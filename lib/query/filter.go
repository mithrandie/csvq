package query

import (
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
	"github.com/mithrandie/csvq/lib/value"
)

type FilterRecord struct {
	View        *View
	RecordIndex int

	fieldReferenceIndices map[string]int
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

	Now time.Time
}

func NewFilter(variablesList VariablesList, tempViewsList TemporaryViewMapList, cursorsList CursorMapList, functionsList UserDefinedFunctionsList) *Filter {
	return &Filter{
		VariablesList: variablesList,
		TempViewsList: tempViewsList,
		CursorsList:   cursorsList,
		FunctionsList: functionsList,
	}
}

func NewEmptyFilter() *Filter {
	return NewFilter(
		VariablesList{{}},
		TemporaryViewMapList{{}},
		CursorMapList{{}},
		UserDefinedFunctionsList{{}},
	)
}

func NewFilterForRecord(view *View, recordIndex int, parentFilter *Filter) *Filter {
	f := &Filter{
		Records: []FilterRecord{
			{
				View:                  view,
				RecordIndex:           recordIndex,
				fieldReferenceIndices: make(map[string]int),
			},
		},
	}
	f.Merge(parentFilter)
	return f
}

func NewFilterForSequentialEvaluation(view *View, parentFilter *Filter) *Filter {
	f := &Filter{
		Records: []FilterRecord{
			{
				View: view,
				fieldReferenceIndices: make(map[string]int),
			},
		},
	}
	f.Merge(parentFilter)
	return f
}

func (f *Filter) Merge(filter *Filter) {
	f.Records = append(f.Records, filter.Records...)
	f.VariablesList = filter.VariablesList
	f.TempViewsList = filter.TempViewsList
	f.CursorsList = filter.CursorsList
	f.FunctionsList = filter.FunctionsList
	f.InlineTablesList = filter.InlineTablesList
	f.AliasesList = filter.AliasesList
	f.Now = filter.Now
}

func (f *Filter) CreateChildScope() *Filter {
	return NewFilter(
		append(VariablesList{{}}, f.VariablesList...),
		append(TemporaryViewMapList{{}}, f.TempViewsList...),
		append(CursorMapList{{}}, f.CursorsList...),
		append(UserDefinedFunctionsList{{}}, f.FunctionsList...),
	)
}

func (f *Filter) CreateNode() *Filter {
	filter := &Filter{
		Records:          f.Records,
		VariablesList:    f.VariablesList,
		TempViewsList:    f.TempViewsList,
		CursorsList:      f.CursorsList,
		FunctionsList:    f.FunctionsList,
		InlineTablesList: append(InlineTablesList{{}}, f.InlineTablesList...),
		AliasesList:      append(AliasMapList{{}}, f.AliasesList...),
		RecursiveTable:   f.RecursiveTable,
		RecursiveTmpView: f.RecursiveTmpView,
		Now:              f.Now,
	}

	if filter.Now.IsZero() {
		filter.Now = cmd.Now()
	}

	return filter
}

func (f *Filter) LoadInlineTable(clause parser.WithClause) error {
	return f.InlineTablesList.Load(clause, f)
}

func (f *Filter) Evaluate(expr parser.QueryExpression) (value.Primary, error) {
	if expr == nil {
		return value.NewTernary(ternary.TRUE), nil
	}

	var val value.Primary
	var err error

	switch expr.(type) {
	case parser.PrimitiveType:
		return expr.(parser.PrimitiveType).Value, nil
	case parser.Parentheses:
		val, err = f.Evaluate(expr.(parser.Parentheses).Expr)
	case parser.FieldReference, parser.ColumnNumber:
		val, err = f.evalFieldReference(expr)
	case parser.Arithmetic:
		val, err = f.evalArithmetic(expr.(parser.Arithmetic))
	case parser.UnaryArithmetic:
		val, err = f.evalUnaryArithmetic(expr.(parser.UnaryArithmetic))
	case parser.Concat:
		val, err = f.evalConcat(expr.(parser.Concat))
	case parser.Comparison:
		val, err = f.evalComparison(expr.(parser.Comparison))
	case parser.Is:
		val, err = f.evalIs(expr.(parser.Is))
	case parser.Between:
		val, err = f.evalBetween(expr.(parser.Between))
	case parser.Like:
		val, err = f.evalLike(expr.(parser.Like))
	case parser.In:
		val, err = f.evalIn(expr.(parser.In))
	case parser.Any:
		val, err = f.evalAny(expr.(parser.Any))
	case parser.All:
		val, err = f.evalAll(expr.(parser.All))
	case parser.Exists:
		val, err = f.evalExists(expr.(parser.Exists))
	case parser.Subquery:
		val, err = f.evalSubqueryForSingleValue(expr.(parser.Subquery))
	case parser.Function:
		val, err = f.evalFunction(expr.(parser.Function))
	case parser.AggregateFunction:
		val, err = f.evalAggregateFunction(expr.(parser.AggregateFunction))
	case parser.ListAgg:
		val, err = f.evalListAgg(expr.(parser.ListAgg))
	case parser.CaseExpr:
		val, err = f.evalCaseExpr(expr.(parser.CaseExpr))
	case parser.Logic:
		val, err = f.evalLogic(expr.(parser.Logic))
	case parser.UnaryLogic:
		val, err = f.evalUnaryLogic(expr.(parser.UnaryLogic))
	case parser.Variable:
		val, err = f.VariablesList.Get(expr.(parser.Variable))
	case parser.VariableSubstitution:
		val, err = f.VariablesList.Substitute(expr.(parser.VariableSubstitution), f)
	case parser.CursorStatus:
		val, err = f.evalCursorStatus(expr.(parser.CursorStatus))
	case parser.CursorAttrebute:
		val, err = f.evalCursorAttribute(expr.(parser.CursorAttrebute))
	default:
		return nil, NewSyntaxErrorFromExpr(expr)
	}

	return val, err
}

func (f *Filter) evalFieldReference(expr parser.QueryExpression) (value.Primary, error) {
	exprStr := expr.String()

	var p value.Primary
	for _, v := range f.Records {
		if v.fieldReferenceIndices != nil {
			if idx, ok := v.fieldReferenceIndices[exprStr]; ok {
				p = v.View.Records[v.RecordIndex][idx].Value()
				break
			}
		}

		idx, err := v.View.FieldIndex(expr)
		if err == nil {
			if v.View.isGrouped && v.View.Header[idx].IsFromTable && !v.View.Header[idx].IsGroupKey {
				return nil, NewFieldNotGroupKeyError(expr)
			}
			p = v.View.Records[v.RecordIndex][idx].Value()
			if v.fieldReferenceIndices != nil {
				v.fieldReferenceIndices[exprStr] = idx
			}
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

func (f *Filter) evalArithmetic(expr parser.Arithmetic) (value.Primary, error) {
	lhs, err := f.Evaluate(expr.LHS)
	if err != nil {
		return nil, err
	}
	if value.IsNull(lhs) {
		return value.NewNull(), nil
	}

	rhs, err := f.Evaluate(expr.RHS)
	if err != nil {
		return nil, err
	}

	return Calculate(lhs, rhs, expr.Operator), nil
}

func (f *Filter) evalUnaryArithmetic(expr parser.UnaryArithmetic) (value.Primary, error) {
	ope, err := f.Evaluate(expr.Operand)
	if err != nil {
		return nil, err
	}

	if pi := value.ToInteger(ope); !value.IsNull(pi) {
		val := pi.(value.Integer).Raw()
		switch expr.Operator.Token {
		case '-':
			val = val * -1
		}
		return value.NewInteger(val), nil
	}

	pf := value.ToFloat(ope)
	if value.IsNull(pf) {
		return value.NewNull(), nil
	}

	val := pf.(value.Float).Raw()

	switch expr.Operator.Token {
	case '-':
		val = val * -1
	}

	return value.ParseFloat64(val), nil
}

func (f *Filter) evalConcat(expr parser.Concat) (value.Primary, error) {
	items := make([]string, len(expr.Items))
	for i, v := range expr.Items {
		s, err := f.Evaluate(v)
		if err != nil {
			return nil, err
		}
		s = value.ToString(s)
		if value.IsNull(s) {
			return value.NewNull(), nil
		}
		items[i] = s.(value.String).Raw()
	}
	return value.NewString(strings.Join(items, "")), nil
}

func (f *Filter) evalComparison(expr parser.Comparison) (value.Primary, error) {
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

		t, err = value.CompareRowValues(lhs, rhs, expr.Operator)
		if err != nil {
			return nil, NewRowValueLengthInComparisonError(expr.RHS.(parser.RowValue), len(lhs))
		}

	default:
		lhs, err := f.Evaluate(expr.LHS)
		if err != nil {
			return nil, err
		}
		if value.IsNull(lhs) {
			t = ternary.UNKNOWN
		} else {
			rhs, err := f.Evaluate(expr.RHS)
			if err != nil {
				return nil, err
			}

			t = value.Compare(lhs, rhs, expr.Operator)
		}
	}
	return value.NewTernary(t), nil
}

func (f *Filter) evalIs(expr parser.Is) (value.Primary, error) {
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
	return value.NewTernary(t), nil
}

func (f *Filter) evalBetween(expr parser.Between) (value.Primary, error) {
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
		lowResult, err := value.CompareRowValues(lhs, low, ">=")
		if err != nil {
			return nil, NewRowValueLengthInComparisonError(expr.Low.(parser.RowValue), len(lhs))
		}

		if lowResult == ternary.FALSE {
			t = ternary.FALSE
		} else {
			high, err := f.evalRowValue(expr.High.(parser.RowValue))
			if err != nil {
				return nil, err
			}

			highResult, err := value.CompareRowValues(lhs, high, "<=")
			if err != nil {
				return nil, NewRowValueLengthInComparisonError(expr.High.(parser.RowValue), len(lhs))
			}

			t = ternary.And(lowResult, highResult)
		}
	default:
		lhs, err := f.Evaluate(expr.LHS)
		if err != nil {
			return nil, err
		}
		if value.IsNull(lhs) {
			t = ternary.UNKNOWN
		} else {
			low, err := f.Evaluate(expr.Low)
			if err != nil {
				return nil, err
			}

			lowResult := value.GreaterOrEqual(lhs, low)
			if lowResult == ternary.FALSE {
				t = ternary.FALSE
			} else {
				high, err := f.Evaluate(expr.High)
				if err != nil {
					return nil, err
				}

				highResult := value.LessOrEqual(lhs, high)
				t = ternary.And(lowResult, highResult)
			}
		}
	}

	if expr.IsNegated() {
		t = ternary.Not(t)
	}
	return value.NewTernary(t), nil
}

func (f *Filter) valuesForRowValueListComparison(lhs parser.QueryExpression, values parser.QueryExpression) (value.RowValue, []value.RowValue, error) {
	var rowValue value.RowValue
	var list []value.RowValue
	var err error

	switch lhs.(type) {
	case parser.RowValue:
		rowValue, err = f.evalRowValue(lhs.(parser.RowValue))
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
		rowValue = value.RowValue{lhs}

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
			list = make([]value.RowValue, len(values))
			for i, v := range values {
				list[i] = value.RowValue{v}
			}
		}
	}
	return rowValue, list, nil
}

func (f *Filter) evalIn(expr parser.In) (value.Primary, error) {
	val, list, err := f.valuesForRowValueListComparison(expr.LHS, expr.Values)
	if err != nil {
		return nil, err
	}

	t, err := Any(val, list, "=")
	if err != nil {
		if subquery, ok := expr.Values.(parser.Subquery); ok {
			return nil, NewSelectFieldLengthInComparisonError(subquery, len(val))
		}

		rvlist, _ := expr.Values.(parser.RowValueList)
		rverr, _ := err.(*RowValueLengthInListError)
		return nil, NewRowValueLengthInComparisonError(rvlist.RowValues[rverr.Index], len(val))
	}

	if expr.IsNegated() {
		t = ternary.Not(t)
	}
	return value.NewTernary(t), nil
}

func (f *Filter) evalAny(expr parser.Any) (value.Primary, error) {
	val, list, err := f.valuesForRowValueListComparison(expr.LHS, expr.Values)
	if err != nil {
		return nil, err
	}

	t, err := Any(val, list, expr.Operator)
	if err != nil {
		if subquery, ok := expr.Values.(parser.Subquery); ok {
			return nil, NewSelectFieldLengthInComparisonError(subquery, len(val))
		}

		rvlist, _ := expr.Values.(parser.RowValueList)
		rverr, _ := err.(*RowValueLengthInListError)
		return nil, NewRowValueLengthInComparisonError(rvlist.RowValues[rverr.Index], len(val))
	}
	return value.NewTernary(t), nil
}

func (f *Filter) evalAll(expr parser.All) (value.Primary, error) {
	val, list, err := f.valuesForRowValueListComparison(expr.LHS, expr.Values)
	if err != nil {
		return nil, err
	}

	t, err := All(val, list, expr.Operator)
	if err != nil {
		if subquery, ok := expr.Values.(parser.Subquery); ok {
			return nil, NewSelectFieldLengthInComparisonError(subquery, len(val))
		}

		rvlist, _ := expr.Values.(parser.RowValueList)
		rverr, _ := err.(*RowValueLengthInListError)
		return nil, NewRowValueLengthInComparisonError(rvlist.RowValues[rverr.Index], len(val))
	}
	return value.NewTernary(t), nil
}

func (f *Filter) evalLike(expr parser.Like) (value.Primary, error) {
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
	return value.NewTernary(t), nil
}

func (f *Filter) evalExists(expr parser.Exists) (value.Primary, error) {
	view, err := Select(expr.Query.Query, f)
	if err != nil {
		return nil, err
	}
	if view.RecordLen() < 1 {
		return value.NewTernary(ternary.FALSE), nil
	}
	return value.NewTernary(ternary.TRUE), nil
}

func (f *Filter) evalFunction(expr parser.Function) (value.Primary, error) {
	name := strings.ToUpper(expr.Name)

	if _, ok := Functions[name]; !ok && name != "NOW" {
		udfn, err := f.FunctionsList.Get(expr, name)
		if err != nil {
			return nil, NewFunctionNotExistError(expr, expr.Name)
		}
		if udfn.IsAggregate {
			aggrdcl := parser.AggregateFunction{
				BaseExpr: expr.BaseExpr,
				Name:     expr.Name,
				Args:     expr.Args,
			}
			return f.evalAggregateFunction(aggrdcl)
		}

		if err = udfn.CheckArgsLen(expr, expr.Name, len(expr.Args)); err != nil {
			return nil, err
		}
	}

	args := make([]value.Primary, len(expr.Args))
	for i, v := range expr.Args {
		arg, err := f.Evaluate(v)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}

	if name == "NOW" {
		return Now(expr, args, f)
	}

	if fn, ok := Functions[name]; ok {
		return fn(expr, args)
	}

	udfn, _ := f.FunctionsList.Get(expr, name)
	return udfn.Execute(args, f)
}

func (f *Filter) evalAggregateFunction(expr parser.AggregateFunction) (value.Primary, error) {
	var aggfn func([]value.Primary) value.Primary
	var udfn *UserDefinedFunction
	var useUserDefined bool
	var err error

	uname := strings.ToUpper(expr.Name)
	if fn, ok := AggregateFunctions[uname]; ok {
		aggfn = fn
	} else {
		if udfn, err = f.FunctionsList.Get(expr, uname); err != nil || !udfn.IsAggregate {
			return nil, NewFunctionNotExistError(expr, expr.Name)
		}
		useUserDefined = true
	}

	if useUserDefined {
		if err = udfn.CheckArgsLen(expr, expr.Name, len(expr.Args)-1); err != nil {
			return nil, err
		}
	} else {
		if len(expr.Args) != 1 {
			return nil, NewFunctionArgumentLengthError(expr, expr.Name, []int{1})
		}
	}

	if len(f.Records) < 1 {
		return nil, NewUnpermittedStatementFunctionError(expr, expr.Name)
	}

	if !f.Records[0].View.isGrouped {
		return nil, NewNotGroupingRecordsError(expr, expr.Name)
	}

	listExpr := expr.Args[0]
	if _, ok := listExpr.(parser.AllColumns); ok {
		listExpr = parser.NewIntegerValue(1)
	}

	if uname == "COUNT" {
		if _, ok := listExpr.(parser.PrimitiveType); ok {
			return value.NewInteger(int64(f.Records[0].View.Records[f.Records[0].RecordIndex].GroupLen())), nil
		}
	}

	view := NewViewFromGroupedRecord(f.Records[0])
	list, err := view.ListValuesForAggregateFunctions(expr, listExpr, expr.IsDistinct(), f)
	if err != nil {
		return nil, err
	}

	if useUserDefined {
		argsExprs := expr.Args[1:]
		args := make([]value.Primary, len(argsExprs))
		for i, v := range argsExprs {
			arg, err := f.Evaluate(v)
			if err != nil {
				return nil, err
			}
			args[i] = arg
		}
		return udfn.ExecuteAggregate(list, args, f)
	}

	return aggfn(list), nil
}

func (f *Filter) evalListAgg(expr parser.ListAgg) (value.Primary, error) {
	if expr.Args == nil || 2 < len(expr.Args) {
		return nil, NewFunctionArgumentLengthError(expr, expr.ListAgg, []int{1, 2})
	}

	if len(f.Records) < 1 {
		return nil, NewUnpermittedStatementFunctionError(expr, expr.ListAgg)
	}

	if !f.Records[0].View.isGrouped {
		return nil, NewNotGroupingRecordsError(expr, expr.ListAgg)
	}

	separator := ""
	if len(expr.Args) == 2 {
		p, err := f.Evaluate(expr.Args[1])
		if err != nil {
			return nil, NewFunctionInvalidArgumentError(expr, expr.ListAgg, "the second argument must be a string")
		}
		s := value.ToString(p)
		if value.IsNull(s) {
			return nil, NewFunctionInvalidArgumentError(expr, expr.ListAgg, "the second argument must be a string")
		}
		separator = s.(value.String).Raw()
	}

	view := NewViewFromGroupedRecord(f.Records[0])
	if expr.OrderBy != nil {
		err := view.OrderBy(expr.OrderBy.(parser.OrderByClause))
		if err != nil {
			return nil, err
		}
	}

	list, err := view.ListValuesForAggregateFunctions(expr, expr.Args[0], expr.IsDistinct(), f)
	if err != nil {
		return nil, err
	}

	return ListAgg(list, separator), nil
}

func (f *Filter) evalCaseExpr(expr parser.CaseExpr) (value.Primary, error) {
	var val value.Primary
	var err error
	if expr.Value != nil {
		val, err = f.Evaluate(expr.Value)
		if err != nil {
			return nil, err
		}
	}

	for _, v := range expr.When {
		when := v.(parser.CaseExprWhen)
		var t ternary.Value

		cond, err := f.Evaluate(when.Condition)
		if err != nil {
			return nil, err
		}

		if val == nil {
			t = cond.Ternary()
		} else {
			t = value.Equal(val, cond)
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
		return value.NewNull(), nil
	}
	result, err := f.Evaluate(expr.Else.(parser.CaseExprElse).Result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f *Filter) evalLogic(expr parser.Logic) (value.Primary, error) {
	lhs, err := f.Evaluate(expr.LHS)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Token {
	case parser.AND:
		if lhs.Ternary() == ternary.FALSE {
			return value.NewTernary(ternary.FALSE), nil
		}
	case parser.OR:
		if lhs.Ternary() == ternary.TRUE {
			return value.NewTernary(ternary.TRUE), nil
		}
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
	return value.NewTernary(t), nil
}

func (f *Filter) evalUnaryLogic(expr parser.UnaryLogic) (value.Primary, error) {
	ope, err := f.Evaluate(expr.Operand)
	if err != nil {
		return nil, err
	}

	var t ternary.Value
	switch expr.Operator.Token {
	case parser.NOT, '!':
		t = ternary.Not(ope.Ternary())
	}
	return value.NewTernary(t), nil
}

func (f *Filter) evalCursorStatus(expr parser.CursorStatus) (value.Primary, error) {
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
	return value.NewTernary(t), nil
}

func (f *Filter) evalCursorAttribute(expr parser.CursorAttrebute) (value.Primary, error) {
	var i int
	var err error

	switch expr.Attrebute.Token {
	case parser.COUNT:
		i, err = f.CursorsList.Count(expr.Cursor)
		if err != nil {
			return nil, err
		}
	}
	return value.NewInteger(int64(i)), nil
}

func (f *Filter) evalRowValue(expr parser.RowValue) (values []value.Primary, err error) {
	switch expr.Value.(type) {
	case parser.Subquery:
		values, err = f.evalSubqueryForRowValue(expr.Value.(parser.Subquery))
	case parser.ValueList:
		values, err = f.evalValueList(expr.Value.(parser.ValueList))
	}

	return
}

func (f *Filter) evalValues(exprs []parser.QueryExpression) ([]value.Primary, error) {
	values := make([]value.Primary, len(exprs))
	for i, v := range exprs {
		val, err := f.Evaluate(v)
		if err != nil {
			return nil, err
		}
		values[i] = val
	}
	return values, nil
}

func (f *Filter) evalValueList(expr parser.ValueList) ([]value.Primary, error) {
	return f.evalValues(expr.Values)
}

func (f *Filter) evalRowValueList(expr parser.RowValueList) ([]value.RowValue, error) {
	list := make([]value.RowValue, len(expr.RowValues))
	for i, v := range expr.RowValues {
		rowValue, err := f.evalRowValue(v.(parser.RowValue))
		if err != nil {
			return nil, err
		}
		list[i] = rowValue
	}
	return list, nil
}

func (f *Filter) evalRowValues(expr parser.QueryExpression) (rowValues []value.RowValue, err error) {
	switch expr.(type) {
	case parser.Subquery:
		rowValues, err = f.evalSubqueryForRowValues(expr.(parser.Subquery))
	case parser.RowValueList:
		rowValues, err = f.evalRowValueList(expr.(parser.RowValueList))
	}

	return
}

func (f *Filter) evalSubqueryForRowValue(expr parser.Subquery) (value.RowValue, error) {
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

	rowValue := make(value.RowValue, view.FieldLen())
	for i, cell := range view.Records[0] {
		rowValue[i] = cell.Value()
	}

	return rowValue, nil
}

func (f *Filter) evalSubqueryForRowValues(expr parser.Subquery) ([]value.RowValue, error) {
	view, err := Select(expr.Query, f)
	if err != nil {
		return nil, err
	}

	if view.RecordLen() < 1 {
		return nil, nil
	}

	list := make([]value.RowValue, view.RecordLen())
	for i, r := range view.Records {
		rowValue := make(value.RowValue, view.FieldLen())
		for j, cell := range r {
			rowValue[j] = cell.Value()
		}
		list[i] = rowValue
	}

	return list, nil
}

func (f *Filter) evalSubqueryForSingleFieldRowValues(expr parser.Subquery) ([]value.RowValue, error) {
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

	list := make([]value.RowValue, view.RecordLen())
	for i, r := range view.Records {
		list[i] = value.RowValue{r[0].Value()}
	}

	return list, nil
}

func (f *Filter) evalSubqueryForSingleValue(expr parser.Subquery) (value.Primary, error) {
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
		return value.NewNull(), nil
	}

	return view.Records[0][0].Value(), nil
}
