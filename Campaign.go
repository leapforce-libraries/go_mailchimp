package mailchimp

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
	"time"
)

type Campaign struct {
	Id                string    `json:"id"`
	WebId             int       `json:"web_id"`
	Type              string    `json:"type"`
	CreateTime        time.Time `json:"create_time"`
	ArchiveUrl        string    `json:"archive_url"`
	LongArchiveUrl    string    `json:"long_archive_url"`
	Status            string    `json:"status"`
	EmailsSent        int       `json:"emails_sent"`
	SendTime          time.Time `json:"send_time"`
	ContentType       string    `json:"content_type"`
	NeedsBlockRefresh bool      `json:"needs_block_refresh"`
	Resendable        bool      `json:"resendable"`
	Recipients        struct {
		ListId         string `json:"list_id"`
		ListIsActive   bool   `json:"list_is_active"`
		ListName       string `json:"list_name"`
		SegmentText    string `json:"segment_text"`
		RecipientCount int    `json:"recipient_count"`
		SegmentOpts    struct {
			SavedSegmentId int    `json:"saved_segment_id"`
			Match          string `json:"match"`
			Conditions     []struct {
				ConditionType string `json:"condition_type"`
				Field         string `json:"field"`
				Op            string `json:"op"`
				Value         int    `json:"value"`
			} `json:"conditions"`
		} `json:"segment_opts"`
	} `json:"recipients"`
	Settings struct {
		SubjectLine     string `json:"subject_line"`
		Title           string `json:"title"`
		FromName        string `json:"from_name"`
		ReplyTo         string `json:"reply_to"`
		UseConversation bool   `json:"use_conversation"`
		ToName          string `json:"to_name"`
		FolderId        string `json:"folder_id"`
		Authenticate    bool   `json:"authenticate"`
		AutoFooter      bool   `json:"auto_footer"`
		InlineCss       bool   `json:"inline_css"`
		AutoTweet       bool   `json:"auto_tweet"`
		FbComments      bool   `json:"fb_comments"`
		Timewarp        bool   `json:"timewarp"`
		TemplateId      int    `json:"template_id"`
		DragAndDrop     bool   `json:"drag_and_drop"`
	} `json:"settings"`
	Tracking struct {
		Opens           bool   `json:"opens"`
		HtmlClicks      bool   `json:"html_clicks"`
		TextClicks      bool   `json:"text_clicks"`
		GoalTracking    bool   `json:"goal_tracking"`
		Ecomm360        bool   `json:"ecomm360"`
		GoogleAnalytics string `json:"google_analytics"`
		Clicktale       string `json:"clicktale"`
	} `json:"tracking"`
	ReportSummary struct {
		Opens            int     `json:"opens"`
		UniqueOpens      int     `json:"unique_opens"`
		OpenRate         float64 `json:"open_rate"`
		Clicks           int     `json:"clicks"`
		SubscriberClicks int     `json:"subscriber_clicks"`
		ClickRate        float64 `json:"click_rate"`
		Ecommerce        struct {
			TotalOrders  int `json:"total_orders"`
			TotalSpent   int `json:"total_spent"`
			TotalRevenue int `json:"total_revenue"`
		} `json:"ecommerce"`
	} `json:"report_summary"`
	DeliveryStatus struct {
		Enabled bool `json:"enabled"`
	} `json:"delivery_status"`
	Links []Link `json:"_links"`
}

type ListCampaignsConfig struct {
	Fields           *[]string
	ExcludeFields    *[]string
	Count            *int64
	Type             *string
	Status           *string
	BeforeSendTime   *time.Time
	SinceSendTime    *time.Time
	BeforeCreateTime *time.Time
	SinceCreateTime  *time.Time
	ListId           *string
	FolderId         *string
	MemberId         *string
	SortField        *string
	SortDir          *string
}

type ListCampaignsResponse struct {
	Campaigns  []Campaign `json:"campaigns"`
	TotalItems int        `json:"total_items"`
	Links      []Link     `json:"_links"`
}

func (service *Service) ListCampaigns(cfg *ListCampaignsConfig) (*[]Campaign, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ListCampaignsConfig must not be bil")
	}

	var campaigns []Campaign

	var values = url.Values{}

	for {
		var response ListCampaignsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("campaigns?%s", values.Encode())),
			ResponseModel: &response,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		campaigns = append(campaigns, response.Campaigns...)

		if len(campaigns) >= response.TotalItems {
			break
		}

		values.Set("offset", fmt.Sprintf("%v", len(campaigns)))
	}

	return &campaigns, nil
}