package sumoapp

import (
	"fmt"

	"github.com/imdario/mergo"
)

func NewApplication() *application {
	return &application{
		Name:          "",
		Description:   "",
		Version:       "",
		Type:          FolderType,
		Items:         make(map[string][]string),
		Children:      make([]interface{}, 0),
		panels:        make(map[string]panel),
		dashboards:    make(map[string]dashboard),
		queries:       make(map[string]query),
		folders:       make(map[string]*folder),
		variables:     make(map[string]variable),
		savedSearches: make(map[string]savedSearch),
	}
}

func (a *application) mergeVariables(app *application) error {
	for vname, variableObj := range app.variables {
		//If a panel by the same name already exists, merge it
		//Otherwise, add the new panel to the list of panels
		if _, ok := a.variables[vname]; ok {
			v := a.variables[vname]
			if err := mergo.Merge(&v, variableObj, mergo.WithOverride); err != nil {
				return err
			}

			a.variables[vname] = v
		} else {
			a.variables[vname] = variableObj
		}
	}

	return nil
}

func (a *application) mergePanels(app *application) error {
	for pname, panelObj := range app.panels {
		//If a panel by the same name already exists, merge it
		//Otherwise, add the new panel to the list of panels
		if _, ok := a.panels[pname]; ok {
			p := a.panels[pname]
			if err := mergo.Merge(&p, panelObj, mergo.WithOverride); err != nil {
				return err
			}

			a.panels[pname] = p
		} else {
			a.panels[pname] = panelObj
		}
	}

	return nil
}

func (a *application) mergeDashboards(app *application) error {
	for dname, dashboard := range app.dashboards {
		//If a dashboard by the same name already exists, merge it
		//Otherwise, add the new dashboard to the list of dashboards
		if _, ok := a.dashboards[dname]; ok {
			d := a.dashboards[dname]
			if err := mergo.Merge(&d, dashboard, mergo.WithOverride); err != nil {
				return err
			}

			a.dashboards[dname] = d
		} else {
			a.dashboards[dname] = dashboard
		}
	}

	return nil
}

func (a *application) mergeFolders(app *application) error {
	for fname, folder := range app.folders {
		//If a folder by the same name already exists, merge it
		//Otherwise, add the new folder to the list of folders
		if _, ok := a.folders[fname]; ok {
			f := a.folders[fname]
			if err := mergo.Merge(&f, folder, mergo.WithOverride); err != nil {
				return err
			}

			a.folders[fname] = f
		} else {
			a.folders[fname] = folder
		}
	}

	return nil
}

func (a *application) mergeItems(app *application) error {
	for iname, item := range app.Items {
		//If an item by the same name already exists, merge it
		//Otherwise, add the new item to the list of items
		if _, ok := a.Items[iname]; ok {
			i := a.Items[iname]
			if err := mergo.Merge(&i, item, mergo.WithOverride); err != nil {
				return err
			}

			a.Items[iname] = i
		} else {
			a.Items[iname] = item
		}
	}

	return nil
}

func (a *application) mergeSavedSearches(app *application) error {
	for sname, search := range app.savedSearches {
		//If a search by the same name already exists, merge it
		//Otherwise, add the new search to the list of saved searches
		if _, ok := a.savedSearches[sname]; ok {
			s := a.savedSearches[sname]
			if err := mergo.Merge(&s, search, mergo.WithOverride); err != nil {
				return err
			}

			a.savedSearches[sname] = s
		} else {
			a.savedSearches[sname] = search
		}
	}

	return nil
}

func (a *application) Merge(app *application) error {
	if err := a.mergePanels(app); err != nil {
		return err
	}

	if err := a.mergeVariables(app); err != nil {
		return err
	}

	if err := a.mergeDashboards(app); err != nil {
		return err
	}

	if err := a.mergeFolders(app); err != nil {
		return err
	}

	if err := a.mergeSavedSearches(app); err != nil {
		return err
	}

	if err := a.mergeItems(app); err != nil {
		return err
	}

	if app.Name != "" {
		a.Name = app.Name
	}

	if app.Description != "" {
		a.Description = app.Description
	}

	if app.Version != "" {
		a.Version = app.Version
	}

	return nil
}

func (a *application) populateFolder(f *folder) error {
	for _, folderName := range f.Items["folders"] {
		if folder, ok := a.folders[folderName]; ok {
			if err := a.populateFolder(folder); err != nil {
				return err
			}

			f.Children = append(f.Children, folder)
		} else {
			msg := fmt.Errorf("Unable to populate folder %s. Child folder %s does not exist", folder.Name, folderName)
			return msg
		}
	}

	for _, dashboardName := range f.Items["dashboards"] {
		if dashboard, ok := a.dashboards[dashboardName]; ok {
			f.Children = append(f.Children, dashboard)
		}
	}

	for _, searchName := range f.Items["savedSearches"] {
		if search, ok := a.savedSearches[searchName]; ok {
			f.Children = append(f.Children, search)
		}
	}

	return nil
}

func (a *application) Build() error {
	//Embed panels and variables into dashboards
	for name, dashboard := range a.dashboards {
		for _, layoutPanel := range dashboard.Layout.LayoutStructures {
			if _, ok := a.panels[layoutPanel.Key]; !ok {
				err := fmt.Errorf("Could not find panel '%s'. Referenced in dashboard '%s' layout", layoutPanel.Key, dashboard.Name)
				return err
			}

			p := a.panels[layoutPanel.Key]
			dashboard.Panels = append(dashboard.Panels, p)
		}

		for _, variableName := range dashboard.IncludeVariables {
			v, ok := a.variables[variableName]
			if !ok {
				err := fmt.Errorf("Could not find variable '%s'. Referenced in dashboard '%s'", variableName, dashboard.Name)
				return err
			}

			dashboard.Variables = append(dashboard.Variables, v)
		}

		a.dashboards[name] = dashboard
	}

	//Embed folder items into the folders. This function also recursively iterates
	//through each folder within a folder to populate that
	for _, folderName := range a.Items["folders"] {
		if folder, ok := a.folders[folderName]; ok {
			if err := a.populateFolder(folder); err != nil {
				return err
			}

			a.Children = append(a.Children, folder)
		}
	}

	//Embed dashboarditems into the list of Children
	//TODO: this code is dupliated in the populateFolder() function.
	//  Consider a way to consolidate it
	for _, dashboardName := range a.Items["dashboards"] {
		if dash, ok := a.dashboards[dashboardName]; ok {
			a.Children = append(a.Children, dash)
		}
	}

	//Embed saved search items into the list of Children
	//TODO: this code is dupliated in the populateFolder() function.
	//  Consider a way to consolidate it
	for _, searchName := range a.Items["savedSearches"] {
		if search, ok := a.savedSearches[searchName]; ok {
			a.Children = append(a.Children, search)
		}
	}

	return nil
}
