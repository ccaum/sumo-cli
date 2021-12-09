package sumoapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
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

func processChildren(folder *folder, overlay *appOverlay) error {
	for _, childObj := range folder.Children {
		if err := processChild(childObj, folder, overlay); err != nil {
			return err
		}
	}

	//These fields aren't needed in the user facing code files
	folder.Children = nil
	folder.Type = ""

	return nil
}

func processChild(obj interface{}, parent *folder, overlay *appOverlay) error {
	childType := obj.(map[string]interface{})["type"]

	switch childType {
	case FolderType:
		folderObj := NewFolder()

		if err := mapstructure.Decode(obj, &folderObj); err != nil {
			return err
		}

		safeName := sanitizeName(folderObj.Name)
		overlay.Folders[safeName] = folderObj

		parent.Items["folders"] = append(parent.Items["folders"], safeName)
		processChildren(folderObj, overlay)

	case DashboardType:
		var dashboardObj dashboard

		if err := mapstructure.Decode(obj, &dashboardObj); err != nil {
			return err
		}

		safeName := sanitizeName(dashboardObj.Name)

		//Add each of the panels to the application
		for _, panelObj := range dashboardObj.Panels {
			overlay.Panels[panelObj.Key] = panelObj
		}

		//Add each of the variables to the application
		for _, variableObj := range dashboardObj.Variables {
			overlay.Variables[variableObj.Name] = variableObj
			dashboardObj.IncludeVariables = append(dashboardObj.IncludeVariables, variableObj.Name)
		}

		//Remove these objects from the user facing code files
		dashboardObj.Panels = nil
		dashboardObj.Variables = nil
		dashboardObj.Type = "" //Making this an empty string will cause the yaml marsheler omit it

		overlay.Dashboards[safeName] = &dashboardObj

		parent.Items["dashboards"] = append(parent.Items["dashboards"], safeName)

	case SavedSearchType:
		var searchObj savedSearch

		if err := mapstructure.Decode(obj, &searchObj); err != nil {
			return err
		}

		safeName := sanitizeName(searchObj.Name)

		overlay.SavedSearches[safeName] = &searchObj

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

func (a *application) Import(pathToFileToImport string, appoverlay string) error {
	return a.ImportWithWriteOption(pathToFileToImport, appoverlay, true)
}

func (a *application) ImportWithWriteOption(pathToFileToImport string, appoverlay string, writeObjects bool) error {
	rootFolder := NewFolder()

	if err := a.LoadAppOverlays(); err != nil {
		return fmt.Errorf("Could not load application app overlays: %w", err)
	}

	overlay, err := a.FindAppOverlay(appoverlay)
	if err != nil {
		return fmt.Errorf("Could not find app overlay %s", appoverlay)
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

	overlay.RootFolder = rootFolder

	//Process the application's child objects (dashboards, folders, saved searches, etc.)
	//This function operates recursively. For each folder it finds in the children,
	//it calls itself to process the new folder. That's why the first
	//argument is a folder
	if err := processChildren(rootFolder, overlay); err != nil {
		return err
	}

	//The root folder of the imported file is essentially the application. The root folder's
	//items (dashboards, folders, saved searches), name, and description need to moved to the
	//application object
	a.Name = rootFolder.Name
	a.Description = rootFolder.Description
	a.Items = rootFolder.Items

	if err := overlay.WriteObjects(); err != nil {
		return err
	}

	return nil
}

func (a *application) NewAppOverlay(name string) *appOverlay {
	overlay := NewAppOverlay(name, a)
	return overlay
}

func (a *application) LoadAppOverlays() error {
	baseOverlay := a.NewAppOverlay("base")
	midOverlay := a.NewAppOverlay("middle")
	finalOverlay := a.NewAppOverlay("final")

	baseOverlay.Child = midOverlay
	midOverlay.Parent = baseOverlay
	midOverlay.Child = finalOverlay
	finalOverlay.Parent = midOverlay

	overlays := []*appOverlay{baseOverlay, midOverlay, finalOverlay}
	for _, overlay := range overlays {
		overlay.Application = a

		if err := overlay.Load(); err != nil {
			return err
		}

		a.appOverlays = append(a.appOverlays, overlay)
	}

	return nil
}

func (a *application) FindAppOverlay(name string) (*appOverlay, error) {
	for _, overlay := range a.appOverlays {
		if overlay.Name == name {
			return overlay, nil
		}
	}

	err := fmt.Errorf("Could not find app overlay '%s' in %s", name, a.path)
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
