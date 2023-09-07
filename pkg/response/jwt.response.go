package response

import (
	"coba_01/src/config"
	"coba_01/src/middleware"
	"coba_01/src/models"
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv(config.ENV_SECRET_KEY)) // Harap ganti dengan kunci rahasia Anda

func CreateToken(user *models.User, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration)
	log.Println("Secret Key jwt respon:", os.Getenv(config.ENV_SECRET_KEY))

	claims := &models.CustomClaims{
		ID:             user.ID.Hex(),
		Username:       user.Username,
		Email:          user.Email,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &models.CustomClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.JwtKey)
}
