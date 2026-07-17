package model

import (
	common "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// PotentialDonor представляет питомца, открытого для приглашений реципиентов
// (donor_preference.open_for_contact = true), вместе с настройками донорства его владельца.
// Это проектный (read-side) тип bloodsearch-контекста: собирается из Pet aggregate
// и User.DonorPreference для отображения в списке потенциальных доноров.
type PotentialDonor struct {
	Pet                  *petmodel.Pet
	CompensationType     common.CompensationType
	TaxiCompensation     bool
	RecoveryPeriodMonths int
}
