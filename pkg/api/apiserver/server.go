package apiserver

import (
	"k8s.io/apimachinery/pkg/apimachinery/announced"
	"k8s.io/apimachinery/pkg/apimachinery/registered"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"

	registry_authconfiguration "kope.io/auth/pkg/api/registry/authconfiguration"
	registry_authprovider "kope.io/auth/pkg/api/registry/authprovider"
	registry_user "kope.io/auth/pkg/api/registry/user"
	"kope.io/auth/pkg/apis/auth"
	authinstall "kope.io/auth/pkg/apis/auth/install"
	authv1alpha1 "kope.io/auth/pkg/apis/auth/v1alpha1"
	"kope.io/auth/pkg/apis/componentconfig"
	componentconfiginstall "kope.io/auth/pkg/apis/componentconfig/install"
	componentconfigv1alpha1 "kope.io/auth/pkg/apis/componentconfig/v1alpha1"
)

var (
	groupFactoryRegistry = make(announced.APIGroupFactoryRegistry)
	registry             = registered.NewOrDie("")
	Scheme               = runtime.NewScheme()
	Codecs               = serializer.NewCodecFactory(Scheme)
)

func init() {
	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})

	authinstall.Install(groupFactoryRegistry, registry, Scheme)
	componentconfiginstall.Install(groupFactoryRegistry, registry, Scheme)

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
		&metav1.WatchEvent{},
	)
}

type Config struct {
	GenericConfig *genericapiserver.Config
}

// AuthServer contains state for a Kubernetes cluster master/api server.
type AuthServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type completedConfig struct {
	*Config
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *Config) Complete() completedConfig {
	c.GenericConfig.Complete()

	c.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}

	return completedConfig{c}
}

// SkipComplete provides a way to construct a server instance without config completion.
func (c *Config) SkipComplete() completedConfig {
	return completedConfig{c}
}

// New returns a new instance of AuthServer from the given config.
func (c completedConfig) New() (*AuthServer, error) {
	genericServer, err := c.Config.GenericConfig.SkipComplete().New() // completion is done in Complete, no need for a second time
	if err != nil {
		return nil, err
	}

	s := &AuthServer{
		GenericAPIServer: genericServer,
	}

	{
		apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(auth.GroupName, registry, Scheme, metav1.ParameterCodec, Codecs)
		apiGroupInfo.GroupMeta.GroupVersion = authv1alpha1.SchemeGroupVersion
		v1alpha1storage := map[string]rest.Storage{}
		v1alpha1storage["users"] = registry_user.NewREST(Scheme, c.GenericConfig.RESTOptionsGetter)
		apiGroupInfo.VersionedResourcesStorageMap["v1alpha1"] = v1alpha1storage

		if err := s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
			return nil, err
		}
	}

	{
		apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(componentconfig.GroupName, registry, Scheme, metav1.ParameterCodec, Codecs)
		apiGroupInfo.GroupMeta.GroupVersion = componentconfigv1alpha1.SchemeGroupVersion

		v1alpha1storage := map[string]rest.Storage{}
		v1alpha1storage["authconfigurations"] = registry_authconfiguration.NewREST(Scheme, c.GenericConfig.RESTOptionsGetter)
		v1alpha1storage["authproviders"] = registry_authprovider.NewREST(Scheme, c.GenericConfig.RESTOptionsGetter)
		apiGroupInfo.VersionedResourcesStorageMap["v1alpha1"] = v1alpha1storage

		if err := s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
			return nil, err
		}
	}

	return s, nil
}
