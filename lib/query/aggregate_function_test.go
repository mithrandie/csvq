package query

import (
	"reflect"
	"testing"

	"github.com/mithrandie/csvq/lib/parser"
)

type aggregateTests struct {
	Distinct bool
	List     []parser.Primary
	Result   parser.Primary
}

var countTests = []aggregateTests{
	{
		Distinct: false,
		List: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(1),
			parser.NewInteger(1),
			parser.NewNull(),
			parser.NewInteger(4),
			parser.NewNull(),
		},
		Result: parser.NewInteger(4),
	},
	{
		Distinct: true,
		List: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(1),
			parser.NewInteger(1),
			parser.NewNull(),
			parser.NewInteger(4),
			parser.NewNull(),
		},
		Result: parser.NewInteger(3),
	},
	{
		Distinct: false,
		List: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewInteger(0),
	},
}

func TestCount(t *testing.T) {
	for _, v := range countTests {
		r := Count(v.Distinct, v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("count distinct = %t, list = %s: result = %s, want %s", v.Distinct, v.List, r, v.Result)
		}
	}
}

var maxTests = []aggregateTests{
	{
		Distinct: false,
		List: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(1),
			parser.NewInteger(1),
			parser.NewNull(),
			parser.NewInteger(4),
			parser.NewNull(),
		},
		Result: parser.NewInteger(4),
	},
	{
		Distinct: true,
		List: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(1),
			parser.NewInteger(1),
			parser.NewNull(),
			parser.NewInteger(4),
			parser.NewNull(),
		},
		Result: parser.NewInteger(4),
	},
	{
		Distinct: false,
		List: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
}

func TestMax(t *testing.T) {
	for _, v := range maxTests {
		r := Max(v.Distinct, v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("max distinct = %t, list = %s: result = %s, want %s", v.Distinct, v.List, r, v.Result)
		}
	}
}

var minTests = []aggregateTests{
	{
		Distinct: false,
		List: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(1),
			parser.NewInteger(1),
			parser.NewNull(),
			parser.NewInteger(4),
			parser.NewNull(),
		},
		Result: parser.NewInteger(1),
	},
	{
		Distinct: true,
		List: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(1),
			parser.NewInteger(1),
			parser.NewNull(),
			parser.NewInteger(4),
			parser.NewNull(),
		},
		Result: parser.NewInteger(1),
	},
	{
		Distinct: false,
		List: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
}

func TestMin(t *testing.T) {
	for _, v := range minTests {
		r := Min(v.Distinct, v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("min distinct = %t, list = %s: result = %s, want %s", v.Distinct, v.List, r, v.Result)
		}
	}
}

var sumTests = []aggregateTests{
	{
		Distinct: false,
		List: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(1),
			parser.NewInteger(1),
			parser.NewNull(),
			parser.NewInteger(4),
			parser.NewNull(),
		},
		Result: parser.NewInteger(8),
	},
	{
		Distinct: true,
		List: []parser.Primary{
			parser.NewInteger(2),
			parser.NewInteger(1),
			parser.NewInteger(1),
			parser.NewNull(),
			parser.NewInteger(4),
			parser.NewNull(),
		},
		Result: parser.NewInteger(7),
	},
	{
		Distinct: false,
		List: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
}

func TestSum(t *testing.T) {
	for _, v := range sumTests {
		r := Sum(v.Distinct, v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("sum distinct = %t, list = %s: result = %s, want %s", v.Distinct, v.List, r, v.Result)
		}
	}
}

var avgTests = []aggregateTests{
	{
		Distinct: false,
		List: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewNull(),
			parser.NewInteger(4),
			parser.NewNull(),
		},
		Result: parser.NewInteger(2),
	},
	{
		Distinct: true,
		List: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewNull(),
			parser.NewInteger(4),
			parser.NewNull(),
		},
		Result: parser.NewFloat(7.0 / 3.0),
	},
	{
		Distinct: false,
		List: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
}

func TestAvg(t *testing.T) {
	for _, v := range avgTests {
		r := Avg(v.Distinct, v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("avg distinct = %t, list = %s: result = %s, want %s", v.Distinct, v.List, r, v.Result)
		}
	}
}
