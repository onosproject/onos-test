package codegen

type ResourceOptions struct {
	Package    string
	ImportPath string
	FilePath   string
	Group      string
	Version    string
	Kind       string
	PluralKind string
	ListKind   string
	Singular   string
	Plural     string
}

func generateResource(options ResourceOptions) error {
	file, err := openFile(options.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return getTemplate("resource", "resource.tpl").Execute(file, options)
}
