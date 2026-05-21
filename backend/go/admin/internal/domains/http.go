package domains

import (
	"fmt"
	"go-common/utils/raw_json"
)

var (
	StatusOk                 = 1
	StatusFail               = 2
	StatusRequestExpiredFail = 3
	StatusRequestNonceFail   = 4
	StatusRequestKeyFail     = 5

	// StatusLoginFail --start 登录相关状态
	//StatusAccountDisabledFail    = 101
	//StatusNotLoginFail           = 102
	//StatusTokenNotFoundFail      = 103
	//StatusInvalidTokenDataFail   = 104
	//StatusLoginLimitExceededFail = 105
	//StatusTokenKickoutFail       = 106
	//StatusTokenReplacedFail      = 107
	StatusLoginFail = 100

	// StatusAuthUnauthorized -- start 权限相关状态
	StatusAuthUnauthorized = 200

	OkMsg             = "success"
	ErrMsg            = "服务繁忙"
	ErrRequestMsg     = "请求错误"
	ErrRequestExpired = "请求超时"
	ErrRequestNonce   = "请求重放"
	ErrRequestKeyFail = "请求错误"

	// ErrAccountDisabled --start 登录相关状态
	ErrAccountDisabled    = "账户已禁用"
	ErrNotLogin           = "未登录"
	ErrTokenNotFound      = "账号错误"
	ErrInvalidTokenData   = "账号错误"
	ErrLoginLimitExceeded = "登录次数超过最大限制。"
	ErrTokenKickout       = "请重新登录"
	ErrTokenReplaced      = "请重新登录"

	// ErrAuthUnauthorized -- start 权限相关状态
	ErrAuthUnauthorized = "未授权"

	JsonOk          = raw_json.RawJson(fmt.Sprintf(`{"code":%d,"msg":"%s"}`, StatusOk, OkMsg))
	JsonErr         = raw_json.RawJson(fmt.Sprintf(`{"code":%d,"msg":"%s"}`, StatusFail, ErrMsg))
	JsonEmptyStruct = raw_json.RawJson(fmt.Sprintf(`{"code":%d,"msg":"%s","data":{}}`, StatusOk, OkMsg))
	JsonEmptySlice  = raw_json.RawJson(fmt.Sprintf(`{"code":%d,"msg":"%s","data":[]}`, StatusOk, OkMsg))
	JsonEmpty       = raw_json.RawJson(fmt.Sprintf(`{"code":%d,"msg":"%s","data":null}`, StatusOk, OkMsg))
)
