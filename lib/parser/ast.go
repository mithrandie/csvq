package parser

import (
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/ternary"
)

const TokenUndefined = 0

type Statement interface{}

type Expression interface {
	GetBaseExpr() *BaseExpr
	HasParseInfo() bool
	Line() int
	Char() int
	SourceFile() string
}

type QueryExpression interface {
	String() string

	GetBaseExpr() *BaseExpr
	HasParseInfo() bool
	Line() int
	Char() int
	SourceFile() string
}

type BaseExpr struct {
	line       int
	char       int
	sourceFile string
}

func (e *BaseExpr) Line() int {
	return e.line
}

func (e *BaseExpr) Char() int {
	return e.char
}

func (e *BaseExpr) SourceFile() string {
	return e.sourceFile
}

func (e *BaseExpr) HasParseInfo() bool {
	if e == nil {
		return false
	}
	return true
}

func (e *BaseExpr) GetBaseExpr() *BaseExpr {
	return e
}

func NewBaseExpr(token Token) *BaseExpr {
	return &BaseExpr{
		line:       token.Line,
		char:       token.Char,
		sourceFile: token.SourceFile,
	}
}

type PrimitiveType struct {
	*BaseExpr
	Literal string
	Value   value.Primary
}

func NewStringValue(s string) PrimitiveType {
	return PrimitiveType{
		Literal: s,
		Value:   value.NewString(s),
	}
}

func NewIntegerValueFromString(s string) PrimitiveType {
	return PrimitiveType{
		Literal: s,
		Value:   value.NewIntegerFromString(s),
	}
}

func NewIntegerValue(i int64) PrimitiveType {
	return PrimitiveType{
		Value: value.NewInteger(i),
	}
}

func NewFloatValueFromString(s string) PrimitiveType {
	return PrimitiveType{
		Literal: s,
		Value:   value.NewFloatFromString(s),
	}
}

func NewFloatValue(f float64) PrimitiveType {
	return PrimitiveType{
		Value: value.NewFloat(f),
	}
}

func NewTernaryValueFromString(s string) PrimitiveType {
	return PrimitiveType{
		Literal: s,
		Value:   value.NewTernaryFromString(s),
	}
}

func NewTernaryValue(t ternary.Value) PrimitiveType {
	return PrimitiveType{
		Value: value.NewTernary(t),
	}
}

func NewDatetimeValueFromString(s string, formats []string) PrimitiveType {
	return PrimitiveType{
		Literal: s,
		Value:   value.NewDatetimeFromString(s, formats),
	}
}

func NewDatetimeValue(t time.Time) PrimitiveType {
	return PrimitiveType{
		Value: value.NewDatetime(t),
	}
}

func NewNullValueFromString(s string) PrimitiveType {
	return PrimitiveType{
		Literal: s,
		Value:   value.NewNull(),
	}
}

func NewNullValue() PrimitiveType {
	return PrimitiveType{
		Value: value.NewNull(),
	}
}

func (e PrimitiveType) String() string {
	if 0 < len(e.Literal) {
		switch e.Value.(type) {
		case value.String, value.Datetime:
			return quoteString(e.Literal)
		default:
			return e.Literal
		}
	}
	return e.Value.String()
}

func (e PrimitiveType) IsInteger() bool {
	_, ok := e.Value.(value.Integer)
	return ok
}

type Identifier struct {
	*BaseExpr
	Literal string
	Quoted  bool
}

func (i Identifier) String() string {
	if i.Quoted {
		return cmd.QuoteIdentifier(i.Literal)
	}
	return i.Literal
}

type FieldReference struct {
	*BaseExpr
	View   Identifier
	Column Identifier
}

func (e FieldReference) String() string {
	s := e.Column.String()
	if 0 < len(e.View.Literal) {
		s = e.View.String() + "." + s
	}
	return s
}

type ColumnNumber struct {
	*BaseExpr
	View   Identifier
	Number value.Integer
}

func (e ColumnNumber) String() string {
	return e.View.String() + "." + e.Number.String()
}

type Parentheses struct {
	*BaseExpr
	Expr QueryExpression
}

func (p Parentheses) String() string {
	return putParentheses(p.Expr.String())
}

type RowValue struct {
	*BaseExpr
	Value QueryExpression
}

func (e RowValue) String() string {
	return e.Value.String()
}

type ValueList struct {
	*BaseExpr
	Values []QueryExpression
}

func (e ValueList) String() string {
	return putParentheses(listQueryExpressions(e.Values))
}

type RowValueList struct {
	*BaseExpr
	RowValues []QueryExpression
}

func (e RowValueList) String() string {
	return putParentheses(listQueryExpressions(e.RowValues))
}

type SelectQuery struct {
	*BaseExpr
	WithClause    QueryExpression
	SelectEntity  QueryExpression
	OrderByClause QueryExpression
	LimitClause   QueryExpression
	OffsetClause  QueryExpression
}

func (e SelectQuery) String() string {
	s := make([]string, 0)
	if e.WithClause != nil {
		s = append(s, e.WithClause.String())
	}
	s = append(s, e.SelectEntity.String())
	if e.OrderByClause != nil {
		s = append(s, e.OrderByClause.String())
	}
	if e.LimitClause != nil {
		s = append(s, e.LimitClause.String())
	}
	if e.OffsetClause != nil {
		s = append(s, e.OffsetClause.String())
	}
	return joinWithSpace(s)
}

type SelectSet struct {
	*BaseExpr
	LHS      QueryExpression
	Operator Token
	All      Token
	RHS      QueryExpression
}

func (e SelectSet) String() string {
	s := []string{e.LHS.String(), e.Operator.Literal}
	if !e.All.IsEmpty() {
		s = append(s, e.All.Literal)
	}
	s = append(s, e.RHS.String())
	return joinWithSpace(s)
}

type SelectEntity struct {
	*BaseExpr
	SelectClause  QueryExpression
	FromClause    QueryExpression
	WhereClause   QueryExpression
	GroupByClause QueryExpression
	HavingClause  QueryExpression
}

func (e SelectEntity) String() string {
	s := []string{e.SelectClause.String()}
	if e.FromClause != nil {
		s = append(s, e.FromClause.String())
	}
	if e.WhereClause != nil {
		s = append(s, e.WhereClause.String())
	}
	if e.GroupByClause != nil {
		s = append(s, e.GroupByClause.String())
	}
	if e.HavingClause != nil {
		s = append(s, e.HavingClause.String())
	}
	return joinWithSpace(s)
}

type SelectClause struct {
	*BaseExpr
	Select   string
	Distinct Token
	Fields   []QueryExpression
}

func (sc SelectClause) IsDistinct() bool {
	return !sc.Distinct.IsEmpty()
}

func (sc SelectClause) String() string {
	s := []string{sc.Select}
	if sc.IsDistinct() {
		s = append(s, sc.Distinct.Literal)
	}
	s = append(s, listQueryExpressions(sc.Fields))
	return joinWithSpace(s)
}

type FromClause struct {
	*BaseExpr
	From   string
	Tables []QueryExpression
}

func (f FromClause) String() string {
	s := []string{f.From, listQueryExpressions(f.Tables)}
	return joinWithSpace(s)
}

type WhereClause struct {
	*BaseExpr
	Where  string
	Filter QueryExpression
}

func (w WhereClause) String() string {
	s := []string{w.Where, w.Filter.String()}
	return joinWithSpace(s)
}

type GroupByClause struct {
	*BaseExpr
	GroupBy string
	Items   []QueryExpression
}

func (gb GroupByClause) String() string {
	s := []string{gb.GroupBy, listQueryExpressions(gb.Items)}
	return joinWithSpace(s)
}

type HavingClause struct {
	*BaseExpr
	Having string
	Filter QueryExpression
}

func (h HavingClause) String() string {
	s := []string{h.Having, h.Filter.String()}
	return joinWithSpace(s)
}

type OrderByClause struct {
	*BaseExpr
	OrderBy string
	Items   []QueryExpression
}

func (ob OrderByClause) String() string {
	s := []string{ob.OrderBy, listQueryExpressions(ob.Items)}
	return joinWithSpace(s)
}

type LimitClause struct {
	*BaseExpr
	Limit   string
	Value   QueryExpression
	Percent string
	With    QueryExpression
}

func (e LimitClause) String() string {
	s := []string{e.Limit, e.Value.String()}
	if e.IsPercentage() {
		s = append(s, e.Percent)
	}
	if e.With != nil {
		s = append(s, e.With.String())
	}
	return joinWithSpace(s)
}

func (e LimitClause) IsPercentage() bool {
	return 0 < len(e.Percent)
}

func (e LimitClause) IsWithTies() bool {
	if e.With == nil {
		return false
	}
	return e.With.(LimitWith).Type.Token == TIES
}

type LimitWith struct {
	*BaseExpr
	With string
	Type Token
}

func (e LimitWith) String() string {
	s := []string{e.With, e.Type.Literal}
	return joinWithSpace(s)
}

type OffsetClause struct {
	*BaseExpr
	Offset string
	Value  QueryExpression
}

func (e OffsetClause) String() string {
	s := []string{e.Offset, e.Value.String()}
	return joinWithSpace(s)
}

type WithClause struct {
	*BaseExpr
	With         string
	InlineTables []QueryExpression
}

func (e WithClause) String() string {
	s := []string{e.With, listQueryExpressions(e.InlineTables)}
	return joinWithSpace(s)
}

type InlineTable struct {
	*BaseExpr
	Recursive Token
	Name      Identifier
	Fields    []QueryExpression
	As        string
	Query     SelectQuery
}

func (e InlineTable) String() string {
	s := make([]string, 0)
	if !e.Recursive.IsEmpty() {
		s = append(s, e.Recursive.Literal)
	}
	s = append(s, e.Name.String())
	if e.Fields != nil {
		s = append(s, putParentheses(listQueryExpressions(e.Fields)))
	}
	s = append(s, e.As, putParentheses(e.Query.String()))
	return joinWithSpace(s)
}

func (e InlineTable) IsRecursive() bool {
	return !e.Recursive.IsEmpty()
}

type Subquery struct {
	*BaseExpr
	Query SelectQuery
}

func (sq Subquery) String() string {
	return putParentheses(sq.Query.String())
}

type TableObject struct {
	*BaseExpr
	Type          Identifier
	FormatElement QueryExpression
	Path          Identifier
	Args          []QueryExpression
}

func (e TableObject) String() string {
	allArgs := make([]QueryExpression, 0, len(e.Args)+2)
	if e.FormatElement != nil {
		allArgs = append(allArgs, e.FormatElement)
	}
	allArgs = append(allArgs, e.Path)
	if e.Args != nil {
		allArgs = append(allArgs, e.Args...)
	}
	return e.Type.String() + putParentheses(listQueryExpressions(allArgs))
}

type JsonQuery struct {
	*BaseExpr
	JsonQuery string
	Query     QueryExpression
	JsonText  QueryExpression
}

func (e JsonQuery) String() string {
	return e.JsonQuery + putParentheses(e.Query.String()+", "+e.JsonText.String())
}

type Comparison struct {
	*BaseExpr
	LHS      QueryExpression
	Operator string
	RHS      QueryExpression
}

func (c Comparison) String() string {
	s := []string{c.LHS.String(), c.Operator, c.RHS.String()}
	return joinWithSpace(s)
}

type Is struct {
	*BaseExpr
	Is       string
	LHS      QueryExpression
	RHS      QueryExpression
	Negation Token
}

func (i Is) IsNegated() bool {
	return !i.Negation.IsEmpty()
}

func (i Is) String() string {
	s := []string{i.LHS.String(), i.Is}
	if i.IsNegated() {
		s = append(s, i.Negation.Literal)
	}
	s = append(s, i.RHS.String())
	return joinWithSpace(s)
}

type Between struct {
	*BaseExpr
	Between  string
	And      string
	LHS      QueryExpression
	Low      QueryExpression
	High     QueryExpression
	Negation Token
}

func (b Between) IsNegated() bool {
	return !b.Negation.IsEmpty()
}

func (b Between) String() string {
	s := []string{b.LHS.String()}
	if b.IsNegated() {
		s = append(s, b.Negation.Literal)
	}
	s = append(s, b.Between, b.Low.String(), b.And, b.High.String())
	return joinWithSpace(s)
}

type In struct {
	*BaseExpr
	In       string
	LHS      QueryExpression
	Values   QueryExpression
	Negation Token
}

func (i In) IsNegated() bool {
	return !i.Negation.IsEmpty()
}

func (i In) String() string {
	s := []string{i.LHS.String()}
	if i.IsNegated() {
		s = append(s, i.Negation.Literal)
	}
	s = append(s, i.In, i.Values.String())
	return joinWithSpace(s)
}

type All struct {
	*BaseExpr
	All      string
	LHS      QueryExpression
	Operator string
	Values   QueryExpression
}

func (a All) String() string {
	s := []string{a.LHS.String(), a.Operator, a.All, a.Values.String()}
	return joinWithSpace(s)
}

type Any struct {
	*BaseExpr
	Any      string
	LHS      QueryExpression
	Operator string
	Values   QueryExpression
}

func (a Any) String() string {
	s := []string{a.LHS.String(), a.Operator, a.Any, a.Values.String()}
	return joinWithSpace(s)
}

type Like struct {
	*BaseExpr
	Like     string
	LHS      QueryExpression
	Pattern  QueryExpression
	Negation Token
}

func (l Like) IsNegated() bool {
	return !l.Negation.IsEmpty()
}

func (l Like) String() string {
	s := []string{l.LHS.String()}
	if l.IsNegated() {
		s = append(s, l.Negation.Literal)
	}
	s = append(s, l.Like, l.Pattern.String())
	return joinWithSpace(s)
}

type Exists struct {
	*BaseExpr
	Exists string
	Query  Subquery
}

func (e Exists) String() string {
	s := []string{e.Exists, e.Query.String()}
	return joinWithSpace(s)
}

type Arithmetic struct {
	*BaseExpr
	LHS      QueryExpression
	Operator int
	RHS      QueryExpression
}

func (a Arithmetic) String() string {
	s := []string{a.LHS.String(), string(rune(a.Operator)), a.RHS.String()}
	return joinWithSpace(s)
}

type UnaryArithmetic struct {
	*BaseExpr
	Operand  QueryExpression
	Operator Token
}

func (e UnaryArithmetic) String() string {
	return e.Operator.Literal + e.Operand.String()
}

type Logic struct {
	*BaseExpr
	LHS      QueryExpression
	Operator Token
	RHS      QueryExpression
}

func (l Logic) String() string {
	s := []string{l.LHS.String(), l.Operator.Literal, l.RHS.String()}
	return joinWithSpace(s)
}

type UnaryLogic struct {
	*BaseExpr
	Operand  QueryExpression
	Operator Token
}

func (e UnaryLogic) String() string {
	if e.Operator.Token == NOT {
		s := []string{e.Operator.Literal, e.Operand.String()}
		return joinWithSpace(s)
	}
	return e.Operator.Literal + e.Operand.String()
}

type Concat struct {
	*BaseExpr
	Items []QueryExpression
}

func (c Concat) String() string {
	s := make([]string, len(c.Items))
	for i, v := range c.Items {
		s[i] = v.String()
	}
	return strings.Join(s, " || ")
}

type Function struct {
	*BaseExpr
	Name string
	Args []QueryExpression
}

func (e Function) String() string {
	return e.Name + "(" + listQueryExpressions(e.Args) + ")"
}

type AggregateFunction struct {
	*BaseExpr
	Name     string
	Distinct Token
	Args     []QueryExpression
}

func (e AggregateFunction) String() string {
	s := make([]string, 0)
	if !e.Distinct.IsEmpty() {
		s = append(s, e.Distinct.Literal)
	}
	s = append(s, listQueryExpressions(e.Args))

	return e.Name + "(" + joinWithSpace(s) + ")"
}

func (e AggregateFunction) IsDistinct() bool {
	return !e.Distinct.IsEmpty()
}

type Table struct {
	*BaseExpr
	Object QueryExpression
	As     string
	Alias  QueryExpression
}

func (t Table) String() string {
	s := []string{t.Object.String()}
	if 0 < len(t.As) {
		s = append(s, t.As)
	}
	if t.Alias != nil {
		s = append(s, t.Alias.String())
	}
	return joinWithSpace(s)
}

func (t Table) Name() Identifier {
	if t.Alias != nil {
		return t.Alias.(Identifier)
	}

	if file, ok := t.Object.(Identifier); ok {
		return Identifier{
			BaseExpr: file.BaseExpr,
			Literal:  FormatTableName(file.Literal),
		}
	}

	return Identifier{
		BaseExpr: t.Object.GetBaseExpr(),
		Literal:  t.Object.String(),
	}

}

type Join struct {
	*BaseExpr
	Join      string
	Table     QueryExpression
	JoinTable QueryExpression
	Natural   Token
	JoinType  Token
	Direction Token
	Condition QueryExpression
}

func (j Join) String() string {
	s := []string{j.Table.String()}
	if !j.Natural.IsEmpty() {
		s = append(s, j.Natural.Literal)
	}
	if !j.Direction.IsEmpty() {
		s = append(s, j.Direction.Literal)
	}
	if !j.JoinType.IsEmpty() {
		s = append(s, j.JoinType.Literal)
	}
	s = append(s, j.Join, j.JoinTable.String())
	if j.Condition != nil {
		s = append(s, j.Condition.String())
	}
	return joinWithSpace(s)
}

type JoinCondition struct {
	*BaseExpr
	Literal string
	On      QueryExpression
	Using   []QueryExpression
}

func (jc JoinCondition) String() string {
	var s []string
	if jc.On != nil {
		s = []string{jc.Literal, jc.On.String()}
	} else {
		s = []string{jc.Literal, putParentheses(listQueryExpressions(jc.Using))}
	}
	return joinWithSpace(s)
}

type Field struct {
	*BaseExpr
	Object QueryExpression
	As     string
	Alias  QueryExpression
}

func (f Field) String() string {
	s := []string{f.Object.String()}
	if 0 < len(f.As) {
		s = append(s, f.As)
	}
	if f.Alias != nil {
		s = append(s, f.Alias.String())
	}
	return joinWithSpace(s)
}

func (f Field) Name() string {
	if f.Alias != nil {
		return f.Alias.(Identifier).Literal
	}
	if t, ok := f.Object.(PrimitiveType); ok {
		return t.Literal
	}
	if fr, ok := f.Object.(FieldReference); ok {
		return fr.Column.Literal
	}
	return f.Object.String()
}

type AllColumns struct {
	*BaseExpr
}

func (ac AllColumns) String() string {
	return "*"
}

type Dual struct {
	*BaseExpr
	Dual string
}

func (d Dual) String() string {
	return d.Dual
}

type Stdin struct {
	*BaseExpr
	Stdin string
}

func (si Stdin) String() string {
	return si.Stdin
}

type OrderItem struct {
	*BaseExpr
	Value     QueryExpression
	Direction Token
	Nulls     string
	Position  Token
}

func (e OrderItem) String() string {
	s := []string{e.Value.String()}
	if !e.Direction.IsEmpty() {
		s = append(s, e.Direction.Literal)
	}
	if 0 < len(e.Nulls) {
		s = append(s, e.Nulls, e.Position.Literal)
	}
	return joinWithSpace(s)
}

type CaseExpr struct {
	*BaseExpr
	Case  string
	End   string
	Value QueryExpression
	When  []QueryExpression
	Else  QueryExpression
}

func (e CaseExpr) String() string {
	s := []string{e.Case}
	if e.Value != nil {
		s = append(s, e.Value.String())
	}
	for _, v := range e.When {
		s = append(s, v.String())
	}
	if e.Else != nil {
		s = append(s, e.Else.String())
	}
	s = append(s, e.End)
	return joinWithSpace(s)
}

type CaseExprWhen struct {
	*BaseExpr
	When      string
	Then      string
	Condition QueryExpression
	Result    QueryExpression
}

func (e CaseExprWhen) String() string {
	s := []string{e.When, e.Condition.String(), e.Then, e.Result.String()}
	return joinWithSpace(s)
}

type CaseExprElse struct {
	*BaseExpr
	Else   string
	Result QueryExpression
}

func (e CaseExprElse) String() string {
	s := []string{e.Else, e.Result.String()}
	return joinWithSpace(s)
}

type ListFunction struct {
	*BaseExpr
	Name        string
	Distinct    Token
	Args        []QueryExpression
	WithinGroup string
	OrderBy     QueryExpression
}

func (e ListFunction) String() string {
	option := make([]string, 0)
	if !e.Distinct.IsEmpty() {
		option = append(option, e.Distinct.Literal)
	}
	option = append(option, listQueryExpressions(e.Args))

	s := []string{e.Name + "(" + joinWithSpace(option) + ")"}
	if 0 < len(e.WithinGroup) {
		s = append(s, e.WithinGroup)
		if e.OrderBy != nil {
			s = append(s, "("+e.OrderBy.String()+")")
		} else {
			s = append(s, "()")
		}
	}
	return joinWithSpace(s)
}

func (e ListFunction) IsDistinct() bool {
	return !e.Distinct.IsEmpty()
}

type AnalyticFunction struct {
	*BaseExpr
	Name           string
	Distinct       Token
	Args           []QueryExpression
	IgnoreNulls    bool
	IgnoreNullsLit string
	Over           string
	AnalyticClause AnalyticClause
}

func (e AnalyticFunction) String() string {
	option := make([]string, 0)
	if !e.Distinct.IsEmpty() {
		option = append(option, e.Distinct.Literal)
	}
	if e.Args != nil {
		option = append(option, listQueryExpressions(e.Args))
	}
	if e.IgnoreNulls {
		option = append(option, e.IgnoreNullsLit)
	}

	s := []string{
		e.Name + "(" + joinWithSpace(option) + ")",
		e.Over,
		"(" + e.AnalyticClause.String() + ")",
	}
	return joinWithSpace(s)
}

func (e AnalyticFunction) IsDistinct() bool {
	return !e.Distinct.IsEmpty()
}

type AnalyticClause struct {
	*BaseExpr
	PartitionClause QueryExpression
	OrderByClause   QueryExpression
	WindowingClause QueryExpression
}

func (e AnalyticClause) String() string {
	s := make([]string, 0)
	if e.PartitionClause != nil {
		s = append(s, e.PartitionClause.String())
	}
	if e.OrderByClause != nil {
		s = append(s, e.OrderByClause.String())
	}
	if e.WindowingClause != nil {
		s = append(s, e.WindowingClause.String())
	}
	return joinWithSpace(s)
}

func (e AnalyticClause) PartitionValues() []QueryExpression {
	if e.PartitionClause == nil {
		return nil
	}
	return e.PartitionClause.(PartitionClause).Values
}

type PartitionClause struct {
	*BaseExpr
	PartitionBy string
	Values      []QueryExpression
}

func (e PartitionClause) String() string {
	s := []string{e.PartitionBy, listQueryExpressions(e.Values)}
	return joinWithSpace(s)
}

type WindowingClause struct {
	*BaseExpr
	Rows      string
	FrameLow  QueryExpression
	FrameHigh QueryExpression
	Between   string
	And       string
}

func (e WindowingClause) String() string {
	s := []string{e.Rows}
	if e.FrameHigh == nil {
		s = append(s, e.FrameLow.String())
	} else {
		s = append(s, e.Between, e.FrameLow.String(), e.And, e.FrameHigh.String())
	}
	return joinWithSpace(s)
}

type WindowFramePosition struct {
	*BaseExpr
	Direction int
	Unbounded bool
	Offset    int
	Literal   string
}

func (e WindowFramePosition) String() string {
	return e.Literal
}

type Variable struct {
	*BaseExpr
	Name string
}

func (v Variable) String() string {
	return string(VariableSign) + v.Name
}

type VariableSubstitution struct {
	*BaseExpr
	Variable Variable
	Value    QueryExpression
}

func (vs VariableSubstitution) String() string {
	return joinWithSpace([]string{vs.Variable.String(), SubstitutionOperator, vs.Value.String()})
}

type VariableAssignment struct {
	*BaseExpr
	Variable Variable
	Value    QueryExpression
}

type VariableDeclaration struct {
	*BaseExpr
	Assignments []VariableAssignment
}

type DisposeVariable struct {
	*BaseExpr
	Variable Variable
}

type EnvironmentVariable struct {
	*BaseExpr
	Name   string
	Quoted bool
}

func (e EnvironmentVariable) String() string {
	name := e.Name
	if e.Quoted {
		name = cmd.QuoteIdentifier(name)
	}

	return string(VariableSign) + string(EnvironmentVariableSign) + name
}

type RuntimeInformation struct {
	*BaseExpr
	Name string
}

func (e RuntimeInformation) String() string {
	return string(VariableSign) + string(RuntimeInformationSign) + e.Name
}

type SetEnvVar struct {
	*BaseExpr
	EnvVar EnvironmentVariable
	Value  QueryExpression
}

type UnsetEnvVar struct {
	*BaseExpr
	EnvVar EnvironmentVariable
}

type InsertQuery struct {
	*BaseExpr
	WithClause QueryExpression
	Table      Table
	Fields     []QueryExpression
	ValuesList []QueryExpression
	Query      QueryExpression
}

type UpdateQuery struct {
	*BaseExpr
	WithClause  QueryExpression
	Tables      []QueryExpression
	SetList     []UpdateSet
	FromClause  QueryExpression
	WhereClause QueryExpression
}

type UpdateSet struct {
	*BaseExpr
	Field QueryExpression
	Value QueryExpression
}

type DeleteQuery struct {
	*BaseExpr
	WithClause  QueryExpression
	Tables      []QueryExpression
	FromClause  FromClause
	WhereClause QueryExpression
}

type CreateTable struct {
	*BaseExpr
	Table  Identifier
	Fields []QueryExpression
	Query  QueryExpression
}

type AddColumns struct {
	*BaseExpr
	Table    QueryExpression
	Columns  []ColumnDefault
	Position Expression
}

type ColumnDefault struct {
	*BaseExpr
	Column Identifier
	Value  QueryExpression
}

type ColumnPosition struct {
	*BaseExpr
	Position Token
	Column   QueryExpression
}

type DropColumns struct {
	*BaseExpr
	Table   QueryExpression
	Columns []QueryExpression
}

type RenameColumn struct {
	*BaseExpr
	Table QueryExpression
	Old   QueryExpression
	New   Identifier
}

type SetTableAttribute struct {
	*BaseExpr
	Table     QueryExpression
	Attribute Identifier
	Value     QueryExpression
}

type FunctionDeclaration struct {
	*BaseExpr
	Name       Identifier
	Parameters []VariableAssignment
	Statements []Statement
}

type AggregateDeclaration struct {
	*BaseExpr
	Name       Identifier
	Cursor     Identifier
	Parameters []VariableAssignment
	Statements []Statement
}

type DisposeFunction struct {
	*BaseExpr
	Name Identifier
}

type Return struct {
	*BaseExpr
	Value QueryExpression
}

type Echo struct {
	*BaseExpr
	Value QueryExpression
}

type Print struct {
	*BaseExpr
	Value QueryExpression
}

type Printf struct {
	*BaseExpr
	Format QueryExpression
	Values []QueryExpression
}

type Source struct {
	*BaseExpr
	FilePath QueryExpression
}

type Chdir struct {
	*BaseExpr
	DirPath QueryExpression
}

type Pwd struct {
	*BaseExpr
}

type Reload struct {
	*BaseExpr
	Type Identifier
}

type Execute struct {
	*BaseExpr
	Statements QueryExpression
	Values     []QueryExpression
}

type Syntax struct {
	*BaseExpr
	Keywords []QueryExpression
}

type SetFlag struct {
	*BaseExpr
	Name  string
	Value QueryExpression
}

type AddFlagElement struct {
	*BaseExpr
	Name  string
	Value QueryExpression
}

type RemoveFlagElement struct {
	*BaseExpr
	Name  string
	Value QueryExpression
}

type ShowFlag struct {
	*BaseExpr
	Name string
}

type ShowObjects struct {
	*BaseExpr
	Type Identifier
}

type ShowFields struct {
	*BaseExpr
	Type  Identifier
	Table QueryExpression
}

type If struct {
	*BaseExpr
	Condition  QueryExpression
	Statements []Statement
	ElseIf     []ElseIf
	Else       Else
}

type ElseIf struct {
	*BaseExpr
	Condition  QueryExpression
	Statements []Statement
}

type Else struct {
	*BaseExpr
	Statements []Statement
}

type Case struct {
	*BaseExpr
	Value QueryExpression
	When  []CaseWhen
	Else  CaseElse
}

type CaseWhen struct {
	*BaseExpr
	Condition  QueryExpression
	Statements []Statement
}

type CaseElse struct {
	*BaseExpr
	Statements []Statement
}

type While struct {
	*BaseExpr
	Condition  QueryExpression
	Statements []Statement
}

type WhileInCursor struct {
	*BaseExpr
	WithDeclaration bool
	Variables       []Variable
	Cursor          Identifier
	Statements      []Statement
}

type CursorDeclaration struct {
	*BaseExpr
	Cursor Identifier
	Query  SelectQuery
}

type OpenCursor struct {
	*BaseExpr
	Cursor Identifier
}

type CloseCursor struct {
	*BaseExpr
	Cursor Identifier
}

type DisposeCursor struct {
	*BaseExpr
	Cursor Identifier
}

type FetchCursor struct {
	*BaseExpr
	Position  FetchPosition
	Cursor    Identifier
	Variables []Variable
}

type FetchPosition struct {
	*BaseExpr
	Position Token
	Number   QueryExpression
}

type CursorStatus struct {
	*BaseExpr
	CursorLit string
	Cursor    Identifier
	Is        string
	Negation  Token
	Type      int
	TypeLit   string
}

func (e CursorStatus) String() string {
	s := []string{e.CursorLit, e.Cursor.String(), e.Is}
	if !e.Negation.IsEmpty() {
		s = append(s, e.Negation.Literal)
	}
	s = append(s, e.TypeLit)
	return joinWithSpace(s)
}

type CursorAttrebute struct {
	*BaseExpr
	CursorLit string
	Cursor    Identifier
	Attrebute Token
}

func (e CursorAttrebute) String() string {
	s := []string{e.CursorLit, e.Cursor.String(), e.Attrebute.Literal}
	return joinWithSpace(s)
}

type ViewDeclaration struct {
	*BaseExpr
	View   Identifier
	Fields []QueryExpression
	Query  QueryExpression
}

type DisposeView struct {
	*BaseExpr
	View Identifier
}

type TransactionControl struct {
	*BaseExpr
	Token int
}

type FlowControl struct {
	*BaseExpr
	Token int
}

type Trigger struct {
	*BaseExpr
	Event   Identifier
	Message QueryExpression
	Code    value.Primary
}

type Exit struct {
	*BaseExpr
	Code value.Primary
}

type ExternalCommand struct {
	*BaseExpr
	Command string
}

func putParentheses(s string) string {
	return "(" + s + ")"
}

func joinWithSpace(s []string) string {
	return strings.Join(s, " ")
}

func listQueryExpressions(exprs []QueryExpression) string {
	s := make([]string, len(exprs))
	for i, v := range exprs {
		s[i] = v.String()
	}
	return strings.Join(s, ", ")
}

func quoteString(s string) string {
	return "'" + s + "'"
}

func quoteIdentifier(s string) string {
	return "`" + s + "`"
}
