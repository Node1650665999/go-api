package hash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
)

//HashByMd5 实现 md5 加密
func HashByMd5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

//HashBySha1 实现sha1散列加密
func HashBySha1(str string) string {
	hash := sha1.Sum([]byte(str))
	return fmt.Sprintf("%x", hash)
}

//HashBySha256 实现sha256散列加密
func HashBySha256(str string) string {
	hash := sha256.New()
	io.WriteString(hash, str)
	// 16 进制转字符串
	//return fmt.Sprintf("%x", hash.Sum(nil))
	return hex.EncodeToString(hash.Sum(nil))
}

//HmacMd5 使用md5哈希加密算法生成加密串
func HmacMd5(data string, key string) string  {
	hash := hmac.New(md5.New, []byte(key))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}

//HmacSha256 使用sha256哈希加密算法生成加密串
func HmacSha256(data, key string) string {
	hash:= hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}

//HmacSha1 使用sha1哈希加密算法生成加密串
func HmacSha1(data, key string) string {
	hash:= hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum([]byte("")))
}


//EncodeByBcrypt 密码加密
func EncodeByBcrypt(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

//DecodeByBcrypt 密码比对
func DecodeByBcrypt(password string, hashed string) (match bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false, errors.New("密码比对错误！")
	}
	return true, nil
}

var  commonIV = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
var  key      = "L7et45GElm1M4a9g"  //key 的长度必须为：16/24/32 字节

//newAES 创建加密算法 aes
func newAES(key string) (cipher.Block, error)  {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil,err
	}
	return c,nil
}

//EncodeByAES 使用对称加密算法加密字符串 src
func EncodeByAES(src string) string {
	// 实例化算法aes
	c, _ := newAES(key)

	text := []byte(src)

	//加密字符串
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(text))
	cfb.XORKeyStream(ciphertext, text)

	//十六进制转字符串
	return hex.EncodeToString(ciphertext)
}

//DecodeByAES 使用对称加密算法解密字符串 hash
func DecodeByAES(hash string) string  {
	// 实例化算法aes
	c, _ := newAES(key)

	//字符串转十六进制
	text,_ := hex.DecodeString(hash)

	// 解密字符串
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(text))
	cfbdec.XORKeyStream(plaintextCopy, text)
	return string(plaintextCopy)
}

//EncodeByRsa 实现RSA加密
func EncodeByRsa(origData []byte, publicKey []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

//DecodeByRsa 实现RSA解密
func DecodeByRsa(ciphertext []byte, privateKey []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}


//EncodeByBase64 实现base64编码
func EncodeByBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

//DecodeByBase64 实现base64解码
func DecodeByBase64(str string) (string, error) {
	strDecode,err := base64.StdEncoding.DecodeString(str)
	return string(strDecode), err
}

//UrlEncode 实现Url编码
func UrlEncode(url string) string {
	return base64.URLEncoding.EncodeToString([]byte(url))
}

//UrlDecode 实现Url解码
func UrlDecode(url string) (string,error) {
	urlDecode, err := base64.URLEncoding.DecodeString(url)
	if err == nil {
		return  string(urlDecode), nil
	}
	return "", err
}

