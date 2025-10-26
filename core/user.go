package core

import "time"

type Authentication struct {
	Provider string `json:"provider"`
	Subject  string `json:"provider_user_id"`
	Email    string `json:"email"`
}

type User struct {
	ID          uint            `json:"id"`
	PrivateData *map[string]any `json:"private_data"`
	PublicData  *map[string]any `json:"public_data"`
	LastLogin   time.Time       `json:"last_login"`

	Authentications []Authentication `json:"authentications"`

	// these just shortcuts: the first valid ones from the authentications

	Email string `json:"email"`
}
