package codegen

import "path"

type ResourceOptions struct {
	Location Location
	Package  Package
	Kind     ResourceKind
	Types    ResourceTypes
	Names    ResourceNames
}

type ResourceKind struct {
	Package  Package
	Group    string
	Version  string
	Kind     string
	ListKind string
}

type ResourceTypes struct {
	Kind     string
	Resource string
	Struct   string
}

type ResourceNames struct {
	Singular string
	Plural   string
}

func generateResource(options ResourceOptions) error {
	return generateTemplate(getTemplate("resource.tpl"), path.Join(options.Location.Path, options.Location.File), options)
}
