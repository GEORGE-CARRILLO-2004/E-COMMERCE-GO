package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golan/application/ports/out"
	"golan/domain/customer"
)

type jwtProvider struct {
	secretKey []byte
}

func NewJWTProvider(secretKey string) out.TokenProvider {
	return &jwtProvider{secretKey: []byte(secretKey)}
}

func (p *jwtProvider) GenerateToken(ctx context.Context, c *customer.Customer) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_id": c.ID().String(),
		"email":       c.Email().String(),
		"exp":         jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	return token.SignedString(p.secretKey)
}

func (p *jwtProvider) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return p.secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("token inválido: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("token inválido")
	}

	customerID, ok := claims["customer_id"].(string)
	if !ok || customerID == "" {
		return "", fmt.Errorf("customer_id no encontrado en el token")
	}

	return customerID, nil
}
