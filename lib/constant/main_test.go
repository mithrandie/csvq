package constant

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"
)

func TestGet(t *testing.T) {
	Definition["TEST"] = map[string]interface{}{
		"INVALID_TYPE": uint8(3),
	}

	defer func() {
		delete(Definition, "Unused")
	}()

	expr := parser.Constant{
		Space: "integer",
		Name:  "max",
	}
	var expect value.Primary = value.NewInteger(math.MaxInt64)

	ret, err := Get(expr)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}

	if !reflect.DeepEqual(ret, expect) {
		t.Errorf("result = %s, want %s", ret.String(), expect.String())
	}

	expr = parser.Constant{
		Space: "math",
		Name:  "pi",
	}
	expect = value.NewFloat(math.Pi)

	ret, err = Get(expr)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}

	if !reflect.DeepEqual(ret, expect) {
		t.Errorf("result = %s, want %s", ret.String(), expect.String())
	}

	expr = parser.Constant{
		Space: "NotDefined",
		Name:  "pi",
	}
	expectError := ErrUndefined

	_, err = Get(expr)
	if err == nil {
		t.Errorf("no error, want error %q", expectError.Error())
	}
	if !errors.Is(err, expectError) {
		t.Errorf("error %q, want error %q", err.Error(), expectError.Error())
	}

	expr = parser.Constant{
		Space: "math",
		Name:  "NotDefined",
	}

	_, err = Get(expr)
	if err == nil {
		t.Errorf("no error, want error %q", expectError.Error())
	}
	if !errors.Is(err, expectError) {
		t.Errorf("error %q, want error %q", err.Error(), expectError.Error())
	}

	expr = parser.Constant{
		Space: "test",
		Name:  "invalid_type",
	}
	expectError = ErrInvalidType

	_, err = Get(expr)
	if err == nil {
		t.Errorf("no error, want error %q", expectError.Error())
	}
	if !errors.Is(err, expectError) {
		t.Errorf("error %q, want error %q", err.Error(), expectError.Error())
	}
}

func TestCount(t *testing.T) {
	oldDef := Definition

	Definition = map[string]map[string]interface{}{}
	Definition["CAT1"] = map[string]interface{}{
		"VAL1": 1,
	}
	Definition["CAT2"] = map[string]interface{}{
		"VAL2": 2,
		"VAL3": 3,
	}

	defer func() {
		Definition = oldDef
	}()

	ret := Count()
	if ret != 3 {
		t.Errorf("result = %d, want %d", ret, 3)
	}
}
