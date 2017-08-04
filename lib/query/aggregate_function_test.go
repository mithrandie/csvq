package query

import (
	"reflect"
	"testing"
	"time"

	"github.com/mithrandie/csvq/lib/parser"
)

type aggregateTests struct {
	List   []parser.Primary
	Result parser.Primary
}

var countTests = []aggregateTests{
	{
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
		List: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewInteger(0),
	},
}

func TestCount(t *testing.T) {
	for _, v := range countTests {
		r := Count(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("count list = %s: result = %s, want %s", v.List, r, v.Result)
		}
	}
}

var maxTests = []aggregateTests{
	{
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
		List: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
}

func TestMax(t *testing.T) {
	for _, v := range maxTests {
		r := Max(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("max list = %s: result = %s, want %s", v.List, r, v.Result)
		}
	}
}

var minTests = []aggregateTests{
	{
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
		List: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
}

func TestMin(t *testing.T) {
	for _, v := range minTests {
		r := Min(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("min list = %s: result = %s, want %s", v.List, r, v.Result)
		}
	}
}

var sumTests = []aggregateTests{
	{
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
		List: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
}

func TestSum(t *testing.T) {
	for _, v := range sumTests {
		r := Sum(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("sum list = %s: result = %s, want %s", v.List, r, v.Result)
		}
	}
}

var avgTests = []aggregateTests{
	{
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
		List: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
}

func TestAvg(t *testing.T) {
	for _, v := range avgTests {
		r := Avg(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("avg list = %s: result = %s, want %s", v.List, r, v.Result)
		}
	}
}

var medianTests = []aggregateTests{
	{
		List: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(4),
			parser.NewInteger(6),
			parser.NewNull(),
			parser.NewInteger(1),
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewNull(),
		},
		Result: parser.NewFloat(1.5),
	},
	{
		List: []parser.Primary{
			parser.NewInteger(1),
			parser.NewInteger(4),
			parser.NewInteger(6),
			parser.NewNull(),
			parser.NewInteger(1),
			parser.NewInteger(2),
			parser.NewNull(),
		},
		Result: parser.NewInteger(2),
	},
	{
		List: []parser.Primary{
			parser.NewDatetime(time.Date(2012, 2, 3, 9, 18, 15, 0, GetTestLocation())),
			parser.NewDatetime(time.Date(2012, 2, 5, 9, 18, 15, 0, GetTestLocation())),
		},
		Result: parser.NewInteger(time.Date(2012, 2, 4, 9, 18, 15, 0, GetTestLocation()).Unix()),
	},
	{
		List: []parser.Primary{
			parser.NewNull(),
		},
		Result: parser.NewNull(),
	},
}

func TestMedian(t *testing.T) {
	for _, v := range medianTests {
		r := Median(v.List)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("median list = %s: result = %s, want %s", v.List, r, v.Result)
		}
	}
}

var listAggTests = []struct {
	List      []parser.Primary
	Separator string
	Result    parser.Primary
}{
	{
		List: []parser.Primary{
			parser.NewString("str1"),
			parser.NewString("str3"),
			parser.NewNull(),
			parser.NewString("str2"),
			parser.NewString("str1"),
			parser.NewString("str2"),
		},
		Separator: ",",
		Result:    parser.NewString("str1,str3,str2,str1,str2"),
	},
	{
		List: []parser.Primary{
			parser.NewNull(),
		},
		Separator: ",",
		Result:    parser.NewNull(),
	},
}

func TestListAgg(t *testing.T) {
	for _, v := range listAggTests {
		r := ListAgg(v.List, v.Separator)
		if !reflect.DeepEqual(r, v.Result) {
			t.Errorf("listagg list = %s: separator = %q, result = %s, want %s", v.List, v.Separator, r, v.Result)
		}
	}
}
