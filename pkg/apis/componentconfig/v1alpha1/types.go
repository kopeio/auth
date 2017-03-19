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

type AuthConfiguration struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`

	Spec AuthConfigurationSpec `json:"spec"`
}

type AuthConfigurationSpec struct {
	AuthProviders []AuthProviderSpec `json:"authProviders,omitempty""`

	GenerateKubeconfig *GenerateKubeconfig `json:"generateKubeconfig,omitempty"`
}

type AuthProviderSpec struct {
	// ID is a system-friendly identifier
	ID string `json:"id,omitempty"`

	// Name is a human-friendly name
	Name string `json:"name,omitempty"`

	OAuthConfig *OAuthConfig `json:"oAuthConfig,omitempty"`

	// Email addresses that are allowed to register using this provider
	PermitEmails []string `json:"permitEmails,omitempty"`
}

type OAuthConfig struct {
	ClientID string `json:"clientID"`

	// TODO(componentconfig-q): What do we do about secrets?  We presumably don't want this secret
	// in the configmap, because that might have a fairly permissive RBAC role.  But do we want to
	// do a layerable configuration?  Keep the secret in a second configuration object?  Have the
	// name of the secret here, and just runtime error until the secret is loaded?

	// ClientSecret is the OAuth secret
	ClientSecret string `json:"clientSecret,omitempty"`
}

type GenerateKubeconfig struct {
	Server string `json:"server,omitempty"`
	Name   string `json:"name,omitempty"`
}
