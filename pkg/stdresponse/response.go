package stdresponse

import (
	"transaction/pkg/genericcode"
	"transaction/pkg/richerror"

	"github.com/labstack/echo/v4"
)

type PaginatedMetadata struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
	Total    int64 `json:"total"`
	PerPage  int64 `json:"per_page"`
}

type StdResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    any                `json:"data"`
	Meta    *PaginatedMetadata `json:"meta,omitempty"`
}

func GenericCodeToHttpCode(code genericcode.Code) int {
	switch code {
	case genericcode.InternalServerError:
		return 500
	case genericcode.NotFound:
		return 404
	case genericcode.Unauthorized:
		return 401
	case genericcode.OK:
		return 200
	case genericcode.BadRequest:
		return 400
	case genericcode.Forbidden:
		return 403
	case genericcode.Conflict:
		return 409
	default:
		return 500
	}
}

func SendHttpResponse(c echo.Context, data ...any) error {
	stdResponse := StdResponse{
		Code:    0,
		Message: "",
		Data:    nil,
		Meta:    nil,
	}

	for _, extra := range data {
		switch v := extra.(type) {
		case genericcode.Code:
			stdResponse.Code = GenericCodeToHttpCode(v)
		case string:
			stdResponse.Message = v
		case richerror.RichError:
			stdResponse.Code = GenericCodeToHttpCode(v.GetCode())
			stdResponse.Message = v.GetMessage()
			if v.Data != nil {
				stdResponse.Data = v.Data
			}
		case PaginatedMetadata:
			stdResponse.Meta = &v
		default:
			stdResponse.Data = v
		}
	}

	return c.JSON(stdResponse.Code, stdResponse)
}
