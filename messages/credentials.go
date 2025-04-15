package messages

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

type CredentialDto struct {
	Username string `json:"username"`
	Secret   string `json:"secret"`
}

func NewCredentialDto(username, secret string) *CredentialDto {
	return &CredentialDto{
		Username: username,
		Secret:   secret,
	}
}
