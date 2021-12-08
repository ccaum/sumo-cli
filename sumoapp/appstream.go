package sumoapp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
	"github.com/r3labs/diff"
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

func (s *appStream) Diff(diffStream *appStream) (diff.Changelog, error) {
	changelogVar, err := diff.Diff(s.Variables, diffStream.Variables)
	if err != nil {
		return nil, err
	}

	changelogPanel, err := diff.Diff(s.Panels, diffStream.Panels)
	if err != nil {
		return nil, err
	}

	changelogSavedSearches, err := diff.Diff(s.SavedSearches, diffStream.SavedSearches)
	if err != nil {
		return nil, err
	}

	changelogDashboard, err := diff.Diff(s.Dashboards, diffStream.Dashboards)
	if err != nil {
		return nil, err
	}

	changelogFolder, err := diff.Diff(s.Folders, diffStream.Folders)
	if err != nil {
		return nil, err
	}

	allChanges := []diff.Changelog{
		changelogVar, changelogPanel, changelogSavedSearches, changelogDashboard, changelogFolder,
	}

	var changelogs diff.Changelog
	for _, c := range allChanges {
		changelogs = append(changelogs, c...)
	}

	fmt.Println("Found", len(changelogs), "changes")
	for _, change := range changelogs {
		fmt.Println("")
		fmt.Println("TYPE: ", change.Type)
		fmt.Println("PATH: ", change.Path)
		fmt.Println("FROM: ", change.From)
		fmt.Println("TO: ", change.To)
	}

	return changelogs, nil
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
	dashboards := make(map[string]*dashboard)

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

	s.Dashboards = dashboards

	//Append the dashboards defined in the parent stream that are NOT
	//overwritten in this stream. Dashboards that have overwrites in this
	//stream will be merged with their parent dashboard and the new
	//object will be added to this stream's list of dashboards
	if s.HasParent() {
		for name, pd := range s.Parent.Dashboards {
			//If the parent stream has a dashboard by the same name
			//merge the current dashboard with its parent
			d, ok := s.Dashboards[name]
			if !ok {
				s.Dashboards[name] = pd
			} else {
				if err := d.Merge(pd); err != nil {
					return err
				}

				s.Dashboards[name] = d
			}
		}
	}

	//TODO: This should leverage go functions to parallelize the
	//population of each dashboard. There's no reason for it to be
	//serialized
	for name, dash := range s.Dashboards {
		dash.key = name
		dash.Type = DashboardType

		//Ensure we have clean lists in case the object was inherited
		//from the parent stream. The Populate() function will repopulate
		//the panels and variables
		dash.Panels = make([]*panel, 0)
		dash.Variables = make([]*variable, 0)

		if dash.Populate(s); err != nil {
			return err
		}
	}

	return nil
}

func (s *appStream) loadVariables(basePath string) error {
	variables := make(map[string]*variable)

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

	s.Variables = variables

	//Append the variables defined in the parent stream that are NOT
	//overwritten in this stream. Variables that have overwrites in this
	//stream will be merged with their parent variable and the new
	//object will be added to this stream's list of variables
	if s.HasParent() {
		for name, pv := range s.Parent.Variables {
			//If the parent stream has a variable by the same name
			//merge the current variable with its parent
			v, ok := s.Variables[name]
			if !ok {
				s.Variables[name] = pv
			} else {
				if err := v.Merge(pv); err != nil {
					return err
				}

				s.Variables[name] = v
			}
		}
	}

	for name, varObj := range variables {
		varObj.Name = name
	}

	return nil
}

func (s *appStream) loadPanels(basePath string) error {
	panels := make(map[string]*panel)

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

	s.Panels = panels

	//Append the panels defined in the parent stream that are NOT
	//overwritten in this stream. Panels that have overwrites in this
	//stream will be merged with their parent panel and the new
	//object will be added to this stream's list of dashboards
	if s.HasParent() {
		for name, pp := range s.Parent.Panels {
			//If the parent stream has a panel by the same name
			//merge the current dashboard with its parent
			p, ok := s.Panels[name]
			if !ok {
				s.Panels[name] = pp
			} else {
				if err := p.Merge(pp); err != nil {
					return err
				}

				s.Panels[name] = p
			}
		}
	}

	for name, pan := range s.Panels {
		pan.Key = name
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
	folders := make(map[string]*folder)

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

	s.Folders = folders

	//Append the folders defined in the parent stream that are NOT
	//overwritten in this stream. Folders that have overwrites in this
	//stream will be merged with their parent folder and the new
	//object will be added to this stream's list of folders
	if s.HasParent() {
		for name, pf := range s.Parent.Folders {
			//If the parent stream has a folder by the same name
			//merge the current folder with its parent
			f, ok := s.Folders[name]
			if !ok {
				s.Folders[name] = pf
			} else {
				if err := f.Merge(pf); err != nil {
					return err
				}

				s.Folders[name] = f
			}
		}
	}

	//Ensure all the folder objects are annotated properly.
	//There might be a more efficient way to do this
	for _, foldObj := range s.Folders {
		foldObj.Type = FolderType
	}

	return nil
}

func (s *appStream) loadSavedSearches(basePath string) error {
	searches := make(map[string]*savedSearch)

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

	s.SavedSearches = searches

	//Append the searches defined in the parent stream that are NOT
	//overwritten in this stream. Searches that have overwrites in this
	//stream will be merged with their parent search and the new
	//object will be added to this stream's list of searches
	if s.HasParent() {
		for name, ps := range s.Parent.SavedSearches {
			//If the parent stream has a search by the same name
			//merge the current search with its parent
			ss, ok := s.SavedSearches[name]
			if !ok {
				s.SavedSearches[name] = ps
			} else {
				if err := ss.Merge(ps); err != nil {
					return err
				}

				s.SavedSearches[name] = ss
			}
		}
	}

	//Ensure all the search objects are annotated properly.
	//There might be a more efficient way to do this
	for _, search := range s.Folders {
		search.Type = SavedSearchType
	}

	return nil
}

func (s *appStream) Load() error {
	var err error

	panelBasePath := fmt.Sprintf("%s/panels", s.Path)
	dashboardBasePath := fmt.Sprintf("%s/dashboards", s.Path)
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
