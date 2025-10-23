package ds

const (
	L1 = 1
	L2 = 2
)

type Item struct {
	ID               int64
	BrandID          int64
	ParentCategoryID int64
	CategoryID       int64
	TypeID           int64
}

// SetBrandID обновляет BrandID и возвращает true если значение было изменено
func (i *Item) SetBrandID(newID int64) bool {
	prevID := i.BrandID
	i.BrandID = newID
	return prevID != newID
}

// SetCategory обновляет ParentCategoryID, CategoryID и TypeID и возвращает true если хоть одно значение было изменено
func (i *Item) SetCategory(parentCatID, catID, typeID int64) bool {
	if i.ParentCategoryID == parentCatID &&
		i.CategoryID == catID &&
		i.TypeID == typeID {
		return false
	}
	i.ParentCategoryID = parentCatID
	i.CategoryID = catID
	i.TypeID = typeID
	return true
}
