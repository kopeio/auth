package client

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
	"kope.io/auth/pkg/apis/auth"
	"kope.io/auth/pkg/apis/auth/v1alpha1"
)

const (
	SchemaGroup = "auth.kope.io"
	//SchemaName        = "auth." + SchemaGroup
	//APIResources   = "user"
	//SchemaDescription = "Auth schema"
	SchemaVersion = "v1alpha1"

	//UserResourceName = "user." + SchemaGroup

)

var Versions = []v1beta1.APIVersion{
	{Name: SchemaVersion},
}

//var (
//	config *rest.Config
//)

type AuthClientset struct {
	client *rest.RESTClient
}

type Interface interface {
	Users(namespace string) UserInterface
}

func (c *AuthClientset) Users(namespace string) UserInterface {
	return &userInterface{c.client, namespace}
}

func RegisterResource(k8sClient kubernetes.Interface) error {
	// initialize third party resource if it does not exist
	_, err := k8sClient.ExtensionsV1beta1().ThirdPartyResources().Get("user."+SchemaGroup, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			tpr := &v1beta1.ThirdPartyResource{
				ObjectMeta: metav1.ObjectMeta{
					Name: "user." + SchemaGroup,
				},
				Versions:    Versions,
				Description: "Users managed by kopeio-auth",
			}

			_, err := k8sClient.Extensions().ThirdPartyResources().Create(tpr)
			if err != nil {
				return fmt.Errorf("unable to register %q ThirdPartyResource: %v", "user."+SchemaGroup, err)
			}
			glog.Infof("Created %q ThirdPartyResource", "user."+SchemaGroup)
		} else {
			return fmt.Errorf("error querying for %q ThirdPartyResource", "user."+SchemaGroup, err)
		}
	} else {
		// TODO: Update versions?
		glog.Infof("Found %q ThirdPartyResource", "user."+SchemaGroup)
	}

	return nil

}

// NewForConfig creates a new AuthClientset for the given config.
func NewForConfig(config *rest.Config) (*AuthClientset, error) {
	groupversion := schema.GroupVersion{
		Group:   SchemaGroup,
		Version: SchemaVersion,
	}

	config.GroupVersion = &groupversion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: api.Codecs}

	schemeBuilder := runtime.NewSchemeBuilder(
		func(scheme *runtime.Scheme) error {
			scheme.AddKnownTypes(
				groupversion,
				&auth.User{},
				&auth.UserList{},
				&metav1.ListOptions{},
				&metav1.DeleteOptions{},
			)
			return nil
		})
	schemeBuilder.AddToScheme(api.Scheme)

	tprClient, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, err
	}

	return &AuthClientset{
		//k8sClient: k8sClient,
		client: tprClient,
	}, nil
}

type UserInterface interface {
	Create(*v1alpha1.User) (*v1alpha1.User, error)
	Update(*v1alpha1.User) (*v1alpha1.User, error)
	//Delete(name string, options *v1.DeleteOptions) error
	//DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string) (*v1alpha1.User, error)
	List(opts metav1.ListOptions) (*v1alpha1.UserList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	//Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.ComponentStatus, err error)
	//ComponentStatusExpansion
}

type userInterface struct {
	client    *rest.RESTClient
	namespace string
}

var _ UserInterface = &userInterface{}

func (c *userInterface) Get(name string) (*v1alpha1.User, error) {
	u := &v1alpha1.User{}
	err := c.client.Get().Resource("users").Namespace(c.namespace).Name(name).Do().Into(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (c *userInterface) List(opts metav1.ListOptions) (*v1alpha1.UserList, error) {
	l := &v1alpha1.UserList{}
	err := c.client.Get().Resource("users").Namespace(c.namespace).VersionedParams(&opts, api.ParameterCodec).Do().Into(l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (c *userInterface) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	watch, err := c.client.Get().Prefix("watch").Resource("users").Namespace(c.namespace).VersionedParams(&opts, api.ParameterCodec).Watch()
	if err != nil {
		return nil, err
	}
	return watch, nil
}

func (c *userInterface) Create(u *v1alpha1.User) (*v1alpha1.User, error) {
	if u.Metadata.Namespace == "" {
		return nil, fmt.Errorf("Namespace not specified")
	}
	if u.Metadata.Name == "" {
		return nil, fmt.Errorf("Name not specified")
	}
	result := &v1alpha1.User{}
	err := c.client.Post().
		Resource("users").
		Namespace(u.Metadata.Namespace).
		Body(u).
		Do().Into(result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *userInterface) Update(u *v1alpha1.User) (*v1alpha1.User, error) {
	if u.Metadata.Namespace == "" {
		return nil, fmt.Errorf("Namespace not specified")
	}
	if u.Metadata.Name == "" {
		return nil, fmt.Errorf("Name not specified")
	}
	result := &v1alpha1.User{}
	err := c.client.Put().
		Resource("users").
		Namespace(u.Metadata.Namespace).
		Name(u.Metadata.Name).
		Body(u).
		Do().Into(result)

	if err != nil {
		return nil, err
	}
	return result, nil
}
