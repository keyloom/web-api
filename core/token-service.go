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

func (s *TokenService) ValidateToken(tokenString string) (*token_dtos.JWTPayload, error) {
	config, err := (&EnvManager{}).GetTokenConfig()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	var payload token_dtos.JWTPayload
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	payload.Sub = claims["sub"].(string)
	payload.Iss = claims["iss"].(string)
	payload.Aud = claims["aud"].(string)
	payload.Exp = int64(claims["exp"].(float64))

	payload.JWTHeader.Alg = token.Header["alg"].(string)
	payload.JWTHeader.Typ = token.Header["typ"].(string)

	return &payload, nil
}
