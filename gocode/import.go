package gocode

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

type (
	// ImportAlias は、import の alias を表す。
	ImportAlias string

	// Import は、import の情報を表す。
	Import struct {
		// alias は、import のエイリアス。
		alias ImportAlias
		// pkgSummary は、import に対応するパッケージ情報。
		pkgSummary *PackageSummary
	}

	// ImportList は、import のリストを表す。
	// import されたパッケージの名前( PackageName ) をキーとして Import 情報を格納する。
	ImportList struct {
		imports map[PackageName]*Import
	}
)

func (ia ImportAlias) String() string {
	return string(ia)
}

func (i Import) AliasName() ImportAlias {
	return i.alias
}

// HasAliasName は、import にエイリアスが設定されているかを返す。
func (i Import) HasAliasName() bool {
	return i.alias != ""
}

func (i Import) PackageSummary() *PackageSummary {
	return i.pkgSummary
}

// packages.Package からImport情報を抽出して生成する。
func newImportList(pkg *packages.Package) *ImportList {
	imports := make(map[PackageName]*Import)
	// general imports
	for _, importedPkg := range pkg.Imports {
		pkgSummary := newPackageSummaryFromPackages(importedPkg)
		imports[pkgSummary.Name()] = &Import{
			alias:      "",
			pkgSummary: pkgSummary,
		}
	}
	// alias Imports
	for _, d := range pkg.TypesInfo.Defs {
		pkgName, ok := d.(*types.PkgName)
		if ok {
			pkgSummary := newPackageSummaryFromGoTypes(pkgName.Imported())
			imports[pkgSummary.Name()] = &Import{
				alias:      ImportAlias(pkgName.Name()),
				pkgSummary: pkgSummary,
			}
		}
	}

	return &ImportList{imports: imports}
}

func (il ImportList) Len() int {
	return len(il.imports)
}

func (il ImportList) Get(pkgName PackageName) (res *Import, ok bool) {
	res, ok = il.imports[pkgName]
	return res, ok
}

func (il ImportList) asSlice() []*Import {
	var slice []*Import
	for i := range il.imports {
		slice = append(slice, il.imports[i])
	}
	return slice
}
