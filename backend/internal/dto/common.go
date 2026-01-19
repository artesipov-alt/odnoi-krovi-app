package dto

// MessageBody представляет тело простого текстового ответа
type MessageBody struct {
	Message string `json:"message" doc:"Сообщение об успехе или ошибке"`
}

// MessageResponse представляет простой текстовый ответ
type MessageResponse struct {
	Body MessageBody
}

// IDPath представляет параметры пути с ID объекта
type IDPath struct {
	ID string `path:"id" doc:"ID объекта (пользователя/питомца)" minLength:"1" example:"USR-ABCDEABCDE"`
}
