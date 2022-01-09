package gocode

import (
	"sort"
	"strings"
)

type (
	PackageGraph struct {
		withExternalPackages bool
		relations            *Relations
		graph                map[PackagePath][]*PackageSummary
	}
)

func newPackageGraph(r *Relations, withExternalPackages bool) *PackageGraph {
	pg := &PackageGraph{
		graph:                make(map[PackagePath][]*PackageSummary),
		relations:            r,
		withExternalPackages: withExternalPackages,
	}
	pg.generate()

	return pg
}

func (pg *PackageGraph) generate() {
	for _, pkg := range pg.relations.Packages().AsSlice() {
		pg.makePackageSummaryMapIfNotExist(pkg)
		for _, im := range pkg.Detail().Imports() {
			pg.addIfTarget(pkg, im)
		}
	}
}

func (pg *PackageGraph) makePackageSummaryMapIfNotExist(pkg *Package) {
	if _, ok := pg.graph[pkg.Summary().Path()]; !ok {
		pg.graph[pkg.Summary().Path()] = make([]*PackageSummary, 0)
	}
}

func (pg *PackageGraph) addIfTarget(pkg *Package, im *Import) {
	if pg.isTargetGraph(im.pkgSummary) {
		pg.graph[pkg.Summary().Path()] = append(pg.graph[pkg.Summary().Path()], im.PackageSummary())
	}
}

func (pg *PackageGraph) isTargetGraph(ps *PackageSummary) bool {
	return pg.withExternalPackages || pg.relations.Packages().Contains(ps.Path())
}

func (pg *PackageGraph) WithExternalPackage() bool {
	return pg.withExternalPackages
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
