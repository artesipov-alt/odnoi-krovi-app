package dto

// ReferenceItemCode представляет элемент справочных данных, определенных в коде (строковое значение)
type ReferenceItemCode struct {
	Value string `json:"value" doc:"Значение элемента справочника"`
	Label string `json:"label" doc:"Отображаемое название элемента справочника"`
}

// ReferenceDataCode представляет ответ со справочными данными, определенными в коде
type ReferenceDataCode struct {
	Data []ReferenceItemCode `json:"data" doc:"Список элементов справочника"`
}

// ReferenceItemLocal представляет элемент справочных данных (строковое значение)
type ReferenceItemLocal struct {
	Value string `json:"value" doc:"Значение элемента справочника"`
	Label string `json:"label" doc:"Отображаемое название элемента справочника"`
}

// ReferenceDataLocal представляет ответ со справочными данными
type ReferenceDataLocal struct {
	Data []ReferenceItemLocal `json:"data" doc:"Список элементов справочника"`
}

// ReferenceItemDB представляет элемент справочных данных из базы данных (целочисленное значение)
type ReferenceItemDB struct {
	Value int    `json:"value" doc:"ID элемента справочника"`
	Label string `json:"label" doc:"Название элемента справочника"`
}

// ReferenceDataDB представляет ответ со справочными данными из базы данных
type ReferenceDataDB struct {
	Data []ReferenceItemDB `json:"data" doc:"Список элементов справочника"`
}

// ReferenceCodeResponse представляет обертку для Huma OpenAPI документации для списка строковых справочных элементов, определенных в коде.
type ReferenceCodeResponse struct {
	Body ReferenceDataCode
}

// ReferenceLocalResponse представляет обертку для Huma OpenAPI документации для списка строковых справочных элементов.
type ReferenceLocalResponse struct {
	Body ReferenceDataLocal
}

// ReferenceDBResponse представляет обертку для Huma OpenAPI документации для списка целочисленных справочных элементов из базы данных.
type ReferenceDBResponse struct {
	Body ReferenceDataDB
}

// PetTypePath представляет параметр пути для типа питомца
type PetTypePath struct {
	PetType string `path:"pet_type" doc:"Тип животного (например, dog, cat)" example:"dog"`
}

// PetTypeQuery представляет параметр запроса для типа питомца
type PetTypeQuery struct {
	PetType string `query:"petType" doc:"Тип животного (dog, cat, etc.)" example:"dog"`
}
