package sumoapp

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/SumoLogic-Incubator/sumologic-go-sdk/service/cip/types"
	"github.com/imdario/mergo"
)

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

func (f *folder) UploadWithOverwrite(a *APIClient, folderId string) error {
	return f.Upload(a, folderId, true)
}

func (f *folder) UploadWithoutOverwrite(a *APIClient, folderId string) error {
	return f.Upload(a, folderId, false)
}

func (f *folder) Upload(a *APIClient, folderId string, overwrite bool) error {
	var (
		localVarHttpMethod  = strings.ToUpper("Post")
		localVarPostBody    interface{}
		localVarFileName    string
		localVarFileBytes   []byte
		localVarReturnValue types.BeginAsyncJobResponse
	)

	// create path and map variables
	localVarPath := a.Cfg.BasePath + "/v2/content/folders/{folderId}/import"
	localVarPath = strings.Replace(localVarPath, "{"+"folderId"+"}", fmt.Sprintf("%v", folderId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if overwrite {
		localVarQueryParams.Add("overwrite", "true")
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}

	// body params
	localVarPostBody = f
	r, err := a.prepareRequest(localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return err
	}

	localVarHttpResponse, err := a.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"))
		if err == nil {
			return err
		}
	} else if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body:  localVarBody,
			error: localVarHttpResponse.Status,
		}
		if localVarHttpResponse.StatusCode == 200 {
			var v types.BeginAsyncJobResponse
			err = a.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return newErr
			}
			newErr.model = v
			return newErr
		} else if localVarHttpResponse.StatusCode >= 400 {
			var v types.ErrorResponse
			err = a.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return newErr
			}
			if v.Errors[0].Meta.Reason != "" {
				newErr.error = v.Errors[0].Message + ": " + v.Errors[0].Meta.Reason
			} else {
				newErr.error = v.Errors[0].Message
			}
			return newErr
		}
		return newErr
	}

	return nil
}
