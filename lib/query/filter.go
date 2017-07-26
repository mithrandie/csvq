package query

import (
	"strings"

	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

type FilterRecord struct {
	View        *View
	RecordIndex int
}

type Filter struct {
	Records []FilterRecord

	VariablesList VariablesList
	TempViewsList TemporaryViewMapList
	CursorsList   CursorMapList
	FunctionsList UserDefinedFunctionsList

	InlineTablesList InlineTablesList
	AliasesList      AliasMapList

	RecursiveTable    *parser.InlineTable
	RecursiveTmpView  *View
	tmpViewIsAccessed bool
}

func NewFilter(variablesList VariablesList, tempViewsList TemporaryViewMapList, cursorsList CursorMapList, functionsList UserDefinedFunctionsList) Filter {
	return Filter{
		VariablesList: variablesList,
		TempViewsList: tempViewsList,
		CursorsList:   cursorsList,
		FunctionsList: functionsList,
	}
}

func NewEmptyFilter() Filter {
	return NewFilter(
		VariablesList{{}},
		TemporaryViewMapList{{}},
		CursorMapList{{}},
		UserDefinedFunctionsList{{}},
	)
}

func NewFilterForRecord(view *View, recordIndex int, parentFilter Filter) Filter {
	f := Filter{
		Records: []FilterRecord{
			{
				View:        view,
				RecordIndex: recordIndex,
			},
		},
	}
	return f.Merge(parentFilter)
}

func NewFilterForSequentialEvaluation(view *View, parentFilter Filter) Filter {
	f := Filter{
		Records: []FilterRecord{
			{
				View: view,
			},
		},
	}
	return f.Merge(parentFilter)
}

func (f Filter) Merge(filter Filter) Filter {
	return Filter{
		Records:          append(f.Records, filter.Records...),
		VariablesList:    filter.VariablesList,
		TempViewsList:    filter.TempViewsList,
		CursorsList:      filter.CursorsList,
		FunctionsList:    filter.FunctionsList,
		InlineTablesList: filter.InlineTablesList,
		AliasesList:      filter.AliasesList,
	}
}

func (f Filter) CreateChildScope() Filter {
	return NewFilter(
		append(VariablesList{{}}, f.VariablesList...),
		append(TemporaryViewMapList{{}}, f.TempViewsList...),
		append(CursorMapList{{}}, f.CursorsList...),
		append(UserDefinedFunctionsList{{}}, f.FunctionsList...),
	)
}

func (f Filter) CreateNode() Filter {
	return Filter{
		Records:          f.Records,
		VariablesList:    f.VariablesList,
		TempViewsList:    f.TempViewsList,
		CursorsList:      f.CursorsList,
		FunctionsList:    f.FunctionsList,
		InlineTablesList: append(InlineTablesList{{}}, f.InlineTablesList...),
		AliasesList:      append(AliasMapList{{}}, f.AliasesList...),
		RecursiveTable:   f.RecursiveTable,
		RecursiveTmpView: f.RecursiveTmpView,
	}
}

func (f Filter) LoadInlineTable(clause parser.WithClause) error {
	return f.InlineTablesList.Load(clause, f)
}

func (f Filter) Evaluate(expr parser.Expression) (parser.Primary, error) {
	if expr == nil {
		return parser.NewTernary(ternary.TRUE), nil
	}

	var primary parser.Primary
	var err error

	if parser.IsPrimary(expr) {
		primary = expr.(parser.Primary)
	} else {
		switch expr.(type) {
		case parser.Parentheses:
			primary, err = f.Evaluate(expr.(parser.Parentheses).Expr)
		case parser.FieldReference:
			primary, err = f.evalFieldReference(expr.(parser.FieldReference))
		case parser.ColumnNumber:
			primary, err = f.evalColumnNumber(expr.(parser.ColumnNumber))
		case parser.Arithmetic:
			primary, err = f.evalArithmetic(expr.(parser.Arithmetic))
		case parser.UnaryArithmetic:
			primary, err = f.evalUnaryArithmetic(expr.(parser.UnaryArithmetic))
		case parser.Concat:
			primary, err = f.evalConcat(expr.(parser.Concat))
		case parser.Comparison:
			primary, err = f.evalComparison(expr.(parser.Comparison))
		case parser.Is:
			primary, err = f.evalIs(expr.(parser.Is))
		case parser.Between:
			primary, err = f.evalBetween(expr.(parser.Between))
		case parser.Like:
			primary, err = f.evalLike(expr.(parser.Like))
		case parser.In:
			primary, err = f.evalIn(expr.(parser.In))
		case parser.Any:
			primary, err = f.evalAny(expr.(parser.Any))
		case parser.All:
			primary, err = f.evalAll(expr.(parser.All))
		case parser.Exists:
			primary, err = f.evalExists(expr.(parser.Exists))
		case parser.Subquery:
			primary, err = f.evalSubqueryForSingleValue(expr.(parser.Subquery))
		case parser.Function:
			primary, err = f.evalFunction(expr.(parser.Function))
		case parser.AggregateFunction:
			primary, err = f.evalAggregateFunction(expr.(parser.AggregateFunction))
		case parser.GroupConcat:
			primary, err = f.evalGroupConcat(expr.(parser.GroupConcat))
		case parser.Case:
			primary, err = f.evalCase(expr.(parser.Case))
		case parser.Logic:
			primary, err = f.evalLogic(expr.(parser.Logic))
		case parser.UnaryLogic:
			primary, err = f.evalUnaryLogic(expr.(parser.UnaryLogic))
		case parser.Variable:
			primary, err = f.VariablesList.Get(expr.(parser.Variable))
		case parser.VariableSubstitution:
			primary, err = f.VariablesList.Substitute(expr.(parser.VariableSubstitution), f)
		case parser.CursorStatus:
			primary, err = f.evalCursorStatus(expr.(parser.CursorStatus))
		default:
			return nil, NewSyntaxErrorFromExpr(expr)
		}
	}

	return primary, err
}

func (f Filter) evalFieldReference(expr parser.FieldReference) (parser.Primary, error) {
	var p parser.Primary
	for _, v := range f.Records {
		idx, err := v.View.Header.Contains(expr)
		if err == nil {
			if v.View.isGrouped && !v.View.Header[idx].IsGroupKey {
				return nil, NewFieldNotGroupKeyError(expr)
			}
			p = v.View.Records[v.RecordIndex][idx].Primary()
			break
		}

		if _, ok := err.(*FieldAmbiguousError); ok {
			return nil, err
		}
	}
	if p == nil {
		return nil, NewFieldNotExistError(expr)
	}
	return p, nil
}

func (f Filter) evalColumnNumber(expr parser.ColumnNumber) (parser.Primary, error) {
	for _, v := range f.Records {
		idx, err := v.View.Header.ContainsNumber(expr)
		if err == nil {
			if v.View.isGrouped && !v.View.Header[idx].IsGroupKey {
				return nil, NewFieldNumberNotGroupKeyError(expr)
			}
			return v.View.Records[v.RecordIndex][idx].Primary(), nil
		}
	}
	return nil, NewFieldNumberNotExistError(expr)
}

func (f Filter) evalArithmetic(expr parser.Arithmetic) (parser.Primary, error) {
	lhs, err := f.Evaluate(expr.LHS)
	if err != nil {
		return nil, err
	}
	rhs, err := f.Evaluate(expr.RHS)
	if err != nil {
		return nil, err
	}

	return Calculate(lhs, rhs, expr.Operator), nil
}

func (f Filter) evalUnaryArithmetic(expr parser.UnaryArithmetic) (parser.Primary, error) {
	ope, err := f.Evaluate(expr.Operand)
	if err != nil {
		return nil, err
	}

	pf := parser.PrimaryToFloat(ope)
	if parser.IsNull(pf) {
		return parser.NewNull(), nil
	}

	value := pf.(parser.Float).Value()

	switch expr.Operator.Token {
	case '-':
		value = value * -1
	}

	return parser.Float64ToPrimary(value), nil
}

func (f Filter) evalConcat(expr parser.Concat) (parser.Primary, error) {
	items := make([]string, len(expr.Items))
	for i, v := range expr.Items {
		s, err := f.Evaluate(v)
		if err != nil {
			return nil, err
		}
		s = parser.PrimaryToString(s)
		if parser.IsNull(s) {
			return parser.NewNull(), nil
		}
		items[i] = s.(parser.String).Value()
	}
	return parser.NewString(strings.Join(items, "")), nil
}

func (f Filter) evalComparison(expr parser.Comparison) (parser.Primary, error) {
	var t ternary.Value

	switch expr.LHS.(type) {
	case parser.RowValue:
		lhs, err := f.evalRowValue(expr.LHS.(parser.RowValue))
		if err != nil {
			return nil, err
		}

		rhs, err := f.evalRowValue(expr.RHS.(parser.RowValue))
		if err != nil {
			return nil, err
		}

		t, err = CompareRowValues(lhs, rhs, expr.Operator)
		if err != nil {
			return nil, NewRowValueLengthInComparisonError(expr.RHS.(parser.RowValue), len(lhs))
		}

	default:
		lhs, err := f.Evaluate(expr.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := f.Evaluate(expr.RHS)
		if err != nil {
			return nil, err
		}

		t = Compare(lhs, rhs, expr.Operator)
	}
	return parser.NewTernary(t), nil
}

func (f Filter) evalIs(expr parser.Is) (parser.Primary, error) {
	lhs, err := f.Evaluate(expr.LHS)
	if err != nil {
		return nil, err
	}
	rhs, err := f.Evaluate(expr.RHS)
	if err != nil {
		return nil, err
	}

	t := Is(lhs, rhs)
	if expr.IsNegated() {
		t = ternary.Not(t)
	}
	return parser.NewTernary(t), nil
}

func (f Filter) evalBetween(expr parser.Between) (parser.Primary, error) {
	var t ternary.Value

	switch expr.LHS.(type) {
	case parser.RowValue:
		lhs, err := f.evalRowValue(expr.LHS.(parser.RowValue))
		if err != nil {
			return nil, err
		}

		low, err := f.evalRowValue(expr.Low.(parser.RowValue))
		if err != nil {
			return nil, err
		}

		high, err := f.evalRowValue(expr.High.(parser.RowValue))
		if err != nil {
			return nil, err
		}

		t1, err := CompareRowValues(lhs, low, ">=")
		if err != nil {
			return nil, NewRowValueLengthInComparisonError(expr.Low.(parser.RowValue), len(lhs))
		}
		t2, err := CompareRowValues(lhs, high, "<=")
		if err != nil {
			return nil, NewRowValueLengthInComparisonError(expr.High.(parser.RowValue), len(lhs))
		}

		t = ternary.And(t1, t2)
	default:
		lhs, err := f.Evaluate(expr.LHS)
		if err != nil {
			return nil, err
		}
		low, err := f.Evaluate(expr.Low)
		if err != nil {
			return nil, err
		}
		high, err := f.Evaluate(expr.High)
		if err != nil {
			return nil, err
		}

		t = ternary.And(GreaterThanOrEqualTo(lhs, low), LessThanOrEqualTo(lhs, high))
	}

	if expr.IsNegated() {
		t = ternary.Not(t)
	}
	return parser.NewTernary(t), nil
}

func (f Filter) valuesForRowValueListComparison(lhs parser.Expression, values parser.Expression) ([]parser.Primary, [][]parser.Primary, error) {
	var value []parser.Primary
	var list [][]parser.Primary
	var err error

	switch lhs.(type) {
	case parser.RowValue:
		value, err = f.evalRowValue(lhs.(parser.RowValue))
		if err != nil {
			return nil, nil, err
		}

		list, err = f.evalRowValues(values)
		if err != nil {
			return nil, nil, err
		}

	default:
		lhs, err := f.Evaluate(lhs)
		if err != nil {
			return nil, nil, err
		}
		value = []parser.Primary{lhs}

		rowValue := values.(parser.RowValue)
		switch rowValue.Value.(type) {
		case parser.Subquery:
			list, err = f.evalSubqueryForSingleFieldRowValues(rowValue.Value.(parser.Subquery))
			if err != nil {
				return nil, nil, err
			}
		case parser.ValueList:
			values, err := f.evalValueList(rowValue.Value.(parser.ValueList))
			if err != nil {
				return nil, nil, err
			}
			list = make([][]parser.Primary, len(values))
			for i, v := range values {
				list[i] = []parser.Primary{v}
			}
		}
	}
	return value, list, nil
}

func (f Filter) evalIn(expr parser.In) (parser.Primary, error) {
	value, list, err := f.valuesForRowValueListComparison(expr.LHS, expr.Values)
	if err != nil {
		return nil, err
	}

	t, err := Any(value, list, "=")
	if err != nil {
		if subquery, ok := expr.Values.(parser.Subquery); ok {
			return nil, NewSelectFieldLengthInComparisonError(subquery, len(value))
		}

		rvlist, _ := expr.Values.(parser.RowValueList)
		rverr, _ := err.(*RowValueLengthInListError)
		return nil, NewRowValueLengthInComparisonError(rvlist.RowValues[rverr.Index], len(value))
	}

	if expr.IsNegated() {
		t = ternary.Not(t)
	}
	return parser.NewTernary(t), nil
}

func (f Filter) evalAny(expr parser.Any) (parser.Primary, error) {
	value, list, err := f.valuesForRowValueListComparison(expr.LHS, expr.Values)
	if err != nil {
		return nil, err
	}

	t, err := Any(value, list, expr.Operator)
	if err != nil {
		if subquery, ok := expr.Values.(parser.Subquery); ok {
			return nil, NewSelectFieldLengthInComparisonError(subquery, len(value))
		}

		rvlist, _ := expr.Values.(parser.RowValueList)
		rverr, _ := err.(*RowValueLengthInListError)
		return nil, NewRowValueLengthInComparisonError(rvlist.RowValues[rverr.Index], len(value))
	}
	return parser.NewTernary(t), nil
}

func (f Filter) evalAll(expr parser.All) (parser.Primary, error) {
	value, list, err := f.valuesForRowValueListComparison(expr.LHS, expr.Values)
	if err != nil {
		return nil, err
	}

	t, err := All(value, list, expr.Operator)
	if err != nil {
		if subquery, ok := expr.Values.(parser.Subquery); ok {
			return nil, NewSelectFieldLengthInComparisonError(subquery, len(value))
		}

		rvlist, _ := expr.Values.(parser.RowValueList)
		rverr, _ := err.(*RowValueLengthInListError)
		return nil, NewRowValueLengthInComparisonError(rvlist.RowValues[rverr.Index], len(value))
	}
	return parser.NewTernary(t), nil
}

func (f Filter) evalLike(expr parser.Like) (parser.Primary, error) {
	lhs, err := f.Evaluate(expr.LHS)
	if err != nil {
		return nil, err
	}
	pattern, err := f.Evaluate(expr.Pattern)
	if err != nil {
		return nil, err
	}

	t := Like(lhs, pattern)
	if expr.IsNegated() {
		t = ternary.Not(t)
	}
	return parser.NewTernary(t), nil
}

func (f Filter) evalExists(expr parser.Exists) (parser.Primary, error) {
	view, err := Select(expr.Query.Query, f)
	if err != nil {
		return nil, err
	}
	if view.RecordLen() < 1 {
		return parser.NewTernary(ternary.FALSE), nil
	}
	return parser.NewTernary(ternary.TRUE), nil
}

func (f Filter) evalFunction(expr parser.Function) (parser.Primary, error) {
	name := strings.ToUpper(expr.Name)
	if strings.EqualFold("GROUP_CONCAT", name) {
		gc := parser.GroupConcat{
			BaseExpr:    expr.BaseExpr,
			GroupConcat: expr.Name,
			Option: parser.AggregateOption{
				Args: expr.Args,
			},
		}
		return f.evalGroupConcat(gc)
	} else if _, ok := AggregateFunctions[name]; ok {
		afn := parser.AggregateFunction{
			BaseExpr: expr.BaseExpr,
			Name:     expr.Name,
			Option: parser.AggregateOption{
				Args: expr.Args,
			},
		}
		return f.evalAggregateFunction(afn)
	}

	args := make([]parser.Primary, len(expr.Args))
	for i, v := range expr.Args {
		arg, err := f.Evaluate(v)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}

	fn, ok := Functions[name]
	if !ok {
		if udfn, err := f.FunctionsList.Get(expr); err == nil {
			return udfn.Execute(args, f)
		}

		return nil, NewFunctionNotExistError(expr, expr.Name)
	}

	return fn(expr, args)
}

func (f Filter) evalAggregateFunction(expr parser.AggregateFunction) (parser.Primary, error) {
	name := strings.ToUpper(expr.Name)
	if strings.EqualFold("GROUP_CONCAT", name) {
		gc := parser.GroupConcat{
			BaseExpr:    expr.BaseExpr,
			GroupConcat: expr.Name,
			Option:      expr.Option,
		}
		return f.evalGroupConcat(gc)
	}

	fn, ok := AggregateFunctions[name]
	if !ok {
		return nil, NewFunctionNotExistError(expr, expr.Name)
	}

	if len(f.Records) < 1 {
		return nil, NewUnpermittedStatementFunctionError(expr, expr.Name)
	}

	if !f.Records[0].View.isGrouped {
		return nil, NewNotGroupingRecordsError(expr, expr.Name)
	}

	if len(expr.Option.Args) != 1 {
		return nil, NewFunctionArgumentLengthError(expr, expr.Name, []int{1})
	}

	arg := expr.Option.Args[0]
	if ac, ok := arg.(parser.AllColumns); ok {
		if !strings.EqualFold(expr.Name, "COUNT") {
			return nil, NewUnpermittedWildCardError(ac, expr.Name)
		}
		arg = parser.NewInteger(1)
	}

	view := NewViewFromGroupedRecord(f.Records[0])

	list := make([]parser.Primary, view.RecordLen())
	filter := NewFilterForSequentialEvaluation(view, f)
	for i := 0; i < view.RecordLen(); i++ {
		filter.Records[0].RecordIndex = i
		p, err := filter.Evaluate(arg)
		if err != nil {
			if _, ok := err.(*NotGroupingRecordsError); ok {
				err = NewNestedAggregateFunctionsError(expr)
			}
			return nil, err
		}
		list[i] = p
	}

	return fn(expr.Option.IsDistinct(), list), nil
}

func (f Filter) evalGroupConcat(expr parser.GroupConcat) (parser.Primary, error) {
	if len(f.Records) < 1 {
		return nil, NewUnpermittedStatementFunctionError(expr, expr.GroupConcat)
	}

	if !f.Records[0].View.isGrouped {
		return nil, NewNotGroupingRecordsError(expr, expr.GroupConcat)
	}

	if len(expr.Option.Args) != 1 {
		return nil, NewFunctionArgumentLengthError(expr, expr.GroupConcat, []int{1})
	}

	arg := expr.Option.Args[0]
	if ac, ok := arg.(parser.AllColumns); ok {
		return nil, NewUnpermittedWildCardError(ac, expr.GroupConcat)
	}

	view := NewViewFromGroupedRecord(f.Records[0])
	if expr.OrderBy != nil {
		err := view.OrderBy(expr.OrderBy.(parser.OrderByClause))
		if err != nil {
			return nil, err
		}
	}

	list := []string{}
	filter := NewFilterForSequentialEvaluation(view, f)
	for i := 0; i < view.RecordLen(); i++ {
		filter.Records[0].RecordIndex = i
		p, err := filter.Evaluate(arg)
		if err != nil {
			if _, ok := err.(*NotGroupingRecordsError); ok {
				err = NewNestedAggregateFunctionsError(expr)
			}
			return nil, err
		}
		s := parser.PrimaryToString(p)
		if parser.IsNull(s) {
			continue
		}
		if expr.Option.IsDistinct() && InStrSlice(s.(parser.String).Value(), list) {
			continue
		}
		list = append(list, s.(parser.String).Value())
	}

	if len(list) < 1 {
		return parser.NewNull(), nil
	}
	return parser.NewString(strings.Join(list, expr.Separator)), nil
}

func (f Filter) evalCase(expr parser.Case) (parser.Primary, error) {
	var value parser.Primary
	var err error
	if expr.Value != nil {
		value, err = f.Evaluate(expr.Value)
		if err != nil {
			return nil, err
		}
	}

	for _, v := range expr.When {
		when := v.(parser.CaseWhen)
		var t ternary.Value

		cond, err := f.Evaluate(when.Condition)
		if err != nil {
			return nil, err
		}

		if value == nil {
			t = cond.Ternary()
		} else {
			t = EqualTo(value, cond)
		}

		if t == ternary.TRUE {
			result, err := f.Evaluate(when.Result)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}

	if expr.Else == nil {
		return parser.NewNull(), nil
	}
	result, err := f.Evaluate(expr.Else.(parser.CaseElse).Result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f Filter) evalLogic(expr parser.Logic) (parser.Primary, error) {
	lhs, err := f.Evaluate(expr.LHS)
	if err != nil {
		return nil, err
	}
	rhs, err := f.Evaluate(expr.RHS)
	if err != nil {
		return nil, err
	}

	var t ternary.Value
	switch expr.Operator.Token {
	case parser.AND:
		t = ternary.And(lhs.Ternary(), rhs.Ternary())
	case parser.OR:
		t = ternary.Or(lhs.Ternary(), rhs.Ternary())
	}
	return parser.NewTernary(t), nil
}

func (f Filter) evalUnaryLogic(expr parser.UnaryLogic) (parser.Primary, error) {
	ope, err := f.Evaluate(expr.Operand)
	if err != nil {
		return nil, err
	}

	var t ternary.Value
	switch expr.Operator.Token {
	case parser.NOT, '!':
		t = ternary.Not(ope.Ternary())
	}
	return parser.NewTernary(t), nil
}

func (f Filter) evalCursorStatus(expr parser.CursorStatus) (parser.Primary, error) {
	var t ternary.Value
	var err error

	switch expr.Type {
	case parser.OPEN:
		t, err = f.CursorsList.IsOpen(expr.Cursor)
		if err != nil {
			return nil, err
		}
	case parser.RANGE:
		t, err = f.CursorsList.IsInRange(expr.Cursor)
		if err != nil {
			return nil, err
		}
	}

	if !expr.Negation.IsEmpty() {
		t = ternary.Not(t)
	}
	return parser.NewTernary(t), nil
}

func (f Filter) evalRowValue(expr parser.RowValue) (values []parser.Primary, err error) {
	switch expr.Value.(type) {
	case parser.Subquery:
		values, err = f.evalSubqueryForRowValue(expr.Value.(parser.Subquery))
	case parser.ValueList:
		values, err = f.evalValueList(expr.Value.(parser.ValueList))
	}

	return
}

func (f Filter) evalValues(exprs []parser.Expression) ([]parser.Primary, error) {
	values := make([]parser.Primary, len(exprs))
	for i, v := range exprs {
		value, err := f.Evaluate(v)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return values, nil
}

func (f Filter) evalValueList(expr parser.ValueList) ([]parser.Primary, error) {
	return f.evalValues(expr.Values)
}

func (f Filter) evalRowValueList(expr parser.RowValueList) ([][]parser.Primary, error) {
	list := make([][]parser.Primary, len(expr.RowValues))
	for i, v := range expr.RowValues {
		values, err := f.evalRowValue(v.(parser.RowValue))
		if err != nil {
			return nil, err
		}
		list[i] = values
	}
	return list, nil
}

func (f Filter) evalRowValues(expr parser.Expression) (values [][]parser.Primary, err error) {
	switch expr.(type) {
	case parser.Subquery:
		values, err = f.evalSubqueryForRowValues(expr.(parser.Subquery))
	case parser.RowValueList:
		values, err = f.evalRowValueList(expr.(parser.RowValueList))
	}

	return
}

func (f Filter) evalSubqueryForRowValue(expr parser.Subquery) ([]parser.Primary, error) {
	view, err := Select(expr.Query, f)
	if err != nil {
		return nil, err
	}

	if view.RecordLen() < 1 {
		return nil, nil
	}

	if 1 < view.RecordLen() {
		return nil, NewSubqueryTooManyRecordsError(expr)
	}

	values := make([]parser.Primary, view.FieldLen())
	for i, cell := range view.Records[0] {
		values[i] = cell.Primary()
	}

	return values, nil
}

func (f Filter) evalSubqueryForRowValues(expr parser.Subquery) ([][]parser.Primary, error) {
	view, err := Select(expr.Query, f)
	if err != nil {
		return nil, err
	}

	if view.RecordLen() < 1 {
		return nil, nil
	}

	list := make([][]parser.Primary, view.RecordLen())
	for i, r := range view.Records {
		values := make([]parser.Primary, view.FieldLen())
		for j, cell := range r {
			values[j] = cell.Primary()
		}
		list[i] = values
	}

	return list, nil
}

func (f Filter) evalSubqueryForSingleFieldRowValues(expr parser.Subquery) ([][]parser.Primary, error) {
	view, err := Select(expr.Query, f)
	if err != nil {
		return nil, err
	}

	if 1 < view.FieldLen() {
		return nil, NewSubqueryTooManyFieldsError(expr)
	}

	if view.RecordLen() < 1 {
		return nil, nil
	}

	list := make([][]parser.Primary, view.RecordLen())
	for i, r := range view.Records {
		list[i] = []parser.Primary{r[0].Primary()}
	}

	return list, nil
}

func (f Filter) evalSubqueryForSingleValue(expr parser.Subquery) (parser.Primary, error) {
	view, err := Select(expr.Query, f)
	if err != nil {
		return nil, err
	}

	if 1 < view.FieldLen() {
		return nil, NewSubqueryTooManyFieldsError(expr)
	}

	if 1 < view.RecordLen() {
		return nil, NewSubqueryTooManyRecordsError(expr)
	}

	if view.RecordLen() < 1 {
		return parser.NewNull(), nil
	}

	return view.Records[0][0].Primary(), nil
}
