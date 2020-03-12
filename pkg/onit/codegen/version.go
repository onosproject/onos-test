package codegen

type VersionOptions struct {
	Package    string
	ImportPath string
	FilePath   string
	Group      string
	Version    string
	Resources  map[string]ResourceOptions
}

func generateVersion(options VersionOptions) error {
	for _, resource := range options.Resources {
		if err := generateResource(resource); err != nil {
			return err
		}
	}

	file, err := openFile(options.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return getTemplate("version", "version.tpl").Execute(file, options)
}
