package tokenstore

import (
	crypto_rand "crypto/rand"
	"crypto/subtle"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"golang.org/x/crypto/bcrypt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	authenticationv1beta1 "k8s.io/client-go/pkg/apis/authentication/v1beta1"
	auth "kope.io/auth/pkg/apis/auth/v1alpha1"
	client "kope.io/auth/pkg/client/clientset_generated/clientset"
	"kope.io/auth/pkg/oauth/session"
)

const bcryptCost = bcrypt.DefaultCost

type APITokenStore struct {
	client    client.Interface
	namespace string

	mutex sync.Mutex
	users map[types.UID]*auth.User
}

var _ Interface = &APITokenStore{}

func NewAPITokenStore(client client.Interface, namespace string) *APITokenStore {
	s := &APITokenStore{
		client:    client,
		namespace: namespace,
		users:     make(map[types.UID]*auth.User),
	}
	return s
}

func (s *APITokenStore) LookupToken(tokenString string) (*authenticationv1beta1.UserInfo, error) {
	// TODO: Cache tokens in memory, to avoid bcrypt?

	items := strings.SplitN(tokenString, "/", 3)
	if len(items) != 3 {
		glog.V(2).Infof("Rejecting token with incorrect number of tokens")
		return nil, nil
	}

	secretBytes, err := base64.URLEncoding.DecodeString(items[2])
	if err != nil {
		glog.V(2).Infof("invalid secret; ignoring token")
		return nil, nil
	}

	user := s.findUserByUid(types.UID(items[0]))
	if user == nil {
		glog.V(2).Infof("user %q not found", items[0])
		return nil, nil
	}

	//TODO: token expiry?
	//TODO: token reuse?

	var token *auth.TokenSpec
	for _, t := range user.Spec.Tokens {
		if t.ID == items[1] {
			token = t
			break
		}
	}

	if token == nil {
		glog.V(2).Infof("token %q not found for user %q", items[1], items[0])
		return nil, nil
	}

	if token.HashedSecret != nil {
		if err := bcrypt.CompareHashAndPassword(token.HashedSecret, secretBytes); err != nil {
			glog.V(2).Infof("invalid token for user %q", items[0])
			return nil, nil
		}
	} else if token.RawSecret != nil {
		if subtle.ConstantTimeCompare(token.RawSecret, secretBytes) != 1 {
			glog.V(2).Infof("invalid token for user %q", items[0])
			return nil, nil
		}
	} else {
		glog.V(2).Infof("no secret set for token %q for user %q", items[1], items[0])
		return nil, nil
	}

	userInfo := &authenticationv1beta1.UserInfo{}
	userInfo.UID = string(user.UID)
	userInfo.Username = user.Spec.Username
	userInfo.Groups = user.Spec.Groups
	return userInfo, nil
}

func (s *APITokenStore) CreateToken(u *auth.User, hashSecret bool) (*auth.TokenSpec, error) {
	objectName := u.Name

	t := &auth.TokenSpec{}
	t.ID = strconv.FormatInt(rand.Int63(), 32)

	// TODO: Check that doesn't already exist?

	secret := make([]byte, 32, 32)
	_, err := crypto_rand.Read(secret)
	if err != nil {
		return nil, fmt.Errorf("error generating random token: %v", err)
	}
	if hashSecret {
		hashedSecret, err := bcrypt.GenerateFromPassword(secret, bcryptCost)
		if err != nil {
			return nil, fmt.Errorf("error doing bcrypt: %v", err)
		}
		t.HashedSecret = hashedSecret
	} else {
		t.RawSecret = secret
	}
	//u, err := s.client.Users(objectNamespace).Get(objectName)
	//if err != nil {
	//	if errors.IsNotFound(err) {
	//		glog.V(2).Infof("user %s/%s not found; will create", objectNamespace, objectName)
	//		u = nil
	//	} else {
	//		return nil, fmt.Errorf("error fetching user %s/%s: %v", objectNamespace, objectName, err)
	//	}
	//}

	//create := false
	//if u == nil {
	//	u = &user.User{}
	//	u.Metadata.Name = uid
	//	u.Metadata.Namespace = objectNamespace
	//
	//	create = true
	//}

	u.Spec.Tokens = append(u.Spec.Tokens, t)

	//if create {
	//	_, err := s.client.Users(objectNamespace).Create(u)
	//	if err != nil {
	//		return nil, fmt.Errorf("error creating user %s/%s: %v", objectNamespace, objectName, err)
	//	}
	//	// TODO: Update directly (vs going through watch)?
	//} else {
	_, err = s.client.AuthV1alpha1().Users(s.namespace).Update(u)
	if err != nil {
		return nil, fmt.Errorf("error updating user %s/%s: %v", s.namespace, objectName, err)
	}
	// TODO: Update directly (vs going through watch)?
	//}

	return t, nil
}

func (s *APITokenStore) FindExistingUser(identity *auth.IdentitySpec) (*auth.User, error) {
	u := s.findUserByProviderInfo(identity)
	return u, nil
}

func (s *APITokenStore) MapToUser(userInfo *session.UserInfo, create bool) (*auth.User, error) {
	// TODO: Check that doesn't already exist?

	identity := &auth.IdentitySpec{
		ProviderID: userInfo.ProviderID,
		ID:         userInfo.ProviderUserID,
	}

	u := s.findUserByProviderInfo(identity)
	if u == nil && create {
		uidBytes := make([]byte, 30, 30)
		_, err := crypto_rand.Read(uidBytes)
		if err != nil {
			return nil, fmt.Errorf("error generating random uid: %v", err)
		}

		u = &auth.User{}
		// TODO: Include a prefix based on the username?
		name := base32.HexEncoding.EncodeToString(uidBytes)
		name = strings.Replace(name, "=", "", -1)
		u.Name = strings.ToLower(name)
		u.Namespace = s.namespace

		u.Spec.Username = userInfo.Email

		u.Spec.Identities = []auth.IdentitySpec{*identity}
		u, err = s.client.AuthV1alpha1().Users(s.namespace).Create(u)
		if err != nil {
			return nil, fmt.Errorf("error creating user %s/%s: %v", s.namespace, u.Name, err)
		}
		// TODO: Update directly (vs going through watch)?
	}

	// TODO: Gather extra information from the merge

	return u, nil
}

func (s *APITokenStore) FindUserByUID(uid string) (*auth.User, error) {
	user := s.findUserByUid(types.UID(uid))
	if user == nil {
		return nil, nil
	}
	return user, nil
}

func (s *APITokenStore) findUserByUid(uid types.UID) *auth.User {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user := s.users[uid]
	return user
}

func (s *APITokenStore) findUserByProviderInfo(id *auth.IdentitySpec) *auth.User {
	if id.ProviderID == "" {
		glog.Fatalf("ProviderID not set")
	}
	if id.ID == "" {
		glog.Fatalf("ID not set")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// TODO: Build index, if we ever think this is going to be a problem...

	// TODO: Check for duplicates?
	for _, u := range s.users {
		for _, i := range u.Spec.Identities {
			if id.ProviderID == i.ProviderID && id.ID == i.ID {
				return u
			}
		}
	}
	return nil
}

// updateUser processes an update notification for a user
func (c *APITokenStore) processUserUpdate(u *auth.User) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.users[u.UID] = u
}

func (c *APITokenStore) processUserDelete(u *auth.User) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.users, u.UID)
}

// Run starts the secretsWatcher.
func (s *APITokenStore) RunWatch(stopCh <-chan struct{}) {
	runOnce := func() (bool, error) {
		var listOpts metav1.ListOptions

		// TODO: Filters?

		userList, err := s.client.AuthV1alpha1().Users(s.namespace).List(listOpts)
		if err != nil {
			return false, fmt.Errorf("error listing users: %v", err)
		}
		for i := range userList.Items {
			u := &userList.Items[i]
			glog.Infof("user: %v", u.Spec.Username)
			s.processUserUpdate(u)
		}

		listOpts.Watch = true
		listOpts.ResourceVersion = userList.ResourceVersion
		watcher, err := s.client.AuthV1alpha1().Users(s.namespace).Watch(listOpts)
		if err != nil {
			return false, fmt.Errorf("error watching users: %v", err)
		}
		ch := watcher.ResultChan()
		for {
			select {
			case <-stopCh:
				glog.Infof("Got stop signal")
				return true, nil
			case event, ok := <-ch:
				if !ok {
					glog.Infof("user watch channel closed")
					return false, nil
				}

				u := event.Object.(*auth.User)
				glog.V(4).Infof("user changed: %s %v", event.Type, u.Spec.Username)

				switch event.Type {
				case watch.Added, watch.Modified:
					s.processUserUpdate(u)

				case watch.Deleted:
					s.processUserDelete(u)
				}
			}
		}
	}

	for {
		stop, err := runOnce()
		if stop {
			return
		}

		if err != nil {
			glog.Warningf("Unexpected error in secret watch, will retry: %v", err)
			time.Sleep(10 * time.Second)
		}
	}
}

func (c *APITokenStore) Run(stopCh <-chan struct{}) {
	//c.RunPolling(stopCh)
	//glog.Warning("Using polling while watches of TPRs are broken")
	c.RunWatch(stopCh)
}

// Run starts the secretsWatcher.
func (s *APITokenStore) RunPolling(stopCh <-chan struct{}) {
	runOnce := func() error {
		var listOpts metav1.ListOptions

		// TODO: Filters?

		glog.V(2).Infof("polling users")
		userList, err := s.client.AuthV1alpha1().Users(s.namespace).List(listOpts)
		if err != nil {
			return fmt.Errorf("error listing users: %v", err)
		}
		for i := range userList.Items {
			u := &userList.Items[i]
			glog.Infof("user: %v", u.Spec.Username)
			s.processUserUpdate(u)
		}
		return nil
	}

	for {
		err := runOnce()
		if err != nil {
			glog.Warningf("Unexpected error in user poll, will retry: %v", err)
		}

		timer := time.NewTimer(10 * time.Second)
		for {
			select {
			case <-stopCh:
				glog.Infof("Got stop signal")
				return

			case <-timer.C:
			}
		}
	}
}
