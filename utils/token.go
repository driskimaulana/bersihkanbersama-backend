package utils

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
)

func GenerateToken(userId primitive.ObjectID) (string, error) {

	//tokenLifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	//
	//if err != nil {
	//	return "", err
	//}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}
