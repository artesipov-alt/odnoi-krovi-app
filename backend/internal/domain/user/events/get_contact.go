package events

import (
	"time"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
)

type UserContact struct {
	NotifyProvider authmodel.ProviderName
	SendTo         string
	UserData       ContactData
	CreatedAt      time.Time
	Recipient      usermodel.User
}

type ContactData struct {
	Name             string
	ProviderMaxID    string
	ProviderTelegram string
	Phone            string
}

func (e UserContact) EventName() string     { return "ApplyDonor" }
func (e UserContact) OccurredAt() time.Time { return e.CreatedAt }

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
