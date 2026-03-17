package fsys

import (
	"os"
	"path/filepath"
)

type Path string

func (p Path) Join(pathElems ...string) Path {
	return Path(filepath.Join(
		append([]string{string(p)}, pathElems...)...))
}

func (p Path) Must() Path {
	_, err := os.Stat(string(p))
	if err != nil {
		panic(err)
	}
	return p
}
