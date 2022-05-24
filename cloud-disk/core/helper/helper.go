package helper

import (
	"cloud-disk/core/define"
	"context"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"path"
	"time"
)

func MD5(s string) string{
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateToken(id int64, identity string, name string)(string, error){
	// generator user's token according to id, identity and name
	uc := define.UserClaim{
		Id: id,
		Identity: identity,
		Name: name,
	}
	// use jwt key to encrypt the claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(define.JwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// AnalyzeToken parsing the token
func AnalyzeToken(token string) (*define.UserClaim, error){
	uc := new(define.UserClaim)
	// use jwt_key to parse claim
	claims, err := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(define.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	// check whether the claim is valid
	if !claims.Valid{
		return uc, errors.New("token is invalid")
	}

	return uc, err
}

// MailSendCode Send Validation Code to Mail
func MailSendCode(dstMail, code string) error{
	e := email.NewEmail()
	e.From = "Get <1426887306@qq.com>"
	e.To = []string{dstMail}

	e.Subject = "Validation Code Sending Test"
	e.HTML = []byte("<h1>Your Validation Code is "+code+"</h1>")
	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "1426887306@qq.com",
		define.MailPassword, "smtp.qq.com"),
		&tls.Config{
			InsecureSkipVerify: true,
			ServerName: "smtp.qq.com",
		})
	if err != nil{
		return err
	}
	return nil
}

func GenerateRandomCode() string{
	s := "1234567890"
	code := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0 ; i < define.CodeLength; i++{
		code += string(s[rand.Intn(len(s))])
	}
	return code
}

func GetUUID()string{
	return uuid.NewV4().String()
}

// CosUpload File Upload to Tencent Cloud COS
func CosUpload(r *http.Request) (string ,error){
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID: define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	file, fileHeader, err := r.FormFile("file")

	key := "cloud-disk/" + GetUUID() + path.Ext(fileHeader.Filename)

	_, err = client.Object.Put(
		context.Background(), key, file, nil,
	)
	if err != nil {
		return "", err
	}
	return define.CosBucket + "/" + key, nil
}