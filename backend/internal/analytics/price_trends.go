package analytics

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"context"
	"fmt"
	"sort"
	"time"
)

// PriceTrend represents price trend data for a voyage
type PriceTrend struct {
	VoyageID        string           `json:"voyage_id"`
	VoyageNumber    string           `json:"voyage_number"`
	CruiseName      string           `json:"cruise_name"`
	DepartureDate   string           `json:"departure_date"`
	CurrentPrice    float64          `json:"current_price"`
	OriginalPrice   float64          `json:"original_price"`
	LowestPrice     float64          `json:"lowest_price"`
	HighestPrice    float64          `json:"highest_price"`
	PriceChange     float64          `json:"price_change"`     // Absolute change
	PriceChangePct  float64          `json:"price_change_pct"` // Percentage change
	TrendDirection  string           `json:"trend_direction"`  // up, down, stable
	PriceHistory    []PricePoint     `json:"price_history"`
	CabinTypePrices []CabinPriceInfo `json:"cabin_type_prices"`
	ForecastPrice   *float64         `json:"forecast_price,omitempty"`
	ForecastTrend   string           `json:"forecast_trend,omitempty"` // likely_up, likely_down, likely_stable
}

// PricePoint represents a historical price point
type PricePoint struct {
	Date        string  `json:"date"`
	Price       float64 `json:"price"`
	CabinTypeID string  `json:"cabin_type_id,omitempty"`
}

// CabinPriceInfo represents current price for a cabin type
type CabinPriceInfo struct {
	CabinTypeID   string  `json:"cabin_type_id"`
	CabinTypeName string  `json:"cabin_type_name"`
	CurrentPrice  float64 `json:"current_price"`
	OriginalPrice float64 `json:"original_price"`
	LowestPrice   float64 `json:"lowest_price"`
}

// PriceCalendar represents price calendar data for date selection
type PriceCalendar struct {
	CruiseID    string              `json:"cruise_id"`
	RouteID     string              `json:"route_id"`
	Year        int                 `json:"year"`
	Month       int                 `json:"month"`
	PriceMatrix [][]CalendarCell    `json:"price_matrix"` // 7 columns (days) x N rows (weeks)
	PriceRange  PriceRangeInfo      `json:"price_range"`
	CabinTypes  []CalendarCabinType `json:"cabin_types"`
	Legend      CalendarLegend      `json:"legend"`
}

// CalendarCell represents a single day in the price calendar
type CalendarCell struct {
	Date           string   `json:"date"`
	Day            int      `json:"day"`
	IsCurrentMonth bool     `json:"is_current_month"`
	IsPast         bool     `json:"is_past"`
	HasVoyage      bool     `json:"has_voyage"`
	MinPrice       *float64 `json:"min_price,omitempty"`
	PriceLevel     string   `json:"price_level,omitempty"` // low, medium, high, premium
	VoyageID       *string  `json:"voyage_id,omitempty"`
	IsSoldOut      bool     `json:"is_sold_out"`
}

// PriceRangeInfo represents price range for the calendar
type PriceRangeInfo struct {
	Min        float64 `json:"min"`
	Max        float64 `json:"max"`
	Avg        float64 `json:"avg"`
	LowestDay  string  `json:"lowest_day"`
	HighestDay string  `json:"highest_day"`
}

// CalendarCabinType represents a cabin type in the calendar
type CalendarCabinType struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Selected bool   `json:"selected"`
}

// CalendarLegend represents the calendar legend
type CalendarLegend struct {
	LowPrice    float64 `json:"low_price_threshold"`
	MediumPrice float64 `json:"medium_price_threshold"`
	HighPrice   float64 `json:"high_price_threshold"`
}

// PriceTrendAnalysis provides price trend analysis functionality
type PriceTrendAnalysis interface {
	// GetPriceTrend gets price trend for a specific voyage
	GetPriceTrend(ctx context.Context, voyageID string) (*PriceTrend, error)

	// GetCruisePriceTrends gets price trends for all voyages of a cruise
	GetCruisePriceTrends(ctx context.Context, cruiseID string, limit int) ([]PriceTrend, error)

	// GetBestBookingWindow recommends optimal booking window
	GetBestBookingWindow(ctx context.Context, cruiseID string, departureMonth int) (*BookingWindow, error)

	// GetPriceCalendar generates price calendar for cruise/route
	GetPriceCalendar(ctx context.Context, cruiseID, routeID string, year, month int, cabinTypeID *string) (*PriceCalendar, error)

	// GetPriceForecast forecasts future price for a voyage
	GetPriceForecast(ctx context.Context, voyageID string) (*PriceForecast, error)

	// CompareVoyagePrices compares prices across multiple voyages
	CompareVoyagePrices(ctx context.Context, voyageIDs []string) ([]VoyagePriceComparison, error)
}

// BookingWindow represents optimal booking recommendations
type BookingWindow struct {
	CruiseID           string  `json:"cruise_id"`
	DepartureMonth     int     `json:"departure_month"`
	OptimalBookBy      string  `json:"optimal_book_by"`
	ExpectedSavings    float64 `json:"expected_savings"`
	ExpectedSavingsPct float64 `json:"expected_savings_pct"`
	Urgency            string  `json:"urgency"` // low, medium, high
	Reason             string  `json:"reason"`
}

// PriceForecast represents a price forecast
type PriceForecast struct {
	VoyageID      string           `json:"voyage_id"`
	CurrentPrice  float64          `json:"current_price"`
	ForecastPrice float64          `json:"forecast_price"`
	Confidence    float64          `json:"confidence"` // 0-1
	Trend         string           `json:"trend"`      // up, down, stable
	Factors       []ForecastFactor `json:"factors"`
	ForecastDate  string           `json:"forecast_date"`
}

// ForecastFactor represents a factor affecting the forecast
type ForecastFactor struct {
	Factor string  `json:"factor"`
	Impact float64 `json:"impact"` // -1 to 1
}

// VoyagePriceComparison represents a price comparison entry
type VoyagePriceComparison struct {
	VoyageID      string  `json:"voyage_id"`
	VoyageNumber  string  `json:"voyage_number"`
	DepartureDate string  `json:"departure_date"`
	MinPrice      float64 `json:"min_price"`
	MaxPrice      float64 `json:"max_price"`
	PriceChange   float64 `json:"price_change"`
	Trend         string  `json:"trend"`
}

// priceTrendAnalyzer implements PriceTrendAnalysis
type priceTrendAnalyzer struct {
	voyageRepo    repository.VoyageRepository
	priceRepo     repository.PriceRepository
	cabinRepo     repository.CabinRepository
	cabinTypeRepo repository.CabinTypeRepository
	cruiseRepo    repository.CruiseRepository
}

// NewPriceTrendAnalysis creates a new price trend analyzer
func NewPriceTrendAnalysis(
	voyageRepo repository.VoyageRepository,
	priceRepo repository.PriceRepository,
	cabinRepo repository.CabinRepository,
	cabinTypeRepo repository.CabinTypeRepository,
	cruiseRepo repository.CruiseRepository,
) PriceTrendAnalysis {
	return &priceTrendAnalyzer{
		voyageRepo:    voyageRepo,
		priceRepo:     priceRepo,
		cabinRepo:     cabinRepo,
		cabinTypeRepo: cabinTypeRepo,
		cruiseRepo:    cruiseRepo,
	}
}

// GetPriceTrend gets price trend for a specific voyage
func (a *priceTrendAnalyzer) GetPriceTrend(ctx context.Context, voyageID string) (*PriceTrend, error) {
	// Get voyage details
	voyage, err := a.voyageRepo.GetByID(ctx, voyageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get voyage: %w", err)
	}

	// Get cruise details
	cruise, err := a.cruiseRepo.GetByID(ctx, voyage.CruiseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cruise: %w", err)
	}

	// Get current prices
	prices, err := a.priceRepo.ListByVoyage(ctx, voyageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get prices: %w", err)
	}

	// Build cabin type price info
	cabinPrices := make([]CabinPriceInfo, 0, len(prices))
	var currentPrice, lowestPrice, highestPrice float64

	for _, price := range prices {
		cabinType, err := a.cabinTypeRepo.GetByID(ctx, price.CabinTypeID)
		if err != nil {
			continue
		}

		cp := CabinPriceInfo{
			CabinTypeID:   price.CabinTypeID,
			CabinTypeName: cabinType.Name,
			CurrentPrice:  price.AdultPrice,
			OriginalPrice: price.AdultPrice, // In real implementation, get from price history
			LowestPrice:   price.AdultPrice,
		}

		cabinPrices = append(cabinPrices, cp)

		// Track overall min/max
		if currentPrice == 0 || price.AdultPrice < currentPrice {
			currentPrice = price.AdultPrice
		}
		if lowestPrice == 0 || price.AdultPrice < lowestPrice {
			lowestPrice = price.AdultPrice
		}
		if price.AdultPrice > highestPrice {
			highestPrice = price.AdultPrice
		}
	}

	// Calculate trend direction
	trend := &PriceTrend{
		VoyageID:        voyageID,
		VoyageNumber:    voyage.VoyageNumber,
		CruiseName:      cruise.NameCN,
		DepartureDate:   voyage.DepartureDate,
		CurrentPrice:    currentPrice,
		OriginalPrice:   currentPrice,
		LowestPrice:     lowestPrice,
		HighestPrice:    highestPrice,
		PriceChange:     0,
		PriceChangePct:  0,
		TrendDirection:  "stable",
		CabinTypePrices: cabinPrices,
		PriceHistory:    []PricePoint{}, // Would come from price history table
	}

	return trend, nil
}

// GetCruisePriceTrends gets price trends for all voyages of a cruise
func (a *priceTrendAnalyzer) GetCruisePriceTrends(ctx context.Context, cruiseID string, limit int) ([]PriceTrend, error) {
	// Get voyages for cruise
	voyages, err := a.voyageRepo.ListByCruise(ctx, cruiseID)
	if err != nil {
		return nil, err
	}

	trends := make([]PriceTrend, 0, len(voyages))
	for _, voyage := range voyages {
		trend, err := a.GetPriceTrend(ctx, voyage.ID.String())
		if err != nil {
			continue
		}
		trends = append(trends, *trend)

		if len(trends) >= limit {
			break
		}
	}

	return trends, nil
}

// GetBestBookingWindow recommends optimal booking window
func (a *priceTrendAnalyzer) GetBestBookingWindow(ctx context.Context, cruiseID string, departureMonth int) (*BookingWindow, error) {
	// Simplified logic - recommend booking 60-90 days in advance for best prices
	now := time.Now()
	targetDate := time.Date(now.Year(), time.Month(departureMonth), 1, 0, 0, 0, 0, time.Local)

	// If month has passed, use next year
	if targetDate.Before(now) {
		targetDate = targetDate.AddDate(1, 0, 0)
	}

	// Optimal booking is 60 days before departure
	optimalBookBy := targetDate.AddDate(0, 0, -60)

	return &BookingWindow{
		CruiseID:           cruiseID,
		DepartureMonth:     departureMonth,
		OptimalBookBy:      optimalBookBy.Format("2006-01-02"),
		ExpectedSavings:    200.0,
		ExpectedSavingsPct: 5.0,
		Urgency:            "low",
		Reason:             "提前60天预订可获得较优价格",
	}, nil
}

// GetPriceCalendar generates price calendar for cruise/route
func (a *priceTrendAnalyzer) GetPriceCalendar(ctx context.Context, cruiseID, routeID string, year, month int, cabinTypeID *string) (*PriceCalendar, error) {
	// Get first day of the month
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)

	// Get voyages in this month
	filters := repository.VoyageFilters{
		CruiseID:      cruiseID,
		RouteID:       routeID,
		DepartureFrom: firstDay.Format("2006-01-02"),
		DepartureTo:   lastDay.Format("2006-01-02"),
	}

	voyages, err := a.voyageRepo.List(ctx, filters, nil)
	if err != nil {
		return nil, err
	}

	// Create a map of date -> voyage
	voyageMap := make(map[string]*domain.Voyage)
	priceMap := make(map[string]float64)

	for _, voyage := range voyages {
		voyageMap[voyage.DepartureDate] = voyage

		// Get min price for this voyage
		prices, err := a.priceRepo.ListByVoyage(ctx, voyage.ID.String())
		if err != nil {
			continue
		}

		var minPrice float64
		for _, p := range prices {
			if minPrice == 0 || p.AdultPrice < minPrice {
				minPrice = p.AdultPrice
			}
		}
		priceMap[voyage.DepartureDate] = minPrice
	}

	// Build calendar matrix
	priceMatrix := make([][]CalendarCell, 0)
	var currentRow []CalendarCell

	// Start from the first Sunday before or on the first day
	startDate := firstDay
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	// Build 6 weeks max
	currentDate := startDate
	globalMin, globalMax := 0.0, 0.0
	lowestDay, highestDay := "", ""

	for week := 0; week < 6; week++ {
		currentRow = make([]CalendarCell, 0, 7)

		for day := 0; day < 7; day++ {
			dateStr := currentDate.Format("2006-01-02")
			isCurrentMonth := currentDate.Month() == time.Month(month)
			isPast := currentDate.Before(time.Now())

			cell := CalendarCell{
				Date:           dateStr,
				Day:            currentDate.Day(),
				IsCurrentMonth: isCurrentMonth,
				IsPast:         isPast,
				HasVoyage:      false,
			}

			if voyage, ok := voyageMap[dateStr]; ok {
				cell.HasVoyage = true
				if price, ok := priceMap[dateStr]; ok && price > 0 {
					cell.MinPrice = &price
					voyageID := voyage.ID.String()
					cell.VoyageID = &voyageID

					// Track global min/max
					if globalMin == 0 || price < globalMin {
						globalMin = price
						lowestDay = dateStr
					}
					if price > globalMax {
						globalMax = price
						highestDay = dateStr
					}
				}
			}

			currentRow = append(currentRow, cell)
			currentDate = currentDate.AddDate(0, 0, 1)
		}

		priceMatrix = append(priceMatrix, currentRow)

		// Stop if we've passed the last day of the month and it's the end of the week
		if currentDate.Month() != time.Month(month) && currentDate.Weekday() == time.Sunday {
			break
		}
	}

	// Assign price levels
	if globalMax > globalMin {
		priceRange := globalMax - globalMin
		for i := range priceMatrix {
			for j := range priceMatrix[i] {
				if priceMatrix[i][j].MinPrice != nil {
					price := *priceMatrix[i][j].MinPrice
					ratio := (price - globalMin) / priceRange

					switch {
					case ratio <= 0.33:
						priceMatrix[i][j].PriceLevel = "low"
					case ratio <= 0.66:
						priceMatrix[i][j].PriceLevel = "medium"
					default:
						priceMatrix[i][j].PriceLevel = "high"
					}
				}
			}
		}
	}

	// Calculate average
	var total float64
	count := 0
	for _, price := range priceMap {
		total += price
		count++
	}
	avg := 0.0
	if count > 0 {
		avg = total / float64(count)
	}

	calendar := &PriceCalendar{
		CruiseID:    cruiseID,
		RouteID:     routeID,
		Year:        year,
		Month:       month,
		PriceMatrix: priceMatrix,
		PriceRange: PriceRangeInfo{
			Min:        globalMin,
			Max:        globalMax,
			Avg:        avg,
			LowestDay:  lowestDay,
			HighestDay: highestDay,
		},
		CabinTypes: []CalendarCabinType{}, // Would populate from cabin types
		Legend: CalendarLegend{
			LowPrice:    globalMin + (globalMax-globalMin)*0.33,
			MediumPrice: globalMin + (globalMax-globalMin)*0.66,
			HighPrice:   globalMax,
		},
	}

	return calendar, nil
}

// GetPriceForecast forecasts future price for a voyage
func (a *priceTrendAnalyzer) GetPriceForecast(ctx context.Context, voyageID string) (*PriceForecast, error) {
	// Get current trend
	trend, err := a.GetPriceTrend(ctx, voyageID)
	if err != nil {
		return nil, err
	}

	// Simple heuristic-based forecast
	forecastPrice := trend.CurrentPrice
	trendDirection := "stable"

	// Parse departure date
	voyage, _ := a.voyageRepo.GetByID(ctx, voyageID)
	if voyage != nil {
		depDate, err := time.Parse("2006-01-02", voyage.DepartureDate)
		if err == nil {
			daysUntil := int(depDate.Sub(time.Now()).Hours() / 24)

			// Prices tend to increase as departure approaches
			if daysUntil < 30 {
				trendDirection = "up"
				forecastPrice = trend.CurrentPrice * 1.05
			} else if daysUntil < 60 {
				trendDirection = "stable"
			} else {
				trendDirection = "down"
				forecastPrice = trend.CurrentPrice * 0.95
			}
		}
	}

	return &PriceForecast{
		VoyageID:      voyageID,
		CurrentPrice:  trend.CurrentPrice,
		ForecastPrice: forecastPrice,
		Confidence:    0.7,
		Trend:         trendDirection,
		Factors: []ForecastFactor{
			{Factor: "booking_lead_time", Impact: 0.3},
			{Factor: "seasonality", Impact: 0.2},
			{Factor: "inventory_level", Impact: 0.2},
		},
		ForecastDate: time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
	}, nil
}

// CompareVoyagePrices compares prices across multiple voyages
func (a *priceTrendAnalyzer) CompareVoyagePrices(ctx context.Context, voyageIDs []string) ([]VoyagePriceComparison, error) {
	comparisons := make([]VoyagePriceComparison, 0, len(voyageIDs))

	for _, voyageID := range voyageIDs {
		trend, err := a.GetPriceTrend(ctx, voyageID)
		if err != nil {
			continue
		}

		comp := VoyagePriceComparison{
			VoyageID:      voyageID,
			VoyageNumber:  trend.VoyageNumber,
			DepartureDate: trend.DepartureDate,
			MinPrice:      trend.LowestPrice,
			MaxPrice:      trend.HighestPrice,
			PriceChange:   trend.PriceChange,
			Trend:         trend.TrendDirection,
		}

		comparisons = append(comparisons, comp)
	}

	// Sort by departure date
	sort.Slice(comparisons, func(i, j int) bool {
		return comparisons[i].DepartureDate < comparisons[j].DepartureDate
	})

	return comparisons, nil
}
