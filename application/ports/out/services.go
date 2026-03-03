package out

import (
	"context"
	"golan/domain/customer"
)

type TokenProvider interface {
	GenerateToken(ctx context.Context, c *customer.Customer) (string, error)
	ValidateToken(ctx context.Context, tokenString string) (string, error)
}
