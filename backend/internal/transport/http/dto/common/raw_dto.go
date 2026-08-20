package common

type BodyOutput[T any] struct {
	Body T
}

type PetTypeStats struct {
	Type                string  `json:"type" doc:"Тип питомца"`
	TotalPets           int64   `json:"totalPets" doc:"Общее количество питомцев данного типа"`
	ActiveBloodRequests int64   `json:"activeBloodRequests" doc:"Количество активных запросов крови для данного типа"`
	TotalDonations      int64   `json:"totalDonations" doc:"Общее количество донаций для данного типа"`
	CompletedDonations  int64   `json:"completedDonations" doc:"Количество завершенных донаций для данного типа"`
	Searches            int64   `json:"searches" doc:"Количество поисков для данного типа"`
	SearchVolume        float64 `json:"searchVolume" doc:"Объём поисков для данного типа"`
	DonationVolume      float64 `json:"donationVolume" doc:"Объём донаций для данного типа"`
}

type PortalStats struct {
	TotalUsers             int64        `json:"totalUsers" doc:"Общее количество пользователей"`
	UsersWithPhone         int64        `json:"usersWithPhone" doc:"Количество пользователей с номером телефона"`
	PhoneConversionPercent float64      `json:"phoneConversionPercent" doc:"Процент пользователей с номером телефона"`
	VerifiedUsers          int64        `json:"verifiedUsers" doc:"Количество верифицированных пользователей"`
	UnverifiedUsers        int64        `json:"unverifiedUsers" doc:"Количество неверифицированных пользователей"`
	TotalPets              int64        `json:"totalPets" doc:"Общее количество питомцев"`
	ActiveBloodRequests    int64        `json:"activeBloodRequests" doc:"Количество активных запросов крови"`
	TotalDonations         int64        `json:"totalDonations" doc:"Общее количество донаций"`
	CompletedDonations     int64        `json:"completedDonations" doc:"Количество завершенных донаций"`
	TotalSearches          int64        `json:"totalSearches" doc:"Общее количество поисков"`
	TotalSearchVolume      float64      `json:"totalSearchVolume" doc:"Общий объём поисков"`
	TotalDonationVolume    float64      `json:"totalDonationVolume" doc:"Общий объём донаций"`
	CatStats               PetTypeStats `json:"catStats" doc:"Статистика по кошкам"`
	DogStats               PetTypeStats `json:"dogStats" doc:"Статистика по собакам"`
}
