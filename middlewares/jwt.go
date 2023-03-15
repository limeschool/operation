package middlewares

import (
	"crypto/md5"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/limeschool/gin"
	"github.com/spf13/viper"
	"operation/consts"
	"operation/errors"
	"operation/middlewares/meta"
	model "operation/models/system"
	"strings"
	"time"
)

var (
	JwtAuth = new(JwtAuthApplication)
)

type JwtAuthApplication struct {
	Enable       bool            `json:"enable" mapstructure:"enable"` //是否启用
	Header       string          `json:"header" mapstructure:"header"`
	Captcha      bool            `json:"captcha" mapstructure:"captcha"`
	UniqueDevice bool            `json:"unique_device" mapstructure:"unique_device"` //是否启用唯一设备登陆
	UniqueCache  string          `json:"unique_cache" mapstructure:"unique_cache"`   //唯一设备登陆的缓存器
	Expire       int64           `json:"expire" mapstructure:"expire"`
	MaxExpire    int64           `json:"max_expire" mapstructure:"max_expire"`
	Secret       string          `json:"secret" mapstructure:"secret"`
	Whitelist    map[string]bool `json:"whitelist" mapstructure:"whitelist"`
}

func NewJwt(v *viper.Viper) {
	if err := v.UnmarshalKey("jwt", JwtAuth); err != nil {
		if JwtAuth == nil {
			panic("jwt 配置解析错误" + err.Error())
		}
		return
	}
	if JwtAuth.Header == "" {
		JwtAuth.Header = "Authorization"
	}
	return
}

func (ja *JwtAuthApplication) CompareToken(ctx *gin.Context, userId int64, token string) error {
	if !ja.UniqueDevice {
		return nil
	}
	storeToken, err := ctx.Redis(ja.UniqueCache).Get(ctx, ja.Md5UUID(fmt.Sprint(userId))).Result()
	if err != nil {
		return errors.TokenDataError
	}

	if storeToken != ja.Md5UUID(token) {
		return errors.DulDeviceLoginError
	}
	return nil
}

func (ja *JwtAuthApplication) StoreToken(ctx *gin.Context, userId int64, token string) error {
	if !ja.UniqueDevice {
		return nil
	}
	key := ja.Md5UUID(fmt.Sprint(userId))
	t := time.Duration(ja.Expire) * time.Second
	return ctx.Redis(ja.UniqueCache).Set(ctx, key, ja.Md5UUID(token), t).Err()
}

func (ja *JwtAuthApplication) ClearToken(ctx *gin.Context, userId int64) error {
	if !ja.UniqueDevice {
		return nil
	}
	key := ja.Md5UUID(fmt.Sprint(userId))
	return ctx.Redis(ja.UniqueCache).Del(ctx, key).Err()
}

func (ja *JwtAuthApplication) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !ja.Enable {
			return
		}

		path := ctx.FullPath()
		method := ctx.Request.Method

		// 判断是否在白名单内
		if ja.isWhitelist(path, method) {
			return
		}

		// 获取token
		token := ctx.Request.Header.Get(ja.Header)
		if token == "" {
			ctx.RespError(errors.TokenEmptyError)
			return
		}

		// 解析jwt信息
		mapClaims, err := ja.parseToken(token)
		if err != nil {
			ctx.RespError(err)
			return
		}

		metadata, err := meta.Parse(mapClaims[consts.JwtMapClaims])
		if err != nil {
			ctx.RespError(errors.TokenDataError)
			return
		}

		// 是否开启唯一设备登陆
		if ja.UniqueDevice {
			if err = ja.CompareToken(ctx, metadata.UserID, token); err != nil {
				ctx.RespError(err)
				return
			}
		}

		// 忽略超级管理员
		if metadata.RoleKey != consts.JwtSuperAdmin && !ja.isBaseApi(ctx, method, path) {
			// 进行权限认证
			if is, _ := ctx.Rbac().Enforce(metadata.RoleKey, path, method); !is {
				ctx.RespError(errors.NotResourcePower)
				return
			}
		}

		// 设置解析数据到上下文
		ctx.Set(consts.JwtMapClaims, metadata)
	}
}

func (ja *JwtAuthApplication) isBaseApi(ctx *gin.Context, method, path string) bool {
	menu := model.Menu{}
	return menu.GetBaseApiPath(ctx)[method+":"+path]
}

func (ja *JwtAuthApplication) isWhitelist(path, method string) bool {
	if path == "/healthy" {
		return true
	}
	key := strings.ToLower(method + ":" + path)
	return ja.Whitelist[key] == true
}

func (ja *JwtAuthApplication) parseToken(token string) (map[string]any, error) {
	var m jwt.MapClaims
	parser, err := jwt.ParseWithClaims(token, &m, func(token *jwt.Token) (interface{}, error) {
		return []byte(ja.Secret), nil
	})

	if err != nil || !parser.Valid {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.TokenExpiredError
		}
		return nil, errors.TokenValidateError
	}

	return m, nil
}

func (ja *JwtAuthApplication) Md5UUID(id string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(id)))
}

func (ja *JwtAuthApplication) CreateToken(ctx *gin.Context, metadata *meta.Metadata) (string, error) {
	// 生成token携带信息
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Unix() + ja.Expire
	claims["iat"] = time.Now().Unix()
	claims[consts.JwtMapClaims] = metadata
	tokenJwt := jwt.New(jwt.SigningMethodHS256)
	tokenJwt.Claims = claims
	token, err := tokenJwt.SignedString([]byte(ja.Secret))
	if err != nil {
		return "", err
	}

	return token, ja.StoreToken(ctx, metadata.UserID, token)
}

func (ja *JwtAuthApplication) MapClaimsAndExpired(ctx *gin.Context) (any, bool, bool) {
	var claims jwt.MapClaims
	token := ctx.Request.Header.Get(ja.Header)
	_, _ = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(ja.Secret), nil
	})
	if claims == nil {
		return nil, true, true
	}
	isMaxExpired := time.Now().Unix()-int64(claims["iat"].(float64)) > ja.MaxExpire
	isExpired := time.Now().Unix()-int64(claims["iat"].(float64)) > ja.Expire
	return claims[consts.JwtMapClaims], isExpired, isMaxExpired
}
