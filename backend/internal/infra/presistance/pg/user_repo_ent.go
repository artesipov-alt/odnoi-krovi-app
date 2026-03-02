package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/donorpreference"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/schema"
	entuser "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/user"
)

// EntUserRepository implements UserRepository using ENT
type EntUserRepository struct {
	client *ent.Client
}

// NewEntUserRepository creates a new ENT user repository
func NewEntUserRepository(client *ent.Client) *EntUserRepository {
	return &EntUserRepository{
		client: client,
	}
}

// Create creates a new user in the database
func (r *EntUserRepository) Create(ctx context.Context, inputuser *usermodel.User) (*usermodel.User, error) {
	if inputuser == nil {
		return nil, errors.New("user cannot be nil")
	}

	// Start transaction
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	builder := tx.User.
		Create().
		SetTelegramID(inputuser.TelegramID).
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

	newUser, err := builder.Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create DonorPreference if provided
	if inputuser.DonorPreference != nil {
		prefBuilder := tx.DonorPreference.Create().
			SetUser(newUser)

		if len(inputuser.DonorPreference.PreferredLocationIDs) > 0 {
			prefBuilder.SetPreferredLocationIds(inputuser.DonorPreference.PreferredLocationIDs)
		}
		if inputuser.DonorPreference.RecoveryPeriodMonths > 0 {
			prefBuilder.SetRecoveryPeriodMonths(inputuser.DonorPreference.RecoveryPeriodMonths)
		}
		if inputuser.DonorPreference.CompensationType != "" {
			prefBuilder.SetCompensationType(donorpreference.CompensationType(inputuser.DonorPreference.CompensationType))
		}
		prefBuilder.SetTaxiCompensation(inputuser.DonorPreference.TaxiCompensation)
		prefBuilder.SetNotificationFrequency(donorpreference.NotificationFrequency(inputuser.DonorPreference.NotificationFrequency))

		_, err = prefBuilder.Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create donor preference: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Re-fetch the created user with relations
	return r.GetByID(ctx, newUser.ID, user.UserPreloadOptions{
		WithPets:            false, // Adjust as needed
		WithDonorPreference: inputuser.DonorPreference != nil,
	})
}

func (r *EntUserRepository) GetByID(ctx context.Context, id string, opts user.UserPreloadOptions) (*usermodel.User, error) {
	quser := r.client.User.Query().Where(entuser.ID(id))

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

func (r *EntUserRepository) GetByTelegram(ctx context.Context, telegramID int64, opts user.UserPreloadOptions) (*usermodel.User, error) {
	quser := r.client.User.Query().Where(entuser.TelegramID(telegramID))

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

// ExistsByID checks if a user with the given ID exists
func (r *EntUserRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, errors.New("invalid user ID")
	}

	exists, err := r.client.User.Query().
		Where(entuser.ID(id)).
		Exist(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence by id %s: %w", id, err)
	}

	return exists, nil
}

// Update updates an existing user in the database
func (r *EntUserRepository) Update(ctx context.Context, id string, input *usermodel.User) error {
	if input == nil {
		return errors.New("user cannot be nil")
	}

	if id == "" {
		return errors.New("invalid user ID")
	}

	builder := r.client.User.UpdateOneID(id)

	if input.FullName != "" {
		builder.SetFullName(input.FullName)
	}
	if input.Phone != "" {
		builder.SetPhone(input.Phone)
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

	_, err := builder.Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user by their ID (soft delete via SoftDeleteMixin)
func (r *EntUserRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	// Soft delete via SoftDeleteMixin hook
	err := r.client.User.DeleteOneID(id).Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user with id %s not found", id)
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ExistsByTelegramID checks if a user with the given Telegram ID exists
func (r *EntUserRepository) ExistsByTelegramID(ctx context.Context, telegramID int64) (bool, error) {
	if telegramID <= 0 {
		return false, errors.New("invalid telegram ID")
	}

	exists, err := r.client.User.Query().
		Where(entuser.TelegramID(telegramID)).
		Exist(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check user existence by telegram id %d: %w", telegramID, err)
	}

	return exists, nil
}

// ResetUser resets user's email and phone number by ID
func (r *EntUserRepository) ResetUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	err := r.client.User.UpdateOneID(id).
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

	err := r.client.User.UpdateOneID(id).
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

	users, err := r.client.User.Query().
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
	err := r.client.User.UpdateOneID(id).
		SetPhotoUrls(newPhotoUrls).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update user photo URLs: %w", err)
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
		Role:             string(e.Role),
		Pets:             nil, // Pets are loaded separately via WithPets
		DonorPreference:  donorPrefToDomain(e.Edges.DonorPreference),
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

	return user
}

func donorPrefToDomain(e *ent.DonorPreference) *usermodel.DonorPreference {
	if e == nil {
		return nil
	}

	return &usermodel.DonorPreference{
		ID:                    e.ID,
		UserID:                e.Edges.User.ID, // Assuming User is loaded
		PreferredLocationIDs:  e.PreferredLocationIds,
		RecoveryPeriodMonths:  e.RecoveryPeriodMonths,
		CompensationType:      usermodel.CompensationType(e.CompensationType.String()),
		TaxiCompensation:      e.TaxiCompensation,
		NotificationFrequency: usermodel.NotificationFrequency(e.NotificationFrequency.String()),
		CreatedAt:             &e.CreatedAt,
		UpdatedAt:             &e.UpdatedAt,
		DeletedAt:             e.DeletedAt,
	}
}
