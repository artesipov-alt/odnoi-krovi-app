package events

import (
	"time"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
)

type ContactData struct {
	Name             string `json:"name"`
	ProviderMaxID    string `json:"providerMaxId"`
	ProviderTelegram string `json:"providerTelegram"`
	Phone            string `json:"phone"`
}

type UserContact struct {
	NotifyProvider authmodel.ProviderName `json:"notifyProvider"`
	SendTo         string                 `json:"sendTo"`
	UserData       ContactData            `json:"userData"`
	CreatedAt      time.Time              `json:"createdAt"`
	Recipient      usermodel.User         `json:"recipient"`
}

func GenerateContact(u *usermodel.User) UserContact {
	event := UserContact{
		UserData: ContactData{
			Name:  u.FullName,
			Phone: u.Phone,
		},
		CreatedAt: time.Now(),
	}
	if len(u.Identities) == 0 {
		return event
	}
	for _, identity := range u.Identities {
		if identity.ProviderName == authmodel.ProviderMax {
			event.UserData.ProviderMaxID = identity.ProviderUserID
		}
		if identity.ProviderName == authmodel.ProviderTelegram {
			event.UserData.ProviderTelegram = identity.ProviderUserID
		}
	}

	return event
}
