package middleware

import (
	"mime"
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	i18nsvc "github.com/mijjjj/gcoll/internal/service/i18n"
)

const (
	contentTypeEventStream  = "text/event-stream"
	contentTypeOctetStream  = "application/octet-stream"
	contentTypeMixedReplace = "multipart/x-mixed-replace"
)

var streamContentTypes = []string{
	contentTypeEventStream,
	contentTypeOctetStream,
	contentTypeMixedReplace,
}

// HandlerResponse 约定控制器统一返回格式。
type HandlerResponse struct {
	Code    int    `json:"code"    dc:"业务错误码"`
	Message string `json:"message" dc:"响应消息"`
	Data    any    `json:"data"    dc:"业务数据"`
}

// Response 统一处理标准路由返回结果和错误。
func (s *Service) Response(r *ghttp.Request) {
	language := i18nsvc.LanguageFromRequest(
		r.GetQuery(i18nsvc.QueryKey()).String(),
		r.Header.Get(i18nsvc.HeaderKey()),
	)
	r.SetCtx(i18nsvc.WithLanguage(r.GetCtx(), language))
	r.Middleware.Next()

	if r.Response.BufferLength() > 0 || r.Response.BytesWritten() > 0 {
		return
	}

	mediaType, _, _ := mime.ParseMediaType(r.Response.Header().Get("Content-Type"))
	for _, contentType := range streamContentTypes {
		if mediaType == contentType {
			return
		}
	}

	var (
		err     = r.GetError()
		res     = r.GetHandlerResponse()
		code    = gerror.Code(err)
		message string
	)

	if err != nil {
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		message = i18nsvc.T(r.GetCtx(), err.Error())
	} else {
		if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
			switch r.Response.Status {
			case http.StatusNotFound:
				code = gcode.CodeNotFound
			case http.StatusForbidden:
				code = gcode.CodeNotAuthorized
			default:
				code = gcode.CodeUnknown
			}
			err = gerror.NewCode(code, code.Message())
			r.SetError(err)
		} else {
			code = gcode.CodeOK
		}
		message = i18nsvc.T(r.GetCtx(), code.Message())
	}

	r.Response.WriteJson(HandlerResponse{
		Code:    code.Code(),
		Message: message,
		Data:    res,
	})
}
