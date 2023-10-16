package mailchimp

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"github.com/leapforce-libraries/go_mailchimp/types"
	"net/http"
	"net/url"
)

type SurveyQuestionAnswer struct {
	Id           string               `json:"id"`
	Value        string               `json:"value"`
	ResponseId   string               `json:"response_id"`
	SubmittedAt  types.DateTimeString `json:"submitted_at"`
	Contact      SurveyContact        `json:"contact"`
	IsNewContact bool                 `json:"is_new_contact"`
}

type ListSurveyQuestionAnswersConfig struct {
	SurveyId         string
	SurveyQuestionId string
}

type ListSurveyQuestionAnswersResponse struct {
	SurveyQuestionAnswers []SurveyQuestionAnswer `json:"answers"`
	TotalItems            int                    `json:"total_items"`
}

func (service *Service) ListSurveyQuestionAnswers(cfg *ListSurveyQuestionAnswersConfig) (*[]SurveyQuestionAnswer, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ListSurveyQuestionAnswersConfig must not be nil")
	}

	var surveyQuestionAnswers []SurveyQuestionAnswer

	var values = url.Values{}

	for {
		var response ListSurveyQuestionAnswersResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("reporting/surveys/%s/questions/%s/answers?%s", cfg.SurveyId, cfg.SurveyQuestionId, values.Encode())),
			ResponseModel: &response,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		surveyQuestionAnswers = append(surveyQuestionAnswers, response.SurveyQuestionAnswers...)

		if len(surveyQuestionAnswers) >= response.TotalItems {
			break
		}

		values.Set("offset", fmt.Sprintf("%v", len(surveyQuestionAnswers)))
	}

	return &surveyQuestionAnswers, nil
}
