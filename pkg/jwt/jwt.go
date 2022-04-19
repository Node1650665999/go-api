package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"gin-api/pkg/hash"
	"gin-api/pkg/helpers"
	"strings"
	"time"
)

//Key 用作签名
const Key = "MMW4n4slID"

//RefreshExpire 指定多长时间内可以刷新 token(一周)
const RefreshExpire  = 7 * 24 * time.Hour

//Header 定义了JWT的头部信息,由两部分组成：加密算法和类型
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

//PayLoad 定义了JWT中的有效信息
type PayLoad struct {
	//用户自定义的数据
	CustomClaims interface{} `json:"custom_claims"`

	//iss: jwt签发者
	Iss string `json:"iss"`

	//sub: jwt所面向的用户
	Sub string `json:"sub"`

	//aud: 接收jwt的一方
	Aud string `json:"aud"`

	//exp: jwt的过期时间，这个过期时间必须要大于签发时间(时间戳)
	Exp int64 `json:"exp"`

	//nbf: 定义在什么时间之前,该jwt都是不可用的(时间戳)
	Nbf int64 `json:"nbf"`

	//iat: jwt的签发时间(时间戳)
	Iat int64 `json:"iat"`

	//jti: jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击。
	Jti string `json:"jti"`
}

//GenerateToken 生成token. customClaims 为用户自定义数据, expire 为生成的token有效时间
func GenerateToken(customClaims interface{}, expire time.Duration) string {
	currentTime := time.Now().Unix()
	exp := currentTime + int64(expire)
	//header
	h := &Header{Alg: "HS256", Typ: "jwt"}
	//payload
	p := &PayLoad{
		CustomClaims: customClaims,
		Exp: exp,
		Iat: currentTime,
		Nbf: currentTime,
		Jti: helpers.StrUuid(30),
	}

	//签名
	headerAndPayLoad, sign := signature(h, p)
	return headerAndPayLoad + "." + sign
}


//VerifyToken 用来验证token是否合法, 如果合法,则返回用户自定义数据,否则返回error
func VerifyToken(token string) (customClaims interface{}, err error)  {
	payload, header, err := parseToken(token)
	if err != nil {
		return nil, err
	}
	_, sg := signature(header, payload)
	sign  := strings.Split(token, ".")[2]
	if sign != sg {
		return nil, fmt.Errorf("token 无效")
	}

	return payload.CustomClaims,nil
}

//RefreshToken 用来刷新token
func RefreshToken(token string, expire time.Duration) (string, error) {
	payload, _, err := parseToken(token)

	if err != nil {
		return "", err
	}

	if time.Unix(payload.Iat, 0).Add(RefreshExpire).Before(time.Now()) {
		return "", errors.New("令牌已过刷新时间")
	}

	newToken := GenerateToken(payload.CustomClaims, expire)
	return newToken, nil
}

//parseToken 解析 token
func parseToken(token string) (*PayLoad, *Header, error) {
	currentTime := time.Now().Unix()
	tokens := strings.Split(token, ".")
	if len(tokens) != 3 {
		return nil, nil, fmt.Errorf("token 格式有误")
	}

	h, err := hash.DecodeByBase64(tokens[0])
	if err != nil {
		return nil,nil,err
	}

	p, err := hash.DecodeByBase64(tokens[1])
	if err != nil {
		return nil, nil, err
	}

	var header  Header
	var payload PayLoad

	json.Unmarshal([]byte(h), &header)
	json.Unmarshal([]byte(p), &payload)

	if payload.Iat > currentTime {
		return nil, nil, fmt.Errorf("签发时间大于当前时间")
	}

	if payload.Exp < currentTime {
		return nil, nil, fmt.Errorf("token 已过期")
	}

	if payload.Nbf > currentTime {
		return nil, nil, fmt.Errorf("token 还未生效")
	}

	return &payload , &header, nil
}

//signature 用来生成签名信息
func signature(h *Header, p *PayLoad) (hp string, sign string) {
	header, _ := json.Marshal(h)
	payLoad, _ := json.Marshal(p)
	headerAndPayLoad := hash.EncodeByBase64(string(header)) + "." + hash.EncodeByBase64(string(payLoad))
	return headerAndPayLoad, hash.HmacSha256(headerAndPayLoad, Key)
}
