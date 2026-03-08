package model

import (
	"time"
)

type ProviderName string

const (
	ProviderTelegram ProviderName = "telegram_bot"
	ProviderMax      ProviderName = "max_bot"
	ProviderService  ProviderName = "service"
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

func NewUserMetadata(metadata map[string]any) *Metadata {
	source, medium, campaign, content, term := extractUTMFromMetadata(metadata)
	return &Metadata{
		UTMData: &UTM{
			Source:   source,
			Medium:   medium,
			Campaign: campaign,
			Content:  content,
			Term:     term,
		},
	}
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
