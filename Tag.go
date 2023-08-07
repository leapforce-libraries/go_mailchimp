package mailchimp

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
)

type Tag struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SearchTagsConfig struct {
	ListId string
	Name   *string
}

type SearchTagsResponse struct {
	Tags       []Tag `json:"tags"`
	TotalItems int   `json:"total_items"`
}

func (service *Service) SearchTags(cfg *SearchTagsConfig) (*[]Tag, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("SearchTagsConfig must not be nil")
	}

	var tags []Tag

	var values = url.Values{}

	if cfg.Name != nil {
		values.Set("name", *cfg.Name)
	}

	for {
		var response SearchTagsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("lists/%s/tag-search?%s", cfg.ListId, values.Encode())),
			ResponseModel: &response,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		tags = append(tags, response.Tags...)

		if len(tags) >= response.TotalItems {
			break
		}

		values.Set("offset", fmt.Sprintf("%v", len(tags)))
	}

	return &tags, nil
}
