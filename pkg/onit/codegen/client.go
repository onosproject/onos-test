package codegen

import "path"

type ClientOptions struct {
	Location Location
	Package  Package
	Groups   map[string]GroupOptions
}

func generateClient(options ClientOptions) error {
	if err := generateTemplate(getTemplate("client.tpl"), path.Join(options.Location.Path, options.Location.File), options); err != nil {
		return err
	}

	for _, group := range options.Groups {
		if err := generateGroupClient(group); err != nil {
			return err
		}
	}
	return nil
}
