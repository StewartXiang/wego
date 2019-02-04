package app

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/godcong/wego/cache"
	"github.com/godcong/wego/core"
	"github.com/godcong/wego/log"
	"github.com/godcong/wego/util"
	"github.com/json-iterator/go"
	"golang.org/x/exp/xerrors"
	"time"
)

// AccessTokenCredential ...
type AccessTokenCredential struct {
	GrantType string
	AppID     string
	Secret    string
}

// ToJSON ...
func (c *AccessTokenCredential) ToJSON() []byte {
	bytes, err := jsoniter.Marshal(c)
	if err != nil {
		return nil
	}
	return bytes
}

// AccessTokenOption ...
type AccessTokenOption struct {
	Credential *AccessTokenCredential
	RemoteHost string
	TokenKey   string
	TokenURL   string
}

/*AccessToken AccessToken */
type AccessToken struct {
	*AccessTokenProperty
	property *Property
	option   *AccessTokenOption
}

/*accessTokenKey 键值 */
const accessTokenKey = "access_token"
const accessTokenURLSuffix = "/cgi-bin/token"

/*AccessTokenSafeSeconds token安全时间 */
const AccessTokenSafeSeconds = 500

// RemoteHost ...
func (obj *AccessToken) RemoteHost() string {
	if obj.option.RemoteHost != "" {
		return obj.option.RemoteHost
	}
	return apiWeixin
}

func (obj *AccessToken) sendRequest(s string) []byte {
	//return core.Get(core.Splice(obj.prefix, obj.URL), obj.credentials).Bytes()
	return nil
}

func newAccessToken(property *Property, opts ...*AccessTokenOption) *AccessToken {
	var opt *AccessTokenOption
	if opts != nil {
		opt = opts[0]
	}
	return &AccessToken{
		AccessTokenProperty: property.AccessToken,
		property:            property,
		option:              opt,
	}
}

/*NewAccessToken NewAccessToken*/
func NewAccessToken(property *Property, opts ...*AccessTokenOption) *AccessToken {
	return newAccessToken(property, opts...)
}

/*Refresh 刷新AccessToken */
func (obj *AccessToken) Refresh() *AccessToken {
	log.Debug("AccessToken|Refresh")
	obj.getToken(true)
	return obj
}

/*GetRefreshedToken 获取刷新token */
func (obj *AccessToken) GetRefreshedToken() *core.Token {
	log.Debug("AccessToken|GetRefreshedToken")
	return obj.getToken(true)
}

/*GetToken 获取token */
func (obj *AccessToken) GetToken() *core.Token {
	log.Debug("AccessToken|GetToken")
	return obj.getToken(false)
}

// KeyMap ...
func (obj *AccessToken) KeyMap() util.Map {
	log.Debug("AccessToken|KeyMap")
	token := obj.getToken(false)
	return token.KeyMap()
}

/*GetTokenWithRefresh 重新获取token */
func (obj *AccessToken) GetTokenWithRefresh() *core.Token {
	log.Debug("AccessToken|GetTokenWithRefresh")
	return obj.getToken(true)
}

func (obj *AccessToken) getToken(refresh bool) *core.Token {
	key := obj.getCacheKey()

	if !refresh && cache.Has(key) {
		log.Debug("cached accessToken", key)
		if v, b := cache.Get(key).(*core.Token); b {
			if v.ExpiresIn > time.Now().Unix() {
				return v
			}
		}
	}

	token := obj.RequestToken(obj.getCredentials())
	if token == nil {
		return nil
	}
	log.Debug("AccessToken|getToken", *token)
	if v := token.ExpiresIn; v != 0 {
		obj.SetTokenWithLife(token.AccessToken, v-AccessTokenSafeSeconds)
	} else {
		obj.SetToken(token.AccessToken)
	}

	return token

}

/*RequestToken 请求获取token */
func (obj *AccessToken) RequestToken(credentials string) *core.Token {
	var token core.Token
	tokenByte := obj.sendRequest(credentials)
	if tokenByte == nil {
		return nil
	}
	err := json.Unmarshal(tokenByte, &token)
	if err != nil {
		log.Error(err)
		return nil
	}
	return &token
}

/*SetTokenWithLife set string accessToken with life time */
func (obj *AccessToken) SetTokenWithLife(token string, lifeTime int64) *AccessToken {
	return obj.setToken(token, lifeTime)
}

/*SetToken set string accessToken */
func (obj *AccessToken) SetToken(token string) *AccessToken {
	return obj.setToken(token, 7200)
}

func (obj *AccessToken) setToken(token string, lifeTime int64) *AccessToken {
	cache.SetWithTTL(obj.getCacheKey(), &core.Token{
		AccessToken: token,
		ExpiresIn:   time.Now().Add(time.Duration(lifeTime)).Unix(),
	}, lifeTime)
	return obj
}

func (obj *AccessToken) getCredentials() string {
	c := md5.Sum(obj.credentials.ToJSON())
	return fmt.Sprintf("%x", c[:])
}

func (obj *AccessToken) getCacheKey() string {
	return "godcong.wego.access_token." + obj.getCredentials()
}

const accessTokenNil = "nil point  accessToken"
const tokenNil = "nil point token"

/*MustKeyMap get accessToken's key,value with map when nil or error return nil map */
func MustKeyMap(at *AccessToken) util.Map {
	m := util.Map{}
	if m, e := KeyMap(at); e != nil {
		return m
	}
	return m
}

/*KeyMap get accessToken's key,value with map */
func KeyMap(at *AccessToken) (util.Map, error) {
	if at == nil {
		return nil, xerrors.New(accessTokenNil)
	}
	if token := at.GetToken(); token != nil {
		return token.KeyMap(), nil
	}
	return nil, xerrors.New(tokenNil)
}
