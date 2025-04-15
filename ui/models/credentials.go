package models

// what is sent to the /ws endpoint
// to login to the chat room
type Credentials struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Secret   string `json:"secret"`
}

func NewCredentials(address, username, secret string) *Credentials {
	return &Credentials{
		Address:  address,
		Username: username,
		Secret:   secret,
	}
}
