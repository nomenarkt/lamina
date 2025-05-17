//go:build tools

package tools

import (
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
)

var _ = jwt.NewWithClaims
