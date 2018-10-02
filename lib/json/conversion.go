package json

import (
	"errors"
	"fmt"
	"github.com/mithrandie/csvq/lib/value"
	"github.com/mithrandie/ternary"
	"time"
)

func ConvertToValue(structure Structure) value.Primary {
	var p value.Primary

	switch structure.(type) {
	case Number:
		p = value.ParseFloat64(float64(structure.(Number)))
	case String:
		p = value.NewString(string(structure.(String)))
	case Boolean:
		p = value.NewBoolean(bool(structure.(Boolean)))
	case Null:
		p = value.NewNull()
	default:
		p = value.NewString(structure.String())
	}

	return p
}

func ConvertToArray(array Array) []value.Primary {
	row := make([]value.Primary, 0, len(array))
	for _, v := range array {
		row = append(row, ConvertToValue(v))
	}

	return row
}

func ConvertToTableValue(array Array) ([]string, [][]value.Primary, error) {
	exists := func(s string, list []string) bool {
		for _, v := range list {
			if s == v {
				return true
			}
		}
		return false
	}

	var header []string
	for _, elem := range array {
		obj, ok := elem.(Object)
		if !ok {
			return nil, nil, errors.New("rows loaded from json must be objects")
		}
		if header == nil {
			header = make([]string, 0, obj.Len())
		}

		for _, m := range obj.Members {
			if !exists(m.Key, header) {
				header = append(header, m.Key)
			}
		}
	}

	rows := make([][]value.Primary, 0, len(array))
	for _, elem := range array {
		row := make([]value.Primary, 0, len(header))

		obj, _ := elem.(Object)
		for _, column := range header {
			if obj.Exists(column) {
				row = append(row, ConvertToValue(obj.Value(column)))
			} else {
				row = append(row, ConvertToValue(Null{}))
			}
		}

		rows = append(rows, row)
	}

	return header, rows, nil
}

func ConvertTableValueToJsonStructure(fields []string, rows [][]value.Primary) (Structure, error) {
	pathes, err := ParsePathes(fields)
	if err != nil {
		return nil, err
	}

	structure := make(Array, 0, len(rows))
	for _, row := range rows {
		rowStructure, err := ConvertRecordValueToJsonStructure(pathes, row)
		if err != nil {
			return nil, err
		}
		structure = append(structure, rowStructure)
	}

	return structure, nil
}

func ParsePathes(fields []string) ([]PathExpression, error) {
	var err error
	pathes := make([]PathExpression, len(fields))
	for i, field := range fields {
		pathes[i], err = Path.Parse(field)
		if err != nil {
			if perr, ok := err.(*PathSyntaxError); ok {
				err = errors.New(fmt.Sprintf("%s at column %d in %q", perr.Error(), perr.Column, field))
			}
			return nil, err
		}
	}
	return pathes, nil
}

func ConvertRecordValueToJsonStructure(pathes []PathExpression, row []value.Primary) (Structure, error) {
	var structure Structure

	fieldLen := len(pathes)

	if len(row) != fieldLen {
		return nil, errors.New("field length does not match")
	}

	for i, path := range pathes {
		structure = addPathValueToRowStructure(structure, path.(ObjectPath), row[i], fieldLen)
	}

	return structure, nil
}

func addPathValueToRowStructure(parent Structure, path ObjectPath, val value.Primary, fieldLen int) Structure {
	var obj Object
	if parent == nil {
		obj = NewObject(fieldLen)
	} else {
		obj = parent.(Object)
	}

	if path.Child == nil {
		obj.Add(path.Name, ParseValueToStructure(val))
	} else {
		valueStructure := addPathValueToRowStructure(obj.Value(path.Name), path.Child.(ObjectPath), val, fieldLen)
		if obj.Exists(path.Name) {
			obj.Update(path.Name, valueStructure)
		} else {
			obj.Add(path.Name, valueStructure)
		}
	}

	return obj
}

func ParseValueToStructure(val value.Primary) Structure {
	var s Structure

	switch val.(type) {
	case value.String:
		s = String(val.(value.String).Raw())
	case value.Integer:
		s = Number(val.(value.Integer).Raw())
	case value.Float:
		s = Number(val.(value.Float).Raw())
	case value.Boolean:
		s = Boolean(val.(value.Boolean).Raw())
	case value.Ternary:
		t := val.(value.Ternary)
		if t.Ternary() == ternary.UNKNOWN {
			s = Null{}
		} else {
			s = Boolean(t.Ternary().ParseBool())
		}
	case value.Datetime:
		s = String(val.(value.Datetime).Format(time.RFC3339Nano))
	case value.Null:
		s = Null{}
	}

	return s
}
