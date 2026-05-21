package res

import (
	"admin/internal/domains"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
)

type Response struct {
	Code int    `json:"code" example:"1"`
	Msg  string `json:"msg" example:"success"`
	Data any    `json:"data"`
}

func (r Response) Error() string {
	var sb strings.Builder
	sb.WriteString(`{"code":`)
	sb.WriteString(strconv.Itoa(r.Code))

	hasMsg := r.Msg != ""
	hasData := r.Data != nil

	if hasMsg {
		sb.WriteString(`,"msg":"`)
		for _, ch := range r.Msg {
			if ch == '"' {
				sb.WriteString(`\"`)
			} else if ch == '\\' {
				sb.WriteString(`\\`)
			} else {
				sb.WriteRune(ch)
			}
		}
		sb.WriteByte('"')
	}

	if hasData {
		sb.WriteString(`,"data":`)
		dataBytes, err := sonic.Marshal(r.Data)
		if err != nil {
			sb.WriteString(`null`)
		} else {
			sb.Write(dataBytes)
		}
	}

	sb.WriteByte('}')

	return sb.String()
}

func OkRes(data any) Response {
	return Response{
		Code: domains.StatusOk,
		Msg:  domains.OkMsg,
		Data: data,
	}
}

func FailMsg(msg string) Response {
	return Response{
		Code: domains.StatusFail,
		Msg:  msg,
		Data: nil,
	}
}

func FailCodeMsg(code int, msg string) Response {
	return Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

var (
	FailDefault            = Response{Code: domains.StatusFail, Msg: domains.ErrMsg}
	FailRequest            = Response{Code: domains.StatusFail, Msg: domains.ErrRequestMsg}
	FailRequestExpired     = Response{Code: domains.StatusRequestExpiredFail, Msg: domains.ErrRequestExpired}
	FailRequestNonce       = Response{Code: domains.StatusRequestNonceFail, Msg: domains.ErrRequestNonce}
	FailRequestKey         = Response{Code: domains.StatusRequestKeyFail, Msg: domains.ErrRequestKeyFail}
	FailAccountDisabled    = Response{Code: domains.StatusLoginFail, Msg: domains.ErrAccountDisabled}
	FailNotLogin           = Response{Code: domains.StatusLoginFail, Msg: domains.ErrNotLogin}
	FailTokenNotFound      = Response{Code: domains.StatusLoginFail, Msg: domains.ErrTokenNotFound}
	FailInvalidTokenData   = Response{Code: domains.StatusLoginFail, Msg: domains.ErrInvalidTokenData}
	FailLoginLimitExceeded = Response{Code: domains.StatusLoginFail, Msg: domains.ErrLoginLimitExceeded}
	FailTokenKickout       = Response{Code: domains.StatusLoginFail, Msg: domains.ErrTokenKickout}
	FailTokenReplaced      = Response{Code: domains.StatusLoginFail, Msg: domains.ErrTokenReplaced}
	FailAuthUnauthorized   = Response{Code: domains.StatusAuthUnauthorized, Msg: domains.ErrAuthUnauthorized}
)
