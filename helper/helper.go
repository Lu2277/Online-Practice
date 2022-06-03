package helper

import (
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

// GetMd5 生成Md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

var myKey = []byte("project - key")

// GenerateToken 生成token
func GenerateToken(identity, name string) (string, error) {
	Userclaim := &UserClaims{
		Identity:       identity,
		Name:           name,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Userclaim)
	signedString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	//fmt.Println(signedString)
	return signedString, nil
}

// AnalyseToken 解析token
func AnalyseToken(signedString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(signedString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("anlyse Token Error:%v", err)
	}
	//fmt.Println(userClaim)
	return userClaim, nil
}
