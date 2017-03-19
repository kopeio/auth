package keystore

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"kope.io/auth/pkg/keystore/pb"
)

type KubernetesKeyStore struct {
	client kubernetes.Interface

	namespace string
	name      string

	mutex            sync.Mutex
	sharedSecretSets map[string]*KubernetesSharedSecretSet
	resourceVersion  int64
}

var _ KeyStore = &KubernetesKeyStore{}

type KubernetesSharedSecretSet struct {
	keystore  *KubernetesKeyStore
	generator SharedSecretGenerator

	name     string
	versions map[int32]*KubernetesSharedSecret
	active   int32
}

var _ SharedSecretSet = &KubernetesSharedSecretSet{}

type KubernetesSharedSecret struct {
	data pb.Key
}

var _ SharedSecret = &KubernetesSharedSecret{}

func NewKubernetesKeyStore(client kubernetes.Interface, namespace string, name string) (*KubernetesKeyStore, error) {
	s := &KubernetesKeyStore{
		client:    client,
		namespace: namespace,
		name:      name,
	}
	return s, nil
}

func (k *KubernetesKeyStore) mutateSecret(mutator func(secret *v1.Secret) error) error {
	secret, err := k.client.CoreV1().Secrets(k.namespace).Get(k.name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			glog.V(2).Infof("secret %s/%s not found; will create", k.namespace, k.name)
			secret = nil
		} else {
			return fmt.Errorf("error fetching secret %s/%s: %v", k.namespace, k.name, err)
		}
	}

	create := false
	if secret == nil {
		secret = &v1.Secret{}
		secret.Name = k.name
		secret.Namespace = k.namespace
		create = true
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	err = mutator(secret)
	if err != nil {
		return err
	}

	if create {
		created, err := k.client.CoreV1().Secrets(k.namespace).Create(secret)
		if err != nil {
			// TODO: Handle concurrent create - retry?
			return fmt.Errorf("error creating secret %s/%s: %v", k.namespace, k.name, err)
		}

		k.updateSecret(created)
	} else {
		// TODO: Make sure this is a conditional update
		// https://github.com/kubernetes/client-go/issues/150
		updated, err := k.client.CoreV1().Secrets(k.namespace).Update(secret)
		if err != nil {
			// TODO: Handle condition update - retry?
			return fmt.Errorf("error updating secret %s/%s: %v", k.namespace, k.name, err)
		}

		k.updateSecret(updated)
	}

	// TODO: Update directly, before watch returns?

	return nil
}

func (k *KubernetesKeyStore) ensureSharedSecretInKubernetes(name string, generator SharedSecretGenerator) error {
	err := k.mutateSecret(func(secret *v1.Secret) error {
		sharedSecretSets := k.decodeSecret(secret)
		sharedSecretSet := sharedSecretSets[name]
		if sharedSecretSet == nil {
			sharedSecretSet = &KubernetesSharedSecretSet{
				keystore:  k,
				generator: generator,
				name:      name,
				versions:  make(map[int32]*KubernetesSharedSecret),
			}
			sharedSecretSets[name] = sharedSecretSet
		}

		sharedSecret := sharedSecretSet.versions[sharedSecretSet.active]
		if sharedSecret == nil {
			maxId := int32(0)
			for id := range sharedSecretSet.versions {
				if id > maxId {
					maxId = id
				}
			}

			secretData, err := generator()
			if err != nil {
				return fmt.Errorf("error generating secret: %v", err)
			}

			sharedSecret := &KubernetesSharedSecret{
				data: pb.Key{
					Id:     maxId + 1,
					Secret: secretData,
				},
			}

			sharedSecretSet.active = sharedSecret.data.Id
			sharedSecretSet.versions[sharedSecret.data.Id] = sharedSecret
		}

		keyPrefix := "secret." + sharedSecretSet.name + "."
		for k := range secret.Data {
			if strings.HasPrefix(k, keyPrefix) {
				delete(secret.Data, k)
			}
		}

		keySet := &pb.KeySet{}
		keySet.ActiveId = sharedSecretSet.active

		for _, sharedSecret := range sharedSecretSet.versions {
			keySet.Keys = append(keySet.Keys, &sharedSecret.data)
		}

		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}
		bytes, err := proto.Marshal(keySet)
		if err != nil {
			return fmt.Errorf("error serializing keyset: %v", err)
		}

		secret.Data["secret."+name] = bytes

		return nil
	})
	return err
}

func int32ToString(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

func (k *KubernetesKeyStore) EnsureSharedSecretSet(name string, generator SharedSecretGenerator) (SharedSecretSet, error) {
	sharedSecretSet := k.sharedSecretSets[name]
	if sharedSecretSet == nil {
		err := k.ensureSharedSecretInKubernetes(name, generator)
		if err != nil {
			return nil, fmt.Errorf("error creating shared secret: %v", err)
		}

		sharedSecretSet = k.sharedSecretSets[name]
		if sharedSecretSet == nil {
			return nil, fmt.Errorf("created secret set was not found")
		}
	}

	if sharedSecretSet.generator == nil {
		sharedSecretSet.generator = generator
	}

	return sharedSecretSet, nil
}

func (s *KubernetesSharedSecretSet) EnsureSharedSecret() (SharedSecret, error) {
	sharedSecret := s.versions[s.active]
	if sharedSecret == nil {
		err := s.keystore.ensureSharedSecretInKubernetes(s.name, s.generator)
		if err != nil {
			return nil, fmt.Errorf("error creating shared secret: %v", err)
		}

		sharedSecret = s.versions[s.active]
		if sharedSecret == nil {
			return nil, fmt.Errorf("created secret was not found")
		}
	}

	return sharedSecret, nil
}

func (s *KubernetesSharedSecret) SecretData() []byte {
	// TODO: Defensive copy?
	return s.data.Secret
}

func (s *KubernetesKeyStore) decodeSecret(secret *v1.Secret) map[string]*KubernetesSharedSecretSet {
	sharedSecretSets := make(map[string]*KubernetesSharedSecretSet)
	for k, v := range secret.Data {
		tokens := strings.Split(k, ".")

		// secret.<name>=<value>
		if len(tokens) == 2 && tokens[0] == "secret" {
			name := tokens[1]
			sharedSecretSet := &KubernetesSharedSecretSet{
				keystore: s,
				name:     name,
				versions: make(map[int32]*KubernetesSharedSecret),
			}
			data := &pb.KeySet{}
			err := proto.Unmarshal(v, data)
			if err != nil {
				glog.Warningf("error parsing key %v", k)
				continue
			}

			sharedSecretSet.active = data.ActiveId
			for _, key := range data.Keys {
				sharedSecretSet.versions[key.Id] = &KubernetesSharedSecret{
					data: *key,
				}
			}

			sharedSecretSets[name] = sharedSecretSet
		} else {
			glog.Warningf("ignoring unrecognized key %v", k)
		}
	}

	return sharedSecretSets
}

// updateSecret parses and updates the specified secret
func (k *KubernetesKeyStore) updateSecret(secret *v1.Secret) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	resourceVersion, err := strconv.ParseInt(secret.ObjectMeta.ResourceVersion, 10, 64)
	if err != nil {
		glog.Warningf("Unable to parse ResourceVersion=%q", secret.ObjectMeta.ResourceVersion)
	} else if resourceVersion <= k.resourceVersion {
		glog.V(2).Infof("Ignoring out of sequence secret update: %d vs %d", resourceVersion, k.resourceVersion)
		return
	}

	sharedSecretSets := k.decodeSecret(secret)
	k.sharedSecretSets = sharedSecretSets

	k.resourceVersion = resourceVersion
}

func (k *KubernetesKeyStore) deleteSecret(resourceVersionString string) {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	resourceVersion, err := strconv.ParseInt(resourceVersionString, 10, 64)
	if err != nil {
		glog.Warningf("Unable to parse ResourceVersion=%q", resourceVersionString)
	} else if resourceVersion <= k.resourceVersion {
		glog.V(2).Infof("Ignoring out of sequence secret update: %d vs %d", resourceVersion, k.resourceVersion)
		return
	}

	sharedSecretSets := make(map[string]*KubernetesSharedSecretSet)
	k.sharedSecretSets = sharedSecretSets

	k.resourceVersion = resourceVersion
}

// Run starts the secretsWatcher.
func (c *KubernetesKeyStore) Run(stopCh <-chan struct{}) {
	runOnce := func() (bool, error) {
		var listOpts metav1.ListOptions

		// TODO: How to watch a single object by name?
		// https://github.com/kubernetes/kubernetes/issues/43299

		secretList, err := c.client.CoreV1().Secrets(c.namespace).List(listOpts)
		if err != nil {
			return false, fmt.Errorf("error watching secrets: %v", err)
		}

		for i := range secretList.Items {
			if secretList.Items[i].Name != c.name {
				continue
			}
			c.updateSecret(&secretList.Items[i])
			// TODO: If this is a multi-item scan, we need to delete any items not present
		}

		listOpts.Watch = true
		listOpts.ResourceVersion = secretList.ResourceVersion
		watcher, err := c.client.CoreV1().Secrets(c.namespace).Watch(listOpts)
		if err != nil {
			return false, fmt.Errorf("error watching secrets: %v", err)
		}
		ch := watcher.ResultChan()
		for {
			select {
			case <-stopCh:
				glog.Infof("Got stop signal")
				return true, nil
			case event, ok := <-ch:
				if !ok {
					glog.Infof("secret watch channel closed")
					return false, nil
				}

				secret := event.Object.(*v1.Secret)
				if secret.Name == c.name {
					glog.V(4).Infof("secret changed: %s %v", event.Type, secret.Name)

					switch event.Type {
					case watch.Added, watch.Modified:
						c.updateSecret(secret)

					case watch.Deleted:
						c.deleteSecret(secret.ResourceVersion)
					}
				} else {
					glog.V(4).Infof("ignoring secret change with wrong name: %s %v", event.Type, secret.Name)
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
