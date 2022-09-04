package jwt

import "github.com/google/uuid"

type Payload struct {
	Iss  string    `json:"iss"`
	Sub  uuid.UUID `json:"sub"`
	Name string    `json:"name"`
	Iat  int       `json:"iat"`
}
