package projfs

import (
	"os/exec"
	"path/filepath"

	"forum/internal/common/system/fsys"
)

var staticDir = projectRoot.Join("website", "static").Must()

func StaticDir() fsys.Path {
	return staticDir
}

var templatesDir = projectRoot.Join("website", "templates").Must()

func TemplatesDir() fsys.Path {
	return templatesDir
}

func StorageDir() fsys.Path {
	return projectRoot.Join("storage")
}

var projectRoot fsys.Path = func() fsys.Path {
	output, err := exec.Command("go", "env", "GOWORK").Output()
	if err != nil {
		panic(err)
	}
	if string(output) == "" {
		panic("GOWORK not set")
	}
	return fsys.Path(filepath.Dir(string(output)))
}()
