package envmanager_dtos

type TokenConfig struct {
	SecretKey     string
	Issuer        string
	Audience      string
	TokenDuration int // in minutes
}
