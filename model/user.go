package model

type User struct {
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	PhotoURL    string `json:"photoURL"`
	ProviderId  string `json:"providerId"`
	UID         string `json:"uid"`
}
