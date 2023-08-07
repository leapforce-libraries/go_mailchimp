package mailchimp

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"github.com/leapforce-libraries/go_mailchimp/types"
	"net/http"
	"net/url"
)

type Survey struct {
	Id             string                `json:"id"`
	ListId         string                `json:"list_id"`
	WebId          string                `json:"web_id"`
	Title          string                `json:"title"`
	Status         string                `json:"status"`
	HostedUrl      string                `json:"hosted_url"`
	IsPipedToInbox bool                  `json:"is_piped_to_inbox"`
	QuestionCount  int                   `json:"question_count"`
	ResponseCount  int                   `json:"response_count"`
	Questions      []string              `json:"questions"`
	CreatedAt      types.DateTimeString  `json:"created_at"`
	UpdatedAt      *types.DateTimeString `json:"updated_at"`
	PublishedAt    *types.DateTimeString `json:"published_at"`
}

type ListSurveysConfig struct {
	ListId string
}

type ListSurveysResponse struct {
	Surveys    []Survey `json:"surveys"`
	TotalItems int      `json:"total_items"`
}

func (service *Service) ListSurveys(cfg *ListSurveysConfig) (*[]Survey, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ListSurveysConfig must not be nil")
	}

	var surveys []Survey

	var values = url.Values{}

	for {
		var response ListSurveysResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("lists/%s/surveys?%s", cfg.ListId, values.Encode())),
			ResponseModel: &response,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		surveys = append(surveys, response.Surveys...)

		if len(surveys) >= response.TotalItems {
			break
		}

		values.Set("offset", fmt.Sprintf("%v", len(surveys)))
	}

	return &surveys, nil
}
