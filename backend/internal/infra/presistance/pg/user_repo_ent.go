package pg

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/donorpreference"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/schema"
	entuser "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/useridentity"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/utmhistory"
	// расширение для апсерта
)

// EntUserRepository implements UserRepository using ENT
type EntUserRepository struct {
	db *ent.Client
}

// client returns the ent.Client from the context if a transaction is active,
// otherwise returns the default client
func (r *EntUserRepository) client(ctx context.Context) *ent.Client {
	if tx := ent.TxFromContext(ctx); tx != nil {
		return tx.Client()
	}
	return r.db
}

// NewEntUserRepository creates a new ENT user repository
func NewEntUserRepository(client *ent.Client) *EntUserRepository {
	return &EntUserRepository{
		db: client,
	}
}

// CreateUser creates a new user in the database along with identity
func (r *EntUserRepository) CreateUser(ctx context.Context, inputuser *usermodel.User) (*usermodel.User, error) {
	if inputuser == nil {
		return nil, errors.New("user cannot be nil")
	}

	c := r.client(ctx)

	builder := c.User.
		Create().
		SetFullName(inputuser.FullName).
		SetRole(entuser.Role(inputuser.Role))

	if inputuser.Phone != "" {
		builder.SetPhone(inputuser.Phone)
	}
	if inputuser.Email != "" {
		builder.SetEmail(inputuser.Email)
	}
	if inputuser.OrganizationName != "" {
		builder.SetOrganizationName(inputuser.OrganizationName)
	}
	if inputuser.LocationID != nil && *inputuser.LocationID != "" {
		builder.SetLocationID(*inputuser.LocationID)
	}
	if len(inputuser.PhotoURLs) > 0 {
		builder.SetPhotoUrls(inputuser.PhotoURLs)
	}
	if len(inputuser.OnBoarding) > 0 {
		builder.SetOnBoarding(inputuser.OnBoarding)
	}
	if inputuser.ConsentPd {
		builder.SetConsentPd(inputuser.ConsentPd)
	}
	if inputuser.AllowGeo {
		builder.SetAllowGeo(inputuser.AllowGeo)
	}
	if inputuser.OriginSource != "" {
		builder.SetOriginSource(inputuser.OriginSource)
	}

	newUser, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Re-fetch the created user without relations
	return r.GetByID(ctx, newUser.ID, user.UserPreloadOptions{
		WithPets:            false,
		WithDonorPreference: false,
	})
}

// UpsertUserIdentity creates or updates a user identity
func (r *EntUserRepository) UpsertUserIdentity(ctx context.Context, userID string, input *authmodel.Identity, metadata *authmodel.Metadata) error {
	if input == nil {
		return errors.New("user identity cannot be nil")
	}
	if userID == "" {
		return errors.New("invalid user ID")
	}

	c := r.client(ctx)

	builder := c.UserIdentity.Create().
		SetUserID(userID).
		SetProvider(useridentity.Provider(input.ProviderName)).
		SetProviderUserID(input.ProviderUserID)

	if input.PartnerID != "" {
		builder.SetPartnerID(input.PartnerID)
	}

	if metadata != nil && metadata.UTMData != nil {
		metadataMap := make(map[string]any)
		metadataMap["utm_source"] = metadata.UTMData.Source
		metadataMap["utm_medium"] = metadata.UTMData.Medium
		metadataMap["utm_campaign"] = metadata.UTMData.Campaign
		metadataMap["utm_content"] = metadata.UTMData.Content
		metadataMap["utm_term"] = metadata.UTMData.Term
		builder.SetMetadata(metadataMap)
	}

	err := builder.
		OnConflict(sql.ConflictColumns("provider_user_id", "provider")).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to upsert user identity: %w", err)
	}

	return nil
}

// CreateDonorPreference creates donor preference for a user
func (r *EntUserRepository) CreateDonorPreference(ctx context.Context, userID string, inputprefs *usermodel.DonorPreference) error {
	if inputprefs == nil {
		return errors.New("donor preference cannot be nil")
	}
	if userID == "" {
		return errors.New("invalid user ID")
	}

	c := r.client(ctx)

	prefBuilder := c.DonorPreference.Create().
		SetUserID(userID)

	prefBuilder.SetPreferredLocationIds(inputprefs.PreferredLocationIDs)
	prefBuilder.SetRecoveryPeriodMonths(inputprefs.RecoveryPeriodMonths)
	if inputprefs.CompensationType != "" {
		prefBuilder.SetCompensationType(donorpreference.CompensationType(inputprefs.CompensationType))
	}
	prefBuilder.SetTaxiCompensation(inputprefs.TaxiCompensation)
	prefBuilder.SetNotificationFrequency(donorpreference.NotificationFrequency(inputprefs.NotificationFrequency))

	_, err := prefBuilder.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create donor preference: %w", err)
	}

	return nil
}

func (r *EntUserRepository) GetByID(ctx context.Context, id string, opts user.UserPreloadOptions) (*usermodel.User, error) {
	quser := r.client(ctx).User.Query().Where(entuser.ID(id))

	if opts.WithPets {
		quser = quser.WithPets()
	}
	if opts.WithDonorPreference {
		quser = quser.WithDonorPreference()
	}

	user, err := quser.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.Internal(err, "failed to get user by ID")
	}

	return EntToModel(user), nil
}

// DEPRECATED
func (r *EntUserRepository) GetByTelegram(ctx context.Context, telegramID int64, opts user.UserPreloadOptions) (*usermodel.User, error) {
	quser := r.client(ctx).User.Query().Where(entuser.TelegramID(telegramID))

	if opts.WithPets {
		quser = quser.WithPets()
	}
	if opts.WithDonorPreference {
		quser = quser.WithDonorPreference()
	}

	user, err := quser.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.Internal(err, "failed to get user by Telegram ID")
	}

	return EntToModel(user), nil
}

func (r *EntUserRepository) GetByProvider(ctx context.Context, providerID string, providerName authmodel.ProviderName) (*authmodel.Identity, error) {
	identity, err := r.client(ctx).UserIdentity.Query().
		Where(useridentity.ProviderUserID(providerID),
			useridentity.ProviderEQ(useridentity.Provider(providerName))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.Internal(err, "failed to get user by provider ID")
	}

	return EntIdentityToModel(identity), nil
}

// ExistsByID checks if a user with the given ID exists
func (r *EntUserRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, errors.New("invalid user ID")
	}

	exists, err := r.client(ctx).User.Query().
		Where(entuser.ID(id)).
		Exist(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence by id %s: %w", id, err)
	}

	return exists, nil
}

// UpdateUserFields updates user fields (simple update without transaction handling)
func (r *EntUserRepository) UpdateUserFields(ctx context.Context, id string, input *usermodel.User) error {
	if input == nil {
		return errors.New("user cannot be nil")
	}

	if id == "" {
		return errors.New("invalid user ID")
	}

	c := r.client(ctx)

	builder := c.User.UpdateOneID(id)

	if input.FullName != "" {
		builder.SetFullName(input.FullName)
	}
	if input.Email != "" {
		builder.SetEmail(input.Email)
	}
	if input.OrganizationName != "" {
		builder.SetOrganizationName(input.OrganizationName)
	}
	if input.LocationID != nil {
		builder.SetLocationID(*input.LocationID)
	}
	if len(input.PhotoURLs) > 0 {
		builder.SetPhotoUrls(input.PhotoURLs)
	}
	if len(input.OnBoarding) > 0 {
		builder.SetOnBoarding(input.OnBoarding)
	}
	builder.SetConsentPd(input.ConsentPd)
	builder.SetAllowGeo(input.AllowGeo)
	if input.Phone != "" {
		builder.SetPhone(input.Phone)
	}

	_, err := builder.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// TransferUserIdentity transfers all identities from one user to another
func (r *EntUserRepository) TransferUserIdentity(ctx context.Context, fromUserID, toUserID string) error {
	if fromUserID == "" || toUserID == "" {
		return errors.New("invalid user IDs")
	}

	c := r.client(ctx)

	_, err := c.UserIdentity.Update().
		Where(useridentity.UserID(fromUserID)).
		SetUserID(toUserID).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to transfer user identity: %w", err)
	}

	return nil
}

// DeleteDonorPreferenceByUserID deletes donor preference for a user
func (r *EntUserRepository) DeleteDonorPreferenceByUserID(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("invalid user ID")
	}

	c := r.client(ctx)
	ctxWithSkip := schema.SkipSoftDelete(ctx)

	_, err := c.DonorPreference.Delete().
		Where(donorpreference.HasUserWith(entuser.ID(userID))).
		Exec(ctxWithSkip)
	if err != nil {
		return fmt.Errorf("failed to delete donor preference: %w", err)
	}

	return nil
}

// DeleteUTMHistoryByUserID deletes UTM history for a user
func (r *EntUserRepository) DeleteUTMHistoryByUserID(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("invalid user ID")
	}

	c := r.client(ctx)
	ctxWithSkip := schema.SkipSoftDelete(ctx)

	_, err := c.UtmHistory.Delete().
		Where(utmhistory.HasUserWith(entuser.ID(userID))).
		Exec(ctxWithSkip)
	if err != nil {
		return fmt.Errorf("failed to delete UTM history: %w", err)
	}

	return nil
}

// TransferUTMHistory transfers UTM history from one user to another
func (r *EntUserRepository) TransferUTMHistory(ctx context.Context, fromUserID, toUserID string) error {
	if fromUserID == "" || toUserID == "" {
		return errors.New("invalid user IDs")
	}

	c := r.client(ctx)

	_, err := c.UtmHistory.Update().
		Where(utmhistory.HasUserWith(entuser.ID(fromUserID))).
		SetUserID(toUserID).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to transfer UTM history: %w", err)
	}

	return nil
}

// DeleteUserHard permanently deletes a user (bypasses soft delete)
func (r *EntUserRepository) DeleteUserHard(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	c := r.client(ctx)
	ctxWithSkip := schema.SkipSoftDelete(ctx)

	// First delete donor preferences (foreign key constraint)
	if err := r.DeleteDonorPreferenceByUserID(ctx, id); err != nil {
		return err
	}

	// Then delete the user
	err := c.User.DeleteOneID(id).Exec(ctxWithSkip)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// UpsertDonorPreference creates or updates donor preference for a user
func (r *EntUserRepository) UpsertDonorPreference(ctx context.Context, userID string, prefs *usermodel.DonorPreference) error {
	if prefs == nil {
		return errors.New("donor preference cannot be nil")
	}
	if userID == "" {
		return errors.New("invalid user ID")
	}

	c := r.client(ctx)

	prefBuilder := c.DonorPreference.Create().
		SetUserID(userID)

	prefBuilder.SetPreferredLocationIds(prefs.PreferredLocationIDs)
	prefBuilder.SetRecoveryPeriodMonths(prefs.RecoveryPeriodMonths)
	if prefs.CompensationType != "" {
		prefBuilder.SetCompensationType(donorpreference.CompensationType(prefs.CompensationType))
	}
	prefBuilder.SetTaxiCompensation(prefs.TaxiCompensation)
	prefBuilder.SetNotificationFrequency(donorpreference.NotificationFrequency(prefs.NotificationFrequency))

	err := prefBuilder.
		OnConflict(sql.ConflictColumns("user_donor_preference")).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to upsert donor preference: %w", err)
	}

	return nil
}

// Delete deletes a user by their ID (soft delete via SoftDeleteMixin)
func (r *EntUserRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	// Soft delete via SoftDeleteMixin hook
	err := r.client(ctx).User.DeleteOneID(id).Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user with id %s not found", id)
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ExistsByProvider checks if a user with the given Provider ID exists
func (r *EntUserRepository) ExistsByProvider(ctx context.Context, providerID string, providerName authmodel.ProviderName) (bool, error) {
	if providerID == "" {
		return false, errors.New("invalid provider ID")
	}

	exists, err := r.client(ctx).UserIdentity.Query().
		Where(useridentity.ProviderUserID(providerID),
			useridentity.ProviderEQ(useridentity.Provider(providerName))).
		Exist(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return false, apperrors.ErrUserNotFound
		}
		return false, apperrors.Internal(err, "failed to get user by provider ID")
	}

	return exists, nil
}

// GetByPhone retrieves a user by phone number
func (r *EntUserRepository) GetByPhone(ctx context.Context, phone string) (string, error) {
	if phone == "" {
		return "", errors.New("invalid phone number")
	}

	user, err := r.client(ctx).User.Query().
		Where(entuser.Phone(phone)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return "", apperrors.ErrUserNotFound
		}
		return "", apperrors.Internal(err, "failed to get user by phone")
	}

	return user.ID, nil
}

// ResetUser resets user's email and phone number by ID
func (r *EntUserRepository) ResetUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	err := r.client(ctx).User.UpdateOneID(id).
		SetEmail("").
		SetPhone("").
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user with id %s not found", id)
		}
		return fmt.Errorf("failed to reset user data: %w", err)
	}

	return nil
}

// RestoreUser restores a soft-deleted user by setting deleted_at to NULL
func (r *EntUserRepository) RestoreUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	// Use SkipSoftDelete context to find the deleted record
	ctxWithSkip := schema.SkipSoftDelete(ctx)

	err := r.client(ctx).User.UpdateOneID(id).
		ClearDeletedAt().
		Exec(ctxWithSkip)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user with id %s not found", id)
		}
		return fmt.Errorf("failed to restore user: %w", err)
	}

	return nil
}

// GetDeletedUsers retrieves all soft-deleted users
func (r *EntUserRepository) GetDeletedUsers(ctx context.Context) ([]*usermodel.User, error) {
	// Use SkipSoftDelete context to see deleted records
	ctxWithSkip := schema.SkipSoftDelete(ctx)

	users, err := r.client(ctx).User.Query().
		Where(entuser.DeletedAtNotNil()).
		All(ctxWithSkip)

	if err != nil {
		return nil, fmt.Errorf("failed to get deleted users: %w", err)
	}

	result := make([]*usermodel.User, len(users))
	for i, u := range users {
		result[i] = EntToModel(u)
	}

	return result, nil
}

// AddPhotoURLs adds new photo paths to the user's PhotoUrls array
func (r *EntUserRepository) AddPhotoURLs(ctx context.Context, id string, paths []string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	// Replace photo URLs with new paths
	newPhotoUrls := paths

	// Update user
	err := r.client(ctx).User.UpdateOneID(id).
		SetPhotoUrls(newPhotoUrls).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update user photo URLs: %w", err)
	}

	return nil
}

// UpsertUTM upserts UTM data for an existing user
func (r *EntUserRepository) UpsertUTM(ctx context.Context, userID string, metadata *authmodel.Metadata) error {
	if metadata == nil || userID == "" {
		return errors.New("invalid user ID")
	}

	builder := r.client(ctx).UtmHistory.Create().
		SetUserID(userID)

	if metadata.UTMData.Source != "" {
		builder.SetUtmSource(metadata.UTMData.Source)
	}
	if metadata.UTMData.Medium != "" {
		builder.SetUtmMedium(metadata.UTMData.Medium)
	}
	if metadata.UTMData.Campaign != "" {
		builder.SetUtmCampaign(metadata.UTMData.Campaign)
	}
	if metadata.UTMData.Content != "" {
		builder.SetUtmContent(metadata.UTMData.Content)
	}
	if metadata.UTMData.Term != "" {
		builder.SetUtmTerm(metadata.UTMData.Term)
	}

	err := builder.
		OnConflict(sql.ConflictColumns("user_id", "utm_campaign", "utm_source")).
		UpdateNewValues().
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to upsert UTM data: %w", err)
	}

	return nil
}

// EntToModel converts ent.User to domain model User
func EntToModel(e *ent.User) *usermodel.User {
	if e == nil {
		return nil
	}

	user := &usermodel.User{
		ID:               e.ID,
		TelegramID:       e.TelegramID,
		FullName:         e.FullName,
		Phone:            e.Phone,
		Email:            e.Email,
		PhotoURLs:        e.PhotoUrls,
		OrganizationName: e.OrganizationName,
		ConsentPd:        e.ConsentPd,
		OnBoarding:       e.OnBoarding,
		AllowGeo:         e.AllowGeo,
		Role:             usermodel.UserRole(e.Role),
		OriginSource:     e.OriginSource,
		Pets:             nil, // Pets are loaded separately via WithPets
		CreatedAt:        &e.CreatedAt,
		UpdatedAt:        &e.UpdatedAt,
		DeletedAt:        e.DeletedAt,
	}

	if e.LocationID != "" {
		user.LocationID = &e.LocationID
	}

	var pets []*petmodel.Pet
	if len(e.Edges.Pets) > 0 {
		for _, p := range e.Edges.Pets {
			pets = append(pets, petToDomain(p))
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
			CompensationType:      usermodel.CompensationType(dp.CompensationType.String()),
			TaxiCompensation:      dp.TaxiCompensation,
			NotificationFrequency: usermodel.NotificationFrequency(dp.NotificationFrequency),
			CreatedAt:             &dp.CreatedAt,
			UpdatedAt:             &dp.UpdatedAt,
			DeletedAt:             dp.DeletedAt,
		}
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
