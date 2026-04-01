package common

// DefaultMessageOutput представляет ответ с общим сообщением
type DefaultMessageOutput struct {
	Body ResultMessage
}

// ResultMessage представляет сообщение о результате операции
type ResultMessage struct {
	Message string `json:"message" doc:"Сообщение о результате операции"`
}

// Основные пути Path

// UserIDPath представляет параметр пути с ID пользователя
type UserIDPath struct {
	ID string `path:"user_id" doc:"ID пользователя" minLength:"1" example:"USR-ABCDEABCDE"`
}

// TelegramIDPath представляет параметр пути с Telegram ID
type TelegramIDPath struct {
	ID int64 `path:"tg_id" doc:"Telegram ID пользователя" minimum:"1" example:"123456789"`
}

// PetPathParam представляет параметр пути с ID питомца
type PetIDPath struct {
	ID string `path:"pet_id" doc:"ID питомца" minLength:"1" example:"PET-aBcDeF1234"`
}

// BloodRequestIDPath представляет параметр пути с ID заявки
type BloodRequestIDPath struct {
	ID string `path:"req_id" doc:"ID заявки на поиск крови" minLength:"1" example:"BLS-ABCDEABCDE"`
}

// DonorApplicationIDPath представляет параметр пути с ID отклика
type DonorApplicationIDPath struct {
	ID string `path:"res_id" doc:"ID отклика донора" minLength:"1" example:"RES-ABCDEABCDE"`
}
