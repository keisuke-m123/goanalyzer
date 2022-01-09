package testdata

import f "fmt"

type (
	ExportedStruct struct {
		Name string
		num  int
	}

	internalStruct struct {
		Name *string
		num  *int
	}

	ExportedInterface interface {
		Test(name string)
	}

	internalInterface interface {
		InternalTest(name string)
	}

	DefinedTypeString string

	AliasInt = int

	TestingSupport struct {
		ExportedStruct
		internalStruct
		ExportedInterface
		internalInterface

		es ExportedStruct
		is internalStruct
		ei ExportedInterface
		ii internalInterface
		Es ExportedStruct
		Is internalStruct
		Ei ExportedInterface
		Ii internalInterface

		Int int
	}
)

func (es *ExportedStruct) Test(name string) {
	f.Println("ExportedStruct.Test", name)
}

func (is *internalStruct) InternalTest(name string) {
	f.Println("internalStruct.InternalTest", name)
}
