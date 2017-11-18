package install

import (
	"k8s.io/apimachinery/pkg/apimachinery/announced"
	"k8s.io/apimachinery/pkg/apimachinery/registered"
	"k8s.io/apimachinery/pkg/runtime"
	"kope.io/auth/pkg/apis/componentconfig"
	"kope.io/auth/pkg/apis/componentconfig/v1alpha1"
	"k8s.io/apimachinery/pkg/util/sets"
)

//func init() {
//	Install(api.GroupFactoryRegistry, api.Registry, api.Scheme)
//}

// Install registers the API group and adds types to a scheme
func Install(groupFactoryRegistry announced.APIGroupFactoryRegistry, registry *registered.APIRegistrationManager, scheme *runtime.Scheme) {
	if err := announced.NewGroupMetaFactory(
		&announced.GroupMetaFactoryArgs{
			GroupName:              componentconfig.GroupName,
			VersionPreferenceOrder: []string{v1alpha1.SchemeGroupVersion.Version},
			//ImportPrefix:               "kope.io/auth/pkg/apis/componentconfig",
			// RootScopedKinds are resources that are not namespaced.
			RootScopedKinds: sets.NewString(
				"AuthConfiguration",
				"AuthConfigurationList",
				"AuthProvider",
				"AuthProviderList",
			),
			AddInternalObjectsToScheme: componentconfig.AddToScheme,
		},
		announced.VersionToSchemeFunc{
			v1alpha1.SchemeGroupVersion.Version: v1alpha1.AddToScheme,
		},
	).Announce(groupFactoryRegistry).RegisterAndEnable(registry, scheme); err != nil {
		panic(err)
	}
}
