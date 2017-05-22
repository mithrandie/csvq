package query

import (
	"errors"

	"github.com/mithrandie/csvq/lib/parser"
)

var Functions = map[string]func([]parser.Primary) (parser.Primary, error){
	"COALESCE": Coalesce,
}

func Coalesce(args []parser.Primary) (parser.Primary, error) {
	if len(args) < 1 {
		return nil, errors.New("function COALESCE is required at least 1 argument")
	}

	for _, arg := range args {
		if !parser.IsNull(arg) {
			return arg, nil
		}
	}
	return parser.NewNull(), nil
}
