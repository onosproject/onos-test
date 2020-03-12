package codegen

import "path"

type GroupOptions struct {
	Location Location
	Package  Package
	Group    string
	Types    GroupTypes
	Names    GroupNames
	Versions map[string]VersionOptions
}

type GroupTypes struct {
	Interface string
	Struct    string
}

type GroupNames struct {
	Natural string
	Proper  string
}

func generateGroupClient(options GroupOptions) error {
	if err := generateTemplate(getTemplate("groupclient.tpl"), path.Join(options.Location.Path, options.Location.File), options); err != nil {
		return err
	}

	for _, version := range options.Versions {
		if err := generateVersionClient(version); err != nil {
			return err
		}
	}
	return nil
}
