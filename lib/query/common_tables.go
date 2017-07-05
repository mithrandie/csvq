package query

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
)

type CommonTables map[string]*View

func (ct CommonTables) Set(commonTable parser.CommonTable) error {
	uname := strings.ToUpper(commonTable.Name.Literal)
	if _, err := ct.Get(uname); err == nil {
		return errors.New(fmt.Sprintf("common table %s already exists", commonTable.Name.Literal))
	}

	filter := Filter{}
	filter.CommonTables = ct
	if commonTable.IsRecursive() {
		filter.RecursiveTable = commonTable
	}
	view, err := SelectAsSubquery(commonTable.Query, filter)
	if err != nil {
		return err
	}

	err = view.UpdateHeader(commonTable.Name.Literal, commonTable.Columns)
	if err != nil {
		return err
	}

	ct[uname] = view
	return nil
}

func (ct CommonTables) Get(name string) (*View, error) {
	uname := strings.ToUpper(name)
	if view, ok := ct[uname]; ok {
		return view.Copy(), nil
	}
	return nil, errors.New(fmt.Sprintf("common table %s does not exist", name))
}

func (ct CommonTables) Copy() CommonTables {
	table := CommonTables{}
	for k, v := range ct {
		table[k] = v
	}
	return table
}

func (ct CommonTables) Merge(tables CommonTables) CommonTables {
	table := ct.Copy()
	for k, v := range tables {
		table[k] = v
	}
	return table
}

func (ct CommonTables) Load(clause parser.CommonTableClause) error {
	for _, v := range clause.CommonTables {
		commonTable := v.(parser.CommonTable)
		err := ct.Set(commonTable)
		if err != nil {
			return err
		}

		view, _ := ct.Get(commonTable.Name.Literal)
		err = view.UpdateHeader(commonTable.Name.Literal, commonTable.Columns)
		if err != nil {
			return err
		}
	}

	return nil
}
