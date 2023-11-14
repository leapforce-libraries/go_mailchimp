package mailchimp

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"github.com/leapforce-libraries/go_mailchimp/types"
	"net/http"
	"net/url"
)

type SurveyResponse struct {
	ResponseId   string               `json:"response_id"`
	SubmittedAt  types.DateTimeString `json:"submitted_at"`
	Contact      SurveyContact        `json:"contact"`
	IsNewContact bool                 `json:"is_new_contact"`
	Results      []SurveyResult       `json:"results"`
}

type SurveyContact struct {
	EmailId                     string `json:"email_id"`
	ContactId                   string `json:"contact_id"`
	Status                      string `json:"status"`
	Email                       string `json:"email"`
	FullName                    string `json:"full_name"`
	ConsentsToOneToOneMessaging bool   `json:"consents_to_one_to_one_messaging"`
	AvatarUrl                   string `json:"avatar_url"`
}

type SurveyResult struct {
	QuestionId   string `json:"question_id"`
	QuestionType string `json:"question_type"`
	Query        string `json:"query"`
	Answer       string `json:"answer"`
}

type ListSurveyResponsesConfig struct {
	SurveyId string
}

type ListSurveyResponsesResponse struct {
	SurveyResponses []SurveyResponse `json:"responses"`
	TotalItems      int              `json:"total_items"`
}

func (service *Service) ListSurveyResponses(cfg *ListSurveyResponsesConfig) (*[]SurveyResponse, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ListSurveyResponsesConfig must not be nil")
	}

	var surveyResponses []SurveyResponse

	var values = url.Values{}

	for {
		var response ListSurveyResponsesResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("reporting/surveys/%s/responses?%s", cfg.SurveyId, values.Encode())),
			ResponseModel: &response,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		surveyResponses = append(surveyResponses, response.SurveyResponses...)

		if len(surveyResponses) >= response.TotalItems {
			break
		}

		values.Set("offset", fmt.Sprintf("%v", len(surveyResponses)))
	}

	return &surveyResponses, nil
}

type GetSurveyResponseConfig struct {
	SurveyId  string
	ReponseId string
}

func (service *Service) GetSurveyResponse(cfg *GetSurveyResponseConfig) (*SurveyResponse, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("GetSurveyResponseConfig must not be nil")
	}

	var surveyResponse SurveyResponse

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("reporting/surveys/%s/responses/%s", cfg.SurveyId, cfg.ReponseId)),
		ResponseModel: &surveyResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &surveyResponse, nil
}
