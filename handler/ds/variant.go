package ds

const (
	brandAttrID = 1
	typeAttrID  = 2
)

type VariantChangeEvent struct {
	Variants []Variant `json:"variants"` // Список товаров
	Diff     []Diff    `json:"diff"`     // Отличия в сравнении
}

type Variant struct {
	ItemIDs    []int64     `json:"ids"`        // Список идентификаторов товаров
	Attributes []Attribute `json:"attributes"` // Атрибуты на товаре
	Categories []Category  `json:"categories"` // Категории на товаре
}

// GetIDs пытается получить brandID и typeID
func (v *Variant) GetIDs(brandID, typeID *int64) {
	for _, attr := range v.Attributes {
		if brandID != nil && attr.getBrandID(brandID) && typeID == nil {
			return
		}
		if typeID != nil && attr.getDescTypeID(typeID) && brandID == nil {
			return
		}
		if brandID != nil && *brandID > 0 && typeID != nil && *typeID > 0 {
			return
		}
	}
}

// GetCategories возвращает категории варианта { lvl: id }
func (v *Variant) GetCategories() map[int64]int64 {
	cats := make(map[int64]int64, L2)
	for _, cat := range v.Categories {
		cats[cat.Level] = cat.ID
	}
	return cats
}

type Attribute struct {
	ID     int64   `json:"id"`     // Идентификатор атрибута
	Values []int64 `json:"values"` // Значения атрибута
}

// getBrandID пытается получить brandID и возвращает признак успеха
func (a *Attribute) getBrandID(out *int64) bool {
	if a.ID == brandAttrID {
		for _, val := range a.Values {
			*out = val
			return true
		}
	}
	return false
}

// getDescTypeID пытается получить descTypeID и возвращает признак успеха
func (a *Attribute) getDescTypeID(out *int64) bool {
	if a.ID == typeAttrID {
		for _, val := range a.Values {
			*out = val
			return true
		}
	}
	return false
}

type Category struct {
	ID    int64 `json:"id"`    // Идентификатор категории
	Level int64 `json:"level"` // Уровень категории
}

type Diff struct {
	Path string `json:"path"` // Что изменилось
}
