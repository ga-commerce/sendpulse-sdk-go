package sendpulse_sdk_go

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type BotsLiveChatService struct {
	client *Client
}

func newBotsLiveChatService(cl *Client) *BotsLiveChatService {
	return &BotsLiveChatService{client: cl}
}

type LiveChatAccount struct {
	Tariff struct {
		Code         string    `json:"code"`
		MaxBots      int       `json:"max_bots"`
		MaxContacts  int       `json:"max_contacts"`
		MaxMessages  int       `json:"max_messages"`
		MaxTags      int       `json:"max_tags"`
		MaxVariables int       `json:"max_variables"`
		Branding     bool      `json:"branding"`
		IsExceeded   bool      `json:"is_exceeded"`
		IsExpired    bool      `json:"is_expired"`
		ExpiredAt    time.Time `json:"expired_at"`
	} `json:"tariff"`
	Statistics struct {
		Messages  int `json:"messages"`
		Bots      int `json:"bots"`
		Contacts  int `json:"contacts"`
		Variables int `json:"variables"`
	} `json:"statistics"`
}

func (service *BotsLiveChatService) GetAccount(ctx context.Context) (*LiveChatAccount, error) {
	path := "/live-chat/account"

	var respData struct {
		Success bool             `json:"success"`
		Data    *LiveChatAccount `json:"data"`
	}
	_, err := service.client.newRequest(ctx, http.MethodGet, path, nil, &respData, true)
	return respData.Data, err
}

type LiveChatBot struct {
	ID          string `json:"id"`
	ChannelData struct {
		ID         string `json:"id"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Name       string `json:"name"`
		NameFormat string `json:"name_format"`
		ShortName  string `json:"short_name"`
		Picture    struct {
			Data struct {
				Height       int    `json:"height"`
				IsSilhouette bool   `json:"is_silhouette"`
				Url          string `json:"url"`
				Width        int    `json:"width"`
			} `json:"data"`
		} `json:"picture"`
	} `json:"channel_data"`
	LiveChatUser struct {
		ID                int    `json:"id"`
		LcID              int    `json:"lc_id"`
		ProfilePictureUrl string `json:"profile_picture_url"`
		Username          string `json:"username"`
		Website           string `json:"website"`
	} `json:"ig_user"`
	Inbox struct {
		Total  int `json:"total"`
		Unread int `json:"unread"`
	} `json:"inbox"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (service *BotsLiveChatService) GetBots(ctx context.Context) ([]*LiveChatBot, error) {
	path := "/live-chat/bots"

	var respData struct {
		Success bool           `json:"success"`
		Data    []*LiveChatBot `json:"data"`
	}
	_, err := service.client.newRequest(ctx, http.MethodGet, path, nil, &respData, true)
	return respData.Data, err
}

type LiveChatBotContact struct {
	ID          string `json:"id"`
	BotID       string `json:"bot_id"`
	Status      int    `json:"status"`
	ChannelData struct {
		ID         string `json:"id"`
		UserName   string `json:"user_name"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Name       string `json:"name"`
		ProfilePic string `json:"profile_pic"`
	} `json:"channel_data"`
	Tags                  []string               `json:"tags"`
	Variables             map[string]interface{} `json:"variables"`
	IsChatOpened          bool                   `json:"is_chat_opened"`
	LastActivityAt        time.Time              `json:"last_activity_at"`
	AutomationPausedUntil time.Time              `json:"automation_paused_until"`
	CreatedAt             time.Time              `json:"created_at"`
}

func (service *BotsLiveChatService) GetContact(ctx context.Context, contactID string) (*LiveChatBotContact, error) {
	path := fmt.Sprintf("/live-chat/contacts/get?id=%s", contactID)

	var respData struct {
		Success bool                `json:"success"`
		Data    *LiveChatBotContact `json:"data"`
	}
	_, err := service.client.newRequest(ctx, http.MethodGet, path, nil, &respData, true)
	return respData.Data, err
}

func (service *BotsLiveChatService) GetContactsByTag(ctx context.Context, tag, botID string) ([]*LiveChatBotContact, error) {
	path := fmt.Sprintf("/live-chat/contacts/getByTag?tag=%s&bot_id=%s", tag, botID)

	var respData struct {
		Success bool                  `json:"success"`
		Data    []*LiveChatBotContact `json:"data"`
	}
	_, err := service.client.newRequest(ctx, http.MethodGet, path, nil, &respData, true)
	return respData.Data, err
}

func (service *BotsLiveChatService) GetContactsByVariable(ctx context.Context, params BotContactsByVariableParams) ([]*LiveChatBotContact, error) {
	urlParams := url.Values{}
	urlParams.Add("variable_value", params.VariableValue)
	if params.VariableID != "" {
		urlParams.Add("variable_id", params.VariableID)
	}
	if params.VariableName != "" {
		urlParams.Add("variable_name", params.VariableName)
	}
	if params.BotID != "" {
		urlParams.Add("bot_id", params.BotID)
	}
	path := "/live-chat/contacts/getByVariable?" + urlParams.Encode()

	var respData struct {
		Success bool                  `json:"success"`
		Data    []*LiveChatBotContact `json:"data"`
	}
	_, err := service.client.newRequest(ctx, http.MethodGet, path, nil, &respData, true)
	return respData.Data, err
}

type LiveChatBotSendMessagesParams struct {
	ContactID string `json:"contact_id"`
	Messages  []struct {
		Type    string `json:"type"`
		Message struct {
			Text string `json:"text"`
		} `json:"message"`
	} `json:"messages"`
}

func (service *BotsLiveChatService) SendTextByContact(ctx context.Context, params LiveChatBotSendMessagesParams) error {
	path := "/live-chat/contacts/send"

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, params, &respData, true)
	return err
}

type LiveChatBotSendMessageAttachmentPayload struct {
	IsExternalAttachment bool   `json:"is_external_attachment"`
	Url                  string `json:"url"`
}

type LiveChatBotSendMessageAttachment struct {
	Type    string                                  `json:"type"`
	Payload LiveChatBotSendMessageAttachmentPayload `json:"payload"`
}

type LiveChatBotSendMessage struct {
	Attachment LiveChatBotSendMessageAttachment `json:"attachment"`
}

type LiveChatBotSendMessages struct {
	Type    string                 `json:"type"`
	Message LiveChatBotSendMessage `json:"message"`
}

type LiveChatBotSendImageMessagesParams struct {
	ContactID string                    `json:"contact_id"`
	Messages  []LiveChatBotSendMessages `json:"messages"`
}

func (service *BotsLiveChatService) SendImageByContact(ctx context.Context, params LiveChatBotSendImageMessagesParams) error {
	path := "/live-chat/contacts/send"

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, params, &respData, true)
	return err
}

func (service *BotsLiveChatService) SetVariableToContact(ctx context.Context, contactID string, variableID string, variableName string, variableValue interface{}) error {
	path := "/live-chat/contacts/setVariable"

	type bodyFormat struct {
		ContactID     string      `json:"contact_id"`
		VariableID    string      `json:"variable_id"`
		VariableName  string      `json:"variable_name"`
		VariableValue interface{} `json:"variable_value"`
	}
	body := bodyFormat{
		ContactID:     contactID,
		VariableID:    variableID,
		VariableName:  variableName,
		VariableValue: variableValue,
	}

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, body, &respData, true)
	return err
}

func (service *BotsLiveChatService) SetTagsToContact(ctx context.Context, contactID string, tags []string) error {
	path := "/live-chat/contacts/setTag"

	type bodyFormat struct {
		ContactID string   `json:"contact_id"`
		Tags      []string `json:"tags"`
	}
	body := bodyFormat{
		ContactID: contactID,
		Tags:      tags,
	}

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, body, &respData, true)
	return err
}

func (service *BotsLiveChatService) DeleteTagFromContact(ctx context.Context, contactID string, tag string) error {
	path := "/live-chat/contacts/deleteTag"

	type bodyFormat struct {
		ContactID string `json:"contact_id"`
		Tag       string `json:"tag"`
	}
	body := bodyFormat{
		ContactID: contactID,
		Tag:       tag,
	}

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, body, &respData, true)
	return err
}

func (service *BotsLiveChatService) DisableContact(ctx context.Context, contactID string) error {
	path := "/live-chat/contacts/disable"

	type bodyFormat struct {
		ContactID string `json:"contact_id"`
	}
	body := bodyFormat{
		ContactID: contactID,
	}

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, body, &respData, true)
	return err
}

func (service *BotsLiveChatService) EnableContact(ctx context.Context, contactID string) error {
	path := "/live-chat/contacts/enable"

	type bodyFormat struct {
		ContactID string `json:"contact_id"`
	}
	body := bodyFormat{
		ContactID: contactID,
	}

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, body, &respData, true)
	return err
}

func (service *BotsLiveChatService) DeleteContact(ctx context.Context, contactID string) error {
	path := "/live-chat/contacts/delete"

	type bodyFormat struct {
		ContactID string `json:"contact_id"`
	}
	body := bodyFormat{
		ContactID: contactID,
	}

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, body, &respData, true)
	return err
}

func (service *BotsLiveChatService) GetPauseAutomation(ctx context.Context, contactID string) (int, error) {
	path := fmt.Sprintf("/live-chat/contacts/getPauseAutomation?contact_id=%s", contactID)

	var respData struct {
		Success bool `json:"success"`
		Data    struct {
			Minutes int `json:"minutes"`
		} `json:"data"`
	}
	_, err := service.client.newRequest(ctx, http.MethodGet, path, nil, &respData, true)
	return respData.Data.Minutes, err
}

func (service *BotsLiveChatService) SetPauseAutomation(ctx context.Context, contactID string, minutes int) error {
	path := "/live-chat/contacts/setPauseAutomation"
	type bodyFormat struct {
		ContactID string `json:"contact_id"`
		Minutes   int    `json:"minutes"`
	}
	body := bodyFormat{
		ContactID: contactID,
		Minutes:   minutes,
	}

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, body, &respData, true)
	return err
}

func (service *BotsLiveChatService) DeletePauseAutomation(ctx context.Context, contactID string) error {
	path := "/live-chat/contacts/deletePauseAutomation"
	type bodyFormat struct {
		ContactID string `json:"contact_id"`
	}
	body := bodyFormat{
		ContactID: contactID,
	}

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, body, &respData, true)
	return err
}

func (service *BotsLiveChatService) GetBotVariables(ctx context.Context, botID string) ([]*BotVariable, error) {
	path := fmt.Sprintf("/live-chat/variables?bot_id=%s", botID)

	var respData struct {
		Success bool           `json:"success"`
		Data    []*BotVariable `json:"data"`
	}
	_, err := service.client.newRequest(ctx, http.MethodGet, path, nil, &respData, true)
	return respData.Data, err
}

type BotLiveChatFlow struct {
	ID     string `json:"id"`
	BotID  string `json:"bot_id"`
	Name   string `json:"name"`
	Status struct {
		Active   int `json:"ACTIVE"`
		Inactive int `json:"INACTIVE"`
		Draft    int `json:"DRAFT"`
	} `json:"status"`
	Triggers []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type int    `json:"type"`
	} `json:"triggers"`
	CreatedAt time.Time `json:"created_at"`
}

func (service *BotsLiveChatService) GetFlows(ctx context.Context, botID string) ([]*BotIgFlow, error) {
	path := fmt.Sprintf("/live-chat/flows?bot_id=%s", botID)

	var respData struct {
		Success bool         `json:"success"`
		Data    []*BotIgFlow `json:"data"`
	}
	_, err := service.client.newRequest(ctx, http.MethodGet, path, nil, &respData, true)
	return respData.Data, err
}

func (service *BotsLiveChatService) RunFlow(ctx context.Context, contactID, flowID string, externalData map[string]interface{}) error {
	path := "/live-chat/flows/run"

	type bodyFormat struct {
		ContactID    string                 `json:"contact_id"`
		FlowID       string                 `json:"flow_id"`
		ExternalData map[string]interface{} `json:"external_data,omitempty"`
	}
	body := bodyFormat{
		ContactID:    contactID,
		FlowID:       flowID,
		ExternalData: externalData,
	}

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, body, &respData, true)
	return err
}

func (service *BotsLiveChatService) RunFlowByTrigger(ctx context.Context, contactID, triggerKeyword string, externalData map[string]interface{}) error {
	path := "/live-chat/flows/runByTrigger"

	type bodyFormat struct {
		ContactID      string                 `json:"contact_id"`
		TriggerKeyword string                 `json:"trigger_keyword"`
		ExternalData   map[string]interface{} `json:"external_data,omitempty"`
	}
	body := bodyFormat{
		ContactID:      contactID,
		TriggerKeyword: triggerKeyword,
		ExternalData:   externalData,
	}

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, body, &respData, true)
	return err
}

func (service *BotsLiveChatService) GetBotTriggers(ctx context.Context, botID string) ([]*BotTrigger, error) {
	path := fmt.Sprintf("/live-chat/triggers?bot_id=%s", botID)

	var respData struct {
		Success bool          `json:"success"`
		Data    []*BotTrigger `json:"data"`
	}
	_, err := service.client.newRequest(ctx, http.MethodGet, path, nil, &respData, true)
	return respData.Data, err
}

type LiveChatBotMessage struct {
	ID         string                 `json:"id"`
	ContactID  string                 `json:"contact_id"`
	BotID      string                 `json:"bot_id"`
	CampaignID string                 `json:"campaign_id"`
	Data       map[string]interface{} `json:"data"`
	Direction  int                    `json:"direction"`
	Status     int                    `json:"status"`
	CreatedAt  time.Time              `json:"created_at"`
	Type       string                 `json:"type"`
}

type LiveChatBotChat struct {
	Contact          *IgBotContact `json:"contact"`
	InboxLastMessage *IgBotMessage `json:"inbox_last_message"`
	InboxUnread      int           `json:"inbox_unread"`
}

func (service *BotsLiveChatService) GetBotChats(ctx context.Context, botID string) ([]*IgBotChat, error) {
	path := fmt.Sprintf("/live-chat/chats?bot_id=%s", botID)

	var respData struct {
		Success bool         `json:"success"`
		Data    []*IgBotChat `json:"data"`
	}
	_, err := service.client.newRequest(ctx, http.MethodGet, path, nil, &respData, true)
	return respData.Data, err
}

func (service *BotsLiveChatService) GetContactMessages(ctx context.Context, contactID string, size *int, skip *int, order *string) ([]*IgBotMessage, error) {
	defaultSize := 20
	defaultSkip := 0
	defaultOrder := "desc"

	// 检查传入的参数，如果为nil则使用默认值
	if size == nil {
		size = &defaultSize
	}
	if skip == nil {
		skip = &defaultSkip
	}
	if order == nil {
		order = &defaultOrder
	}
	path := fmt.Sprintf("/live-chat/chats/messages?contact_id=%s&size=%d&skip=%d&order=%s", contactID, *size, *skip, *order)

	var respData struct {
		Success bool            `json:"success"`
		Data    []*IgBotMessage `json:"data"`
	}
	_, err := service.client.newRequest(ctx, http.MethodGet, path, nil, &respData, true)
	return respData.Data, err
}

type LiveChatBotSendCampaignParams struct {
	Title    string                 `json:"title"`
	BotID    string                 `json:"bot_id"`
	SendAt   time.Time              `json:"send_at"`
	Messages []IgBotCampaignMessage `json:"messages"`
}

type LiveChatBotCampaignMessage struct {
	Type    string `json:"type"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}

func (service *BotsLiveChatService) SendCampaign(ctx context.Context, params IgBotSendCampaignParams) error {
	path := "/live-chat/campaigns/send"

	var respData struct {
		Success bool `json:"success"`
	}
	_, err := service.client.newRequest(ctx, http.MethodPost, path, params, &respData, true)
	return err
}
