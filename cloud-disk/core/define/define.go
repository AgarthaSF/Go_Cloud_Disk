package define

import "github.com/golang-jwt/jwt/v4"

type UserClaim struct {
	Id       int64
	Identity string
	Name     string
	jwt.StandardClaims
}

var JwtKey = ""
var MailPassword = ""

// CodeLength 验证码长度
var CodeLength = 6

// CodeExpirationTime 验证码过期时间（s）
var CodeExpirationTime = 300

var TencentSecretID = ""

var TencentSecretKey = ""

// CosBucket Tencent COS Bucket Path
var CosBucket = ""

// PageSize default page size
var PageSize = 20

var DateTime = "2006-01-02 15:04:05"

var TokenExpireTime = 3600

var RefreshTokenExpireTime = 7200
