/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TODO(authprovider-q): Is the Auth in AuthConfiguration redundant?
type AuthConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	GenerateKubeconfig GenerateKubeconfig `json:"generateKubeconfig,omitempty"`
}

type GenerateKubeconfig struct {
	Server string `json:"server,omitempty"`
	Name   string `json:"name,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AuthConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []AuthConfiguration `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AuthProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// Description is a human-friendly name
	Description string `json:"description,omitempty"`

	OAuthConfig OAuthConfig `json:"oAuthConfig,omitempty"`

	// Email addresses that are allowed to register using this provider
	PermitEmails []string `json:"permitEmails,omitempty"`
}

type OAuthConfig struct {
	ClientID string `json:"clientID,omitempty"`

	// TODO(authprovider-q): What do we do about secrets?  We presumably don't want this secret
	// in the configmap, because that might have a fairly permissive RBAC role.  But do we want to
	// do a layerable configuration?  Keep the secret in a second configuration object?  Have the
	// name of the secret here, and just runtime error until the secret is loaded?

	// ClientSecret is the OAuth secret
	ClientSecret string `json:"clientSecret,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AuthProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []AuthProvider `json:"items"`
}
