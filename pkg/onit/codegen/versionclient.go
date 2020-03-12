package codegen

import "path"

type VersionOptions struct {
	Location  Location
	Package   Package
	Group     string
	Version   string
	Types     VersionTypes
	Names     VersionNames
	Resources map[string]ResourceClientOptions
}

type VersionTypes struct {
	Interface string
	Struct    string
}

type VersionNames struct {
	Natural string
	Proper  string
}

func generateVersionClient(options VersionOptions) error {
	if err := generateTemplate(getTemplate("versionclient.tpl"), path.Join(options.Location.Path, options.Location.File), options); err != nil {
		return err
	}

	for _, resource := range options.Resources {
		if err := generateResourceClient(resource); err != nil {
			return err
		}
	}
	return nil
}
