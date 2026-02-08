package service

import (
	_ "embed"

	"github.com/goccy/go-yaml"
)

//go:embed metadata/category-tree.yaml
var CategoryTreeYAML []byte

type Node struct {
	Key      string `yaml:"key"`
	NameKo   string `yaml:"name_ko"`
	NameEn   string `yaml:"name_en"`
	Children []Node `yaml:"children"`
	Leaf     bool   `yaml:"leaf"`
}

type CategoryTree struct {
	// TODO : consider case that not use int for versioning
	Version int              `yaml:"version"`
	Nodes   []Node           `yaml:"nodes"`
	ByKey   map[string]*Node `yaml:"-"`
}

var CategoryTreeData CategoryTree

func init() {
	err := yaml.Unmarshal(CategoryTreeYAML, &CategoryTreeData)
	if err != nil {
		return
	}
	// TODO : build ByKey map
}

// TODO : 근데 이걸 어떻게 분류 태그로 만들지? -> 나중에 고민
