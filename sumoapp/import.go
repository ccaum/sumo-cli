package sumoapp

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

func sanitizeName(name string) string {
	replacer := strings.NewReplacer("/", "_", ";", "-", " ", "-")
	sanitizedString := replacer.Replace(name)
	sanitizedString = strings.ToLower(sanitizedString)

	return sanitizedString
}

func processChildren(folder *folder, app *application) error {
	for _, childObj := range folder.Children {
		if err := processChild(childObj, folder, app); err != nil {
			return err
		}
	}

	//These fields aren't needed in the user facing code files
	folder.Children = nil
	folder.Type = ""

	return nil
}

func processChild(obj interface{}, parent *folder, app *application) error {
	childType := obj.(map[string]interface{})["type"]

	switch childType {
	case FolderType:
		folderObj := NewFolder()

		if err := mapstructure.Decode(obj, &folderObj); err != nil {
			return err
		}

		safeName := sanitizeName(folderObj.Name)
		app.folders[safeName] = folderObj

		parent.Items["folders"] = append(parent.Items["folders"], safeName)
		processChildren(folderObj, app)

	case DashboardType:
		var dashboardObj dashboard

		if err := mapstructure.Decode(obj, &dashboardObj); err != nil {
			return err
		}

		safeName := sanitizeName(dashboardObj.Name)

		//Add each of the panels to the application
		for _, panelObj := range dashboardObj.Panels {
			app.panels[panelObj.Key] = panelObj
		}

		//Add each of the variables to the application
		for _, variableObj := range dashboardObj.Variables {
			app.variables[variableObj.Name] = variableObj
			dashboardObj.IncludeVariables = append(dashboardObj.IncludeVariables, variableObj.Name)
		}

		//Remove these objects from the user facing code files
		dashboardObj.Panels = nil
		dashboardObj.Variables = nil
		dashboardObj.Type = "" //Making this an empty string will cause the yaml marsheler omit it

		app.dashboards[safeName] = dashboardObj

		parent.Items["dashboards"] = append(parent.Items["dashboards"], safeName)

	case SavedSearchType:
		var searchObj savedSearch

		if err := mapstructure.Decode(obj, &searchObj); err != nil {
			return err
		}

		safeName := sanitizeName(searchObj.Name)

		app.savedSearches[safeName] = searchObj

		parent.Items["savedSearches"] = append(parent.Items["savedSearches"], safeName)

	default:
		errMessage := fmt.Sprintf("Unknown child type: %s", childType)
		return errors.New(errMessage)
	}

	return nil
}

func writeObjects(app *application, appstream string) error {

	//Write the folder objects to the app stream
	for fName, folderObj := range app.folders {
		folderMap := make(map[string]*folder)
		folderMap[fName] = folderObj

		y, err := yaml.Marshal(folderMap)
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("%s/folders/%s.yaml", appstream, fName)
		if err := os.WriteFile(filePath, y, 0644); err != nil {
			return err
		}
	}

	//Write the dashboard objects to the app stream
	for dName, dashboardObj := range app.dashboards {
		dashMap := make(map[string]dashboard)
		dashMap[dName] = dashboardObj

		y, err := yaml.Marshal(dashMap)
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("%s/dashboards/%s.yaml", appstream, dName)
		if err := os.WriteFile(filePath, y, 0644); err != nil {
			return err
		}
	}

	//Write the panel objects to the app stream
	for pName, panelObj := range app.panels {
		panelMap := make(map[string]panel)
		panelMap[pName] = panelObj

		p, err := yaml.Marshal(panelMap)
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("%s/panels/%s.yaml", appstream, pName)
		if err := os.WriteFile(filePath, p, 0644); err != nil {
			return err
		}
	}

	//Write the variable objects to the app stream
	for vName, variableObj := range app.variables {
		variableMap := make(map[string]variable)
		variableMap[vName] = variableObj

		v, err := yaml.Marshal(variableMap)
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("%s/variables/%s.yaml", appstream, vName)
		if err := os.WriteFile(filePath, v, 0644); err != nil {
			return err
		}
	}

	//Write the saved search objects to the app stream
	for sName, searchObj := range app.savedSearches {
		searchMap := make(map[string]savedSearch)
		searchMap[sName] = searchObj

		v, err := yaml.Marshal(searchMap)
		if err != nil {
			return err
		}

		filePath := fmt.Sprintf("%s/saved-searches/%s.yaml", appstream, sName)
		if err := os.WriteFile(filePath, v, 0644); err != nil {
			return err
		}
	}

	//Write the application's definition to the init file in the stream
	a, err := yaml.Marshal(app)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/init.yaml", appstream)
	if err := os.WriteFile(filePath, a, 0644); err != nil {
		return err
	}

	return nil
}

func Import(path string, appstream string) error {
	rootFolder := NewFolder()
	app := NewApplication()

	data, err := os.ReadFile(path)
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

	//Process the application's child objects (dashboards, folders, saved searches, etc.)
	if err := processChildren(rootFolder, app); err != nil {
		return err
	}

	//The root folder of the imported file is essentially the application. The root folder's
	//items (dashboards, folders, saved searches), name, and description need to moved to the
	//application object
	app.Name = rootFolder.Name
	app.Description = rootFolder.Description
	app.Items = rootFolder.Items
	app.Type = ""

	if err := writeObjects(app, appstream); err != nil {
		return err
	}

	return nil
}
