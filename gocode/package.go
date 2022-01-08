package gocode

import (
	"go/types"

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
)

func (pn PackageName) String() string {
	return string(pn)
}

func (pp PackagePath) String() string {
	return string(pp)
}

func newPackageSummaryFromPackages(pkg *packages.Package) *PackageSummary {
	return &PackageSummary{
		name: PackageName(pkg.Name),
		path: PackagePath(pkg.PkgPath),
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

func newPackageDetail(pkg *packages.Package) *PackageDetail {
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

func newPackage(pkg *packages.Package) *Package {
	return &Package{
		summary: newPackageSummaryFromPackages(pkg),
		detail:  newPackageDetail(pkg),
	}
}

func (p *Package) Summary() *PackageSummary {
	return p.summary
}

func (p *Package) Detail() *PackageDetail {
	return p.detail
}
