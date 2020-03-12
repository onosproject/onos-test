package codegen

type ClientOptions struct {
	Package    string
	ImportPath string
	FilePath   string
	Groups     map[string]GroupOptions
}

func generateClient(options ClientOptions) error {
	for _, group := range options.Groups {
		if err := generateGroup(group); err != nil {
			return err
		}
	}

	file, err := openFile(options.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return getTemplate("client", "client.tpl").Execute(file, options)
}
