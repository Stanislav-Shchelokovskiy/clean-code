package ds

type Item struct {
	ID               int64
	BrandID          int64
	ParentCategoryID int64
	CategoryID       int64
	TypeID           int64
}

// Items ...
type Items []*Item

// MapBySkuKey ...
func (items Items) AsMap() map[int64]*Item {
	itemsMap := make(map[int64]*Item, len(items))
	for _, item := range items {
		itemsMap[item.ID] = item
	}
	return itemsMap
}
