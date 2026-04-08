package setting

type GlobalSetting struct {
	Path           string
	CreateDir      bool
	EndpointFilter Filter
}

type Filter struct {
	Include []string
	Exclude []string
}

type GoSetting struct {
	GlobalSetting
	PackageName string
	IsMod       bool
}

type PythonSetting struct {
	GlobalSetting
	PackageName string
}
