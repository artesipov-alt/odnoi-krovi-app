package domainmapper

import (
	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

// EntToModel converts ent.User to domain model User
func EntToModel(e *ent.User) *usermodel.User {
	if e == nil {
		return nil
	}

	user := &usermodel.User{
		ID:                  e.ID,
		FullName:            e.FullName,
		Phone:               e.Phone,
		Verified:            e.Verified,
		Email:               e.Email,
		PhotoURLs:           e.PhotoUrls,
		OrganizationName:    e.OrganizationName,
		ConsentPd:           e.ConsentPd,
		OnBoarding:          e.OnBoarding,
		AllowGeo:            e.AllowGeo,
		Role:                usermodel.UserRole(e.Role),
		OriginSource:        e.OriginSource,
		PrioritySearchCount: e.PrioritySearchCount,
		CreatedAt:           &e.CreatedAt,
		UpdatedAt:           &e.UpdatedAt,
		DeletedAt:           e.DeletedAt,
		LastSeenAt:          e.LastSeenAt,
	}

	if e.LocationID != "" {
		user.LocationID = &e.LocationID
	}

	var pets []*petmodel.Pet
	if len(e.Edges.Pets) > 0 {
		for _, p := range e.Edges.Pets {
			pets = append(pets, PetToDomain(p))
		}
	}
	user.Pets = pets

	// Map DonorPreference with UserID from user
	if e.Edges.DonorPreference != nil {
		dp := e.Edges.DonorPreference

		user.DonorPreference = &usermodel.DonorPreference{
			ID:                    dp.ID,
			UserID:                e.ID,
			PreferredLocationIDs:  dp.PreferredLocationIds,
			RecoveryPeriodMonths:  dp.RecoveryPeriodMonths,
			CompensationType:      common.CompensationType(dp.CompensationType.String()),
			TaxiCompensation:      dp.TaxiCompensation,
			NotificationFrequency: usermodel.NotificationFrequency(dp.NotificationFrequency),
			OpenForContact:        dp.OpenForContact,
			CreatedAt:             &dp.CreatedAt,
			UpdatedAt:             &dp.UpdatedAt,
			DeletedAt:             dp.DeletedAt,
		}
	}
	if len(e.Edges.Identities) != 0 {
		var identities []*authmodel.Identity
		for _, i := range e.Edges.Identities {
			identities = append(identities, EntIdentityToModel(i))
		}
		user.Identities = identities
	}

	return user
}

func EntIdentityToModel(identity *ent.UserIdentity) *authmodel.Identity {
	return &authmodel.Identity{
		ID:             identity.ID,
		UserID:         identity.UserID,
		ProviderName:   authmodel.ProviderName(identity.Provider),
		ProviderUserID: identity.ProviderUserID,
		Metadata:       &identity.Metadata,
		CreatedAt:      identity.CreatedAt,
		UpdatedAt:      identity.UpdatedAt,
		DeletedAt:      identity.DeletedAt,
	}
}
