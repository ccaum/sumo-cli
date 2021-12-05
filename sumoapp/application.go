package sumoapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/r3labs/diff/v2"
)

func NewApplicationWithPath(path string) *application {
	app := NewApplication()
	app.path = path

	return app
}

func NewApplication() *application {
	return &application{
		Name:        "",
		Description: "",
		Version:     "",
		Type:        FolderType,
		Items:       make(map[string][]string),
		Children:    make([]interface{}, 0),
		path:        "",
	}
}

func sanitizeName(name string) string {
	replacer := strings.NewReplacer("/", "_", ";", "-", " ", "-")
	sanitizedString := replacer.Replace(name)
	sanitizedString = strings.ToLower(sanitizedString)

	return sanitizedString
}

func processChildren(folder *folder, stream *appStream) error {
	for _, childObj := range folder.Children {
		if err := processChild(childObj, folder, stream); err != nil {
			return err
		}
	}

	//These fields aren't needed in the user facing code files
	folder.Children = nil
	folder.Type = ""

	return nil
}

func processChild(obj interface{}, parent *folder, stream *appStream) error {
	childType := obj.(map[string]interface{})["type"]

	switch childType {
	case FolderType:
		folderObj := NewFolder()

		if err := mapstructure.Decode(obj, &folderObj); err != nil {
			return err
		}

		safeName := sanitizeName(folderObj.Name)
		stream.Folders[safeName] = folderObj

		parent.Items["folders"] = append(parent.Items["folders"], safeName)
		processChildren(folderObj, stream)

	case DashboardType:
		var dashboardObj dashboard

		if err := mapstructure.Decode(obj, &dashboardObj); err != nil {
			return err
		}

		safeName := sanitizeName(dashboardObj.Name)

		//Add each of the panels to the application
		for _, panelObj := range dashboardObj.Panels {
			stream.Panels[panelObj.Key] = panelObj
		}

		//Add each of the variables to the application
		for _, variableObj := range dashboardObj.Variables {
			stream.Variables[variableObj.Name] = variableObj
			dashboardObj.IncludeVariables = append(dashboardObj.IncludeVariables, variableObj.Name)
		}

		//Remove these objects from the user facing code files
		dashboardObj.Panels = nil
		dashboardObj.Variables = nil
		dashboardObj.Type = "" //Making this an empty string will cause the yaml marsheler omit it

		stream.Dashboards[safeName] = &dashboardObj

		parent.Items["dashboards"] = append(parent.Items["dashboards"], safeName)

	case SavedSearchType:
		var searchObj savedSearch

		if err := mapstructure.Decode(obj, &searchObj); err != nil {
			return err
		}

		safeName := sanitizeName(searchObj.Name)

		stream.SavedSearches[safeName] = &searchObj

		parent.Items["savedSearches"] = append(parent.Items["savedSearches"], safeName)

	default:
		errMessage := fmt.Sprintf("Unknown child type: %s", childType)
		return errors.New(errMessage)
	}

	return nil
}
func (a *application) BasePath() string {
	return a.path
}

func (a *application) Import(pathToFileToImport string, appstream string) error {
	rootFolder := NewFolder()

	if err := a.LoadAppStreams(); err != nil {
		return fmt.Errorf("Could not load application app streams: %w", err)
	}

	stream, err := a.FindAppStream(appstream)
	if err != nil {
		return fmt.Errorf("Could not find app stream %s", appstream)
	}

	//Read the JSON file and load it into objects
	data, err := os.ReadFile(pathToFileToImport)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, rootFolder); err != nil {
		if jsonErr, ok := err.(*json.SyntaxError); ok {
			problemPart := data[jsonErr.Offset-10 : jsonErr.Offset+10]
			err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
			return err
		}
	}

	stream.RootFolder = rootFolder

	//Process the application's child objects (dashboards, folders, saved searches, etc.)
	//This function operates recursively. For each folder it finds in the children,
	//it calls itself to process the new folder. That's why the first
	//argument is a folder
	if err := processChildren(rootFolder, stream); err != nil {
		return err
	}

	//The root folder of the imported file is essentially the application. The root folder's
	//items (dashboards, folders, saved searches), name, and description need to moved to the
	//application object
	a.Name = rootFolder.Name
	a.Description = rootFolder.Description
	a.Items = rootFolder.Items

	if err := stream.WriteObjects(); err != nil {
		return err
	}

	return nil
}
func (a *application) NewAppStream(name string) *appStream {
	stream := NewAppStream(name, a)
	return stream
}

func (a *application) LoadAppStreams() error {
	upstream := a.NewAppStream("upstream")
	midstream := a.NewAppStream("midstream")
	downstream := a.NewAppStream("downstream")

	upstream.Child = midstream
	midstream.Parent = upstream
	midstream.Child = downstream
	downstream.Parent = midstream

	streams := []*appStream{upstream, midstream, downstream}
	for _, stream := range streams {
		stream.Application = a

		if err := stream.Load(); err != nil {
			return err
		}

		a.appStreams = append(a.appStreams, stream)
	}

	return nil
}

//func (a *application) Copy() *application {
//	newApp := NewApplicationWithPath(a.path)
//
//	newApp.Name = a.Name
//	newApp.Description = a.Description
//	newApp.Version = a.Version
//	newApp.Type = a.Type
//	newApp.Items = a.Items
//	newApp.Children = a.Children
//
//	newApp.panels = copyPanelList(a.panels)
//	newApp.dashboards = copyDashboardList(a.dashboards)
//	newApp.queries = a.queries             //TODO: implement copy function for queries
//	newApp.variables = a.variables         //TODO: implement copy function for variables
//	newApp.savedSearches = a.savedSearches //TODO: implement copy function for saved searches
//
//	//Folders are a slice of pointers so they
//	//need special handling to not point to the original
//	//application's folder objects
//	folderList := make(map[string]*folder)
//	for folderName, folderToCopy := range a.folders {
//		folderList[folderName] = folderToCopy.Copy()
//	}
//	newApp.folders = folderList
//
//	return newApp
//}

func (a *application) Diff(app *application) error {
	changelog, err := diff.Diff(a.Children, app.Children)
	if err != nil {
		return err
	}

	//TODO: Filter out diffing on list of available panels, dashboards,
	//saved searches, variables, and queries. Those aren't relevant as
	//they may not end up in the built app artifact. Only diff the objects
	//under items list
	for _, change := range changelog {
		fmt.Println("")
		fmt.Println("TYPE: ", change.Type)
		fmt.Println("PATH: ", change.Path)
		fmt.Println("FROM: ", change.From)
		fmt.Println("TO: ", change.To)
	}

	return nil
}

//func (a *application) mergeVariables(app *application) error {
//	for vname, variableObj := range app.variables {
//		//If a panel by the same name already exists, merge it
//		//Otherwise, add the new panel to the list of panels
//		if _, ok := a.variables[vname]; ok {
//			v := a.variables[vname]
//			if err := mergo.Merge(&v, variableObj, mergo.WithOverride); err != nil {
//				return err
//			}
//
//			a.variables[vname] = v
//		} else {
//			a.variables[vname] = variableObj
//		}
//	}
//
//	return nil
//}
//
//func (a *application) mergePanels(app *application) error {
//	for pname, panelObj := range app.panels {
//		//If a panel by the same name already exists, merge it
//		//Otherwise, add the new panel to the list of panels
//		if _, ok := a.panels[pname]; ok {
//			p := a.panels[pname]
//			if err := mergo.Merge(&p, panelObj, mergo.WithOverride); err != nil {
//				return err
//			}
//
//			a.panels[pname] = p
//		} else {
//			a.panels[pname] = panelObj
//		}
//	}
//
//	return nil
//}
//
//func (a *application) mergeDashboards(app *application) error {
//	for dname, dashboard := range app.dashboards {
//		//If a dashboard by the same name already exists, merge it
//		//Otherwise, add the new dashboard to the list of dashboards
//		if _, ok := a.dashboards[dname]; ok {
//			d := a.dashboards[dname]
//			if err := mergo.Merge(&d, dashboard, mergo.WithOverride); err != nil {
//				return err
//			}
//
//			a.dashboards[dname] = d
//		} else {
//			a.dashboards[dname] = dashboard
//		}
//	}
//
//	return nil
//}
//
//func (a *application) mergeFolders(app *application) error {
//	for fname, folder := range app.folders {
//		//If a folder by the same name already exists, merge it
//		//Otherwise, add the new folder to the list of folders
//		if _, ok := a.folders[fname]; ok {
//			f := a.folders[fname]
//			if err := mergo.Merge(&f, folder, mergo.WithOverride); err != nil {
//				return err
//			}
//
//			a.folders[fname] = f
//		} else {
//			a.folders[fname] = folder
//		}
//	}
//
//	return nil
//}
//
//func (a *application) mergeItems(app *application) error {
//	for iname, item := range app.Items {
//		//If an item by the same name already exists, merge it
//		//Otherwise, add the new item to the list of items
//		if _, ok := a.Items[iname]; ok {
//			i := a.Items[iname]
//			if err := mergo.Merge(&i, item, mergo.WithOverride); err != nil {
//				return err
//			}
//
//			a.Items[iname] = i
//		} else {
//			a.Items[iname] = item
//		}
//	}
//
//	return nil
//}
//
//func (a *application) mergeSavedSearches(app *application) error {
//	for sname, search := range app.savedSearches {
//		//If a search by the same name already exists, merge it
//		//Otherwise, add the new search to the list of saved searches
//		if _, ok := a.savedSearches[sname]; ok {
//			s := a.savedSearches[sname]
//			if err := mergo.Merge(&s, search, mergo.WithOverride); err != nil {
//				return err
//			}
//
//			a.savedSearches[sname] = s
//		} else {
//			a.savedSearches[sname] = search
//		}
//	}
//
//	return nil
//}
//
//func (a *application) Merge(app *application) error {
//	if err := a.mergePanels(app); err != nil {
//		return err
//	}
//
//	if err := a.mergeVariables(app); err != nil {
//		return err
//	}
//
//	if err := a.mergeDashboards(app); err != nil {
//		return err
//	}
//
//	if err := a.mergeFolders(app); err != nil {
//		return err
//	}
//
//	if err := a.mergeSavedSearches(app); err != nil {
//		return err
//	}
//
//	if err := a.mergeItems(app); err != nil {
//		return err
//	}
//
//	if app.Name != "" {
//		a.Name = app.Name
//	}
//
//	if app.Description != "" {
//		a.Description = app.Description
//	}
//
//	if app.Version != "" {
//		a.Version = app.Version
//	}
//
//	return nil
//}

func (a *application) FindAppStream(name string) (*appStream, error) {
	for _, stream := range a.appStreams {
		if stream.Name == name {
			return stream, nil
		}
	}

	err := fmt.Errorf("Could not find app stream '%s' in %s", name, a.path)
	return nil, err
}

func (a *application) ToJSON() ([]byte, error) {
	//Compile into JSON return
	jsonByteString, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return jsonByteString, nil
}
