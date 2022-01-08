package gocode

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type (
	// InterfaceName はインターフェース名を表す。
	InterfaceName string

	// PackageInterfaceName はパッケージ名付きインターフェース名。
	PackageInterfaceName string

	// Interface はinterfaceを表す。
	Interface struct {
		goInterface *types.Interface
		name        InterfaceName
		pkgSummary  *PackageSummary
		methods     *FunctionList
		embeds      *EmbedList
	}

	// InterfaceList はinterfaceのリストを表す。
	InterfaceList struct {
		interfaces []*Interface
	}
)

func (in InterfaceName) String() string {
	return string(in)
}

func NewPackageInterfaceName(pName PackageName, iName InterfaceName) PackageInterfaceName {
	return PackageInterfaceName(strings.Join([]string{pName.String(), iName.String()}, "."))
}

func (pin PackageInterfaceName) String() string {
	return string(pin)
}

func newInterfaceList(pkg *packages.Package) *InterfaceList {
	var interfaces []*Interface
	scope := pkg.Types.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if i, ok := newInterfaceIfInterfaceType(obj); ok {
			interfaces = append(interfaces, i)
		}
	}
	return &InterfaceList{interfaces: interfaces}
}

func (il *InterfaceList) asSlice() []*Interface {
	var slice []*Interface
	for i := range il.interfaces {
		slice = append(slice, il.interfaces[i])
	}
	return slice
}

func newInterfaceIfInterfaceType(obj types.Object) (res *Interface, ok bool) {
	interfaceType, ok := obj.Type().Underlying().(*types.Interface)
	if !ok {
		return &Interface{}, false
	}

	pkgSummary := newPackageSummaryFromGoTypes(obj.Pkg())

	return &Interface{
		goInterface: interfaceType,
		pkgSummary:  pkgSummary,
		name:        InterfaceName(obj.Name()),
		methods:     newFunctionListFromInterface(interfaceType),
		embeds:      newEmbedListFromInterfaceType(pkgSummary, interfaceType),
	}, true
}

func (i *Interface) PackageSummary() *PackageSummary {
	return i.pkgSummary
}

func (i *Interface) Name() InterfaceName {
	return i.name
}

func (i *Interface) PackageInterfaceName() PackageInterfaceName {
	return NewPackageInterfaceName(i.PackageSummary().Name(), i.Name())
}

func (i *Interface) Methods() []*Function {
	return i.methods.asSlice()
}

func (i *Interface) Embeds() []*Embed {
	return i.embeds.asSlice()
}
