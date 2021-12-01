package sumoapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
)

func loadVariables(basePath string) (map[string]variable, error) {
	files, _ := ioutil.ReadDir(basePath)
	vList := make(map[string]variable)

	for _, file := range files {
		curList := make(map[string]variable)

		path := fmt.Sprintf("%s/%s", basePath, file.Name())
		extension := filepath.Ext(path)
		if extension != ".json" {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			err := fmt.Errorf("Unable to read file ", path, ": ", err)
			return nil, err
		}

		if err := json.Unmarshal(data, &curList); err != nil {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				problemPart := data[jsonErr.Offset-10 : jsonErr.Offset+10]
				err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
				return nil, err
			}
		}

		for variableName, v := range curList {
			v.Name = variableName
			vList[variableName] = v
		}
	}

	return vList, nil
}

func loadPanels(basePath string) (map[string]panel, error) {
	files, _ := ioutil.ReadDir(basePath)
	pList := make(map[string]panel)

	for _, file := range files {
		var curList map[string]panel

		path := fmt.Sprintf("%s/%s", basePath, file.Name())
		extension := filepath.Ext(path)
		if extension != ".json" {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			err := fmt.Errorf("Unable to read file ", path, ": ", err)
			return nil, err
		}

		if err := json.Unmarshal(data, &curList); err != nil {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				problemPart := data[jsonErr.Offset-10 : jsonErr.Offset+10]
				err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
				return nil, err
			}
		}

		for panelName, p := range curList {
			p.Key = panelName
			pList[panelName] = p
		}
	}

	return pList, nil
}

func loadRootFolder(appFilePath string) (folder, error) {
	var app folder

	data, err := os.ReadFile(appFilePath)
	if err != nil {
		//If the root folder file doesn't exist, that's OK. At least one stream
		//needs a root folder, but not all of them. If one isn't found, return an
		//empty root folder
		noSuchFileError := fmt.Sprintf("open %s: no such file or directory", appFilePath)
		if err.Error() == noSuchFileError {
			return folder{}, nil
		}

		msg := fmt.Errorf("Unable to read file %s: %s", appFilePath, err.Error())
		//This means we found the root folder file, but can't read it. Since
		//that's a legitimate error, return the error object
		return folder{}, msg
	}

	if err := json.Unmarshal(data, &app); err != nil {
		if jsonErr, ok := err.(*json.SyntaxError); ok {
			problemPart := data[jsonErr.Offset-10 : jsonErr.Offset+10]
			err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
			return folder{}, err
		}
	}

	app.Type = FolderType

	return app, nil
}

func loadFolders(basePath string) (map[string]folder, error) {
	var folders map[string]folder

	ffiles, _ := ioutil.ReadDir(basePath)
	for _, file := range ffiles {
		var curList map[string]folder
		path := fmt.Sprintf("%s/%s", basePath, file.Name())
		extension := filepath.Ext(path)
		if extension != ".json" {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			err := fmt.Errorf("Unable to read file ", path, ": ", err)
			return nil, err
		}

		if err := json.Unmarshal(data, &curList); err != nil {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				problemPart := data[jsonErr.Offset-10 : jsonErr.Offset+10]
				err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
				return nil, err
			}
		}

		if err := mergo.Merge(&folders, curList); err != nil {
			return nil, err
		}
	}

	for name, folder := range folders {
		folder.Name = name
		folder.Type = FolderType
	}

	return folders, nil
}

func loadDashboards(basePath string) (map[string]dashboard, error) {
	var dashboards map[string]dashboard

	dfiles, _ := ioutil.ReadDir(basePath)
	for _, file := range dfiles {
		var curList map[string]dashboard

		path := fmt.Sprintf("%s/%s", basePath, file.Name())
		extension := filepath.Ext(path)
		if extension != ".json" {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			err := fmt.Errorf("Unable to read file ", path, ": ", err)
			return nil, err
		}

		if err := json.Unmarshal(data, &curList); err != nil {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				problemPart := data[jsonErr.Offset-10 : jsonErr.Offset+10]
				err = fmt.Errorf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
				return nil, err
			}
		}

		if err := mergo.Merge(&dashboards, curList); err != nil {
			return nil, err
		}
	}

	for name, dashboard := range dashboards {
		dashboard.Name = name
	}

	return dashboards, nil
}

func loadAppStreams(basePath string) ([]appStream, error) {
	var appStreams []appStream

	//The order here is important! Always start with the highest
	//order stream and move to the lowest
	streams := [3]string{"upstream", "midstream", "downstream"}

	for _, streamName := range streams {
		stream := appStream{
			Name:        streamName,
			Path:        fmt.Sprintf("%s/%s", basePath, streamName),
			Application: InitApplication(),
		}

		panelBasePath := fmt.Sprintf("%s/panels", stream.Path)
		dashboardBasePath := fmt.Sprintf("%s/dashboards", stream.Path)
		folderBasePath := fmt.Sprintf("%s/folders", stream.Path)
		variableBasePath := fmt.Sprintf("%s/variables", stream.Path)

		variables, err := loadVariables(variableBasePath)
		if err != nil {
			err := fmt.Errorf("Could not load variables at %s: %w", variableBasePath, err)
			return nil, err
		}

		panels, err := loadPanels(panelBasePath)
		if err != nil {
			err := fmt.Errorf("Could not load panels at %s: %w", panelBasePath, err)
			return nil, err
		}

		dashboards, err := loadDashboards(dashboardBasePath)
		if err != nil {
			err := fmt.Errorf("Could not load dashboards at %s: %w", dashboardBasePath, err)
			return nil, err
		}

		folders, err := loadFolders(folderBasePath)
		if err != nil {
			err := fmt.Errorf("Could not load folders at %s: %w", folderBasePath, err)
			return nil, err
		}

		rootPath := fmt.Sprintf("%s/init.json", stream.Path)
		rootFolder, err := loadRootFolder(rootPath)
		if err != nil {
			err := fmt.Errorf("Could not load root application at %s: %w", rootPath, err)
			return nil, err
		}

		stream.Application.panels = panels
		stream.Application.dashboards = dashboards
		stream.Application.variables = variables
		stream.Application.folders = folders
		stream.Application.Name = rootFolder.Name
		stream.Application.Description = rootFolder.Description
		stream.Application.items = rootFolder.Items

		appStreams = append(appStreams, stream)
	}

	return appStreams, nil
}

func CompileApp(basePath string) ([]byte, error) {
	app := InitApplication()

	//Load each app stream and all objects within the stream's app
	streams, err := loadAppStreams(basePath)
	if err != nil {
		return nil, err
	}

	//For each app stream, merge with the previous one. The order here is important!
	for _, stream := range streams {
		app.Merge(stream.Application)
	}

	if err := app.Build(); err != nil {
		return nil, err
	}

	//Compile into JSON return
	jsonByteString, err := json.Marshal(app)
	if err != nil {
		return nil, err
	}

	return jsonByteString, nil
}
