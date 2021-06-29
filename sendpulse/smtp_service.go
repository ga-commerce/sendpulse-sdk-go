package sendpulse

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/dimuska139/sendpulse-sdk-go/sendpulse/models"
	"net/http"
	"strings"
	"time"
)

type SmtpService struct {
	client *Client
}

func newSmtpService(cl *Client) *SmtpService {
	return &SmtpService{client: cl}
}

type SendEmailParams struct {
	Html     string `json:"html,omitempty"`
	Text     string `json:"text,omitempty"`
	Template *struct {
		ID        string                 `json:"id"`
		Variables map[string]interface{} `json:"variables"`
	} `json:"template"`
	AutoPlainText bool   `json:"auto_plain_text"`
	Subject       string `json:"subject"`
	From          struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"from"`
	To struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"to"`
	Attachments map[string]string `json:"attachments"`
}

func (service *SmtpService) Send(params SendEmailParams) (string, error) {
	path := "/smtp/emails"

	type paramsFormat struct {
		Email SendEmailParams `json:"email"`
	}

	if params.Html != "" {
		html := b64.StdEncoding.EncodeToString([]byte(params.Html))
		params.Html = html
	}

	data := paramsFormat{Email: params}

	var response struct {
		Result bool   `json:"result"`
		ID     string `json:"id"`
	}
	_, err := service.client.NewRequest(http.MethodPost, fmt.Sprintf(path), data, &response, true)
	return response.ID, err
}

type SmtpMessage struct {
	ID                    string              `json:"id"`
	Sender                string              `json:"sender"`
	TotalSize             int                 `json:"total_size"`
	SenderIP              string              `json:"sender_ip"`
	SmtpAnswerCode        int                 `json:"smtp_answer_code"`
	SmtpAnswerCodeExplain string              `json:"smtp_answer_code_explain"`
	SmtpAnswerSubcode     string              `json:"smtp_answer_subcode"`
	SmtpAnswerData        string              `json:"smtp_answer_data"`
	UsedIP                string              `json:"used_ip"`
	Recipient             string              `json:"recipient"`
	Subject               string              `json:"subject"`
	SendDate              models.DateTimeType `json:"subject"`
	Tracking              struct {
		Click int `json:"click"`
		Open  int `json:"open"`
		Link  []*struct {
			Url              string              `json:"url"`
			Browser          string              `json:"browser"`
			Os               string              `json:"os"`
			ScreenResolution string              `json:"screen_resolution"`
			IP               string              `json:"ip"`
			ActionDate       models.DateTimeType `json:"action_date"`
		} `json:"link"`
		ClientInfo []*struct {
			Browser    string              `json:"browser"`
			Os         string              `json:"os"`
			IP         string              `json:"ip"`
			ActionDate models.DateTimeType `json:"action_date"`
		} `json:"client_info"`
	} `json:"tracking"`
}

type SmtpListParams struct {
	Limit     int
	Offset    int
	From      time.Time
	To        time.Time
	Sender    string
	Recipient string
}

func (service *SmtpService) List(params SmtpListParams) ([]*SmtpMessage, error) {
	path := "/smtp/emails"
	var urlParts []string
	urlParts = append(urlParts, fmt.Sprintf("offset=%d", params.Offset))
	if params.Limit != 0 {
		urlParts = append(urlParts, fmt.Sprintf("limit=%d", params.Limit))
	}
	if !params.From.IsZero() {
		urlParts = append(urlParts, fmt.Sprintf("from=%s", params.From.Format("2006-01-02")))
	}
	if !params.To.IsZero() {
		urlParts = append(urlParts, fmt.Sprintf("to=%s", params.From.Format("2006-01-02")))
	}
	if params.Sender != "" {
		urlParts = append(urlParts, fmt.Sprintf("sender=%s", params.Sender))
	}
	if params.Recipient != "" {
		urlParts = append(urlParts, fmt.Sprintf("recipient=%s", params.Recipient))
	}

	if len(urlParts) != 0 {
		path += "?" + strings.Join(urlParts, "&")
	}

	var respData []*SmtpMessage
	_, err := service.client.NewRequest(http.MethodGet, fmt.Sprintf(path), nil, &respData, true)
	return respData, err
}

func (service *SmtpService) Total() (int, error) {
	path := "/smtp/emails/total"
	var respData struct {
		Total int `json:"total"`
	}
	_, err := service.client.NewRequest(http.MethodGet, fmt.Sprintf(path), nil, &respData, true)
	return respData.Total, err
}

func (service *SmtpService) Get(id int) (*SmtpMessage, error) {
	path := fmt.Sprintf("/smtp/emails/%d", id)
	var respData *SmtpMessage
	_, err := service.client.NewRequest(http.MethodGet, fmt.Sprintf(path), nil, &respData, true)
	return respData, err
}

type BouncesList struct {
	Total  int `json:"total"`
	Emails []struct {
		EmailTo           string              `json:"email_to"`
		Sender            string              `json:"sender"`
		SendDate          models.DateTimeType `json:"send_date"`
		Subject           string              `json:"subject"`
		SmtpAnswerCode    int                 `json:"smtp_answer_code"`
		SmtpAnswerSubcode string              `json:"smtp_answer_subcode"`
		SmtpAnswerData    string              `json:"smtp_answer_data"`
	} `json:"emails"`
	RequestLimit int `json:"request_limit"`
	Found        int `json:"found"`
}

func (service *SmtpService) DailyBounces(limit, offset int, date time.Time) (*BouncesList, error) {
	path := fmt.Sprintf("/smtp/bounces/day?limit=%d&offset=%d", limit, offset)
	if !date.IsZero() {
		path += "&date=" + date.Format("2006-01-02")
	}

	var respData *BouncesList
	_, err := service.client.NewRequest(http.MethodGet, fmt.Sprintf(path), nil, &respData, true)
	return respData, err
}

func (service *SmtpService) TotalBounces() (int, error) {
	path := "/smtp/bounces/day/total"

	var respData struct {
		Total int `json:"total"`
	}
	_, err := service.client.NewRequest(http.MethodGet, fmt.Sprintf(path), nil, &respData, true)
	return respData.Total, err
}

type SmtpUnsubscribeEmail struct {
	Email   string `json:"email"`
	Comment string `json:"comment"`
}

func (service *SmtpService) Unsubscribe(emails []*SmtpUnsubscribeEmail) error {
	path := "/smtp/unsubscribe"

	type paramsFormat struct {
		Emails []*SmtpUnsubscribeEmail `json:"emails"`
	}

	data := paramsFormat{Emails: emails}

	var respData struct {
		Result bool `json:"true"`
	}
	_, err := service.client.NewRequest(http.MethodPost, fmt.Sprintf(path), data, &respData, true)
	return err
}

func (service *SmtpService) DeleteUnsubscribed(emails []string) error {
	path := "/smtp/unsubscribe"

	type paramsFormat struct {
		Emails []string `json:"emails"`
	}

	data := paramsFormat{Emails: emails}

	var respData struct {
		Result bool `json:"true"`
	}
	_, err := service.client.NewRequest(http.MethodDelete, fmt.Sprintf(path), data, &respData, true)
	return err
}

type UnsubscribedListParams struct {
	Limit  int
	Offset int
	Date   time.Time
}

type Unsubscribed struct {
	Email             string              `json:"email"`
	UnsubscribeByLink int                 `json:"unsubscribe_by_link"`
	UnsubscribeByUser int                 `json:"unsubscribe_by_user"`
	SpamComplaint     int                 `json:"spam_complaint"`
	Date              models.DateTimeType `json:"date"`
}

func (service *SmtpService) UnsubscribedList(params UnsubscribedListParams) ([]Unsubscribed, error) {
	path := fmt.Sprintf("/smtp/unsubscribe?offset=%d", params.Offset)
	if params.Limit != 0 {
		path += fmt.Sprintf("&limit=%d", params.Limit)
	}

	if !params.Date.IsZero() {
		path += "&date=" + params.Date.Format("2006-01-02")
	}

	var respData []Unsubscribed

	_, err := service.client.NewRequest(http.MethodGet, fmt.Sprintf(path), nil, &respData, true)
	return respData, err
}

func (service *SmtpService) SendersIPs() ([]string, error) {
	path := "/smtp/ips"

	var respData []string
	_, err := service.client.NewRequest(http.MethodGet, fmt.Sprintf(path), nil, &respData, true)
	return respData, err
}

func (service *SmtpService) SendersEmails() ([]string, error) {
	path := "/smtp/senders"

	var respData []string
	_, err := service.client.NewRequest(http.MethodGet, fmt.Sprintf(path), nil, &respData, true)
	return respData, err
}

func (service *SmtpService) AllowedDomains() ([]string, error) {
	path := "/smtp/domains"

	var respData []string
	_, err := service.client.NewRequest(http.MethodGet, fmt.Sprintf(path), nil, &respData, true)
	return respData, err
}

func (service *SmtpService) AddDomain(email string) error {
	path := "/smtp/domains"

	type data struct {
		Email string `json:"email"`
	}

	var respData struct {
		Result bool `json:"result"`
	}
	params := data{Email: email}

	_, err := service.client.NewRequest(http.MethodPost, fmt.Sprintf(path), params, &respData, true)
	return err
}

func (service *SmtpService) VerifyDomain(email string) error {
	path := fmt.Sprintf("/domains/%s", email)

	var respData struct {
		Result bool `json:"result"`
	}

	_, err := service.client.NewRequest(http.MethodGet, fmt.Sprintf(path), nil, &respData, true)
	return err
}
