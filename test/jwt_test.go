package test

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"testing"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

var myKey = []byte("project - key")

//生成token
func TestGenerateToken(t *testing.T) {
	Userclaim := &UserClaims{
		Identity:       "user1",
		Name:           "luke",
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Userclaim)
	signedString, err := token.SignedString(myKey)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(signedString)
}

//解析token
func TestAnalyseToken(t *testing.T) {
	signedString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXIxIiwibmFtZSI6Imx1a2UifQ.G0SdFob3lJz0gBY_tlJYlj2e2kSr4lDGkZUcn51WNRU"
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(signedString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		log.Fatalln(err)
		return
	}
	if claims.Valid {
		fmt.Println(userClaim)
	}
}
