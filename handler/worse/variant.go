package worse

type VariantChangeEvent struct {
	Variants []Variant `json:"variants"` // Список товаров
	Diff     []Diff    `json:"diff"`     // Отличия в сравнении
}

type Variant struct {
	ItemIDs    []int64     `json:"ids"`        // Список идентификаторов товаров
	Attributes []Attribute `json:"attributes"` // Атрибуты на товаре
	Categories []Category  `json:"categories"` // Категории на товаре
}

type Attribute struct {
	ID     int64   `json:"id"`     // Идентификатор атрибута
	Values []int64 `json:"values"` // Значения атрибута
}

type Category struct {
	ID    int64 `json:"id"`    // Идентификатор категории
	Level int64 `json:"level"` // Уровень категории
}

type Diff struct {
	Path string `json:"path"` // Что изменилось
}
