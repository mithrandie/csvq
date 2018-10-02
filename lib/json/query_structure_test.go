package json

import (
	"testing"
)

var elementFieldLabelTests = []struct {
	Element Element
	Expect  string
}{
	{
		Element: Element{
			Label: "key",
		},
		Expect: "key",
	},
	{
		Element: Element{
			Label: "ke.y",
		},
		Expect: "ke\\.y",
	},
	{
		Element: Element{
			Label: "key",
			Child: Element{
				Label: "key2",
			},
		},
		Expect: "key.key2",
	},
	{
		Element: Element{
			Label: "key",
			Child: ArrayItem{
				Index: 2,
			},
		},
		Expect: "key[2]",
	},
}

func TestElement_FieldLabel(t *testing.T) {
	for _, v := range elementFieldLabelTests {
		result := v.Element.FieldLabel()
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %#v", result, v.Expect, v.Element)
		}
	}
}

var arrayItemFieldLabelTests = []struct {
	ArrayItem ArrayItem
	Expect    string
}{
	{
		ArrayItem: ArrayItem{
			Index: 0,
		},
		Expect: "[0]",
	},
	{
		ArrayItem: ArrayItem{
			Index: 0,
			Child: Element{
				Label: "key2",
			},
		},
		Expect: "[0].key2",
	},
	{
		ArrayItem: ArrayItem{
			Index: 0,
			Child: ArrayItem{
				Index: 2,
			},
		},
		Expect: "[0][2]",
	},
}

func TestArrayItem_FieldLabel(t *testing.T) {
	for _, v := range arrayItemFieldLabelTests {
		result := v.ArrayItem.FieldLabel()
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %#v", result, v.Expect, v.ArrayItem)
		}
	}
}

var columnExprFieldLabelTests = []struct {
	ColumnExpr FieldExpr
	Expect     string
}{
	{
		ColumnExpr: FieldExpr{
			Element: Element{Label: "key"},
		},
		Expect: "key",
	},
	{
		ColumnExpr: FieldExpr{
			Element: Element{Label: "key"},
			Alias:   "alias",
		},
		Expect: "alias",
	},
}

func TestColumnExpr_FieldLabel(t *testing.T) {
	for _, v := range columnExprFieldLabelTests {
		result := v.ColumnExpr.FieldLabel()
		if result != v.Expect {
			t.Errorf("result = %q, want %q for %#v", result, v.Expect, v.ColumnExpr)
		}
	}
}
