package datasources

import (
	"io/fs"
	"os"
	"wintercloak/internal"
)

func ListDataSources() ([]fs.DirEntry, error) {
	return os.ReadDir(internal.DataSourcesDir)
}
