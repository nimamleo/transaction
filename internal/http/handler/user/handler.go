package user

import (
	"transaction/internal/user/application"
	"transaction/pkg/genericcode"
	"transaction/pkg/stdresponse"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	userService *application.Service
}

func NewHandler(userService *application.Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) CreateUser(c echo.Context) error {
	var req CreateUserRequest

	if err := c.Bind(&req); err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	if err := req.Validate(); err != nil {
		return stdresponse.SendHttpResponse(c, err.Error())
	}

	result, err := h.userService.CreateUser(c.Request().Context(), req.Name, req.Email)
	if err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	return stdresponse.SendHttpResponse(c, genericcode.OK, ToCreateResponse(result.User, result.APIKey))
}

func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	return stdresponse.SendHttpResponse(c, genericcode.OK, ToResponse(user))
}
