package gocode_test

import (
	"testing"

	"github.com/keisuke-m123/goanalyzer/gocode/testdata"
)

func TestType_EqualReflectionType(t *testing.T) {
	s, ok := testingSupportPackages.Structs().Get("testdata", "ExportedStruct")
	if !ok {
		t.Fatal("expected to find struct")
	}
	if !s.Type().EqualReflectionType((*testdata.ExportedStruct)(nil)) {
		t.Error("EqualReflectionType failed")
	}
}
