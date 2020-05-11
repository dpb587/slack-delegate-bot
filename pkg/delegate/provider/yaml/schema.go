package yaml

type Schema struct {
	DelegateBot SchemaDelegateBot `yaml:"delegatebot"`
}

type SchemaDelegateBot struct {
	Watch    []interface{}                `yaml:"watch"`
	Delegate interface{}                  `yaml:"delegate"`
	Options  SchemaDelegateBotWithOptions `yaml:"options"`
}

type SchemaDelegateBotWithOptions struct {
	EmptyMessage string `yaml:"empty_message"`
}
