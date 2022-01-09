package gocode_test

import (
	"testing"

	"github.com/keisuke-m123/goanalyzer/gocode"
	"github.com/spf13/afero"
)

var testingSupportPackages *gocode.Relations

func TestMain(m *testing.M) {
	if testingSupportPackages == nil {
		r, err := gocode.LoadRelations(&gocode.LoadOptions{
			FileSystem:  afero.NewOsFs(),
			Directories: []string{"./testdata/"},
			Recursive:   true,
		})
		if err != nil {
			panic(err)
		}
		testingSupportPackages = r
	}

	m.Run()
}
