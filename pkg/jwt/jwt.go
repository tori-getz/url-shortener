package jwt

import "github.com/golang-jwt/jwt/v5"

type JwtPayload struct {
	Email string
}

type Jwt struct {
	secret string
}

func NewJwt(secret string) *Jwt {
	return &Jwt{
		secret: secret,
	}
}

func (j *Jwt) Create(payload JwtPayload) (string, error) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": payload.Email,
	})

	signedToken, err := unsignedToken.SignedString([]byte(j.secret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *Jwt) Parse(token string) (bool, *JwtPayload) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return false, nil
	}

	email := parsedToken.Claims.(jwt.MapClaims)["email"]

	return parsedToken.Valid, &JwtPayload{
		Email: email.(string),
	}
}
