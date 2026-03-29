package dto

// DefaultMessageOutput представляет ответ с общим сообщением
type DefaultMessageOutput struct {
	Body ResultMessage
}

// ResultMessage представляет сообщение о результате операции
type ResultMessage struct {
	Message string `json:"message" doc:"Сообщение о результате операции"`
}
