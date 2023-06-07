package metrics

type CreateMetricRequest struct {
	MetricType string `json:"type"`
	SubType    string `json:"subtype"`
}
