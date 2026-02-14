package dto

// ReferenceItemLocal представляет элемент справочных данных (строковое значение)
type ReferenceItem struct {
	Value string `json:"value" doc:"Значение элемента справочника"`
	Label string `json:"label" doc:"Отображаемое название элемента справочника"`
}

// ReferenceDataLocal представляет ответ со справочными данными
type ReferenceData struct {
	Data []ReferenceItem `json:"data" doc:"Список элементов справочника"`
}

// ReferenceLocalResponse представляет обертку для Huma OpenAPI документации для списка строковых справочных элементов.
type ReferenceResponse struct {
	Body ReferenceData
}

// PetTypePath представляет параметр пути для типа питомца
type PetTypePath struct {
	PetType string `path:"pet_type" doc:"Тип животного (например, dog, cat)" example:"dog"`
}

// PetTypeQuery представляет параметр запроса для типа питомца
type PetTypeQuery struct {
	PetType string `query:"petType" doc:"Тип животного (dog, cat, etc.)" example:"dog"`
}
