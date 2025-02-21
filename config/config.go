package config

//GokuConfig goku根配置
type GokuConfig struct {
	Version string `json:"version"`
	Cluster string `json:"cluster"`
	//Port                int                        `json:"port"`
	DiscoverConfig      map[string]*DiscoverConfig `json:"discover,omitempty"`
	Balance             map[string]*BalanceConfig  `json:"balance,omitempty"`
	Plugins             *GatewayPluginConfig       `json:"plugins,omitempty"`
	APIS                []*APIContent              `json:"apis,omitempty"`
	Strategy            []*StrategyConfig          `json:"strategy,omitempty"`
	AnonymousStrategyID string                     `json:"anonymousStrategyID,omitempty"`
	AuthPlugin          map[string]string          `json:"authPlugin,omitempty"`

	Log       *LogConfig       `json:"log,omitempty"`
	AccessLog *AccessLogConfig `json:"access_log,omitempty"`
}

//AccessLogConfig access日志配置
type AccessLogConfig struct {
	Name   string   `json:"name"`
	Enable int      `json:"enable"`
	Dir    string   `json:"dir"`
	File   string   `json:"file"`
	Period string   `json:"period"`
	Expire int      `json:"expire"`
	Fields []string `json:"fields"`
}

//LogConfig log日志配置
type LogConfig struct {
	Name   string `json:"name"`
	Enable int    `json:"enable"`
	Dir    string `json:"dir"`
	File   string `json:"file"`
	Level  string `json:"level"`
	Period string `json:"period"`
	Expire int    `json:"expire"`
}

//GatewayPluginConfig 网关插件配置
type GatewayPluginConfig struct {
	BeforePlugins []*PluginConfig `json:"before"`
	GlobalPlugins []*PluginConfig `json:"global"`
}

//DiscoverConfig 服务发现配置
type DiscoverConfig struct {
	Name        string             `json:"name"`
	Driver      string             `json:"driver"`
	Config      string             `json:"config"`
	HealthCheck *HealthCheckConfig `json:"healthCheck"` // nil 表示不启用，非nil表示启用
}

//HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	IsHealthCheck bool   `json:"is_health_check"`
	URL           string `json:"url"`
	Second        int    `json:"second"`
	TimeOutMill   int    `json:"timeoutMill"`
	StatusCode    string `json:"statusCode"`
}

//BalanceConfig 负载配置
type BalanceConfig struct {
	Name         string `json:"name"`
	DiscoverName string `json:"discover"`
	Config       string `json:"config"` // appName(for discovery) or  address (for static)
}

//PluginConfig 插件配置
type PluginConfig struct {
	Name      string `json:"name"`
	IsStop    bool   `json:"stop"`
	Config    string `json:"config"`
	UpdateTag string `json:"updateTag"`
}

//APIContent api详情
type APIContent struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	OutPutEncoder  string           `json:"output"`
	RequestURL     string           `json:"requestUrl"`
	Methods        []string         `json:"methods"`
	TimeOutTotal   int              `json:"timeoutTotal"`
	AlertThreshold int              `json:"alert_threshold"`
	Steps          []*APIStepConfig `json:"steps"`

	StaticResponseStrategy string `json:"static_respone_strategy"`
	StaticResponse         string `json:"staticResponse"`
}

//APIStepConfig 链路配置
type APIStepConfig struct {
	Proto   string   `json:"proto"`
	Balance string   `json:"balance"`
	Method  string   `json:"method"` // follow | get | post | put ...
	Path    string   `json:"path"`
	Body    string   `json:"body"`
	Headers []string `json:"headers,omitempty"`
	Decode  string   `json:"decode"` // origin | json
	Encode  string   `json:"encode"` // origin | form | json

	Actions   []*ActionConfig `json:"actions"`
	BlackList []string        `json:"blackList"`
	WhiteList []string        `json:"whiteList"`

	Target  string `json:"target"`
	Group   string `json:"group"`
	Retry   int    `json:"retry"`
	TimeOut int    `json:"timeout"`
}

//APIStepUIConfig 链路UI配置
type APIStepUIConfig struct {
	Proto   string   `json:"proto"`
	Balance string   `json:"balance"`
	Method  string   `json:"method"` // follow | get | post | put ...
	Path    string   `json:"path"`
	Body    string   `json:"body"`
	Headers []string `json:"headers,omitempty"`
	Decode  string   `json:"decode"` // origin | json
	Encode  string   `json:"encode"` // origin | form | json

	BlackList []string `json:"blackList"`
	WhiteList []string `json:"whiteList"`

	Move    []MoveConfig   `json:"move"`
	Delete  []DeleteConfig `json:"delete"`
	Rename  []RenameConfig `json:"rename"`
	Target  string         `json:"target"`
	Group   string         `json:"group"`
	Retry   int            `json:"retry"`
	TimeOut int            `json:"timeout"`
}

//MoveConfig move配置
type MoveConfig struct {
	Origin string `json:"origin"`
	Target string `json:"target"`
}

//DeleteConfig delete配置
type DeleteConfig struct {
	Origin string `json:"origin"`
}

//RenameConfig rename配置
type RenameConfig struct {
	Origin string `json:"origin"`
	Target string `json:"target"`
}

//ActionConfig action配置
type ActionConfig struct {
	ActionType string `json:"type"`
	Original   string `json:"original"`
	Target     string `json:"target"`
}

//StrategyConfig 策略配置
type StrategyConfig struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Enable  bool              `json:"enable"`
	APIS    []*APIOfStrategy  `json:"apis"`
	AUTH    map[string]string `json:"auth"`
	Plugins []*PluginConfig   `json:"plugins"`
}

//APIOfStrategy 策略接口配置
type APIOfStrategy struct {
	ID      int             `json:"id"`
	Balance string          `json:"balance"` // 单step有效
	Plugins []*PluginConfig `json:"plugins"`
}

//VersionConfig 版本配置
type VersionConfig struct {
	VersionID     int    `json:"versionID"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	Remark        string `json:"remark"`
	CreateTime    string `json:"createTime"`
	PublishStatus int    `json:"publishStatus"`
	PublishTime   string `json:"publishTime"`
}

//Project 项目
type Project struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	CreateTime string `json:"createTime"`
}
