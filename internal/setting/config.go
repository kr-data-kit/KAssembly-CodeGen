package setting

type GlobalSetting struct {
	Path      string
	CreateDir bool
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
