package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/MoWan-inc/aqua/pkg/util/object"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
)

type GinServerHandler func(ctx *gin.Context) (any, error)

type APIHandler interface {
	RegisterTo(group *gin.RouterGroup)
}

func JSONRequestError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
}

func JSONInternalError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
}

func JSONSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data, "msg": "success"})
}

func ParseParamID(c *gin.Context) (uint, error) {
	id := c.Param("id")
	parseID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, errors.New("id must be an unsigned integer")
	}
	return uint(parseID), nil
}

func HandleAPIWithLimiter(handler GinServerHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := handler(c)
		if err == nil {
			JSONSuccess(c, data)
		} else {
			if errors.As(err, &RequestError{}) {
				JSONRequestError(c, err)
			} else {
				JSONInternalError(c, err)
			}
			// todo 日志
		}
	}
}

func HandleAPIDefaultNotQuery(handler GinServerHandler) GinServerHandler {
	return func(c *gin.Context) (any, error) {
		if c.Query("not") == "" {
			u := c.Request.URL
			queryValues := u.Query()
			queryValues.Set("not", "{}")
			u.RawQuery = queryValues.Encode()
		}
		return handler(c)
	}
}

func DefaultHandlers(handler GinServerHandler) gin.HandlerFunc {
	return HandleAPIWithLimiter(HandleAPIDefaultNotQuery(handler))
}

type RequestError struct {
	error
}

type InternalError struct {
	error
}

type HTTPRsp[T any] struct {
	Data T      `json:"data"`
	Msg  string `json:"msg"`
}

type ErrorHTTPRsp HTTPRsp[any]

func NewRequestError(err error) *RequestError {
	return &RequestError{error: err}
}

func NewInternalError(err error) *InternalError {
	return &InternalError{error: err}
}

// 构造gin context，用于测试
func getTestGinContest(method, path string, forms map[string]string, body any) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	buffer := &bytes.Buffer{}
	// 表单数据
	url := path
	if !object.IsEmpty(forms) {
		var params []string
		for k, v := range forms {
			params = append(params, fmt.Sprintf("%s=%s", k, v))
			c.Params = append(c.Params, gin.Param{Key: k, Value: v})
		}
		url += "?" + strings.Join(params, "&")
	}
	// json body
	if !object.IsEmpty(body) {
		switch reflect.TypeOf(body).Kind() {
		case reflect.String:
			buffer.WriteString(body.(string))
		default:
			buffer.WriteString(object.MustMarshalJSON(body))
		}
	}
	request, err := http.NewRequest(method, url, buffer)
	if err != nil {
		panic(fmt.Sprintf("create request error %v", err))
	}
	c.Request = request
	return c
}
