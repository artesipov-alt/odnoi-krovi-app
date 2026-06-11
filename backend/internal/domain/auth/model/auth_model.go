package model

import (
	"net/url"
	"time"
)

type ProviderName string

const (
	ProviderTelegram ProviderName = "telegram_bot"
	ProviderMax      ProviderName = "max_bot"
	ProviderService  ProviderName = "service"
)

const (
	TelegramBotURL = "https://t.me/app1krovi_bot?start="
	MaxBotURL      = "https://max.ru/id3200014662_bot?start="
)

// Identity представляет доменную модель пользователя
type Identity struct {
	ID             string
	UserID         string
	ProviderName   ProviderName
	ProviderUserID string
	ServiceKey     string
	AppInitData    string
	AccessToken    string
	PartnerID      string
	Metadata       *map[string]any
	ExpiresAt      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// NewUser creates a new User aggregate with validation
func NewIdentity(providerName ProviderName, appInitData string, metadata *map[string]any) (*Identity, error) {

	idn := &Identity{
		ProviderName: providerName,
		AppInitData:  appInitData,
		Metadata:     metadata,
	}

	return idn, nil
}

// NewUser creates a new User aggregate with validation
func NewServiceIdentity(providerID, providerName, apiKey string, metadata *map[string]any) (*Identity, error) {
	idn := &Identity{
		ProviderName:   ProviderName(providerName),
		ProviderUserID: providerID,
		ServiceKey:     apiKey,
		Metadata:       metadata,
	}

	return idn, nil
}

// NewUser creates a new User aggregate with validation
func NewMiniAppIdentity(providerName ProviderName, appInitData string, metadata *map[string]any) (*Identity, error) {
	idn := &Identity{
		ProviderName: providerName,
		AppInitData:  appInitData,
		Metadata:     metadata,
	}

	return idn, nil
}

func (i *Identity) SetPartnerID(partnerID string) {
	i.PartnerID = partnerID
}

func (i *Identity) SetSystemUserID(userID string) {
	i.UserID = userID
}

func (i *Identity) SetJWTData(accessToken string, expiresAt time.Time) {
	i.AccessToken = accessToken
	i.ExpiresAt = expiresAt
}

func (i *Identity) SetProviderID(providerID string) {
	i.ProviderUserID = providerID
}

func (i *Identity) GenerateRefURL() string {
	utm := url.Values{}
	utm.Set("utm_source", "ref")
	utm.Set("utm_medium", string(i.ProviderName))
	utm.Set("utm_campaign", "ref_"+i.UserID)
	startParam := utm.Encode()
	switch i.ProviderName {
	case ProviderTelegram:
		return TelegramBotURL + startParam
	case ProviderMax:
		return MaxBotURL + startParam
	default:
		return ""
	}
}

// ====================================================================================================
//
//                                             UTM-Метки
//
// ====================================================================================================

// Metadata represents user metadata with UTM and other fields
type UTM struct {
	Source   string
	Medium   string
	Campaign string
	Content  string
	Term     string
}

type Metadata struct {
	UTMData *UTM
}

func NewUserMetadata(metadata map[string]any) (*Metadata, error) {
	source, medium, campaign, content, term := extractUTMFromMetadata(metadata)
	return &Metadata{
		UTMData: &UTM{
			Source:   source,
			Medium:   medium,
			Campaign: campaign,
			Content:  content,
			Term:     term,
		},
	}, nil
}

func (m *Metadata) GetCampaign() string {
	if m == nil || m.UTMData == nil {
		return ""
	}
	return m.UTMData.Campaign
}

func extractUTMFromMetadata(metadata map[string]any) (string, string, string, string, string) {
	source := ""
	if s, ok := metadata["utm_source"].(string); ok {
		source = s
	}
	medium := ""
	if s, ok := metadata["utm_medium"].(string); ok {
		medium = s
	}
	campaign := ""
	if s, ok := metadata["utm_campaign"].(string); ok {
		campaign = s
	}
	content := ""
	if s, ok := metadata["utm_content"].(string); ok {
		content = s
	}
	term := ""
	if s, ok := metadata["utm_term"].(string); ok {
		term = s
	}
	return source, medium, campaign, content, term
}
