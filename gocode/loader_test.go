package gocode_test

import (
	"testing"

	"github.com/keisuke-m123/goanalyzer/gocode"
	"github.com/spf13/afero"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLoadRelations(t *testing.T) {
	tests := []struct {
		name            string
		directories     []string
		numPackages     int
		numStructs      int
		numInterfaces   int
		numDefinedTypes int
		numTypeAliases  int
	}{
		{
			name:            "testingsupport-recursive",
			directories:     []string{"./testdata/"},
			numPackages:     1,
			numStructs:      3,
			numInterfaces:   2,
			numDefinedTypes: 1,
			numTypeAliases:  1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := gocode.LoadRelations(&gocode.LoadOptions{
				FileSystem:  afero.NewOsFs(),
				Directories: test.directories,
				Recursive:   true,
			})
			if err != nil {
				t.Fatalf("failed to load relations: %s", err)
			}

			if r.Packages().NumPackages() != test.numPackages {
				t.Errorf("failed to load packages: %d", r.Packages().NumPackages())
			}

			if len(r.Structs().StructAll()) != test.numStructs {
				t.Errorf("failed to load structs: %d", len(r.Structs().StructAll()))
			}

			if len(r.Interfaces().InterfaceAll()) != test.numInterfaces {
				t.Errorf("failed to load interfaces: %d", len(r.Interfaces().InterfaceAll()))
			}

			if len(r.DefinedTypes().DefinedTypeAll()) != test.numDefinedTypes {
				t.Errorf("failed to load defined types: %d", len(r.DefinedTypes().DefinedTypeAll()))
			}

			if len(r.TypeAliases().AliasAll()) != test.numTypeAliases {
				t.Errorf("failed to load type aliases: %d", len(r.TypeAliases().AliasAll()))
			}
		})
	}
}

func TestLoadRelationsFromAnalysis(t *testing.T) {
	tests := []struct {
		name            string
		numPackages     int
		numStructs      int
		numInterfaces   int
		numDefinedTypes int
		numTypeAliases  int
	}{
		{
			name:            "testingsupport",
			numPackages:     1,
			numStructs:      3,
			numInterfaces:   2,
			numDefinedTypes: 1,
			numTypeAliases:  1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			analyzer := &analysis.Analyzer{Name: "test"}
			analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
				r := gocode.LoadRelationsFromAnalysis(pass)

				if r.Packages().NumPackages() != test.numPackages {
					t.Errorf("failed to load packages: %d", r.Packages().NumPackages())
				}

				if len(r.Structs().StructAll()) != test.numStructs {
					t.Errorf("failed to load structs: %d", len(r.Structs().StructAll()))
				}

				if len(r.Interfaces().InterfaceAll()) != test.numInterfaces {
					t.Errorf("failed to load interfaces: %d", len(r.Interfaces().InterfaceAll()))
				}

				if len(r.DefinedTypes().DefinedTypeAll()) != test.numDefinedTypes {
					t.Errorf("failed to load defined types: %d", len(r.DefinedTypes().DefinedTypeAll()))
				}

				if len(r.TypeAliases().AliasAll()) != test.numTypeAliases {
					t.Errorf("failed to load type aliases: %d", len(r.TypeAliases().AliasAll()))
				}

				return nil, nil
			}
			analysistest.Run(t, analysistest.TestData(), analyzer)
		})
	}
}
