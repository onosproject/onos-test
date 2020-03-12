package codegen

type GroupOptions struct {
	Package    string
	ImportPath string
	FilePath   string
	Group      string
	Versions   map[string]VersionOptions
}

func generateGroup(options GroupOptions) error {
	for _, version := range options.Versions {
		if err := generateVersion(version); err != nil {
			return err
		}
	}

	file, err := openFile(options.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return getTemplate("group", "group.tpl").Execute(file, options)
}
