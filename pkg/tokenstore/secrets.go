package tokenstore

//
//import (
//	authenticationv1beta1 "k8s.io/client-go/pkg/apis/authentication/v1beta1"
//	"k8s.io/client-go/kubernetes"
//	"github.com/golang/glog"
//	"fmt"
//	"k8s.io/client-go/pkg/api/v1"
//	"k8s.io/client-go/pkg/watch"
//	"time"
//	"sync"
//	"strings"
//	"golang.org/x/crypto/bcrypt"
//	"encoding/json"
//	"encoding/base64"
//	"strconv"
//	"math/rand"
//	crypto_rand "crypto/rand"
//	"k8s.io/client-go/pkg/api/errors"
//	"kope.io/auth/pkg/apis/auth"
//)
//
//const bcryptCost = bcrypt.DefaultCost
//const userSecretNamePrefix = "user-"
//const secretDataTokenPrefix = "token_"
//
//type Secrets struct {
//	kubeClient *kubernetes.Clientset
//	namespace  string
//
//	mutex      sync.Mutex
//	users      map[string]*user.User
//}
//
//func NewSecrets(kubeClient *kubernetes.Clientset, namespace string) (*Secrets) {
//	s := &Secrets{
//		kubeClient: kubeClient,
//		namespace:namespace,
//		users: make(map[string]*user.User),
//	}
//	return s
//}
//
//func (s*Secrets) LookupToken(tokenString string) (*authenticationv1beta1.UserInfo, error) {
//	items := strings.SplitN(tokenString, "/", 3)
//	if len(items) != 3 {
//		return nil, nil
//	}
//
//	secretBytes, err := base64.URLEncoding.DecodeString(items[2])
//	if err != nil {
//		glog.V(2).Infof("invalid secret; ignoring token")
//		return nil, nil
//	}
//
//	user := s.findUser(items[0])
//	if user == nil {
//		glog.V(2).Infof("user %q not found", items[0])
//		return nil, nil
//	}
//
//	var token *TokenRecord
//	for _, t := range user.Tokens {
//		if t.ID == items[1] {
//			token = t
//			break
//		}
//	}
//
//	if token == nil {
//		glog.V(2).Infof("token %q not found for user %q", items[1], items[0])
//		return nil, nil
//	}
//
//	if err := bcrypt.CompareHashAndPassword(token.Secret, secretBytes); err != nil {
//		glog.V(2).Infof("invalid token for user %q", items[0])
//		return nil, nil
//	}
//
//	return user.User.User, nil
//}
//
//func (s*Secrets) CreateToken(userInfo *authenticationv1beta1.UserInfo) (*TokenInfo, error) {
//	t := &TokenInfo{}
//
//	t.UserID = userInfo.UID
//	if t.UserID == "" {
//		return nil, fmt.Errorf("UID is required")
//	}
//
//	t.TokenID = strconv.FormatInt(rand.Int63(), 32)
//
//	// TODO: Check that doesn't already exist?
//
//	token := make([]byte, 32, 32)
//	_, err := crypto_rand.Read(token)
//	if err != nil {
//		return nil, fmt.Errorf("error generating random token: %v", err)
//	}
//	secretBytes, err := bcrypt.GenerateFromPassword(token, bcryptCost)
//	if err != nil {
//		return nil, fmt.Errorf("error doing bcrypt: %v", err)
//	}
//	t.Secret = secretBytes
//
//	secretName := userSecretNamePrefix + t.UserID
//	secret, err := s.kubeClient.Secrets(s.namespace).Get(secretName)
//	if err != nil {
//		if errors.IsNotFound(err) {
//			glog.V(2).Infof("secret %s/%s not found; will create", s.namespace, secretName)
//			secret = nil
//		} else {
//			return nil, fmt.Errorf("error fetching secret %s/%s: %v", s.namespace, secretName, err)
//		}
//	}
//
//	create := false
//	if secret == nil {
//		secret = &v1.Secret{}
//		secret.Name = secretName
//		secret.Namespace = s.namespace
//		create = true
//	}
//
//	if secret.Data == nil {
//		secret.Data = make(map[string][]byte)
//
//		// TODO: Always update?
//		user := &UserRecord{
//			User: userInfo,
//		}
//		userJson, err := json.Marshal(user)
//		if err != nil {
//			return nil, fmt.Errorf("error building json for user: %v", err)
//		}
//
//		secret.Data["user"] = userJson
//	}
//
//	{
//		record := &TokenRecord{
//			ID: t.TokenID,
//			Secret: t.Secret,
//		}
//		recordJson, err := json.Marshal(record)
//		if err != nil {
//			return nil, fmt.Errorf("error building json for token: %v", err)
//		}
//
//		secret.Data[secretDataTokenPrefix + t.TokenID] = recordJson
//	}
//
//	if create {
//		_, err := s.kubeClient.Secrets(s.namespace).Create(secret)
//		if err != nil {
//			return nil, fmt.Errorf("error creating secret %s/%s: %v", s.namespace, secretName, err)
//		}
//		// TODO: Update directly?
//	} else {
//		_, err := s.kubeClient.Secrets(s.namespace).Update(secret)
//		if err != nil {
//			return nil, fmt.Errorf("error updating secret %s/%s: %v", s.namespace, secretName, err)
//		}
//		// TODO: Update directly?
//	}
//
//	return t, nil
//}
//
//func (s*Secrets) ListTokens(uid string) ([]*TokenInfo, error) {
//	user := s.findUser(uid)
//
//	if user == nil {
//		return nil, nil
//	}
//
//	var tokens []*TokenInfo
//	for _, t := range user.Tokens {
//		tokens = append(tokens, &TokenInfo{
//			UserID: uid,
//			TokenID: t.ID,
//			Secret: t.Secret,
//		})
//	}
//	return tokens, nil
//}
//
//func (s *Secrets) findUser(uid string) (*user) {
//	s.mutex.Lock()
//	defer s.mutex.Unlock()
//
//	user := s.users[uid]
//	return user
//}
//
//// updateSecret parses and updates the specified secret
//func (c *Secrets) updateSecret(secret *v1.Secret) {
//	uid := strings.TrimPrefix(secret.Name, userSecretNamePrefix)
//
//	u := &user{}
//	{
//		v := secret.Data["user"]
//		if v == nil {
//			return
//		}
//
//		err := json.Unmarshal(v, &u.User)
//		if err != nil {
//			glog.Warningf("error decoding user %s/user: %v", uid, err)
//			return
//		}
//	}
//
//	for k, v := range secret.Data {
//		if strings.HasPrefix(k, secretDataTokenPrefix) {
//			data := &TokenRecord{}
//			err := json.Unmarshal(v, data)
//			if err != nil {
//				glog.Warningf("error decoding token %s/%s: %v", uid, k, err)
//				continue
//			}
//			u.Tokens = append(u.Tokens, data)
//		}
//	}
//
//	c.mutex.Lock()
//	defer c.mutex.Unlock()
//
//	c.users[uid] = u
//}
//
//func (c *Secrets) deleteSecret(secret *v1.Secret) {
//	uid := strings.TrimPrefix(secret.Name, userSecretNamePrefix)
//
//	c.mutex.Lock()
//	defer c.mutex.Unlock()
//
//	delete(c.users, uid)
//}
//
//
//// Run starts the secretsWatcher.
//func (c *Secrets) Run(stopCh <-chan struct{}) {
//	runOnce := func() (bool, error) {
//		var listOpts v1.ListOptions
//
//		//listOpts.LabelSelector = labels.Everything()
//		//glog.Warningf("querying without field filter")
//		//listOpts.FieldSelector = fields.Everything()
//
//		c.kubeClient.ThirdPartyResources().L
//		secretList, err := c.kubeClient.Secrets(c.namespace).List(listOpts)
//		if err != nil {
//			return false, fmt.Errorf("error listing secrets: %v", err)
//		}
//		for i := range secretList.Items {
//			secret := &secretList.Items[i]
//			glog.Infof("secret: %v", secret.Name)
//			c.updateSecret(secret)
//		}
//
//		//listOpts.LabelSelector = labels.Everything()
//		//glog.Warningf("querying without field filter")
//		//listOpts.FieldSelector = fields.Everything()
//
//		listOpts.Watch = true
//		listOpts.ResourceVersion = secretList.ResourceVersion
//		watcher, err := c.kubeClient.Secrets(c.namespace).Watch(listOpts)
//		if err != nil {
//			return false, fmt.Errorf("error watching secrets: %v", err)
//		}
//		ch := watcher.ResultChan()
//		for {
//			select {
//			case <-stopCh:
//				glog.Infof("Got stop signal")
//				return true, nil
//			case event, ok := <-ch:
//				if !ok {
//					glog.Infof("secret watch channel closed")
//					return false, nil
//				}
//
//				secret := event.Object.(*v1.Secret)
//				glog.V(4).Infof("secret changed: %s %v", event.Type, secret.Name)
//
//				switch event.Type {
//				case watch.Added, watch.Modified:
//					c.updateSecret(secret)
//
//				case watch.Deleted:
//					c.deleteSecret(secret)
//				}
//			}
//		}
//	}
//
//	for {
//		stop, err := runOnce()
//		if stop {
//			return
//		}
//
//		if err != nil {
//			glog.Warningf("Unexpected error in secret watch, will retry: %v", err)
//			time.Sleep(10 * time.Second)
//		}
//	}
//}
