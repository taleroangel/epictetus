package security

import (
	"context"
	"dev/taleroangel/epictetus/internal/env"
	"dev/taleroangel/epictetus/internal/types"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT token user claim
type UserClaim struct {
	jwt.RegisteredClaims
	Name string
	Sudo bool
}

func NewUserClaim(user types.User, exp time.Duration) *UserClaim {
	return &UserClaim{
		jwt.RegisteredClaims{
			Subject:   user.User,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
		user.Name,
		user.Sudo,
	}
}

func (uc UserClaim) NewUser() *types.User {
	return &types.User{
		User: uc.Subject,
		Name: uc.Name,
		Sudo: uc.Sudo,
	}
}

// Calculate password hash with BCrypt algorithm for secure storage
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Generate an authentication token JWT
func GenerateToken(ctx context.Context, user types.User) (string, error) {
	// Obtain token valid time
	tvf, err := strconv.Atoi(ctx.Value(env.TokenValidFor).(string))
	if err != nil {
		return "", err
	}

	// Generate the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, NewUserClaim(user, time.Duration(tvf)*time.Hour))

	// Sign it with HMAC
	return token.SignedString(ctx.Value(env.SecretKey))
}

// Validate generated token, returns the username or error if token is no longer valid
func ValidateToken(ctx context.Context, token string) (*types.User, error) {
	// Parse the claims
	var claim UserClaim

	// Parse the token
	tkn, err := jwt.ParseWithClaims(token, &claim, func(t *jwt.Token) (interface{}, error) {
		// Check if method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method: %v", t.Header["alg"])
		}

		// Return the signing key
		return ctx.Value(env.SecretKey), nil
	})

	// Check if token could be parsed
	if err != nil {
		return nil, err
	} else if !tkn.Valid {
		return nil, errors.New("token is not valid")
	}

	// Return the claimed user
	return claim.NewUser(), nil
}
