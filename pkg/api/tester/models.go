package tester

import (
	"github.com/go-resty/resty/v2"
	"time"
)

// LoginRequest 定义登录请求的数据结构
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 定义登录响应的数据结构
type LoginResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    struct {
		Successful bool   `json:"successful"`
		Detail     string `json:"detail"`
		Token      string `json:"token"`
	} `json:"data"`
}

// ApiResponse 定义登录响应的数据结构
type ApiResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

// Report 测试报告
type Report struct {
	ApiName  string
	ApiPath  string
	IsPassed bool
	Elapsed  time.Duration
	Response *resty.Response
	Data     *ApiResponse
}

// ApiInfo 接口信息
type ApiInfo struct {
	ApiName     string            `json:"apiName"`
	ApiPath     string            `json:"apiPath"`
	Method      string            `json:"method"`
	QueryParams map[string]string `json:"queryParams"`
	Body        map[string]string `json:"body"`
}
