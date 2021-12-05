package sumoapp

import "github.com/imdario/mergo"

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

func (f *folder) Merge(folderObj *folder) error {
	if err := mergo.Merge(f, folderObj, mergo.WithOverride); err != nil {
		return err
	}

	return nil
}

func (f *folder) Copy() *folder {
	folderList := make(map[string]*folder)

	newFolder := NewFolder()

	newFolder.Type = f.Type
	newFolder.Name = f.Name
	newFolder.Description = f.Description
	newFolder.Children = f.Children
	newFolder.Items = f.Items
	newFolder.dashboards = f.dashboards
	newFolder.savedSearches = f.savedSearches

	for folderName, folderToCopy := range f.folders {
		folderList[folderName] = folderToCopy.Copy()
	}

	newFolder.folders = folderList

	return newFolder
}
