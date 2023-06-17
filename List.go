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

type List struct {
	Id                   string               `json:"id"`
	WebId                int                  `json:"web_id"`
	Name                 string               `json:"name"`
	Contact              Contact              `json:"contact"`
	PermissionReminder   string               `json:"permission_reminder"`
	UseArchiveBar        bool                 `json:"use_archive_bar"`
	CampaignDefaults     CampaignDefaults     `json:"campaign_defaults"`
	NotifyOnSubscribe    string               `json:"notify_on_subscribe"`
	NotifyOnUnsubscribe  string               `json:"notify_on_unsubscribe"`
	DateCreated          types.DateTimeString `json:"date_created"`
	ListRating           int                  `json:"list_rating"`
	EmailTypeOption      bool                 `json:"email_type_option"`
	SubscribeUrlShort    string               `json:"subscribe_url_short"`
	SubscribeUrlLong     string               `json:"subscribe_url_long"`
	BeamerAddress        string               `json:"beamer_address"`
	Visibility           string               `json:"visibility"`
	DoubleOptin          bool                 `json:"double_optin"`
	HasWelcome           bool                 `json:"has_welcome"`
	MarketingPermissions bool                 `json:"marketing_permissions"`
	Modules              []string             `json:"modules"`
	Stats                ListStats            `json:"stats"`
	Links                []Link               `json:"_links"`
}

type CampaignDefaults struct {
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	Subject   string `json:"subject"`
	Language  string `json:"language"`
}

type ListStats struct {
	MemberCount               int                   `json:"member_count"`
	UnsubscribeCount          int                   `json:"unsubscribe_count"`
	CleanedCount              int                   `json:"cleaned_count"`
	MemberCountSinceSend      int                   `json:"member_count_since_send"`
	UnsubscribeCountSinceSend int                   `json:"unsubscribe_count_since_send"`
	CleanedCountSinceSend     int                   `json:"cleaned_count_since_send"`
	CampaignCount             int                   `json:"campaign_count"`
	CampaignLastSent          *types.DateTimeString `json:"campaign_last_sent"`
	MergeFieldCount           int                   `json:"merge_field_count"`
	AvgSubRate                int                   `json:"avg_sub_rate"`
	AvgUnsubRate              int                   `json:"avg_unsub_rate"`
	TargetSubRate             int                   `json:"target_sub_rate"`
	OpenRate                  float64               `json:"open_rate"`
	ClickRate                 float64               `json:"click_rate"`
	LastSubDate               *types.DateTimeString `json:"last_sub_date"`
	LastUnsubDate             *types.DateTimeString `json:"last_unsub_date"`
}

type ListListsConfig struct {
	Fields                 *[]string
	ExcludeFields          *[]string
	Count                  *int64
	BeforeDateCreated      *time.Time
	SinceDateCreated       *time.Time
	BeforeCampaignLastSent *time.Time
	SinceCampaignLastSent  *time.Time
	Email                  *string
	SortField              *string
	SortDir                *string
	HasEcommerceStore      *bool
	IncludeTotalContacts   *bool
}

type ListListsResponse struct {
	Lists       []List `json:"lists"`
	TotalItems  int    `json:"total_items"`
	Constraints struct {
		MayCreate             bool `json:"may_create"`
		MaxInstances          int  `json:"max_instances"`
		CurrentTotalInstances int  `json:"current_total_instances"`
	} `json:"constraints"`
	Links []Link `json:"_links"`
}

func (service *Service) ListLists(cfg *ListListsConfig) (*[]List, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ListListsConfig must not be bil")
	}

	var lists []List

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

	if cfg.BeforeDateCreated != nil {
		values.Set("before_date_created", (*cfg.BeforeDateCreated).Format(types.DateTimeFormat))
	}

	if cfg.SinceDateCreated != nil {
		values.Set("since_date_created", (*cfg.SinceDateCreated).Format(types.DateTimeFormat))
	}

	if cfg.BeforeCampaignLastSent != nil {
		values.Set("before_campaign_last_sent", (*cfg.BeforeCampaignLastSent).Format(types.DateTimeFormat))
	}

	if cfg.SinceCampaignLastSent != nil {
		values.Set("since_campaign_last_sent", (*cfg.SinceCampaignLastSent).Format(types.DateTimeFormat))
	}

	if cfg.Email != nil {
		values.Set("email", *cfg.Email)
	}

	if cfg.SortField != nil {
		values.Set("sort_field", *cfg.SortField)
	}

	if cfg.SortDir != nil {
		values.Set("sort_dir", *cfg.SortDir)
	}

	if cfg.HasEcommerceStore != nil {
		values.Set("has_ecommerce_store", fmt.Sprintf("%v", *cfg.HasEcommerceStore))
	}

	if cfg.IncludeTotalContacts != nil {
		values.Set("include_total_contacts", fmt.Sprintf("%v", *cfg.IncludeTotalContacts))
	}

	for {
		var response ListListsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("lists?%s", values.Encode())),
			ResponseModel: &response,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		lists = append(lists, response.Lists...)

		if len(lists) >= response.TotalItems {
			break
		}

		values.Set("offset", fmt.Sprintf("%v", len(lists)))
	}

	return &lists, nil
}
