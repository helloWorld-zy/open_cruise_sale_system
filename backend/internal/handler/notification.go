package handler

import (
	"backend/internal/pagination"
	"backend/internal/service"
	"net/http"
	"strconv"

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
// @Success 200 {object} Response{data=pagination.Result}
// @Failure 401 {object} Response
// @Router /notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	userID := c.GetUint64("user_id")

	paginator := pagination.NewPaginator(
		c.Query("page"),
		c.Query("page_size"),
	)

	unreadOnly := c.Query("unread_only") == "true"

	result, err := h.notificationService.List(c.Request.Context(), userID, paginator, unreadOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to fetch notifications",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Data:    result,
		Message: "Success",
	})
}

// GetUnreadCount godoc
// @Summary Get unread notification count
// @Description Get the count of unread notifications for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=map[string]int64}
// @Failure 401 {object} Response
// @Router /notifications/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint64("user_id")

	count, err := h.notificationService.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to get unread count",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: map[string]int64{
			"count": count,
		},
		Message: "Success",
	})
}

// MarkAsRead godoc
// @Summary Mark notification as read
// @Description Mark a specific notification as read
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Router /notifications/{id}/read [post]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid notification ID",
		})
		return
	}

	if err := h.notificationService.MarkAsRead(c.Request.Context(), id); err != nil {
		if err == service.ErrNotificationNotFound {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "Notification not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to mark notification as read",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "Notification marked as read",
	})
}

// MarkAllAsRead godoc
// @Summary Mark all notifications as read
// @Description Mark all notifications for the current user as read
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Router /notifications/read-all [post]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint64("user_id")

	if err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to mark all notifications as read",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "All notifications marked as read",
	})
}

// Archive godoc
// @Summary Archive notification
// @Description Archive a specific notification
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Router /notifications/{id}/archive [post]
func (h *NotificationHandler) Archive(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid notification ID",
		})
		return
	}

	if err := h.notificationService.Archive(c.Request.Context(), id); err != nil {
		if err == service.ErrNotificationNotFound {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "Notification not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to archive notification",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "Notification archived",
	})
}

// Delete godoc
// @Summary Delete notification
// @Description Delete a specific notification
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Router /notifications/{id} [delete]
func (h *NotificationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid notification ID",
		})
		return
	}

	if err := h.notificationService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to delete notification",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "Notification deleted",
	})
}

// GetSettings godoc
// @Summary Get notification settings
// @Description Get notification settings for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=domain.NotificationSetting}
// @Failure 401 {object} Response
// @Router /notifications/settings [get]
func (h *NotificationHandler) GetSettings(c *gin.Context) {
	userID := c.GetUint64("user_id")

	settings, err := h.notificationService.GetSettings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to get notification settings",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Data:    settings,
		Message: "Success",
	})
}

// UpdateSettings godoc
// @Summary Update notification settings
// @Description Update notification settings for the current user
// @Tags notifications
// @Accept json
// @Produce json
// @Param settings body service.UpdateSettingsRequest true "Notification settings"
// @Success 200 {object} Response{data=domain.NotificationSetting}
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /notifications/settings [put]
func (h *NotificationHandler) UpdateSettings(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req service.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	settings, err := h.notificationService.UpdateSettings(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to update notification settings",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Data:    settings,
		Message: "Settings updated successfully",
	})
}
