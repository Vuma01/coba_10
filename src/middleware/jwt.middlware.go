package middleware

import (
	"coba_01/src/config"
	"coba_01/src/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

var JwtKey = []byte(os.Getenv(config.ENV_SECRET_KEY))

func Authenticate(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")

	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}

	claims := &models.CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Login Tidak Valdi"})
		c.Abort()
		return
	}

	c.Set("username", claims.Username)
	c.Next()
}
