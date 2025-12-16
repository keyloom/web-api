package core

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct{}

func (s *TokenService) GenerateToken(
	sub string,
) (string, error) {
	config, err := (&EnvManager{}).GetTokenConfig()
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(time.Duration(config.TokenDuration) * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"iss": config.Issuer,
		"aud": config.Audience,
		"exp": expirationTime.Unix(),
	})

	signedToken, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
