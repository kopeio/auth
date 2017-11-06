package auth

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
