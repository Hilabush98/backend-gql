package auth

import (
	"backend-gql/graph/model"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(payload model.PayloadData, accessExp, refreshExp time.Duration) (model.AuthPayload, error) {
	access, err := generateToken(payload, "access", accessExp)
	if err != nil {
		return model.AuthPayload{}, err
	}

	refresh, err := generateToken(payload, "refresh", refreshExp)
	if err != nil {
		return model.AuthPayload{}, err
	}

	return model.AuthPayload{AccessToken: access, RefreshToken: refresh}, nil
}

func DecodeJWT(tokenStr string) (jwt.MapClaims, error) {
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}
	return claims, nil
}

func RefreshAccessToken(refreshTokenStr string, accessExp, refreshExp time.Duration) (model.AuthPayload, error) {
	claims, err := parseAndVerify(refreshTokenStr)
	if err != nil {
		return model.AuthPayload{}, err
	}

	if claims["token_type"] != "refresh" {
		return model.AuthPayload{}, errors.New("provided token is not a refresh token")
	}

	payload := claimsToPayload(claims)
	return Login(payload, accessExp, refreshExp)
}

// ── internals ────────────────────────────────────────────────────────────────

func generateToken(payload model.PayloadData, tokenType string, expiration time.Duration) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(expiration).Unix(),
		"token_type": tokenType,
		"name":       payload.Name,
		"last_name":  payload.LastName,
		"profiles":   payload.Profiles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func parseAndVerify(tokenStr string) (jwt.MapClaims, error) {
	secret, err := jwtSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}
	return claims, nil
}

func claimsToPayload(claims jwt.MapClaims) model.PayloadData {
	payload := model.PayloadData{
		Name:     stringClaim(claims, "name"),
		LastName: stringClaim(claims, "last_name"),
	}
	if raw, ok := claims["profiles"].([]any); ok {
		for _, v := range raw {
			if f, ok := v.(float64); ok {
				n := int(f)
				payload.Profiles = append(payload.Profiles, &n)
			}
		}
	}
	return payload
}

func stringClaim(claims jwt.MapClaims, key string) string {
	if v, ok := claims[key].(string); ok {
		return v
	}
	return ""
}

func jwtSecret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return nil, errors.New("JWT_SECRET env var not set")
	}
	return []byte(s), nil
}

// ── password ─────────────────────────────────────────────────────────────────

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
