package ds

type Item struct {
	ID   int64
	Name string
}

type RawItem struct {
	ID   int64
	Name string
}

func CreateItems(ids []int64, names map[int64]string) []*Item {
	items := make([]*Item, 0, len(ids))
	for _, id := range ids {
		items = append(items, &Item{ID: id, Name: names[id]})
	}
	return items
}
