package mailchimp

// ErrorResponse stores general API error response
type ErrorResponse struct {
	Status        string            `json:"status"`
	Message       string            `json:"message"`
	CorrelationId string            `json:"correlationId"`
	Category      string            `json:"category"`
	Links         map[string]string `json:"links"`
}
