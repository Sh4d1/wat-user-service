package main

import (
	"log"
	"time"

	pb "github.com/Sh4d1/wat-user-service/proto/user"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	key = []byte("kjasdfjhaejrgadjkfaweiogidavnzbxnmxcbjveoghshlaalgkzdfmngreggkfzg")
)

type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims
}

type Authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

type TokenService struct {
	repo Repository
}

func (s *TokenService) Decode(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		log.Println(err)
		return nil, err
	}
}

func (s *TokenService) Encode(user *pb.User) (string, error) {
	expireToken := time.Now().Add(time.Hour * 72).Unix()

	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "wat.user",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}
