package auth

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient=true

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

type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []User `json:"items"`
}

//// Required to satisfy Object interface
//func (e *User) GetObjectKind() schema.ObjectKind {
//	return &e.TypeMeta
//}

//// Required to satisfy ObjectMetaAccessor interface
//func (e *User) GetObjectMeta() metav1.Object {
//	return &e.Metadata
//}
//
//// Required to satisfy Object interface
//func (el *UserList) GetObjectKind() schema.ObjectKind {
//	return &el.TypeMeta
//}
//
//// Required to satisfy ListMetaAccessor interface
//func (el *UserList) GetListMeta() metav1.List {
//	return &el.ListMeta
//}

//// The code below is used only to work around a known problem with third-party
//// resources and ugorji. If/when these issues are resolved, the code below
//// should no longer be required.
//
//type UserListCopy UserList
//type UserCopy User
//
//func (e *User) UnmarshalJSON(data []byte) error {
//	tmp := UserCopy{}
//	err := json.Unmarshal(data, &tmp)
//	if err != nil {
//		return err
//	}
//	tmp2 := User(tmp)
//	*e = tmp2
//	return nil
//}
//
//func (el *UserList) UnmarshalJSON(data []byte) error {
//	tmp := UserListCopy{}
//	err := json.Unmarshal(data, &tmp)
//	if err != nil {
//		return err
//	}
//	tmp2 := UserList(tmp)
//	*el = tmp2
//	return nil
//}
