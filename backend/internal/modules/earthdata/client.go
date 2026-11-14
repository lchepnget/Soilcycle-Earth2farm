package earthdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

const (
	nasaPOWERBaseURL = "https://power.larc.nasa.gov/api/temporal/daily/point"
	nasaPOWERSource  = "NASA POWER"
)

// POWERClient communicates with the NASA POWER API.
type POWERClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewPOWERClient creates a NASA POWER API client.
func NewPOWERClient() *POWERClient {
	return &POWERClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: nasaPOWERBaseURL,
	}
}

// NASA POWER response structure.
type powerResponse struct {
	Properties struct {
		Parameter map[string]map[string]float64 `json:"parameter"`
	} `json:"properties"`
}

// FetchDaily retrieves daily environmental observations.
func (c *POWERClient) FetchDaily(
	ctx context.Context,
	req EarthdataRequest,
) ([]ClimateObservation, error) {

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("parse start date: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("parse end date: %w", err)
	}

	params := url.Values{}

	params.Set(
		"parameters",
		"PRECTOTCORR,T2M,T2M_MAX,T2M_MIN,RH2M",
	)
	params.Set("community", "AG")
	params.Set("longitude", strconv.FormatFloat(req.Longitude, 'f', 4, 64))
	params.Set("latitude", strconv.FormatFloat(req.Latitude, 'f', 4, 64))
	params.Set("start", startDate.Format("20060102"))
	params.Set("end", endDate.Format("20060102"))
	params.Set("format", "JSON")
	params.Set("time-standard", "UTC")

	apiURL := c.baseURL + "?" + params.Encode()

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		apiURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create NASA request: %w", err)
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "SoilCycle-Earth2Farm/1.0")

	response, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("NASA POWER request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return nil, fmt.Errorf(
			"NASA POWER returned status %d: %s",
			response.StatusCode,
			string(body),
		)
	}

	var result powerResponse

	if err := json.NewDecoder(
		io.LimitReader(response.Body, 10<<20),
	).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode NASA response: %w", err)
	}

	parameters := result.Properties.Parameter

	if len(parameters) == 0 {
		return nil, fmt.Errorf("NASA POWER returned no parameter data")
	}

	rainfall := parameters["PRECTOTCORR"]
	temperature := parameters["T2M"]
	maxTemperature := parameters["T2M_MAX"]
	minTemperature := parameters["T2M_MIN"]
	humidity := parameters["RH2M"]

	// Collect all observation dates.
	dateSet := make(map[string]struct{})

	for date := range rainfall {
		dateSet[date] = struct{}{}
	}
	for date := range temperature {
		dateSet[date] = struct{}{}
	}
	for date := range maxTemperature {
		dateSet[date] = struct{}{}
	}
	for date := range minTemperature {
		dateSet[date] = struct{}{}
	}
	for date := range humidity {
		dateSet[date] = struct{}{}
	}

	dates := make([]string, 0, len(dateSet))

	for date := range dateSet {
		dates = append(dates, date)
	}

	sort.Strings(dates)

	observations := make([]ClimateObservation, 0, len(dates))

	for _, dateKey := range dates {
		observationDate, err := time.Parse("20060102", dateKey)
		if err != nil {
			continue
		}

		observation := ClimateObservation{
			Date:            observationDate,
			RainfallMM:      validMeasurement(rainfall, dateKey),
			TemperatureC:    validMeasurement(temperature, dateKey),
			MaxTemperatureC: validMeasurement(maxTemperature, dateKey),
			MinTemperatureC: validMeasurement(minTemperature, dateKey),
			HumidityPercent: validMeasurement(humidity, dateKey),
			Source:          nasaPOWERSource,
		}

		observations = append(observations, observation)
	}

	return observations, nil
}

// validMeasurement converts NASA missing-value indicators to nil.
func validMeasurement(
	values map[string]float64,
	date string,
) *float64 {

	value, exists := values[date]

	if !exists || value == -999 {
		return nil
	}

	return &value
}
