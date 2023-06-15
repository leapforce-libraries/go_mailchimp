package mailchimp

type Link struct {
	Rel          string `json:"rel"`
	Href         string `json:"href"`
	Method       string `json:"method"`
	TargetSchema string `json:"targetSchema,omitempty"`
	Schema       string `json:"schema,omitempty"`
}
