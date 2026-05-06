package auth
import(
	"crypto/rand"
	"encoding/hex"
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashed_password, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err 
	}

	return hashed_password, nil 
}

func CheckHashedPassword(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	return match, err
}

func MakeRefreshToken() string {
	encoded_str := make([]byte, 32)
	rand.Read(encoded_str)

	str := hex.EncodeToString(encoded_str)

	return str
}