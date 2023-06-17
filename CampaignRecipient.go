package mailchimp

import (
	"encoding/json"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"github.com/leapforce-libraries/go_mailchimp/types"
	"net/http"
	"net/url"
	"strings"
)

type CampaignRecipient struct {
	EmailId      string `json:"email_id"`
	EmailAddress string `json:"email_address"`
	MergeFields  struct {
		FNAME   string                   `json:"FNAME"`
		LNAME   string                   `json:"LNAME"`
		ADDRESS CampaignRecipientAddress `json:"ADDRESS"`
		PHONE   string                   `json:"PHONE"`
		AGE     string                   `json:"AGE"`
	} `json:"merge_fields"`
	Vip          bool                  `json:"vip"`
	Status       string                `json:"status"`
	OpenCount    int                   `json:"open_count"`
	LastOpen     *types.DateTimeString `json:"last_open"`
	AbsplitGroup string                `json:"absplit_group"`
	GmtOffset    int                   `json:"gmt_offset"`
	CampaignId   string                `json:"campaign_id"`
	ListId       string                `json:"list_id"`
	ListIsActive bool                  `json:"list_is_active"`
	Links        []Link                `json:"_links"`
}

type campaignRecipientAddress struct {
	Addr1   string `json:"addr1"`
	Addr2   string `json:"addr2"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}

type CampaignRecipientAddress campaignRecipientAddress

func (a *CampaignRecipientAddress) UnmarshalJSON(b []byte) error {
	if string(b) == `""` {
		return nil
	}

	var c campaignRecipientAddress
	err := json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	*a = CampaignRecipientAddress(c)

	return nil
}

type ListCampaignRecipientsConfig struct {
	CampaignId    string
	Fields        *[]string
	ExcludeFields *[]string
	Count         *int64
}

type ListCampaignRecipientsResponse struct {
	CampaignRecipients []CampaignRecipient `json:"sent_to"`
	TotalItems         int                 `json:"total_items"`
	Links              []Link              `json:"_links"`
}

func (service *Service) ListCampaignRecipients(cfg *ListCampaignRecipientsConfig) (*[]CampaignRecipient, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ListCampaignRecipientsConfig must not be bil")
	}

	var campaignReports []CampaignRecipient

	var values = url.Values{}

	if cfg.Fields != nil {
		values.Set("fields", strings.Join(*cfg.Fields, ","))
	}

	if cfg.ExcludeFields != nil {
		values.Set("exclude_fields", strings.Join(*cfg.ExcludeFields, ","))
	}

	var count = countDefault
	if cfg.Count != nil {
		count = *cfg.Count
	}
	values.Set("count", fmt.Sprintf("%v", count))

	for {
		var response ListCampaignRecipientsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("reports/%s/sent-to?%s", cfg.CampaignId, values.Encode())),
			ResponseModel: &response,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		campaignReports = append(campaignReports, response.CampaignRecipients...)

		if len(campaignReports) >= response.TotalItems {
			break
		}

		values.Set("offset", fmt.Sprintf("%v", len(campaignReports)))
	}

	return &campaignReports, nil
}
