package sumoapp

type childType string

const (
	FolderType      string = "FolderSyncDefinition"
	DashboardType          = "DashboardV2SyncDefinition"
	SavedSearchType        = "SavedSearchWithScheduleSyncDefinition"
)

type timeBoundary struct {
	Type         string `yaml:"type,omitempty"`
	RelativeTime string `yaml:"relativeTime,omitempty"`
}

type timerange struct {
	Type string        `yaml:"type"`
	From *timeBoundary `yaml:"from"`
	To   *timeBoundary `yaml:"to,omitempty"`
}

type query struct {
	QueryString      string `yaml:"queryString"`
	QueryType        string `yaml:"queryType"`
	QueryKey         string `yaml:"queryKey"`
	MetricsQueryMode string `yaml:"metricsQueryMode,omitempty"`
	MetricsQueryData string `yaml:"metricsQueryData,omitempty"`
	TracesQueryData  string `yaml:"tracesQueryData,omitempty"`
	ParseMode        string `yaml:"parseMode,omitempty"`
	TimeSource       string `yaml:"timeSource,omitempty"`
}

type queryParameter struct{}

type layoutStructure struct {
	Key       string `json:"key"`
	Structure string `json:"structure"`
}

type search struct {
	QueryText        string        `yaml:"queryText"`
	DefaultTimeRange string        `yaml:"defaultTimeRange"`
	ByReceiptTime    bool          `yaml:"byReceiptTime"`
	ViewName         string        `yaml:"ViewName"`
	ViewStartTime    string        `yaml:"ViewStartTime"`
	QueryParameters  []interface{} `yaml:"queryParmaeters"`
	ParsingMode      string        `yaml:"parsingMode"`
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
	Type        string `yaml:"type"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Search      search `yaml:"search"`
	//SearchSchedule searchSchedule `json:"searchSchedule"` TODO: Add this back in. The searchSchedule type needs to be defined first
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
	LayoutType             string            `yaml:"layoutType"`
	LayoutStructures       []layoutStructure `yaml:"layoutStructures"`
	AppendLayoutStructures []layoutStructure `yaml:"appendLayoutStructures,omitempty" yaml:"appendLayoutStructures,omitempty"`
}

type panel struct {
	Id                                     string     `yaml:"id,omitempty"`
	Key                                    string     `yaml:"key"`
	Title                                  string     `yaml:"title"`
	VisualSettings                         string     `yaml:"visualSettings"`
	KeepVisualSettingsConsistentWithParent bool       `yaml:"keepVisualSettingsConsistentWithParent"`
	PanelType                              string     `yaml:"panelType"`
	Queries                                []query    `yaml:"queries"`
	Description                            string     `yaml:"descriptions"`
	TimeRange                              *timerange `yaml:"timeRange"`
	ColoringRules                          []string   `yaml:"coloringRules"`
	LinkedDashboards                       []string   `yaml:"linkedDashboards"`
	Text                                   string     `yaml:"text,omitempty"`
}

type sourceDefinition struct {
	VariableSourceType string `yaml:"variableSourceType"`
	Query              string `yaml:"query"`
	Field              string `yaml:"field"`
}

type variable struct {
	Id               string           `yaml:"id,omitempty"`
	Name             string           `yaml:"name"`
	DisplayName      string           `yaml:"displayName"`
	DefaultValue     string           `yaml:"defaultValue"`
	SourceDefinition sourceDefinition `yaml:"sourceDefinition"`
	AllowMultiSelect bool             `yaml:"allowMultiSelect"`
	IncludeAllOption bool             `yaml:"includeAllOption"`
	HideFromUI       bool             `yaml:"hideFromUI"`
	ValueType        string           `yaml:"valueType"`
}
