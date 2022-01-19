package gocode

import (
	"go/token"
	"go/types"
	"strings"
)

type (
	// InterfaceName はインターフェース名を表す。
	InterfaceName string

	// PackageInterfaceName はパッケージ名付きインターフェース名。
	PackageInterfaceName string

	// Interface はinterfaceを表す。
	Interface struct {
		definedPos  token.Pos
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

func newInterfaceList(pkg packageIn) *InterfaceList {
	var interfaces []*Interface
	scope := pkg.Scope()
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
		definedPos:  obj.Pos(),
		goInterface: interfaceType,
		pkgSummary:  pkgSummary,
		name:        InterfaceName(obj.Name()),
		methods:     newFunctionListFromInterface(interfaceType),
		embeds:      newEmbedListFromInterfaceType(pkgSummary, interfaceType),
	}, true
}

func (i *Interface) DefinedPos() token.Pos {
	return i.definedPos
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

func implements(typ types.Type, i *Interface) bool {
	if len(i.Methods()) == 0 {
		return false
	}

	switch t := typ.(type) {
	case *types.Pointer:
	default:
		// pointerにしておかないとtypes.Implementsで正しく判定されない
		typ = types.NewPointer(t)
	}
	return types.Implements(typ, i.goInterface)
}
