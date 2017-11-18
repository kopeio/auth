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

type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec UserSpec `json:"spec"`
}

type UserSpec struct {
	// The name that uniquely identifies this user among all active users.
	// +optional
	Username string `json:"username,omitempty"`

	// The names of groups this user is a part of.
	// +optional
	Groups []string `json:"groups,omitempty"`

	Tokens []*TokenSpec `json:"tokens,omitempty"`

	Identities []IdentitySpec `json:"identities,omitempty"`
}

type TokenSpec struct {
	ID           string `json:"id,omitempty"`
	HashedSecret []byte `json:"hashedSecret,omitempty"`
	RawSecret    []byte `json:"rawSecret,omitempty"`
}

type IdentitySpec struct {
	ID         string `json:"id,omitempty"`
	ProviderID string `json:"provider,omitempty"`

	Username string `json:"username,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []User `json:"items"`
}
