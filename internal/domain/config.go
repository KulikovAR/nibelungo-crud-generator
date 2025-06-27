package domain

type ProjectConfig struct {
	Name         string   `json:"name"`
	Module       string   `json:"module"`
	Entities     []Entity `json:"entities"`
	Repositories []string `json:"repositories"`
	Features     Features `json:"features"`
	Port         int      `json:"port,omitempty"`
}

type Entity struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Tags     []string `json:"tags,omitempty"`
	Required bool     `json:"required,omitempty"`
	Unique   bool     `json:"unique,omitempty"`
}

type Features struct {
	GRPC       bool `json:"grpc"`
	REST       bool `json:"rest"`
	Events     bool `json:"events"`
	Tests      bool `json:"tests"`
	Docker     bool `json:"docker"`
	Migrations bool `json:"migrations"`
	Swagger    bool `json:"swagger"`
}
