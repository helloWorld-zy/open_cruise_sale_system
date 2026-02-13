package handler

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"backend/internal/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	wechatSvc *service.WechatAuthService
	smsSvc    *service.SMSService
	repo      repository.UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	wechatSvc *service.WechatAuthService,
	smsSvc *service.SMSService,
	repo repository.UserRepository,
) *UserHandler {
	return &UserHandler{
		wechatSvc: wechatSvc,
		smsSvc:    smsSvc,
		repo:      repo,
	}
}

// WechatLogin godoc
// @Summary WeChat Mini Program Login
// @Description Login using WeChat mini program code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body WechatLoginRequest true "Login request"
// @Success 200 {object} response.Response{data=service.WechatLoginResult}
// @Failure 400 {object} response.Response
// @Router /auth/wechat/login [post]
func (h *UserHandler) WechatLogin(c *gin.Context) {
	var req WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.wechatSvc.WechatLogin(c.Request.Context(), req.Code)
	if err != nil {
		if err == service.ErrWechatAuthFailed {
			response.Error(c, http.StatusUnauthorized, "微信登录失败")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// WechatPhoneLogin godoc
// @Summary WeChat Phone Login
// @Description Login with WeChat and phone number
// @Tags auth
// @Accept json
// @Produce json
// @Param request body WechatPhoneLoginRequest true "Login request"
// @Success 200 {object} response.Response{data=service.WechatLoginResult}
// @Failure 400 {object} response.Response
// @Router /auth/wechat/phone-login [post]
func (h *UserHandler) WechatPhoneLogin(c *gin.Context) {
	var req WechatPhoneLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.wechatSvc.WechatPhoneLogin(
		c.Request.Context(),
		req.Code,
		req.EncryptedData,
		req.IV,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// SendSMSCode godoc
// @Summary Send SMS verification code
// @Description Send verification code to phone number
// @Tags auth
// @Accept json
// @Produce json
// @Param request body SendSMSRequest true "Send SMS request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/sms/send [post]
func (h *UserHandler) SendSMSCode(c *gin.Context) {
	var req SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.smsSvc.SendVerificationCode(c.Request.Context(), req.Phone); err != nil {
		if err == service.ErrInvalidPhone {
			response.BadRequest(c, "无效的手机号码")
			return
		}
		if err == service.ErrTooManyRequests {
			response.Error(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "验证码已发送"})
}

// SMSLogin godoc
// @Summary SMS Login
// @Description Login using phone number and SMS code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body SMSLoginRequest true "Login request"
// @Success 200 {object} response.Response{data=service.UserLoginResult}
// @Failure 400 {object} response.Response
// @Router /auth/sms/login [post]
func (h *UserHandler) SMSLogin(c *gin.Context) {
	var req SMSLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.smsSvc.VerifyAndLogin(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		if err == service.ErrInvalidCode {
			response.Error(c, http.StatusUnauthorized, "验证码错误")
			return
		}
		if err == service.ErrCodeExpired {
			response.Error(c, http.StatusUnauthorized, "验证码已过期")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile information
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=domain.User}
// @Failure 401 {object} response.Response
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	user, err := h.repo.GetByID(c.Request.Context(), userID.(string))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// Hide sensitive info
	user.PasswordHash = ""

	response.Success(c, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user profile
// @Tags user
// @Accept json
// @Produce json
// @Param request body UpdateProfileRequest true "Update request"
// @Success 200 {object} response.Response{data=domain.User}
// @Failure 400 {object} response.Response
// @Router /user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.repo.GetByID(c.Request.Context(), userID.(string))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// Update fields
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.BirthDate != "" {
		user.BirthDate = req.BirthDate
	}

	if err := h.repo.Update(c.Request.Context(), user); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	user.PasswordHash = ""
	response.Success(c, user)
}

// ListFrequentPassengers godoc
// @Summary List frequent passengers
// @Description Get user's frequent passengers list
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]domain.FrequentPassenger}
// @Failure 401 {object} response.Response
// @Router /user/passengers [get]
func (h *UserHandler) ListFrequentPassengers(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	passengers, err := h.repo.ListFrequentPassengersByUser(c.Request.Context(), userID.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, passengers)
}

// CreateFrequentPassenger godoc
// @Summary Create frequent passenger
// @Description Add a new frequent passenger
// @Tags user
// @Accept json
// @Produce json
// @Param request body CreatePassengerRequest true "Passenger info"
// @Success 201 {object} response.Response{data=domain.FrequentPassenger}
// @Failure 400 {object} response.Response
// @Router /user/passengers [post]
func (h *UserHandler) CreateFrequentPassenger(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req CreatePassengerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.ValidateStruct(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	passenger := &domain.FrequentPassenger{
		UserID:              userID.(string),
		Name:                req.Name,
		Surname:             req.Surname,
		GivenName:           req.GivenName,
		Gender:              req.Gender,
		BirthDate:           req.BirthDate,
		Nationality:         req.Nationality,
		PassportNumber:      req.PassportNumber,
		PassportExpiry:      req.PassportExpiry,
		IDNumber:            req.IDNumber,
		Phone:               req.Phone,
		Email:               req.Email,
		DietaryRequirements: req.DietaryRequirements,
		MedicalNotes:        req.MedicalNotes,
		IsDefault:           req.IsDefault,
	}

	if err := h.repo.CreateFrequentPassenger(c.Request.Context(), passenger); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(c, passenger)
}

// UpdateFrequentPassenger godoc
// @Summary Update frequent passenger
// @Description Update a frequent passenger
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "Passenger ID"
// @Param request body UpdatePassengerRequest true "Passenger info"
// @Success 200 {object} response.Response{data=domain.FrequentPassenger}
// @Failure 400 {object} response.Response
// @Router /user/passengers/{id} [put]
func (h *UserHandler) UpdateFrequentPassenger(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	passengerID := c.Param("id")

	// Verify ownership
	passenger, err := h.repo.GetFrequentPassengerByID(c.Request.Context(), passengerID)
	if err != nil {
		response.NotFound(c, "乘客不存在")
		return
	}

	if passenger.UserID != userID.(string) {
		response.Forbidden(c, "无权修改此乘客信息")
		return
	}

	var req UpdatePassengerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Update fields
	if req.Name != "" {
		passenger.Name = req.Name
	}
	if req.Surname != "" {
		passenger.Surname = req.Surname
	}
	if req.GivenName != "" {
		passenger.GivenName = req.GivenName
	}
	if req.Gender != "" {
		passenger.Gender = req.Gender
	}
	if req.BirthDate != "" {
		passenger.BirthDate = req.BirthDate
	}
	if req.PassportNumber != "" {
		passenger.PassportNumber = req.PassportNumber
	}
	if req.IDNumber != "" {
		passenger.IDNumber = req.IDNumber
	}

	if err := h.repo.UpdateFrequentPassenger(c.Request.Context(), passenger); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, passenger)
}

// DeleteFrequentPassenger godoc
// @Summary Delete frequent passenger
// @Description Delete a frequent passenger
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "Passenger ID"
// @Success 200 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /user/passengers/{id} [delete]
func (h *UserHandler) DeleteFrequentPassenger(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	passengerID := c.Param("id")

	// Verify ownership
	passenger, err := h.repo.GetFrequentPassengerByID(c.Request.Context(), passengerID)
	if err != nil {
		response.NotFound(c, "乘客不存在")
		return
	}

	if passenger.UserID != userID.(string) {
		response.Forbidden(c, "无权删除此乘客信息")
		return
	}

	if err := h.repo.DeleteFrequentPassenger(c.Request.Context(), passengerID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Request types

type WechatLoginRequest struct {
	Code string `json:"code" validate:"required"`
}

type WechatPhoneLoginRequest struct {
	Code          string `json:"code" validate:"required"`
	EncryptedData string `json:"encrypted_data" validate:"required"`
	IV            string `json:"iv" validate:"required"`
}

type SendSMSRequest struct {
	Phone string `json:"phone" validate:"required,len=11"`
}

type SMSLoginRequest struct {
	Phone string `json:"phone" validate:"required,len=11"`
	Code  string `json:"code" validate:"required,len=6"`
}

type UpdateProfileRequest struct {
	Nickname  string `json:"nickname,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	RealName  string `json:"real_name,omitempty"`
	Gender    string `json:"gender,omitempty" validate:"omitempty,oneof=male female unknown"`
	BirthDate string `json:"birth_date,omitempty"`
}

type CreatePassengerRequest struct {
	Name                string `json:"name" validate:"required"`
	Surname             string `json:"surname" validate:"required"`
	GivenName           string `json:"given_name,omitempty"`
	Gender              string `json:"gender" validate:"required,oneof=male female"`
	BirthDate           string `json:"birth_date" validate:"required"`
	Nationality         string `json:"nationality,omitempty"`
	PassportNumber      string `json:"passport_number,omitempty"`
	PassportExpiry      string `json:"passport_expiry,omitempty"`
	IDNumber            string `json:"id_number,omitempty"`
	Phone               string `json:"phone,omitempty"`
	Email               string `json:"email,omitempty"`
	DietaryRequirements string `json:"dietary_requirements,omitempty"`
	MedicalNotes        string `json:"medical_notes,omitempty"`
	IsDefault           bool   `json:"is_default,omitempty"`
}

type UpdatePassengerRequest struct {
	Name           string `json:"name,omitempty"`
	Surname        string `json:"surname,omitempty"`
	GivenName      string `json:"given_name,omitempty"`
	Gender         string `json:"gender,omitempty" validate:"omitempty,oneof=male female"`
	BirthDate      string `json:"birth_date,omitempty"`
	PassportNumber string `json:"passport_number,omitempty"`
	IDNumber       string `json:"id_number,omitempty"`
}
