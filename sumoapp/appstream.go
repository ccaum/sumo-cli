package sumoapp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

func NewAppStream(name string, app *application) *appStream {
	path := fmt.Sprintf("%s/%s", app.BasePath(), name)

	return &appStream{
		Name:          name,
		Application:   app,
		Path:          path,
		Dashboards:    make(map[string]*dashboard),
		Variables:     make(map[string]*variable),
		Panels:        make(map[string]*panel),
		SavedSearches: make(map[string]*savedSearch),
		Folders:       make(map[string]*folder),
		Queries:       make(map[string]*query),
	}
}

func (s *appStream) HasParent() bool {
	if s.Parent == nil {
		return false
	}

	return true
}

func (s *appStream) WriteObjects() error {
	//Write the folder objects to the app stream
	for fName, folderObj := range s.Folders {
		folderMap := make(map[string]*folder)
		folderMap[fName] = folderObj

		y, err := yaml.Marshal(folderMap)
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("%s/folders/%s.yaml", s.Path, fName)
		if err := os.WriteFile(filePath, y, 0644); err != nil {
			return err
		}
	}

	//Write the dashboard objects to the app stream
	for dName, dashboardObj := range s.Dashboards {
		dashMap := make(map[string]*dashboard)
		dashMap[dName] = dashboardObj

		y, err := yaml.Marshal(dashMap)
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("%s/dashboards/%s.yaml", s.Path, dName)
		if err := os.WriteFile(filePath, y, 0644); err != nil {
			return err
		}
	}

	//Write the panel objects to the app stream
	for pName, panelObj := range s.Panels {
		panelMap := make(map[string]*panel)
		panelMap[pName] = panelObj

		p, err := yaml.Marshal(panelMap)
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("%s/panels/%s.yaml", s.Path, pName)
		if err := os.WriteFile(filePath, p, 0644); err != nil {
			return err
		}
	}

	//Write the variable objects to the app stream
	for vName, variableObj := range s.Variables {
		variableMap := make(map[string]*variable)
		variableMap[vName] = variableObj

		v, err := yaml.Marshal(variableMap)
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("%s/variables/%s.yaml", s.Path, vName)
		if err := os.WriteFile(filePath, v, 0644); err != nil {
			return err
		}
	}

	//Write the saved search objects to the app stream
	for sName, searchObj := range s.SavedSearches {
		searchMap := make(map[string]*savedSearch)
		searchMap[sName] = searchObj

		v, err := yaml.Marshal(searchMap)
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("%s/saved-searches/%s.yaml", s.Path, sName)
		if err := os.WriteFile(filePath, v, 0644); err != nil {
			return err
		}
	}

	//Write the application's definition to the init file in the stream
	a, err := yaml.Marshal(s.Application)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/init.yaml", s.Path)
	if err := os.WriteFile(filePath, a, 0644); err != nil {
		return err
	}

	return nil
}

func (s *appStream) populateFolder(f *folder) error {
	//Children was inherited from the parent stream. It needs to be cleared
	//or objects will be appear in the list multiple times.
	//TODO: this shouldn't be necessary. Can we avoid this in the
	//Merge() function?
	f.Children = make([]interface{}, 0)

	for _, folderName := range f.Items["folders"] {
		if folder, ok := s.Folders[folderName]; ok {
			if err := s.populateFolder(folder); err != nil {
				return err
			}

			f.Children = append(f.Children, folder)
		} else {
			msg := fmt.Errorf("Unable to populate folder %s. Child folder %s does not exist", folderName, folderName)
			return msg
		}
	}

	for _, dashboardName := range f.Items["dashboards"] {
		if dashboard, ok := s.Dashboards[dashboardName]; ok {
			f.Children = append(f.Children, dashboard)
		}
	}

	for _, searchName := range f.Items["savedSearches"] {
		if search, ok := s.SavedSearches[searchName]; ok {
			f.Children = append(f.Children, search)
		}
	}

	return nil
}

func (s *appStream) loadDashboards(basePath string) error {
	var dashboards map[string]*dashboard

	dfiles, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range dfiles {
		var curList map[string]*dashboard

		path := fmt.Sprintf("%s/%s", basePath, file.Name())
		extension := filepath.Ext(path)
		if extension != ".yaml" {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			err := fmt.Errorf("Unable to read file ", path, ": ", err)
			return err
		}

		if err := yaml.Unmarshal(data, &curList); err != nil {
			return err
		}

		if err := mergo.Merge(&dashboards, curList); err != nil {
			return err
		}
	}

	//First copy the parent's panel list into a local list that
	//can be overwritten with new panel values. Each dashboard
	//will need to be iterated over and populated since this app
	//stream may have overwritten panels or variables in dashboards
	//defined in previous app streams
	if s.HasParent() {
		s.Dashboards = s.Parent.Dashboards
	}

	ds := s.Dashboards
	if err := mergo.Merge(&ds, dashboards, mergo.WithOverride); err != nil {
		return err
	}

	//TODO: This should leverage go functions to parallelize the
	//population of each dashboard. There's no reason for it to be
	//serialized
	for name, dash := range s.Dashboards {
		dash.key = name

		if dash.Populate(s); err != nil {
			return err
		}
	}

	return nil
}

func (s *appStream) loadVariables(basePath string) error {
	var variables map[string]*variable

	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		var curList map[string]*variable

		path := fmt.Sprintf("%s/%s", basePath, file.Name())
		extension := filepath.Ext(path)
		if extension != ".yaml" {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			err := fmt.Errorf("Unable to read file ", path, ": ", err)
			return err
		}

		if err := yaml.Unmarshal(data, &curList); err != nil {
			return err
		}

		if err := mergo.Merge(&variables, curList); err != nil {
			return err
		}
	}

	//First copy the parent's variables list into a local list that
	//can be overwritten with new panel values
	//Later, variables found in this stream will need to be merged
	//with the parent variables's object, if it exists, before updating
	//this stream's variable pointer to the new merged object
	if s.HasParent() {
		s.Variables = s.Parent.Variables
	}

	for name, varObj := range variables {
		//If the parent stream already has a variable in
		//its component library, merge it with the one in this
		//app stream and save the new merged object in this app stream's
		//component library
		if s.HasParent() {
			appVar, err := s.Parent.FindVariable(name)
			if err == nil {
				varObj.Merge(appVar)
			}
		}

		varObj.Name = name

		s.Variables[name] = varObj
	}

	return nil
}

func (s *appStream) loadPanels(basePath string) error {
	var panels map[string]*panel

	//Load the panel files from this stream's file system
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		var curList map[string]*panel

		path := fmt.Sprintf("%s/%s", basePath, file.Name())
		extension := filepath.Ext(path)
		if extension != ".yaml" {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			err := fmt.Errorf("Unable to read file ", path, ": ", err)
			return err
		}

		if err := yaml.Unmarshal(data, &curList); err != nil {
			return err
		}

		if err := mergo.Merge(&panels, curList); err != nil {
			return err
		}
	}

	//First copy the parent's panel list into a local list that
	//can be overwritten with new panel values
	//Later, panels found in this stream will need to be merged
	//with the parent panel's object, if it exists, before updating
	//this stream's panel pointer to the new merged object
	if s.HasParent() {
		s.Panels = s.Parent.Panels
	}

	for name, panel := range panels {
		//If the parent stream already has a panel in
		//its component library, merge it with the one in this
		//app stream and save the new merged object in this app stream's
		//component library
		if s.HasParent() {
			appPanel, err := s.Parent.FindPanel(name)
			if err == nil {
				panel.Merge(appPanel)
			}
		}

		panel.Key = name
		s.Panels[name] = panel
	}

	return nil
}

func (s *appStream) FindDashboard(name string) (*dashboard, error) {
	dash, ok := s.Dashboards[name]
	if !ok {
		err := fmt.Errorf("Could not find dashboard '%s'", name)
		return nil, err
	}

	return dash, nil
}

func (s *appStream) FindVariable(name string) (*variable, error) {
	varObj, ok := s.Variables[name]
	if !ok {
		err := fmt.Errorf("Could not find variable '%s'", name)
		return nil, err
	}

	return varObj, nil
}

func (s *appStream) FindSavedSearch(name string) (*savedSearch, error) {
	search, ok := s.SavedSearches[name]
	if !ok {
		err := fmt.Errorf("Could not find saved search '%s'", name)
		return nil, err
	}

	return search, nil
}

func (s *appStream) FindFolder(name string) (*folder, error) {
	folderObj, ok := s.Folders[name]
	if !ok {
		err := fmt.Errorf("Could not find folder '%s'", name)
		return nil, err
	}

	return folderObj, nil
}

func (s *appStream) FindPanel(name string) (*panel, error) {
	pan, ok := s.Panels[name]
	if !ok {
		err := fmt.Errorf("Could not find panel '%s'", name)
		return nil, err
	}

	return pan, nil
}

func (s *appStream) loadRootFolder(appFilePath string) error {
	var root folder

	data, err := os.ReadFile(appFilePath)
	if err == nil {
		if err := yaml.Unmarshal(data, &root); err != nil {
			return err
		}
	}

	root.Type = FolderType

	//Before the folder is populated, the items needs to be merged
	//with the parent's items, if a parent exists
	if s.HasParent() {
		err := mergo.Merge(&root.Items, s.Parent.RootFolder.Items)
		if err != nil {
			return err
		}
	}

	s.populateFolder(&root)

	s.RootFolder = &root

	if root.Description != "" {
		s.Application.Description = root.Description
	}

	if root.Name != "" {
		s.Application.Name = root.Name
	}

	//Update the application's children to be this stream's children
	s.Application.Children = root.Children

	return nil
}

func (s *appStream) loadFolders(basePath string) error {
	var folders map[string]*folder

	ffiles, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range ffiles {
		var curList map[string]*folder

		path := fmt.Sprintf("%s/%s", basePath, file.Name())
		extension := filepath.Ext(path)
		if extension != ".yaml" {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			err := fmt.Errorf("Unable to read file ", path, ": ", err)
			return err
		}

		if err := yaml.Unmarshal(data, &curList); err != nil {
			return err
		}

		if err := mergo.Merge(&folders, curList); err != nil {
			return err
		}
	}

	//First copy the parent's folder list into a local list that
	//can be overwritten with new folder values
	//Later, folders found in this stream will need to be merged
	//with the parent folders's object, if it exists, before updating
	//this stream's folder pointer to the new merged object
	if s.HasParent() {
		s.Folders = s.Parent.Folders
	}

	for name, foldObj := range folders {
		//If the parent stream already has a folder with the same name in
		//its component library, merge it with the one in this
		//app stream and save the new merged object in this app stream's
		//component library
		if s.HasParent() {
			appFolder, err := s.Parent.FindFolder(name)
			if err == nil {
				foldObj.Merge(appFolder)
			}
		}

		foldObj.Type = FolderType
		s.Folders[name] = foldObj
	}

	return nil
}

func (s *appStream) loadSavedSearches(basePath string) error {
	var searches map[string]*savedSearch

	sfiles, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range sfiles {
		var curList map[string]*savedSearch

		path := fmt.Sprintf("%s/%s", basePath, file.Name())
		extension := filepath.Ext(path)
		if extension != ".yaml" {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			err := fmt.Errorf("Unable to read file ", path, ": ", err)
			return err
		}

		if err := yaml.Unmarshal(data, &curList); err != nil {
			return err
		}

		if err := mergo.Merge(&searches, curList); err != nil {
			return err
		}
	}

	//First copy the parent's saved search list into a local list that
	//can be overwritten with new saved search values
	//Later, saved searches found in this stream will be merged
	//with the parent search's object, if it exists, before updating
	//this stream's search pointer to the new merged object
	if s.HasParent() {
		s.SavedSearches = s.Parent.SavedSearches
	}

	for name, search := range searches {
		//If the parent stream already has a saved search with the same name in
		//its component library, merge it with the one in this
		//app stream and save the new merged object in this app stream's
		//component library
		if s.HasParent() {
			appSearch, err := s.Parent.FindSavedSearch(name)
			if err == nil {
				search.Merge(appSearch)
			}
		}

		search.Type = SavedSearchType
		s.SavedSearches[name] = search
	}

	return nil
}

func (s *appStream) Load() error {
	var err error

	panelBasePath := fmt.Sprintf("%s/panels", s.Path)
	dashboardBasePath := fmt.Sprintf("%s//dashboards", s.Path)
	folderBasePath := fmt.Sprintf("%s/folders", s.Path)
	variableBasePath := fmt.Sprintf("%s/variables", s.Path)
	searchesBasePath := fmt.Sprintf("%s/saved-searches", s.Path)

	//It's important the components be loaded in
	//the correct order. Variables and panels should
	//be loaded before dashboards since dashboards
	//reference variables and panels. Folders should
	//be loaded last since they can contain
	//saved searches and dashboards

	err = s.loadVariables(variableBasePath)
	if err != nil {
		err := fmt.Errorf("Could not load variables at %s: %w", variableBasePath, err)
		return err
	}

	err = s.loadPanels(panelBasePath)
	if err != nil {
		err := fmt.Errorf("Could not load panels at %s: %w", panelBasePath, err)
		return err
	}

	err = s.loadDashboards(dashboardBasePath)
	if err != nil {
		err := fmt.Errorf("Could not load dashboards at %s: %w", dashboardBasePath, err)
		return err
	}

	err = s.loadSavedSearches(searchesBasePath)
	if err != nil {
		err := fmt.Errorf("Could not load saved searches at %s: %w", folderBasePath, err)
		return err
	}

	err = s.loadFolders(folderBasePath)
	if err != nil {
		err := fmt.Errorf("Could not load folders at %s: %w", folderBasePath, err)
		return err
	}

	rootPath := fmt.Sprintf("%s/init.yaml", s.Path)
	err = s.loadRootFolder(rootPath)
	if err != nil {
		noSuchFileError := fmt.Sprintf("open %s: no such file or directory", rootPath)
		if err.Error() == noSuchFileError {
			return nil
		} else {
			err := fmt.Errorf("Could not load root application at %s: %w", rootPath, err)
			return err
		}
	}

	return nil
}
