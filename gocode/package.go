package gocode

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

type (
	// PackageName はパッケージ名を表す。
	PackageName string

	// PackagePath はパッケージへのパスを表す。
	PackagePath string

	// PackageSummary はパッケージのサマリ。
	PackageSummary struct {
		name PackageName
		path PackagePath
	}

	// PackageDetail はパッケージの詳細情報。
	PackageDetail struct {
		// imports はパッケージのインポート情報の一覧。
		imports *ImportList
		// structs はパッケージ内の struct の一覧。
		structs *StructList
		// interfaces はパッケージ内の interface の一覧。
		interfaces *InterfaceList
		// typeAliases はパッケージ内の type alias の一覧。
		typeAliases *TypeAliasList
		// definedTypes はパッケージ内の defined type の一覧。
		definedTypes *DefinedTypeList
	}

	// Package はパッケージ情報を表す。
	Package struct {
		// summary はパッケージのサマリ。
		summary *PackageSummary
		// detail はパッケージ内の詳細情報。
		detail *PackageDetail
	}

	packageIn interface {
		PkgPath() string
		PkgName() string
		Import() []*types.Package
		Defs() []types.Object
		Scope() *types.Scope
	}

	packageInPackagesPackage struct {
		pkg *packages.Package
	}

	packageInAnalysisPass struct {
		pass *analysis.Pass
	}
)

func (pn PackageName) String() string {
	return string(pn)
}

func (pp PackagePath) String() string {
	return string(pp)
}

func newPackageSummary(pkg packageIn) *PackageSummary {
	return &PackageSummary{
		name: PackageName(pkg.PkgName()),
		path: PackagePath(pkg.PkgPath()),
	}
}

func newPackageSummaryFromGoTypes(pkg *types.Package) *PackageSummary {
	return &PackageSummary{
		name: PackageName(pkg.Name()),
		path: PackagePath(pkg.Path()),
	}
}

func (p *PackageSummary) Name() PackageName {
	return p.name
}

func (p *PackageSummary) Path() PackagePath {
	return p.path
}

func (p *PackageSummary) Equal(other *PackageSummary) bool {
	return p.Path() == other.Path()
}

func newPackageDetail(pkg packageIn) *PackageDetail {
	return &PackageDetail{
		imports:      newImportList(pkg),
		structs:      newStructList(pkg),
		interfaces:   newInterfaceList(pkg),
		typeAliases:  newAliasList(pkg),
		definedTypes: newDefinedList(pkg),
	}
}

func (pd *PackageDetail) Imports() []*Import {
	return pd.imports.asSlice()
}

func (pd *PackageDetail) Structs() []*Struct {
	return pd.structs.asSlice()
}

func (pd *PackageDetail) Interfaces() []*Interface {
	return pd.interfaces.asSlice()
}

func (pd *PackageDetail) TypeAliases() []*TypeAlias {
	return pd.typeAliases.asSlice()
}

func (pd *PackageDetail) DefinedTypes() []*DefinedType {
	return pd.definedTypes.asSlice()
}

func newPackage(pkg packageIn) *Package {
	return &Package{
		summary: newPackageSummary(pkg),
		detail:  newPackageDetail(pkg),
	}
}

func newPackageFromPackages(pkg *packages.Package) *Package {
	return newPackage(newPackageInPackages(pkg))
}

func newPackageFromAnalysis(pass *analysis.Pass) *Package {
	return newPackage(newPackageInAnalysis(pass))
}

func (p *Package) Summary() *PackageSummary {
	return p.summary
}

func (p *Package) Detail() *PackageDetail {
	return p.detail
}

func newPackageInPackages(pkg *packages.Package) packageIn {
	return &packageInPackagesPackage{
		pkg: pkg,
	}
}

func (p *packageInPackagesPackage) PkgPath() string {
	return p.pkg.PkgPath
}

func (p *packageInPackagesPackage) PkgName() string {
	return p.pkg.Name
}

func (p *packageInPackagesPackage) Import() []*types.Package {
	var imports []*types.Package
	for _, pkg := range p.pkg.Imports {
		imports = append(imports, pkg.Types)
	}
	return imports
}

func (p *packageInPackagesPackage) Defs() []types.Object {
	var defs []types.Object
	for _, d := range p.pkg.TypesInfo.Defs {
		defs = append(defs, d)
	}
	return defs
}

func (p *packageInPackagesPackage) Scope() *types.Scope {
	return p.pkg.Types.Scope()
}

func newPackageInAnalysis(pass *analysis.Pass) packageIn {
	return &packageInAnalysisPass{
		pass: pass,
	}
}

func (p *packageInAnalysisPass) PkgPath() string {
	return p.pass.Pkg.Path()
}

func (p *packageInAnalysisPass) PkgName() string {
	return p.pass.Pkg.Name()
}

func (p *packageInAnalysisPass) Import() []*types.Package {
	return p.pass.Pkg.Imports()
}

func (p *packageInAnalysisPass) Defs() []types.Object {
	var defs []types.Object
	for _, d := range p.pass.TypesInfo.Defs {
		defs = append(defs, d)
	}
	return defs
}

func (p *packageInAnalysisPass) Scope() *types.Scope {
	return p.pass.Pkg.Scope()
}
