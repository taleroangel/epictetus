package security

import (
	"errors"
	"fmt"
	"time"

	"dev/taleroangel/epictetus/internal/entities"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Calculate password hash with BCrypt algorithm for secure storage
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check if password and hash match
func CheckPassword(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// JWT token user claim
type UserClaim struct {
	jwt.RegisteredClaims
	Name string
	Sudo bool
}

func NewUserClaim(user entities.User, exp time.Duration) *UserClaim {
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

func (uc UserClaim) NewUser() *entities.User {
	return &entities.User{
		User: uc.Subject,
		Name: uc.Name,
		Sudo: uc.Sudo,
	}
}

// Generate an authentication token JWT
func GenerateToken(secretKey []byte, user entities.User) (string, error) {
	// Generate the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, NewUserClaim(user, 8*time.Hour))
	// Sign it with HMAC
	return token.SignedString(secretKey)
}

// Validate generated token, returns the username or error if token is no longer valid
func ValidateToken(secretKey []byte, token string) (*entities.User, error) {
	// Parse the claims
	var claim UserClaim

	// Parse the token
	tkn, err := jwt.ParseWithClaims(token, &claim, func(t *jwt.Token) (interface{}, error) {
		// Check if method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method: %v", t.Header["alg"])
		}

		// Return the signing key
		return secretKey, nil
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
