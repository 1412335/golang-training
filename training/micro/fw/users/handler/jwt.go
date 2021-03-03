package handler

import (
	"fmt"
	"time"

	"fw/users/config"

	"github.com/dgrijalva/jwt-go"
)

type JWTManager struct {
	secretKey     []byte
	tokenDuration time.Duration
	issuer        string
}

func NewJWTManager(cfg *config.JWT) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(cfg.SecretKey),
		tokenDuration: cfg.Duration,
		issuer:        cfg.Issuer,
	}
}

type UserClaims struct {
	jwt.StandardClaims
	ID        string
	FirstName string
	LastName  string
	Email     string
}

func (manager *JWTManager) Generate(user *User) (string, error) {
	current := time.Now()
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: current.Add(manager.tokenDuration).Unix(),
			IssuedAt:  current.Unix(),
			Issuer:    manager.issuer,
		},
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
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
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %+v", err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok || claims.Issuer != manager.issuer {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
