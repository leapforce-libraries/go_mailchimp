package mailchimp

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"github.com/leapforce-libraries/go_mailchimp/types"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CampaignReport struct {
	Id            string                `json:"id"`
	CampaignTitle string                `json:"campaign_title"`
	Type          string                `json:"type"`
	ListId        string                `json:"list_id"`
	ListIsActive  bool                  `json:"list_is_active"`
	ListName      string                `json:"list_name"`
	SubjectLine   string                `json:"subject_line"`
	PreviewText   string                `json:"preview_text"`
	EmailsSent    int                   `json:"emails_sent"`
	AbuseReports  int                   `json:"abuse_reports"`
	Unsubscribed  int                   `json:"unsubscribed"`
	SendTime      *types.DateTimeString `json:"send_time"`
	RssLastSend   *types.DateTimeString `json:"rss_last_send"`
	Bounces       struct {
		HardBounces  int `json:"hard_bounces"`
		SoftBounces  int `json:"soft_bounces"`
		SyntaxErrors int `json:"syntax_errors"`
	} `json:"bounces"`
	Forwards struct {
		ForwardsCount int `json:"forwards_count"`
		ForwardsOpens int `json:"forwards_opens"`
	} `json:"forwards"`
	Opens struct {
		OpensTotal  int                   `json:"opens_total"`
		UniqueOpens int                   `json:"unique_opens"`
		OpenRate    float64               `json:"open_rate"`
		LastOpen    *types.DateTimeString `json:"last_open"`
	} `json:"opens"`
	Clicks struct {
		ClicksTotal            int                   `json:"clicks_total"`
		UniqueClicks           int                   `json:"unique_clicks"`
		UniqueSubscriberClicks int                   `json:"unique_subscriber_clicks"`
		ClickRate              float64               `json:"click_rate"`
		LastClick              *types.DateTimeString `json:"last_click"`
	} `json:"clicks"`
	FacebookLikes struct {
		RecipientLikes int `json:"recipient_likes"`
		UniqueLikes    int `json:"unique_likes"`
		FacebookLikes  int `json:"facebook_likes"`
	} `json:"facebook_likes"`
	IndustryStats struct {
		Type       string  `json:"type"`
		OpenRate   float64 `json:"open_rate"`
		ClickRate  float64 `json:"click_rate"`
		BounceRate float64 `json:"bounce_rate"`
		UnopenRate float64 `json:"unopen_rate"`
		UnsubRate  float64 `json:"unsub_rate"`
		AbuseRate  float64 `json:"abuse_rate"`
	} `json:"industry_stats"`
	ListStats struct {
		SubRate   float64 `json:"sub_rate"`
		UnsubRate float64 `json:"unsub_rate"`
		OpenRate  float64 `json:"open_rate"`
		ClickRate float64 `json:"click_rate"`
	} `json:"list_stats"`
	AbSplit struct {
		A CampaignReportAb `json:"a"`
		B CampaignReportAb `json:"b"`
	} `json:"ab_split"`
	Timewarp    []CampaignReportTimewarpItem   `json:"timewarp"`
	Timeseries  []CampaignReportTimeseriesItem `json:"timeseries"`
	ShareReport struct {
		ShareUrl      string `json:"share_url"`
		SharePassword string `json:"share_password"`
	} `json:"share_report"`
	Ecommerce struct {
		TotalOrders  int     `json:"total_orders"`
		TotalSpent   float64 `json:"total_spent"`
		TotalRevenue float64 `json:"total_revenue"`
		CurrencyCode string  `json:"currency_code"`
	} `json:"ecommerce"`
	DeliveryStatus struct {
		Enabled        bool   `json:"enabled"`
		CanCancel      bool   `json:"can_cancel"`
		Status         string `json:"status"`
		EmailsSent     int    `json:"emails_sent"`
		EmailsCanceled int    `json:"emails_canceled"`
	} `json:"delivery_status"`
	Links []Link `json:"_links"`
}

type CampaignReportAb struct {
	Bounces         int                   `json:"bounces"`
	AbuseReports    int                   `json:"abuse_reports"`
	Unsubs          int                   `json:"unsubs"`
	RecipientClicks int                   `json:"recipient_clicks"`
	Forwards        int                   `json:"forwards"`
	ForwardsOpens   int                   `json:"forwards_opens"`
	Opens           int                   `json:"opens"`
	LastOpen        *types.DateTimeString `json:"last_open"`
	UniqueOpens     int                   `json:"unique_opens"`
}

type CampaignReportTimewarpItem struct {
	GmtOffset    int                   `json:"gmt_offset"`
	Opens        int                   `json:"opens"`
	LastOpen     *types.DateTimeString `json:"last_open"`
	UniqueOpens  int                   `json:"unique_opens"`
	Clicks       int                   `json:"clicks"`
	LastClick    *types.DateTimeString `json:"last_click"`
	UniqueClicks int                   `json:"unique_clicks"`
	Bounces      int                   `json:"bounces"`
}

type CampaignReportTimeseriesItem struct {
	Timestamp        *types.DateTimeString `json:"timestamp"`
	EmailsSent       int                   `json:"emails_sent"`
	UniqueOpens      int                   `json:"unique_opens"`
	RecipientsClicks int                   `json:"recipients_clicks"`
}

type ListCampaignReportsConfig struct {
	Fields         *[]string
	ExcludeFields  *[]string
	Count          *int64
	Type           *CampaignType
	BeforeSendTime *time.Time
	SinceSendTime  *time.Time
}

type ListCampaignReportsResponse struct {
	CampaignReports []CampaignReport `json:"reports"`
	TotalItems      int              `json:"total_items"`
	Links           []Link           `json:"_links"`
}

func (service *Service) ListCampaignReports(cfg *ListCampaignReportsConfig) (*[]CampaignReport, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ListCampaignReportsConfig must not be bil")
	}

	var campaignReports []CampaignReport

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

	if cfg.Type != nil {
		values.Set("type", string(*cfg.Type))
	}

	if cfg.BeforeSendTime != nil {
		values.Set("before_send_time", (*cfg.BeforeSendTime).Format(types.DateTimeFormat))
	}

	if cfg.SinceSendTime != nil {
		values.Set("since_send_time", (*cfg.SinceSendTime).Format(types.DateTimeFormat))
	}

	for {
		var response ListCampaignReportsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("reports?%s", values.Encode())),
			ResponseModel: &response,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		campaignReports = append(campaignReports, response.CampaignReports...)

		if len(campaignReports) >= response.TotalItems {
			break
		}

		values.Set("offset", fmt.Sprintf("%v", len(campaignReports)))
	}

	return &campaignReports, nil
}
