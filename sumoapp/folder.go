package sumoapp

func NewFolder() *folder {
	return &folder{
		Type:          "",
		Name:          "",
		Description:   "",
		Children:      make([]interface{}, 0),
		Items:         make(map[string][]string),
		folders:       make(map[string]*folder),
		dashboards:    make(map[string]dashboard),
		savedSearches: make(map[string]savedSearch),
	}
}
