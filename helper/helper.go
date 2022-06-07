package helper

import (
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"
)

// GetMd5 生成Md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var myKey = []byte("project - key")

// GenerateToken 生成token
func GenerateToken(identity string, name string, isAdmin int) (string, error) {
	UserClaim := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        isAdmin,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
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

// CreateRandCode  生成随机验证码
func CreateRandCode() string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < 6; i++ {
		randNum := rand.Intn(10)
		code = code + strconv.Itoa(randNum)
	}
	return code
}

// SendCode 发送验证码
func SendCode(toEmail, code string) error {
	em := email.NewEmail()
	// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
	em.From = "LU <2591134973@qq.com>"
	// 设置 receiver 接收方的邮箱
	em.To = []string{toEmail}
	// 设置主题
	em.Subject = "发送验证码"
	// 简单设置文件发送的内容，暂时设置成纯文本
	em.HTML = []byte("本次验证码为：<b>" + code + "</b>")
	//设置服务器相关的配置
	err := em.Send("smtp.qq.com:25", smtp.PlainAuth("", "2591134973@qq.com", "uhpitgbjrkirdjed", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("send successfully ... ")
	return nil
}

// GetUUID 获取uuid:唯一标识码
func GetUUID() string {
	return uuid.NewV4().String()
}
