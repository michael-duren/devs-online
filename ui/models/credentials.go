package models

type Credentials struct {
	Address  string
	Username string
	Secret   string
}

func NewCredentials(address, username, secret string) *Credentials {
	return &Credentials{
		Address:  address,
		Username: username,
		Secret:   secret,
	}
}
