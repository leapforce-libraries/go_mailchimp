package mailchimp

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
)

type SurveyQuestion struct {
	Id                       string                 `json:"id"`
	SurveyId                 string                 `json:"survey_id"`
	Query                    string                 `json:"query"`
	Type                     string                 `json:"type"`
	TotalResponses           int                    `json:"total_responses"`
	IsRequired               bool                   `json:"is_required"`
	HasOther                 bool                   `json:"has_other"`
	OtherLabel               string                 `json:"other_label"`
	RangeLowLabel            string                 `json:"range_low_label"`
	RangeHighLabel           string                 `json:"range_high_label"`
	PlaceholderLabel         string                 `json:"placeholder_label"`
	SubscribeCheckboxEnabled bool                   `json:"subscribe_checkbox_enabled"`
	SubscribeCheckboxLabel   string                 `json:"subscribe_checkbox_label"`
	Options                  []SurveyQuestionOption `json:"options"`
}

type SurveyQuestionOption struct {
	Label string `json:"label"`
	Id    string `json:"id"`
	Count int    `json:"count"`
}

type ListSurveyQuestionsConfig struct {
	SurveyId string
}

type ListSurveyQuestionsResponse struct {
	SurveyQuestions []SurveyQuestion `json:"questions"`
	TotalItems      int              `json:"total_items"`
}

func (service *Service) ListSurveyQuestions(cfg *ListSurveyQuestionsConfig) (*[]SurveyQuestion, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ListSurveyQuestionsConfig must not be nil")
	}

	var surveyQuestions []SurveyQuestion

	var values = url.Values{}

	for {
		var response ListSurveyQuestionsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("reporting/surveys/%s/questions?%s", cfg.SurveyId, values.Encode())),
			ResponseModel: &response,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		surveyQuestions = append(surveyQuestions, response.SurveyQuestions...)

		if len(surveyQuestions) >= response.TotalItems {
			break
		}

		values.Set("offset", fmt.Sprintf("%v", len(surveyQuestions)))
	}

	return &surveyQuestions, nil
}
