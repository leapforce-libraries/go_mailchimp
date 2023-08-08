package mailchimp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"github.com/leapforce-libraries/go_mailchimp/types"
)

type ListMember struct {
	Id                          string                     `json:"id"`
	EmailAddress                string                     `json:"email_address"`
	UniqueEmailId               string                     `json:"unique_email_id"`
	ContactId                   string                     `json:"contact_id"`
	FullName                    string                     `json:"full_name"`
	WebId                       int                        `json:"web_id"`
	EmailType                   string                     `json:"email_type"`
	Status                      string                     `json:"status"`
	ConsentsToOneToOneMessaging bool                       `json:"consents_to_one_to_one_messaging"`
	MergeFields                 map[string]json.RawMessage `json:"merge_fields"`
	Interests                   map[string]bool            `json:"interests"`
	Stats                       ListMemberStats            `json:"stats"`
	IpSignup                    string                     `json:"ip_signup"`
	TimestampSignup             string                     `json:"timestamp_signup"`
	IpOpt                       string                     `json:"ip_opt"`
	TimestampOpt                *types.DateTimeString      `json:"timestamp_opt"`
	MemberRating                int                        `json:"member_rating"`
	LastChanged                 *types.DateTimeString      `json:"last_changed"`
	Language                    string                     `json:"language"`
	Vip                         bool                       `json:"vip"`
	EmailClient                 string                     `json:"email_client"`
	Location                    Location                   `json:"location"`
	Source                      string                     `json:"source"`
	TagsCount                   int                        `json:"tags_count"`
	Tags                        []Tag                      `json:"tags"`
	ListId                      string                     `json:"list_id"`
	Links                       []Link                     `json:"_links"`
}

type ListMemberStats struct {
	AvgOpenRate   float64 `json:"avg_open_rate"`
	AvgClickRate  float64 `json:"avg_click_rate"`
	EcommerceData struct {
		TotalRevenue   float64 `json:"total_revenue"`
		NumberOfOrders int     `json:"number_of_orders"`
		CurrencyCode   string  `json:"currency_code"`
	} `json:"ecommerce_data"`
}

type Location struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Gmtoff      int     `json:"gmtoff"`
	Dstoff      int     `json:"dstoff"`
	CountryCode string  `json:"country_code"`
	Timezone    string  `json:"timezone"`
	Region      string  `json:"region"`
}

type ListListMembersConfig struct {
	ListId             string
	Fields             *[]string
	ExcludeFields      *[]string
	Count              *int64
	EmailType          *string
	Status             *string
	SinceTimestampOpt  *time.Time
	BeforeTimestampOpt *time.Time
	SinceLastChanged   *time.Time
	BeforeLastChanged  *time.Time
	UniqueEmailId      *string
	VipOnly            *bool
	InterestCategoryId *string
	InterestIds        *string
	InterestMatch      *string
	SortField          *string
	SortDir            *string
	SinceLastCampaign  *bool
	UnsubscribedSince  *time.Time
}

type ListListMembersResponse struct {
	ListMembers []ListMember `json:"members"`
	TotalItems  int          `json:"total_items"`
	Links       []Link       `json:"_links"`
}

func (service *Service) ListListMembers(cfg *ListListMembersConfig) (*[]ListMember, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ListListMembersConfig must not be nil")
	}

	var listMembers []ListMember

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

	if cfg.EmailType != nil {
		values.Set("email_type", *cfg.EmailType)
	}

	if cfg.Status != nil {
		values.Set("status", *cfg.Status)
	}

	if cfg.SinceTimestampOpt != nil {
		values.Set("since_timestamp_opt", (*cfg.SinceTimestampOpt).Format(types.DateTimeFormat))
	}

	if cfg.BeforeTimestampOpt != nil {
		values.Set("before_timestamp_opt", (*cfg.BeforeTimestampOpt).Format(types.DateTimeFormat))
	}

	if cfg.SinceLastChanged != nil {
		values.Set("since_last_changed", (*cfg.SinceLastChanged).Format(types.DateTimeFormat))
	}

	if cfg.BeforeLastChanged != nil {
		values.Set("before_last_changed", (*cfg.BeforeLastChanged).Format(types.DateTimeFormat))
	}

	if cfg.UniqueEmailId != nil {
		values.Set("unique_email_id", *cfg.UniqueEmailId)
	}

	if cfg.VipOnly != nil {
		values.Set("vip_only", fmt.Sprintf("%v", *cfg.VipOnly))
	}

	if cfg.InterestCategoryId != nil {
		values.Set("interest_category_id", *cfg.InterestCategoryId)
	}

	if cfg.InterestIds != nil {
		values.Set("interest_ids", *cfg.InterestIds)
	}

	if cfg.InterestMatch != nil {
		values.Set("interest_match", *cfg.InterestMatch)
	}

	if cfg.SortField != nil {
		values.Set("sort_field", *cfg.SortField)
	}

	if cfg.SortDir != nil {
		values.Set("sort_dir", *cfg.SortDir)
	}

	if cfg.SinceLastCampaign != nil {
		values.Set("since_last_campaign", fmt.Sprintf("%v", *cfg.SinceLastCampaign))
	}

	if cfg.UnsubscribedSince != nil {
		values.Set("unsubscribed_since", (*cfg.UnsubscribedSince).Format(types.DateTimeFormat))
	}

	for {
		var response ListListMembersResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("lists/%s/members?%s", cfg.ListId, values.Encode())),
			ResponseModel: &response,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		listMembers = append(listMembers, response.ListMembers...)

		if len(listMembers) >= response.TotalItems {
			break
		}

		values.Set("offset", fmt.Sprintf("%v", len(listMembers)))
	}

	return &listMembers, nil
}
