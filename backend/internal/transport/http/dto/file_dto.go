package dto

// ============================================
// Query Parameters
// ============================================

// UploadQuery представляет параметры для получения URL загрузки
type UploadQuery struct {
	PhotosCount   int64 `query:"photos_count" doc:"Количество фотографий для загрузки" minimum:"1" maximum:"10" example:"1"`
	ForPetAvatar  bool  `query:"for_pet_avatar" doc:"URL для аватарки питомца"`
	ForUserAvatar bool  `query:"for_user_avatar" doc:"URL для аватарки пользователя"`
	ForBloodReq   bool  `query:"for_blood_req" doc:"URL для заявки на поиск крови"`
}

// ============================================
// Get Upload URLs
// ============================================

// GetUploadURLsInput представляет запрос на получение URL для загрузки файлов
type GetUploadURLsInput struct {
	ID string `path:"id" doc:"ID сущности" minLength:"1" example:"PET-ABCDEABCDE"`
	UploadQuery
}

// GetUploadURLsOutput представляет ответ с URL для загрузки файлов
type GetUploadURLsOutput struct {
	Body UploadURLsResult
}

// UploadURLsResult представляет результат с URL для загрузки
type UploadURLsResult struct {
	Items []UploadItem `json:"items" doc:"Список подписанных URL для загрузки файлов"`
}

// UploadItem представляет элемент загрузки с подписанным URL
type UploadItem struct {
	URL  string `json:"url" doc:"Подписанная ссылка для загрузки файла"`
	Path string `json:"path" doc:"Путь к файлу в хранилище"`
}

// ============================================
// Confirm Upload
// ============================================

// ConfirmUploadInput представляет запрос на подтверждение загрузки файлов
type ConfirmUploadInput struct {
	Body ConfirmUploadBody
}

// ConfirmUploadBody представляет тело запроса на подтверждение загрузки
type ConfirmUploadBody struct {
	EntityID string   `json:"entityId" doc:"ID сущности (питомец/пользователь/заявка)" minLength:"1" example:"PET-ABCDEABCDE"`
	Paths    []string `json:"paths" doc:"Массив путей к загруженным файлам" validate:"required,dive,min=1,max=255"`
}

// ConfirmUploadOutput представляет ответ на подтверждение загрузки
type ConfirmUploadOutput struct {
	Body ConfirmUploadResult
}

// ConfirmUploadResult представляет результат подтверждения загрузки
type ConfirmUploadResult struct {
	Message string `json:"message" doc:"Сообщение о результате операции"`
}
