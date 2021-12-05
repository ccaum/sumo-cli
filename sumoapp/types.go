package sumoapp

import "github.com/silas/dag"

type childType string

const (
	FolderType      string = "FolderSyncDefinition"
	DashboardType          = "DashboardV2SyncDefinition"
	SavedSearchType        = "SavedSearchWithScheduleSyncDefinition"
)

type timeBoundary struct {
	Type         string `json:"type,omitempty"`
	RelativeTime string `json:"relativeTime,omitempty"`
}

type timerange struct {
	Type string        `json:"type"`
	From *timeBoundary `json:"from"`
	To   *timeBoundary `json:"to,omitempty"`
}

type query struct {
	QueryString      string `json:"queryString"`
	QueryType        string `json:"queryType"`
	QueryKey         string `json:"queryKey"`
	MetricsQueryMode string `json:"metricsQueryMode,omitempty"`
	MetricsQueryData string `json:"metricsQueryData,omitempty"`
	TracesQueryData  string `json:"tracesQueryData,omitempty"`
	ParseMode        string `json:"parseMode,omitempty"`
	TimeSource       string `json:"timeSource,omitempty"`
}

type queryParameter struct{}

type layoutStructure struct {
	Key       string `json:"key"`
	Structure string `json:"structure"`
}

type search struct {
	QueryText        string        `json:"queryText"`
	DefaultTimeRange string        `json:"defaultTimeRange"`
	ByReceiptTime    bool          `json:"byReceiptTime"`
	ViewName         string        `json:"ViewName"`
	ViewStartTime    string        `json:"ViewStartTime"`
	QueryParameters  []interface{} `json:"queryParmaeters"`
	ParsingMode      string        `json:"parsingMode"`
}

type componentLibrary struct {
}

type application struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Version     string        `json:"version"`
	Children    []interface{} `json:"children" yaml:"children,omitempty"`
	Type        string        `json:"type" yaml:"type,omitempty"`
	Items       map[string][]string
	path        string
	appStreams  []*appStream
}

type appStream struct {
	Path          string
	Name          string
	Application   *application
	Parent        *appStream
	Child         *appStream
	Dashboards    map[string]*dashboard
	Panels        map[string]*panel
	SavedSearches map[string]*savedSearch
	Variables     map[string]*variable
	Queries       map[string]*query
	Folders       map[string]*folder
	RootFolder    *folder
	Graph         *dag.AcyclicGraph
}

type searchSchedule struct{}

type folder struct {
	Type          string        `json:"type" yaml:"type,omitempty"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Children      []interface{} `json:"children" yaml:"children,omitempty"`
	Items         map[string][]string
	folders       map[string]*folder
	dashboards    map[string]dashboard
	savedSearches map[string]savedSearch
}

type savedSearch struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Search      search `json:"search"`
	//SearchSchedule searchSchedule `json:"searchSchedule"`
}

type labelMap struct {
	Data map[string]string `json:"data,omitempty"`
}

type dashboard struct {
	Type             string      `json:"type" yaml:"type,omitempty"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	Title            string      `json:"title"`
	Theme            string      `json:"theme"`
	TopologyLabelMap labelMap    `json:"topologyLabelMap,omitempty", yaml:topologylabelmap,omitempty"`
	RefreshInterval  int64       `json:"refreshInterval"`
	TimeRange        *timerange  `json:"timeRange"`
	Layout           layout      `json:"layout"`
	Panels           []*panel    `json:"panels" yaml:"panels,omitempty"`
	Variables        []*variable `json:"variables" yaml:"variables,omitempty"`
	RootPanel        string      `json:"rootPanel,omitempty"`
	IncludeVariables []string
	key              string
}

type layout struct {
	LayoutType             string            `json:"layoutType"`
	LayoutStructures       []layoutStructure `json:"layoutStructures"`
	AppendLayoutStructures []layoutStructure `json:"appendLayoutStructures,omitempty" yaml:"appendLayoutStructures,omitempty"`
}

type panel struct {
	Id                                     string     `json:"id,omitempty"`
	Key                                    string     `json:"key"`
	Title                                  string     `json:"title"`
	VisualSettings                         string     `json:"visualSettings"`
	KeepVisualSettingsConsistentWithParent bool       `json:"keepVisualSettingsConsistentWithParent"`
	PanelType                              string     `json:"panelType"`
	Queries                                []query    `json:"queries"`
	Description                            string     `json:"descriptions"`
	TimeRange                              *timerange `json:"timeRange"`
	ColoringRules                          []string   `json:"coloringRules"`
	LinkedDashboards                       []string   `json:"linkedDashboards"`
	Text                                   string     `json:"text,omitempty"`
}

type sourceDefinition struct {
	VariableSourceType string `json:"variableSourceType"`
	Query              string `json:"query"`
	Field              string `json:"field"`
}

type variable struct {
	Id               string           `json:"id,omitempty"`
	Name             string           `json:"name"`
	DisplayName      string           `json:"displayName"`
	DefaultValue     string           `json:"defaultValue"`
	SourceDefinition sourceDefinition `json:"sourceDefinition"`
	AllowMultiSelect bool             `json:"allowMultiSelect"`
	IncludeAllOption bool             `json:"includeAllOption"`
	HideFromUI       bool             `json:"hideFromUI"`
	ValueType        string           `json:"valueType"`
}
