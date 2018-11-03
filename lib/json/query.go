package json

import (
	"errors"
	"fmt"
	"github.com/mithrandie/csvq/lib/value"
)

func LoadValue(queryString string, json string) (value.Primary, error) {
	structure, _, err := load(queryString, json)
	if err != nil {
		return nil, err
	}

	return ConvertToValue(structure), nil
}

func LoadArray(queryString string, json string) ([]value.Primary, error) {
	structure, _, err := load(queryString, json)
	if err != nil {
		return nil, err
	}

	array, ok := structure.(Array)
	if !ok {
		return nil, errors.New(fmt.Sprintf("json value does not exists for %q", queryString))
	}

	return ConvertToArray(array), nil
}

func LoadTable(queryString string, json string) ([]string, [][]value.Primary, EscapeType, error) {
	structure, et, err := load(queryString, json)
	if err != nil {
		return nil, nil, et, err
	}

	array, ok := structure.(Array)
	if !ok {
		return nil, nil, et, errors.New(fmt.Sprintf("json value does not exists for %q", queryString))
	}

	h, rows, err := ConvertToTableValue(array)
	return h, rows, et, err
}

func load(queryString string, json string) (Structure, EscapeType, error) {
	query, err := Query.Parse(queryString)
	if err != nil {
		return nil, 0, err
	}

	d := NewDecoder()
	data, et, err := d.Decode(json)
	if err != nil {
		return nil, et, err
	}

	st, err := Extract(query, data)
	return st, et, err
}

func Extract(query QueryExpression, data Structure) (Structure, error) {
	var extracted Structure
	var err error

	if query == nil {
		return data, nil
	}

	switch query.(type) {
	case Element:
		switch data.(type) {
		case Object:
			element := query.(Element)

			obj := data.(Object)
			if obj.Exists(element.Label) {
				if element.Child == nil {
					extracted = obj.Value(element.Label)
				} else {
					extracted, err = Extract(element.Child, obj.Value(element.Label))
				}
			} else {
				extracted = Null{}
			}
		default:
			extracted = Null{}
		}
	case ArrayItem:
		switch data.(type) {
		case Array:
			arrayItem := query.(ArrayItem)

			ar := data.(Array)
			if arrayItem.Index < len(ar) {
				if arrayItem.Child == nil {
					extracted = ar[arrayItem.Index]
				} else {
					extracted, err = Extract(arrayItem.Child, ar[arrayItem.Index])
				}
			} else {
				extracted = Null{}
			}
		default:
			extracted = Null{}
		}
	case RowValueExpr:
		switch data.(type) {
		case Array:
			rowValue := query.(RowValueExpr)
			if rowValue.Child == nil {
				extracted = data
			} else {
				ar := data.(Array)
				elems := make(Array, 0, len(ar))
				for _, v := range ar {
					e, err := Extract(rowValue.Child, v)
					if err != nil {
						return extracted, err
					}
					elems = append(elems, e)
				}
				extracted = elems
			}
		default:
			return extracted, errors.New("json value must be an array")
		}
	case TableExpr:
		switch data.(type) {
		case Object:
			table := query.(TableExpr)
			if table.Fields == nil {
				extracted = Array{data}
			} else {
				obj := NewObject(len(table.Fields))
				for _, field := range table.Fields {
					e, err := Extract(field.Element, data)
					if err != nil {
						return extracted, err
					}

					obj.Add(field.FieldLabel(), e)
				}
				extracted = Array{obj}
			}
		case Array:
			table := query.(TableExpr)
			var fields []FieldExpr

			if table.Fields != nil {
				fields = table.Fields
			}

			array := data.(Array)
			for _, v := range array {
				obj, ok := v.(Object)
				if !ok {
					return extracted, errors.New("all elements in array must be objects")
				}

				if table.Fields == nil {
					if fields == nil {
						fields = make([]FieldExpr, 0, obj.Len())
					}
					for _, members := range obj.Members {
						if !existsKeyInFields(members.Key, fields) {
							fields = append(fields, FieldExpr{Element: Element{Label: members.Key}})
						}
					}
				}
			}

			elems := make(Array, 0, len(array))
			for _, v := range array {
				obj := NewObject(len(fields))
				for _, field := range fields {
					e, err := Extract(field.Element, v)
					if err != nil {
						return extracted, err
					}

					obj.Add(field.FieldLabel(), e)
				}
				elems = append(elems, obj)
			}
			extracted = elems
		default:
			return extracted, errors.New("json value must be an array or object")
		}
	default:
		return extracted, errors.New("invalid expression")
	}

	return extracted, err
}

func existsKeyInFields(key string, list []FieldExpr) bool {
	for _, v := range list {
		if key == v.Element.Label {
			return true
		}
	}
	return false
}
