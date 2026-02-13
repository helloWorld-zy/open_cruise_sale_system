package handler

import (
	"backend/internal/pagination"
	"backend/internal/response"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotificationHandler handles notification HTTP requests
type NotificationHandler struct {
	notificationService service.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// RegisterRoutes registers the notification routes
func (h *NotificationHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	notifications := router.Group("/notifications")
	notifications.Use(authMiddleware)
	{
		notifications.GET("", h.List)
		notifications.GET("/unread-count", h.GetUnreadCount)
		notifications.POST("/:id/read", h.MarkAsRead)
		notifications.POST("/read-all", h.MarkAllAsRead)
		notifications.POST("/:id/archive", h.Archive)
		notifications.DELETE("/:id", h.Delete)

		// Settings
		notifications.GET("/settings", h.GetSettings)
		notifications.PUT("/settings", h.UpdateSettings)
	}
}

// List godoc
// @Summary List user notifications
// @Description Get a paginated list of notifications for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20)"
// @Param unread_only query bool false "Show only unread notifications"
// @Success 200 {object} response.Response{data=pagination.Result}
// @Failure 401 {object} response.Response
// @Router /notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	paginator := pagination.NewPaginator(c)

	unreadOnly := c.Query("unread_only") == "true"

	result, err := h.notificationService.List(c.Request.Context(), userID, paginator, unreadOnly)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}

	response.Success(c, result)
}

// GetUnreadCount godoc
// @Summary Get unread notification count
// @Description Get the count of unread notifications for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]int64}
// @Failure 401 {object} response.Response
// @Router /notifications/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetString("user_id")

	count, err := h.notificationService.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get unread count")
		return
	}

	response.Success(c, map[string]int64{
		"count": count,
	})
}

// MarkAsRead godoc
// @Summary Mark notification as read
// @Description Mark a specific notification as read
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /notifications/{id}/read [post]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	if err := h.notificationService.MarkAsRead(c.Request.Context(), id); err != nil {
		if err == service.ErrNotificationNotFound {
			response.Error(c, http.StatusNotFound, "Notification not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to mark notification as read")
		return
	}

	response.Success(c, gin.H{"message": "Notification marked as read"})
}

// MarkAllAsRead godoc
// @Summary Mark all notifications as read
// @Description Mark all notifications for the current user as read
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /notifications/read-all [post]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to mark all notifications as read")
		return
	}

	response.Success(c, gin.H{"message": "All notifications marked as read"})
}

// Archive godoc
// @Summary Archive notification
// @Description Archive a specific notification
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /notifications/{id}/archive [post]
func (h *NotificationHandler) Archive(c *gin.Context) {
	id := c.Param("id")

	if err := h.notificationService.Archive(c.Request.Context(), id); err != nil {
		if err == service.ErrNotificationNotFound {
			response.Error(c, http.StatusNotFound, "Notification not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to archive notification")
		return
	}

	response.Success(c, gin.H{"message": "Notification archived"})
}

// Delete godoc
// @Summary Delete notification
// @Description Delete a specific notification
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /notifications/{id} [delete]
func (h *NotificationHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.notificationService.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete notification")
		return
	}

	response.Success(c, gin.H{"message": "Notification deleted"})
}

// GetSettings godoc
// @Summary Get notification settings
// @Description Get notification settings for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=domain.NotificationSetting}
// @Failure 401 {object} response.Response
// @Router /notifications/settings [get]
func (h *NotificationHandler) GetSettings(c *gin.Context) {
	userID := c.GetString("user_id")

	settings, err := h.notificationService.GetSettings(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get notification settings")
		return
	}

	response.Success(c, settings)
}

// UpdateSettings godoc
// @Summary Update notification settings
// @Description Update notification settings for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Param settings body service.UpdateSettingsRequest true "Notification settings"
// @Success 200 {object} response.Response{data=domain.NotificationSetting}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /notifications/settings [put]
func (h *NotificationHandler) UpdateSettings(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data: "+err.Error())
		return
	}

	settings, err := h.notificationService.UpdateSettings(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update notification settings")
		return
	}

	response.Success(c, settings)
}
