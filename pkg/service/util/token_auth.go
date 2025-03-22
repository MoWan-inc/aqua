package util

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strconv"
	"time"
)

const (
	UsrKey            = "user"
	InternalDeveloper = "internal developer"
	TimeLayout        = "2006102150405"
)

type TokenAuth interface {
	Need(ctx *gin.Context) bool
	CheckToken(ctx *gin.Context, token string) bool
}

type Token struct {
	Token     string `form:"token" json:"token"`
	Timestamp string `form:"timestamp" json:"timestamp"`
	Sign      string `form:"sign" json:"sign"`
}

func (t *Token) Empty() bool {
	return len(t.Token) == 0 && len(t.Timestamp) == 0 && len(t.Sign) == 0
}

func TokenAuthentication(auth TokenAuth) gin.HandlerFunc {
	// todo 日志
	return func(ctx *gin.Context) {
		// 如果不适用token验证，适用别的验证方式
		if !UseTokenAuthentication(ctx) {
			return
		}
		if !auth.Need(ctx) {
			return
		}
		token, err := getTokenFromCtx(ctx)
		if err != nil {
			JSONRequestError(ctx, err)
		}
		if !auth.CheckToken(ctx, token.Token) {
			JSONRequestError(ctx, errors.New("invalid token"))
			return
		}
		if err = checkAuthentication(token); err != nil {
			JSONRequestError(ctx, err)
			return
		}
	}
}

func checkAuthentication(token *Token) error {
	second, err := strconv.ParseInt(token.Timestamp, 10, 64)
	if err != nil {
		// todo 日志
		return err
	}
	ts := time.Unix(second, 0)
	checkSign := TokenSign(ts, token.Token)
	if checkSign != token.Sign {
		return fmt.Errorf("token check sign error, token:%s", token)
	}
	return nil
}

func getTokenFromCtx(ctx *gin.Context) (*Token, error) {
	token := &Token{}
	if err := ctx.ShouldBindWith(token, binding.Query); err != nil {
		return token, err
	}
	return token, nil
}

func GetTokenAuth(tokens []string, whiteList []string) TokenAuth {
	auth := &realTokenAuth{
		WhiteListURL: make(map[string]any),
		Tokens:       make(map[string]any),
	}
	for _, tk := range tokens {
		auth.Tokens[tk] = struct{}{}
	}
	for _, wl := range whiteList {
		auth.WhiteListURL[wl] = struct{}{}
	}
	return auth
}

type realTokenAuth struct {
	WhiteListURL map[string]any
	Tokens       map[string]any
}

func (a *realTokenAuth) Need(ctx *gin.Context) bool {
	_, has := a.WhiteListURL[ctx.FullPath()]
	return !has
}

func (a *realTokenAuth) CheckToken(ctx *gin.Context, token string) bool {
	if _, has := a.Tokens[token]; !has {
		return false
	}
	ctx.Set(UsrKey, InternalDeveloper)
	return true
}

func UseTokenAuthentication(ctx *gin.Context) bool {
	token, err := getTokenFromCtx(ctx)
	return err == nil && !token.Empty()
}

func TokenSign(now time.Time, token string) string {
	vin := fmt.Sprintf("%s-%s", now.Format(TimeLayout), token)
	return fmt.Sprintf("%s", md5.Sum([]byte(vin)))
}
