package mailchimp

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"github.com/leapforce-libraries/go_mailchimp/types"
	"net/http"
	"net/url"
	"strings"
)

type Order struct {
	Id       string `json:"id"`
	Customer struct {
		Id           string  `json:"id"`
		EmailAddress string  `json:"email_address"`
		OptInStatus  bool    `json:"opt_in_status"`
		Company      string  `json:"company"`
		FirstName    string  `json:"first_name"`
		LastName     string  `json:"last_name"`
		OrdersCount  int     `json:"orders_count"`
		TotalSpent   float64 `json:"total_spent"`
		Address      struct {
			Address1     string `json:"address1"`
			Address2     string `json:"address2"`
			City         string `json:"city"`
			Province     string `json:"province"`
			ProvinceCode string `json:"province_code"`
			PostalCode   string `json:"postal_code"`
			Country      string `json:"country"`
			CountryCode  string `json:"country_code"`
		} `json:"address"`
		CreatedAt *types.DateTimeString `json:"created_at"`
		UpdatedAt *types.DateTimeString `json:"updated_at"`
		Links     []Link                `json:"_links"`
	} `json:"customer"`
	StoreId            string                `json:"store_id"`
	CampaignId         string                `json:"campaign_id"`
	LandingSite        string                `json:"landing_site"`
	FinancialStatus    string                `json:"financial_status"`
	FulfillmentStatus  string                `json:"fulfillment_status"`
	CurrencyCode       string                `json:"currency_code"`
	OrderTotal         float64               `json:"order_total"`
	OrderUrl           string                `json:"order_url"`
	DiscountTotal      float64               `json:"discount_total"`
	TaxTotal           float64               `json:"tax_total"`
	ShippingTotal      float64               `json:"shipping_total"`
	TrackingCode       string                `json:"tracking_code"`
	ProcessedAtForeign *types.DateTimeString `json:"processed_at_foreign"`
	CancelledAtForeign *types.DateTimeString `json:"cancelled_at_foreign"`
	UpdatedAtForeign   *types.DateTimeString `json:"updated_at_foreign"`
	ShippingAddress    struct {
		Name         string  `json:"name"`
		Address1     string  `json:"address1"`
		Address2     string  `json:"address2"`
		City         string  `json:"city"`
		Province     string  `json:"province"`
		ProvinceCode string  `json:"province_code"`
		PostalCode   string  `json:"postal_code"`
		Country      string  `json:"country"`
		CountryCode  string  `json:"country_code"`
		Longitude    float64 `json:"longitude"`
		Latitude     float64 `json:"latitude"`
		Phone        string  `json:"phone"`
		Company      string  `json:"company"`
	} `json:"shipping_address"`
	BillingAddress struct {
		Name         string  `json:"name"`
		Address1     string  `json:"address1"`
		Address2     string  `json:"address2"`
		City         string  `json:"city"`
		Province     string  `json:"province"`
		ProvinceCode string  `json:"province_code"`
		PostalCode   string  `json:"postal_code"`
		Country      string  `json:"country"`
		CountryCode  string  `json:"country_code"`
		Longitude    float64 `json:"longitude"`
		Latitude     float64 `json:"latitude"`
		Phone        string  `json:"phone"`
		Company      string  `json:"company"`
	} `json:"billing_address"`
	Promos []struct {
		Code             string  `json:"code"`
		AmountDiscounted float64 `json:"amount_discounted"`
		Type             string  `json:"type"`
	} `json:"promos"`
	Lines []struct {
		Id                  string  `json:"id"`
		ProductId           string  `json:"product_id"`
		ProductTitle        string  `json:"product_title"`
		ProductVariantId    string  `json:"product_variant_id"`
		ProductVariantTitle string  `json:"product_variant_title"`
		ImageUrl            string  `json:"image_url"`
		Quantity            float64 `json:"quantity"`
		Price               float64 `json:"price"`
		Discount            float64 `json:"discount"`
		Links               []Link  `json:"_links"`
	} `json:"lines"`
	Outreach struct {
		Id            string                `json:"id"`
		Name          string                `json:"name"`
		Type          string                `json:"type"`
		PublishedTime *types.DateTimeString `json:"published_time"`
	} `json:"outreach"`
	TrackingNumber  string `json:"tracking_number"`
	TrackingCarrier string `json:"tracking_carrier"`
	TrackingUrl     string `json:"tracking_url"`
}

type ListAccountOrdersConfig struct {
	Fields        *[]string
	ExcludeFields *[]string
	Count         *int64
	CampaignId    *string
	OutreachId    *string
	CustomerId    *string
	HasOutreach   *bool
}

type ListAccountOrdersResponse struct {
	Orders     []Order `json:"orders"`
	TotalItems int     `json:"total_items"`
	Links      []Link  `json:"_links"`
}

func (service *Service) ListAccountOrders(cfg *ListAccountOrdersConfig) (*[]Order, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ListAccountOrdersConfig must not be nil")
	}

	var orders []Order

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

	if cfg.CampaignId != nil {
		values.Set("campaign_id", *cfg.CampaignId)
	}

	if cfg.OutreachId != nil {
		values.Set("outreach_id", *cfg.OutreachId)
	}

	if cfg.CustomerId != nil {
		values.Set("customer_id", *cfg.CustomerId)
	}

	if cfg.CampaignId != nil {
		values.Set("campaign_id", *cfg.CampaignId)
	}

	if cfg.HasOutreach != nil {
		values.Set("has_outreach", fmt.Sprintf("%v", *cfg.HasOutreach))
	}

	for {
		var response ListAccountOrdersResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("ecommerce/orders?%s", values.Encode())),
			ResponseModel: &response,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		orders = append(orders, response.Orders...)

		if len(orders) >= response.TotalItems {
			break
		}

		values.Set("offset", fmt.Sprintf("%v", len(orders)))
	}

	return &orders, nil
}
