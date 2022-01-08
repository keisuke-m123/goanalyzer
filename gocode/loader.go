package gocode

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/afero"
	"golang.org/x/tools/go/packages"
)

type (
	// Relations は解析したgoコードの結果を保持する構造体。
	Relations struct {
		packages     *PackageMap
		structs      *PackageStructureMap
		interfaces   *PackageInterfaceMap
		typeAliases  *PackageTypeAliasMap
		definedTypes *PackageDefinedTypeMap
	}

	PackageGraph struct {
		graph map[PackagePath][]*PackageSummary
	}

	// LoadOptions はgoコード解析時のオプション。
	LoadOptions struct {
		FileSystem         afero.Fs
		Directories        []string
		IgnoredDirectories []string
		Recursive          bool
	}
)

func newRelations() *Relations {
	return &Relations{
		packages:     newPackageMap(),
		structs:      newPackageStructureMap(),
		interfaces:   newPackageInterfaceMap(),
		typeAliases:  newPackageTypeAliasMap(),
		definedTypes: newPackageDefinedTypeMap(),
	}
}

func LoadRelations(options *LoadOptions) (*Relations, error) {
	r := newRelations()
	if err := r.load(options); err != nil {
		return r, err
	}
	return r, nil
}

func (r *Relations) Packages() *PackageMap {
	return r.packages
}

func (r *Relations) Structs() *PackageStructureMap {
	return r.structs
}

func (r *Relations) Interfaces() *PackageInterfaceMap {
	return r.interfaces
}

func (r *Relations) TypeAliases() *PackageTypeAliasMap {
	return r.typeAliases
}

func (r *Relations) DefinedTypes() *PackageDefinedTypeMap {
	return r.definedTypes
}

func (r *Relations) load(options *LoadOptions) error {
	ignoreDirectoryMap := map[string]struct{}{}
	for _, dir := range options.IgnoredDirectories {
		ignoreDirectoryMap[dir] = struct{}{}
	}

	for _, directoryPath := range options.Directories {
		if options.Recursive {
			err := afero.Walk(options.FileSystem, directoryPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					if strings.HasPrefix(info.Name(), ".") || info.Name() == "vendor" {
						return filepath.SkipDir
					}
					if _, ok := ignoreDirectoryMap[path]; ok {
						return filepath.SkipDir
					}
					return r.parseDirectory(path)
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			err := r.parseDirectory(directoryPath)
			if err != nil {
				return err
			}
		}
	}

	r.registerRelations()

	return nil
}

func (r *Relations) parseDirectory(directoryPath string) error {
	loadConfig := &packages.Config{
		Mode: packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedSyntax |
			packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports,
		Dir: directoryPath,
	}
	pkgs, err := packages.Load(loadConfig)
	if err != nil {
		return fmt.Errorf("load packages failed: %w", err)
	}
	for i := range pkgs {
		p := newPackage(pkgs[i])
		r.packages.add(p)
		r.registerStructs(p)
		r.registerInterfaces(p)
		r.registerTypeAliases(p)
		r.registerDefinedTypes(p)
	}

	return nil
}

func (r *Relations) registerRelations() {
	structs := r.structs.StructAll()
	for si := range structs {
		interfaces := r.interfaces.InterfaceAll()
		for i := range interfaces {
			structs[si].addInterfaceIfImplements(interfaces[i])
		}
	}
}

func (r *Relations) registerStructs(pkg *Package) {
	structs := pkg.Detail().Structs()
	for i := range structs {
		r.structs.put(structs[i])
	}
}

func (r *Relations) registerInterfaces(pkg *Package) {
	interfaces := pkg.Detail().Interfaces()
	for i := range interfaces {
		r.interfaces.put(interfaces[i])
	}
}

func (r *Relations) registerTypeAliases(pkg *Package) {
	aliases := pkg.Detail().TypeAliases()
	for i := range aliases {
		r.typeAliases.put(aliases[i])
	}
}

func (r *Relations) registerDefinedTypes(pkg *Package) {
	definedTypes := pkg.Detail().DefinedTypes()
	for i := range definedTypes {
		r.definedTypes.put(definedTypes[i])
	}
}

func (r *Relations) GeneratePackageGraph() *PackageGraph {
	pg := newPackageGraph()
	for _, pkg := range r.Packages().AsSlice() {
		if pkg.Summary().Path() == "." {
			continue
		}
		if _, ok := pg.graph[pkg.Summary().Path()]; !ok {
			pg.graph[pkg.Summary().Path()] = make([]*PackageSummary, 0)
		}
		for _, im := range pkg.Detail().Imports() {
			if r.Packages().Contains(im.PackageSummary().Path()) {
				pg.graph[pkg.Summary().Path()] = append(pg.graph[pkg.Summary().Path()], im.PackageSummary())
			}
		}
	}
	return pg
}

func newPackageGraph() *PackageGraph {
	return &PackageGraph{
		graph: make(map[PackagePath][]*PackageSummary),
	}
}
func (pg *PackageGraph) SortedPackagePaths() []PackagePath {
	var paths []PackagePath
	for path := range pg.graph {
		paths = append(paths, path)
	}
	sort.Slice(paths, func(i, j int) bool {
		return strings.Compare(paths[i].String(), paths[j].String()) < 0
	})
	return paths
}

func (pg *PackageGraph) SortedImportPackagePaths(pkgPath PackagePath) []*PackageSummary {
	summaries := append([]*PackageSummary{}, pg.graph[pkgPath]...)
	sort.Slice(summaries, func(i, j int) bool {
		return strings.Compare(summaries[i].Path().String(), summaries[j].Path().String()) < 0
	})
	return summaries
}
