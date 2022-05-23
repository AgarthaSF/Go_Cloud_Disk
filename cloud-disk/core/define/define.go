package define

import "github.com/golang-jwt/jwt/v4"

type UserClaim struct {
	Id       int64
	Identity string
	Name     string
	jwt.StandardClaims
}

var JwtKey = "cloud-disk-key"
var MailPassword = "wtuuqmkmrhqofhgj"

// CodeLength 验证码长度
var CodeLength = 6

// CodeExpirationTime 验证码过期时间（s）
var CodeExpirationTime = 300

// 腾讯云COS密钥
var TencentSecretID = "AKIDe8NM4VgaDZfU2JK4XjeIfkY6NN6YjdGd"

var TencentSecretKey = "7yYrugBPN4UdZEcjZ9RSChuqVu2tSknZ"

// 腾讯云存储桶路径
var CosBucket = "https://cloud-disk-1304032890.cos.ap-nanjing.myqcloud.com"