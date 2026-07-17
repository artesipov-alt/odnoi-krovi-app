package pet

import (
	"context"
	"time"

	bloodsearchmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	commonmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// PetReadRepository определяет операции чтения для питомцев
type PetReadRepository interface {
	// GetByID получает питомца по ID с опциями загрузки связанных данных
	GetByID(ctx context.Context, id string, opts PetPreloadOptions) (*model.Pet, error)

	// GetByUserID получает всех питомцев пользователя
	GetByUserID(ctx context.Context, userID string, opts PetPreloadOptions) ([]*model.Pet, error)

	// Exists проверяет существование питомца по ID
	Exists(ctx context.Context, id string) (bool, error)

	// ExistsByID алиас для Exists (для обратной совместимости)
	// Deprecated: используйте Exists
	ExistsByID(ctx context.Context, id string) (bool, error)

	//Для уведомлений
	GetPetsByBloodGroupAndRegion(ctx context.Context, petType commonmodel.PetType, bloodGroups, regions []string) ([]*model.Pet, error)

	// GetPetIDsByBloodGroupAndRegion возвращает ID питомцев по группе крови и регионам (raw SQL)
	GetPetIDsByBloodGroupAndRegion(ctx context.Context, bloodGroups, regions []string) ([]string, error)

	// GetByIDs загружает питомцев по слайсу ID с полными данными
	GetByIDs(ctx context.Context, ids []string, opts PetPreloadOptions) ([]*model.Pet, error)

	// FindPotentialDonors возвращает питомцев, открытых для приглашений реципиентов
	// (donor_preference.open_for_contact = true) с подходящей группой крови,
	// регионом и типом, исключая самого реципиента и питомцев, уже откликнувшихся
	// на заявку ExcludeRequestID. Возвращает bloodsearch-модель PotentialDonor —
	// питомца с настройками донорства его владельца.
	FindPotentialDonors(ctx context.Context, criteria PotentialDonorsCriteria) ([]*bloodsearchmodel.PotentialDonor, error)
}

// PetWriteRepository определяет операции записи для питомцев
type PetWriteRepository interface {
	// Create создает нового питомца
	Create(ctx context.Context, pet *model.Pet) (*model.Pet, error)

	// Update обновляет существующего питомца
	Update(ctx context.Context, id string, pet *model.Pet) (*model.Pet, error)

	SetLastDonation(ctx context.Context, petID string, lastDonationDate *time.Time) error

	SetTransfused(ctx context.Context, petID string, transfused bool) error

	// Delete удаляет питомца (soft delete)
	DeleteWithRelations(ctx context.Context, id string) error
}

// PetStatsRepository определяет операции для статистики и поиска доноров
type PetStatsRepository interface {
	// CountSuitableDonors считает количество подходящих доноров для групп крови
	CountSuitableDonors(ctx context.Context, bloodGroups []string) (int, error)
}

// PetPhotoRepository определяет операции для работы с фото питомцев
type PetPhotoRepository interface {
	// AddPhotoURLs добавляет URL фотографий к питомцу
	AddPhotoURLs(ctx context.Context, id string, paths []string) error
}

// PetPreloadOptions определяет опции для загрузки связанных данных
type PetPreloadOptions struct {
	WithHealth       bool
	WithTreatments   bool
	WithAnalyses     bool
	WithBonuses      bool
	WithAll          bool
	IgnoreSoftDelete bool
}

// PotentialDonorsCriteria содержит критерии для поиска потенциальных доноров.
type PotentialDonorsCriteria struct {
	PetType          commonmodel.PetType
	BloodGroups      []string
	Regions          []string
	ExcludeRequestID string
	ExcludePetID     string
	Limit            int
	Offset           int
}

func (pr *PetPreloadOptions) SetIgnoreSoftDelete() {
	pr.IgnoreSoftDelete = true
}

// Repository объединяет все интерфейсы для обратной совместимости
// Deprecated: используйте специализированные интерфейсы
type Repository interface {
	PetReadRepository
	PetWriteRepository
	PetStatsRepository
	PetPhotoRepository
}
