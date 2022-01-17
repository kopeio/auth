package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AuthProviderSpec defines the desired state of AuthProvider
type AuthProviderSpec struct {
	ClientID        string   `json:"clientId,omitempty"`
	ClientSecret    string   `json:"clientSecret,omitempty"`
	PermittedEmails []string `json:"permittedEmails,omitempty"`
}

// AuthProviderStatus defines the observed state of AuthProvider
type AuthProviderStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced

// AuthProvider defines authentication configuration
type AuthProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AuthProviderSpec   `json:"spec,omitempty"`
	Status AuthProviderStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AuthProviderList contains a list of AuthProvider
type AuthProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AuthProvider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AuthProvider{}, &AuthProviderList{})
}

// func (a *AuthProvider) DeepCopyObject() runtime.Object {
// 	return DeepCopyByReflection(a)
// }

// func (a *AuthProviderList) DeepCopyObject() runtime.Object {
// 	return DeepCopyByReflection(a)
// }

// func DeepCopyByReflection(o runtime.Object) runtime.Object {
// 	panic("DeepCopyByReflection not implemented")
// }
