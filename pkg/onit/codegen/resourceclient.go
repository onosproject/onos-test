package codegen

import "path"

type ResourceClientOptions struct {
	Location Location
	Package  Package
	Types    ResourceClientTypes
	Resource ResourceOptions
}

type ResourceClientTypes struct {
	Interface string
	Struct    string
}

func generateResourceClient(options ResourceClientOptions) error {
	if err := generateTemplate(getTemplate("resourceclient.tpl"), path.Join(options.Location.Path, options.Location.File), options); err != nil {
		return err
	}
	return generateResource(options.Resource)
}
