package handlers

import (
	"log/slog"

	"github.com/OxytocinGroup/theca-v3/internal/metrics"
	"github.com/OxytocinGroup/theca-v3/internal/model"
	"github.com/OxytocinGroup/theca-v3/internal/service"
	errors "github.com/OxytocinGroup/theca-v3/internal/utils/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
	log     *slog.Logger
}

func NewHandler(service service.Service, log *slog.Logger) *Handler {
	return &Handler{service: service, log: log}
}

// @Summary Register
// @Description Register a new user
// @Tags user
// @Accept json
// @Produce json
// @Param registerRequest body model.RegisterRequest true "Register request"
// @Success 200
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Router /register [post]
func (h *Handler) Register(c *gin.Context) {
	const op = "handler.register"
	log := h.log.With(slog.String("op", op))

	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("binding json", "err", err, "req", req)
		metrics.RecordError(c.Request.Context(), "validation_error", c.Request.URL.Path, c.Request.Method)
		errors.RespondWithError(c, errors.New(errors.CodeInvalidRequest, "Неверный формат запроса"))
		return
	}

	err := h.service.Register(c.Request.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		metrics.RecordError(c.Request.Context(), "registration_error", c.Request.URL.Path, c.Request.Method)
		errors.RespondWithError(c, err)
		return
	}

	errors.RespondWithSuccess(c, "User registered successfully")
}

// @Summary Login
// @Description Login a user
// @Tags user
// @Accept json
// @Produce json
// @Param loginRequest body model.LoginRequest true "Login request"
// @Success 200
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	const op = "handler.login"
	log := h.log.With(slog.String("op", op))

	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Debug("binding json", "err", err, "req", req)
		metrics.RecordError(c.Request.Context(), "validation_error", c.Request.URL.Path, c.Request.Method)
		errors.RespondWithError(c, errors.New(errors.CodeInvalidRequest, "Неверный формат запроса"))
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		metrics.RecordError(c.Request.Context(), "authentication_error", c.Request.URL.Path, c.Request.Method)
		errors.RespondWithError(c, err)
		return
	}

	c.SetCookie("refreshToken", refreshToken, 0, "/", "", false, true)

	errors.RespondWithSuccess(c, gin.H{
		"access_token": accessToken,
	})
}

// @Summary Logout
// @Description Logout a user
// @Tags user
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {object} errors.Error
// @Failure 500 {object} errors.Error
// @Router /api/logout [delete]
func (h *Handler) Logout(c *gin.Context) {
	const op = "handler.logout"
	log := h.log.With(slog.String("op", op))

	c.SetCookie("refreshToken", "", -1, "/", "", false, true)
	log.Debug("logged out", "user", c.GetUint("user_id"))
	errors.RespondWithSuccess(c, "Logged out successfully")
}
