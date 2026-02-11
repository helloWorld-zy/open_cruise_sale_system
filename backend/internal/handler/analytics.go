package handler

import (
	"backend/internal/analytics"
	"backend/internal/recommendation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AnalyticsHandler handles analytics and recommendation HTTP requests
type AnalyticsHandler struct {
	recommendationEngine recommendation.RecommendationEngine
	priceTrendAnalysis   analytics.PriceTrendAnalysis
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(
	recommendationEngine recommendation.RecommendationEngine,
	priceTrendAnalysis analytics.PriceTrendAnalysis,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		recommendationEngine: recommendationEngine,
		priceTrendAnalysis:   priceTrendAnalysis,
	}
}

// RegisterRoutes registers the analytics routes
func (h *AnalyticsHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	analytics := router.Group("/analytics")
	{
		// Recommendations
		analytics.GET("/recommendations", h.GetRecommendations)
		analytics.GET("/recommendations/personalized", authMiddleware, h.GetPersonalizedRecommendations)
		analytics.GET("/recommendations/popular", h.GetPopularVoyages)
		analytics.GET("/recommendations/similar/:voyageId", h.GetSimilarVoyages)
		analytics.GET("/recommendations/last-minute", h.GetLastMinuteDeals)

		// Price trends
		analytics.GET("/price-trends/:voyageId", h.GetPriceTrend)
		analytics.GET("/price-trends/cruise/:cruiseId", h.GetCruisePriceTrends)
		analytics.GET("/price-calendar", h.GetPriceCalendar)
		analytics.GET("/price-forecast/:voyageId", h.GetPriceForecast)
		analytics.POST("/price-compare", h.CompareVoyagePrices)
	}
}

// GetRecommendations godoc
// @Summary Get recommendations
// @Description Get cruise recommendations based on type
// @Tags analytics
// @Accept json
// @Produce json
// @Param type query string false "Recommendation type (popular, personalized, similar, last_minute)"
// @Param limit query int false "Number of recommendations (default: 10)"
// @Success 200 {object} Response{data=[]recommendation.Recommendation}
// @Failure 400 {object} Response
// @Router /analytics/recommendations [get]
func (h *AnalyticsHandler) GetRecommendations(c *gin.Context) {
	req := recommendation.RecommendationRequest{
		Type: recommendation.RecommendationType(c.DefaultQuery("type", string(recommendation.RecommendationTypePopular))),
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	req.Limit = limit

	recommendations, err := h.recommendationEngine.GetRecommendations(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendations})
}

// GetPersonalizedRecommendations godoc
// @Summary Get personalized recommendations
// @Description Get personalized cruise recommendations for the authenticated user
// @Tags analytics
// @Accept json
// @Produce json
// @Param limit query int false "Number of recommendations"
// @Success 200 {object} Response{data=[]recommendation.Recommendation}
// @Failure 401 {object} Response
// @Router /analytics/recommendations/personalized [get]
func (h *AnalyticsHandler) GetPersonalizedRecommendations(c *gin.Context) {
	userID := c.GetUint64("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	recommendations, err := h.recommendationEngine.GetPersonalizedRecommendations(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendations})
}

// GetPopularVoyages godoc
// @Summary Get popular voyages
// @Description Get currently popular cruise voyages
// @Tags analytics
// @Accept json
// @Produce json
// @Param limit query int false "Number of voyages"
// @Success 200 {object} Response{data=[]recommendation.Recommendation}
// @Router /analytics/recommendations/popular [get]
func (h *AnalyticsHandler) GetPopularVoyages(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	recommendations, err := h.recommendationEngine.GetPopularVoyages(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendations})
}

// GetSimilarVoyages godoc
// @Summary Get similar voyages
// @Description Get voyages similar to the specified voyage
// @Tags analytics
// @Accept json
// @Produce json
// @Param voyageId path string true "Voyage ID"
// @Param limit query int false "Number of voyages"
// @Success 200 {object} Response{data=[]recommendation.Recommendation}
// @Failure 404 {object} Response
// @Router /analytics/recommendations/similar/{voyageId} [get]
func (h *AnalyticsHandler) GetSimilarVoyages(c *gin.Context) {
	voyageID := c.Param("voyageId")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	recommendations, err := h.recommendationEngine.GetSimilarVoyages(c.Request.Context(), voyageID, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voyage not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendations})
}

// GetLastMinuteDeals godoc
// @Summary Get last minute deals
// @Description Get voyages departing soon with available inventory
// @Tags analytics
// @Accept json
// @Produce json
// @Param limit query int false "Number of voyages"
// @Success 200 {object} Response{data=[]recommendation.Recommendation}
// @Router /analytics/recommendations/last-minute [get]
func (h *AnalyticsHandler) GetLastMinuteDeals(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	recommendations, err := h.recommendationEngine.GetLastMinuteDeals(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendations})
}

// GetPriceTrend godoc
// @Summary Get price trend
// @Description Get price trend data for a specific voyage
// @Tags analytics
// @Accept json
// @Produce json
// @Param voyageId path string true "Voyage ID"
// @Success 200 {object} Response{data=analytics.PriceTrend}
// @Failure 404 {object} Response
// @Router /analytics/price-trends/{voyageId} [get]
func (h *AnalyticsHandler) GetPriceTrend(c *gin.Context) {
	voyageID := c.Param("voyageId")

	trend, err := h.priceTrendAnalysis.GetPriceTrend(c.Request.Context(), voyageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voyage not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trend})
}

// GetCruisePriceTrends godoc
// @Summary Get cruise price trends
// @Description Get price trends for all voyages of a cruise
// @Tags analytics
// @Accept json
// @Produce json
// @Param cruiseId path string true "Cruise ID"
// @Param limit query int false "Number of voyages"
// @Success 200 {object} Response{data=[]analytics.PriceTrend}
// @Failure 404 {object} Response
// @Router /analytics/price-trends/cruise/{cruiseId} [get]
func (h *AnalyticsHandler) GetCruisePriceTrends(c *gin.Context) {
	cruiseID := c.Param("cruiseId")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	trends, err := h.priceTrendAnalysis.GetCruisePriceTrends(c.Request.Context(), cruiseID, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cruise not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": trends})
}

// GetPriceCalendar godoc
// @Summary Get price calendar
// @Description Get price calendar for selecting departure dates
// @Tags analytics
// @Accept json
// @Produce json
// @Param cruise_id query string true "Cruise ID"
// @Param route_id query string false "Route ID"
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Param cabin_type_id query string false "Cabin Type ID"
// @Success 200 {object} Response{data=analytics.PriceCalendar}
// @Failure 400 {object} Response
// @Router /analytics/price-calendar [get]
func (h *AnalyticsHandler) GetPriceCalendar(c *gin.Context) {
	cruiseID := c.Query("cruise_id")
	routeID := c.Query("route_id")
	year, _ := strconv.Atoi(c.Query("year"))
	month, _ := strconv.Atoi(c.Query("month"))

	if cruiseID == "" || year == 0 || month == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	var cabinTypeID *string
	if ctid := c.Query("cabin_type_id"); ctid != "" {
		cabinTypeID = &ctid
	}

	calendar, err := h.priceTrendAnalysis.GetPriceCalendar(c.Request.Context(), cruiseID, routeID, year, month, cabinTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": calendar})
}

// GetPriceForecast godoc
// @Summary Get price forecast
// @Description Get price forecast for a specific voyage
// @Tags analytics
// @Accept json
// @Produce json
// @Param voyageId path string true "Voyage ID"
// @Success 200 {object} Response{data=analytics.PriceForecast}
// @Failure 404 {object} Response
// @Router /analytics/price-forecast/{voyageId} [get]
func (h *AnalyticsHandler) GetPriceForecast(c *gin.Context) {
	voyageID := c.Param("voyageId")

	forecast, err := h.priceTrendAnalysis.GetPriceForecast(c.Request.Context(), voyageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voyage not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": forecast})
}

// CompareVoyagePrices godoc
// @Summary Compare voyage prices
// @Description Compare prices across multiple voyages
// @Tags analytics
// @Accept json
// @Produce json
// @Param voyage_ids body []string true "Array of voyage IDs"
// @Success 200 {object} Response{data=[]analytics.VoyagePriceComparison}
// @Failure 400 {object} Response
// @Router /analytics/price-compare [post]
func (h *AnalyticsHandler) CompareVoyagePrices(c *gin.Context) {
	var voyageIDs []string
	if err := c.ShouldBindJSON(&voyageIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	comparisons, err := h.priceTrendAnalysis.CompareVoyagePrices(c.Request.Context(), voyageIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comparisons})
}
