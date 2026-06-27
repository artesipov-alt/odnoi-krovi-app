package common

type BodyOutput[T any] struct {
	Body T
}

type PortalStats struct {
	TotalUsers             int64   `json:"totalUsers" doc:"Общее количество пользователей"`
	UsersWithPhone         int64   `json:"usersWithPhone" doc:"Количество пользователей с номером телефона"`
	PhoneConversionPercent float64 `json:"phoneConversionPercent" doc:"Процент пользователей с номером телефона"`
	VerifiedUsers          int64   `json:"verifiedUsers" doc:"Количество верифицированных пользователей"`
	UnverifiedUsers        int64   `json:"unverifiedUsers" doc:"Количество неверифицированных пользователей"`
	TotalPets              int64   `json:"totalPets" doc:"Общее количество питомцев"`
	ActiveBloodRequests    int64   `json:"activeBloodRequests" doc:"Количество активных запросов крови"`
	TotalDonations         int64   `json:"totalDonations" doc:"Общее количество донаций"`
	CompletedDonations     int64   `json:"completedDonations" doc:"Количество завершенных донаций"`
}
