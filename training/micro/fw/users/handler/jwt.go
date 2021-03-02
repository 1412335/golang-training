package handler

import (
	"fmt"
	"time"

	"fw/configs"

	"github.com/dgrijalva/jwt-go"
)

type JWTManager struct {
	secretKey     []byte
	tokenDuration time.Duration
	issuer        string
}

func NewJWTManager(config *configs.JWT) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(config.SecretKey),
		tokenDuration: config.Duration,
		issuer:        config.Issuer,
	}
}

type UserClaims struct {
	jwt.StandardClaims
	User *User
}

func (manager *JWTManager) Generate(user *User) (string, error) {
	current := time.Now()
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: current.Add(manager.tokenDuration).Unix(),
			IssuedAt:  current.Unix(),
			Issuer:    manager.issuer,
		},
		User: user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(manager.secretKey)
}

func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return manager.secretKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %+v", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok || claims.Issuer != manager.issuer {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
