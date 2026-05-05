package auth 
import (
	"fmt"
	"time"
	"strings"
	"net/http"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	jwt.RegisteredClaims
}

func MakeJWT(userID uuid.UUID, tokenSecretKey string, expiresIn time.Duration) (string, error) {
	claim := MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userID.String(),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claim)

	return token.SignedString([]byte(tokenSecretKey))
}

func ValidateJWT(tokenString, tokenSecretKey string) (uuid.UUID, error) {
	claim := MyClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claim, func(tkn *jwt.Token) (interface{}, error){
		if tkn.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("Token signing method mismatch")
		}
		return []byte(tokenSecretKey), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	user_id_str, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err 
	} 

	user_id, err := uuid.Parse(user_id_str)
	if err != nil {
		return uuid.Nil, err 
	}

	return user_id, nil 
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")

	if token == "" {
		return "", fmt.Errorf("No authorization information found")
	}

	trimmed_token := strings.TrimPrefix(token, "Bearer ")

	tokenString := strings.TrimSpace(trimmed_token)

	return tokenString, nil
}