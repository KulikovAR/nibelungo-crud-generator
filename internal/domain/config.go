package domain

type ProjectConfig struct {
	Name         string   `json:"name"`
	Module       string   `json:"module"`
	Entities     []Entity `json:"entities"`
	Repositories []string `json:"repositories"`
	Features     Features `json:"features"`
}

type Entity struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name string   `json:"name"`
	Type string   `json:"type"`
	Tags []string `json:"tags"`
}

type Features struct {
	GRPC   bool `json:"grpc"`
	REST   bool `json:"rest"`
	Events bool `json:"events"`
}
