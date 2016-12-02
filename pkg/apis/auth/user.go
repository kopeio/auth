package auth

type User struct {
	unversioned.TypeMeta `json:",inline"`
	ObjectMeta           `json:"metadata,omitempty"`

	Spec UserSpec `json:"spec,omitempty"`
}

type UserList struct {
	Items []User `json:"items"`
}

type UserSpec struct {
}
