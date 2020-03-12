package codegen

import "path"

type ResourceClientOptions struct {
	Location Location
	Package  Package
	Types    ResourceClientTypes
	Names    ResourceClientNames
	Resource ResourceOptions
}

type ResourceClientTypes struct {
	Interface string
	Struct    string
}

type ResourceClientNames struct {
	Singular string
	Plural   string
}

func generateResourceClient(options ResourceClientOptions) error {
	if err := generateTemplate(getTemplate("resourceclient.tpl"), path.Join(options.Location.Path, options.Location.File), options); err != nil {
		return err
	}
	return generateResource(options.Resource)
}
