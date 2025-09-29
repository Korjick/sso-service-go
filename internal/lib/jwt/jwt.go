package jwt

import (
	"github.com/Korjick/sso-service-go/internal/domain/model"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user model.User, app model.App, ttl time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["app_id"] = app.ID
	claims["exp"] = time.Now().Add(ttl).Unix()
	claims["email"] = user.Email
	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
