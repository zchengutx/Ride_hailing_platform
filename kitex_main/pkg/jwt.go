package pkg

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"kitex_main/config"
	"time"
)

func TokenHandler(id string) (string, error) {
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": id,
		"exp":  time.Now().Add(time.Hour * time.Duration(10)).Unix(),
		"iat":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWT_TOKEN))

	return tokenString, err
}

func GetToken(tokenString string) (jwt.MapClaims, string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_TOKEN), nil
	})

	switch {
	case token.Valid:
		fmt.Println("You look nice today")
	case errors.Is(err, jwt.ErrTokenMalformed):
		fmt.Println("That's not even a token")
		return nil, "token错误"
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		// Invalid signature
		fmt.Println("Invalid signature")
		return nil, "token过期"
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		// Token is either expired or not active yet
		fmt.Println("Timing is everything")
		return nil, "token错误"
	default:
		fmt.Println("Couldn't handle this token:", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
		return claims, ""
	} else {
		fmt.Println(err)
	}

	return nil, ""
}
