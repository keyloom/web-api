package core

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	token_dtos "github.com/keyloom/web-api/dtos/token"
)

type TokenService struct{}

func (s *TokenService) GenerateToken(
	sub string,
) (token_dtos.AccessTokenResponse, error) {
	config, err := (&EnvManager{}).GetTokenConfig()
	if err != nil {
		return token_dtos.AccessTokenResponse{}, err
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
		return token_dtos.AccessTokenResponse{}, err
	}

	return token_dtos.AccessTokenResponse{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresAt:   expirationTime.Unix(),
	}, nil
}

func (s *TokenService) ValidateToken(tokenString string) (*jwt.Token, error) {
	config, err := (&EnvManager{}).GetTokenConfig()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	return token, err
}
