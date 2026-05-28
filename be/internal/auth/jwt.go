package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const defaultSecret = "hani-dev-secret-change-me"

type Claims struct {
	UserID    int    `json:"uid"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      int    `json:"role"`
	Plan      string `json:"plan"`
	Guest     bool   `json:"guest"`
	GuestID   string `json:"gid,omitempty"`
	SessionID string `json:"sid,omitempty"`
	jwt.RegisteredClaims
}

func jwtSecret() []byte {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return []byte(s)
	}
	return []byte(defaultSecret)
}

func accessTTL() time.Duration {
	if d := os.Getenv("JWT_ACCESS_TTL"); d != "" {
		if parsed, err := time.ParseDuration(d); err == nil {
			return parsed
		}
	}
	return 15 * time.Minute
}

func GenerateAccessToken(userID int, email, name string, role int, plan string, guest bool, guestID, sessionID string) (string, time.Time, error) {
	exp := time.Now().Add(accessTTL())
	if plan == "" {
		if guest {
			plan = "guest"
		} else {
			plan = "free"
		}
	}
	claims := Claims{
		UserID:    userID,
		Email:     email,
		Name:      name,
		Role:      role,
		Plan:      plan,
		Guest:     guest,
		GuestID:   guestID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret())
	return signed, exp, err
}

// GenerateToken issues a legacy 7-day token (avoid for new clients).
func GenerateToken(userID int, email, name string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Name:   name,
		Plan:   "free",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
