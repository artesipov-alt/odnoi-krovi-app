package domain

type PhotoPreloadQuery struct {
	PhotosCount   int64 `query:"photos_count" doc:"Количество фотографий прикрепленных пользователем" example:"1"`
	ForPetAvatar  bool  `query:"for_pet_avatar" doc:"Получение ссылок для аватарки питомца"`
	ForUserAvatar bool  `query:"for_user_avatar" doc:"Получение ссылок для аватарки пользователя"`
	ForPetBlood   bool  `query:"for_blood_req" doc:"Получение ссылок для заявки на поиск крови для питомца"`
}

type PathParam struct {
	Path string `path:"path" doc:"Путь к фото" example:"pets/PET-aBcDeF1234/avatar.jpg"`
}

type UploadItem struct {
	URL  string `json:"url" doc:"Подписанная ссылка для загрузки"`
	Path string `json:"path" doc:"Путь к файлу в хранилище"`
}

type UploadURLResponse struct {
	Body struct {
		Items []UploadItem `json:"items" doc:"Список ссылок для загрузки"`
	}
}

type ConfirmUploadRequest struct {
	EntityID string   `json:"entityId" doc:"ID сущности (питомец/пользователь/заявка)" example:"PET-aBcDeF1234"`
	Paths    []string `json:"paths" doc:"Массив путей к загруженным фото" example:"pets/PET-aBcDeF1234/photos/1.jpg"`
}

type ConfirmUploadResponse struct {
	Body struct {
		Message string `json:"message" doc:"Сообщение об успехе"`
	}
}
