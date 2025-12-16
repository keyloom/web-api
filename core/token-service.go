package core

type TokenService struct{}

func (s *TokenService) GenerateToken(
	sub string,
) (error, string) {
	return nil, "mocked-token"
}
