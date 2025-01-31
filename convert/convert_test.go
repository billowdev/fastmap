package convert_test

import (
	"testing"

	"github.com/billowdev/fastmap/convert"
)

func TestConvert(t *testing.T) {
	type Source struct {
		Name  string
		Age   int
		Email string
	}

	type Dest struct {
		Name    string
		Age     int
		Contact string
	}

	src := &Source{
		Name:  "John",
		Age:   30,
		Email: "john@example.com",
	}

	dst := &Dest{}

	opts := &convert.FastmapConvertOptions{
		MatchCase: true,
		CustomMatchers: map[string]string{
			"Email": "Contact",
		},
	}

	err := convert.ConvertStruct(src, dst, opts)
	if err != nil {
		t.Errorf("Convert failed: %v", err)
	}

	if dst.Name != src.Name || dst.Age != src.Age || dst.Contact != src.Email {
		t.Error("Field conversion failed")
	}
}

func TestConvertCaseInsensitive(t *testing.T) {
	type Source struct {
		NAME string
		AGE  int
	}

	type Dest struct {
		Name string
		Age  int
	}

	src := &Source{
		NAME: "Jane",
		AGE:  25,
	}

	dst := &Dest{}

	opts := &convert.FastmapConvertOptions{MatchCase: false}

	err := convert.ConvertStruct(src, dst, opts)
	if err != nil {
		t.Errorf("Convert failed: %v", err)
	}

	if dst.Name != src.NAME || dst.Age != src.AGE {
		t.Error("Case-insensitive conversion failed")
	}
}

func TestConvertErrors(t *testing.T) {
	type TestStruct struct {
		Field string
	}

	src := TestStruct{}
	dst := TestStruct{}

	if err := convert.ConvertStruct(src, &dst, nil); err != convert.FastmapErrNotPointer {
		t.Error("Should fail on non-pointer source")
	}

	var srcInt, dstInt int
	if err := convert.ConvertStruct(&srcInt, &dstInt, nil); err != convert.FastmapErrNotStruct {
		t.Error("Should fail on non-struct types")
	}
}
