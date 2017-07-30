package parser

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/ternary"
)

const TOKEN_UNDEFINED = 0
const DEFAULT_DATETIME_FORMAT = "2006-01-02 15:04:05.999999999"

func IsPrimary(e Expression) bool {
	if e == nil {
		return false
	}

	t := reflect.TypeOf(e)
	v := reflect.TypeOf((*Primary)(nil)).Elem()
	return t.Implements(v)
}

func IsNull(v Primary) bool {
	_, ok := v.(Null)
	return ok
}

type Statement interface{}

type ProcExpr interface {
	GetBaseExpr() *BaseExpr
	HasParseInfo() bool
	Line() int
	Char() int
	SourceFile() string
}

type Expression interface {
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

type Primary interface {
	String() string
	Ternary() ternary.Value

	GetBaseExpr() *BaseExpr
	HasParseInfo() bool
	Line() int
	Char() int
	SourceFile() string
}

type String struct {
	*BaseExpr
	literal string
}

func (s String) String() string {
	return quoteString(s.literal)
}

func NewString(s string) String {
	return String{
		literal: s,
	}
}

func (s String) Value() string {
	return s.literal
}

func (s String) Ternary() ternary.Value {
	if b, err := strconv.ParseBool(s.Value()); err == nil {
		return ternary.ParseBool(b)
	}
	return ternary.UNKNOWN
}

type Integer struct {
	*BaseExpr
	literal string
	value   int64
}

func NewIntegerFromString(s string) Integer {
	i, _ := strconv.ParseInt(s, 10, 64)
	return Integer{
		literal: s,
		value:   i,
	}
}

func NewInteger(i int64) Integer {
	return Integer{
		literal: Int64ToStr(i),
		value:   i,
	}
}

func (i Integer) String() string {
	return i.literal
}

func (i Integer) Value() int64 {
	return i.value
}

func (i Integer) Ternary() ternary.Value {
	switch i.Value() {
	case 0:
		return ternary.FALSE
	case 1:
		return ternary.TRUE
	default:
		return ternary.UNKNOWN
	}
}

type Float struct {
	*BaseExpr
	literal string
	value   float64
}

func NewFloatFromString(s string) Float {
	f, _ := strconv.ParseFloat(s, 64)
	return Float{
		literal: s,
		value:   f,
	}
}

func NewFloat(f float64) Float {
	return Float{
		literal: Float64ToStr(f),
		value:   f,
	}
}

func (f Float) String() string {
	return f.literal
}

func (f Float) Value() float64 {
	return f.value
}

func (f Float) Ternary() ternary.Value {
	switch f.Value() {
	case 0:
		return ternary.FALSE
	case 1:
		return ternary.TRUE
	default:
		return ternary.UNKNOWN
	}
}

type Boolean struct {
	*BaseExpr
	literal string
	value   bool
}

func NewBoolean(b bool) Boolean {
	return Boolean{
		literal: strconv.FormatBool(b),
		value:   b,
	}
}

func (b Boolean) String() string {
	return b.literal
}

func (b Boolean) Value() bool {
	return b.value
}

func (b Boolean) Ternary() ternary.Value {
	return ternary.ParseBool(b.Value())
}

type Ternary struct {
	*BaseExpr
	literal string
	value   ternary.Value
}

func NewTernaryFromString(s string) Ternary {
	t, _ := ternary.Parse(s)
	return Ternary{
		literal: s,
		value:   t,
	}
}

func NewTernary(t ternary.Value) Ternary {
	return Ternary{
		literal: t.String(),
		value:   t,
	}
}

func (t Ternary) String() string {
	return t.literal
}

func (t Ternary) Ternary() ternary.Value {
	return t.value
}

type Datetime struct {
	*BaseExpr
	literal string
	value   time.Time
	format  string
}

func NewDatetimeFromString(s string) Datetime {
	t, _ := StrToTime(s)
	return Datetime{
		literal: s,
		value:   t,
		format:  DEFAULT_DATETIME_FORMAT,
	}
}

func NewDatetime(t time.Time) Datetime {
	return Datetime{
		literal: t.Format(time.RFC3339Nano),
		value:   t,
		format:  DEFAULT_DATETIME_FORMAT,
	}
}

func (dt Datetime) String() string {
	return quoteString(dt.literal)
}

func (dt Datetime) Value() time.Time {
	return dt.value
}

func (dt Datetime) Ternary() ternary.Value {
	return ternary.UNKNOWN
}

func (dt Datetime) Format() string {
	return dt.value.Format(dt.format)
}

func (dt *Datetime) SetFormat(format string) {
	dt.format = ConvertDatetimeFormat(format)
}

type Null struct {
	*BaseExpr
	literal string
}

func NewNullFromString(s string) Null {
	return Null{
		literal: s,
	}
}

func NewNull() Null {
	return Null{}
}

func (n Null) String() string {
	if len(n.literal) < 1 {
		return "NULL"
	}
	return n.literal
}

func (n Null) Ternary() ternary.Value {
	return ternary.UNKNOWN
}

type Identifier struct {
	*BaseExpr
	Literal string
	Quoted  bool
}

func (i Identifier) String() string {
	if i.Quoted {
		return quoteIdentifier(i.Literal)
	}
	return i.Literal
}

type FieldReference struct {
	*BaseExpr
	View   Expression
	Column Identifier
}

func (e FieldReference) String() string {
	s := e.Column.String()
	if e.View != nil {
		s = e.View.(Identifier).String() + "." + s
	}
	return s
}

type ColumnNumber struct {
	*BaseExpr
	View   Identifier
	Number Integer
}

func (e ColumnNumber) String() string {
	return e.View.String() + "." + e.Number.String()
}

type Parentheses struct {
	*BaseExpr
	Expr Expression
}

func (p Parentheses) String() string {
	return putParentheses(p.Expr.String())
}

type RowValue struct {
	*BaseExpr
	Value Expression
}

func (e RowValue) String() string {
	return e.Value.String()
}

type ValueList struct {
	*BaseExpr
	Values []Expression
}

func (e ValueList) String() string {
	return putParentheses(listExpressions(e.Values))
}

type RowValueList struct {
	*BaseExpr
	RowValues []Expression
}

func (e RowValueList) String() string {
	return putParentheses(listExpressions(e.RowValues))
}

type SelectQuery struct {
	*BaseExpr
	WithClause    Expression
	SelectEntity  Expression
	OrderByClause Expression
	LimitClause   Expression
	OffsetClause  Expression
}

func (e SelectQuery) String() string {
	s := []string{}
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
	LHS      Expression
	Operator Token
	All      Token
	RHS      Expression
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
	SelectClause  Expression
	FromClause    Expression
	WhereClause   Expression
	GroupByClause Expression
	HavingClause  Expression
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
	Fields   []Expression
}

func (sc SelectClause) IsDistinct() bool {
	return !sc.Distinct.IsEmpty()
}

func (sc SelectClause) String() string {
	s := []string{sc.Select}
	if sc.IsDistinct() {
		s = append(s, sc.Distinct.Literal)
	}
	s = append(s, listExpressions(sc.Fields))
	return joinWithSpace(s)
}

type FromClause struct {
	*BaseExpr
	From   string
	Tables []Expression
}

func (f FromClause) String() string {
	s := []string{f.From, listExpressions(f.Tables)}
	return joinWithSpace(s)
}

type WhereClause struct {
	*BaseExpr
	Where  string
	Filter Expression
}

func (w WhereClause) String() string {
	s := []string{w.Where, w.Filter.String()}
	return joinWithSpace(s)
}

type GroupByClause struct {
	*BaseExpr
	GroupBy string
	Items   []Expression
}

func (gb GroupByClause) String() string {
	s := []string{gb.GroupBy, listExpressions(gb.Items)}
	return joinWithSpace(s)
}

type HavingClause struct {
	*BaseExpr
	Having string
	Filter Expression
}

func (h HavingClause) String() string {
	s := []string{h.Having, h.Filter.String()}
	return joinWithSpace(s)
}

type OrderByClause struct {
	*BaseExpr
	OrderBy string
	Items   []Expression
}

func (ob OrderByClause) String() string {
	s := []string{ob.OrderBy, listExpressions(ob.Items)}
	return joinWithSpace(s)
}

type LimitClause struct {
	*BaseExpr
	Limit   string
	Value   Expression
	Percent string
	With    Expression
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
	Value  Expression
}

func (e OffsetClause) String() string {
	s := []string{e.Offset, e.Value.String()}
	return joinWithSpace(s)
}

type WithClause struct {
	*BaseExpr
	With         string
	InlineTables []Expression
}

func (e WithClause) String() string {
	s := []string{e.With, listExpressions(e.InlineTables)}
	return joinWithSpace(s)
}

type InlineTable struct {
	*BaseExpr
	Recursive Token
	Name      Identifier
	Fields    []Expression
	As        string
	Query     SelectQuery
}

func (e InlineTable) String() string {
	s := []string{}
	if !e.Recursive.IsEmpty() {
		s = append(s, e.Recursive.Literal)
	}
	s = append(s, e.Name.String())
	if e.Fields != nil {
		s = append(s, putParentheses(listExpressions(e.Fields)))
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

type Comparison struct {
	*BaseExpr
	LHS      Expression
	Operator string
	RHS      Expression
}

func (c Comparison) String() string {
	s := []string{c.LHS.String(), c.Operator, c.RHS.String()}
	return joinWithSpace(s)
}

type Is struct {
	*BaseExpr
	Is       string
	LHS      Expression
	RHS      Expression
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
	LHS      Expression
	Low      Expression
	High     Expression
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
	LHS      Expression
	Values   Expression
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
	LHS      Expression
	Operator string
	Values   Expression
}

func (a All) String() string {
	s := []string{a.LHS.String(), a.Operator, a.All, a.Values.String()}
	return joinWithSpace(s)
}

type Any struct {
	*BaseExpr
	Any      string
	LHS      Expression
	Operator string
	Values   Expression
}

func (a Any) String() string {
	s := []string{a.LHS.String(), a.Operator, a.Any, a.Values.String()}
	return joinWithSpace(s)
}

type Like struct {
	*BaseExpr
	Like     string
	LHS      Expression
	Pattern  Expression
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
	LHS      Expression
	Operator int
	RHS      Expression
}

func (a Arithmetic) String() string {
	s := []string{a.LHS.String(), string(rune(a.Operator)), a.RHS.String()}
	return joinWithSpace(s)
}

type UnaryArithmetic struct {
	*BaseExpr
	Operand  Expression
	Operator Token
}

func (e UnaryArithmetic) String() string {
	return e.Operator.Literal + e.Operand.String()
}

type Logic struct {
	*BaseExpr
	LHS      Expression
	Operator Token
	RHS      Expression
}

func (l Logic) String() string {
	s := []string{l.LHS.String(), l.Operator.Literal, l.RHS.String()}
	return joinWithSpace(s)
}

type UnaryLogic struct {
	*BaseExpr
	Operand  Expression
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
	Items []Expression
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
	Args []Expression
}

func (e Function) String() string {
	return e.Name + "(" + listExpressions(e.Args) + ")"
}

type AggregateFunction struct {
	*BaseExpr
	Name     string
	Distinct Token
	Args     []Expression
}

func (e AggregateFunction) String() string {
	s := []string{}
	if !e.Distinct.IsEmpty() {
		s = append(s, e.Distinct.Literal)
	}
	s = append(s, listExpressions(e.Args))

	return e.Name + "(" + joinWithSpace(s) + ")"
}

func (e AggregateFunction) IsDistinct() bool {
	return !e.Distinct.IsEmpty()
}

type Table struct {
	*BaseExpr
	Object Expression
	As     string
	Alias  Expression
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
	Table     Table
	JoinTable Table
	Natural   Token
	JoinType  Token
	Direction Token
	Condition Expression
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
	On      Expression
	Using   []Expression
}

func (jc JoinCondition) String() string {
	var s []string
	if jc.On != nil {
		s = []string{jc.Literal, jc.On.String()}
	} else {
		s = []string{jc.Literal, putParentheses(listExpressions(jc.Using))}
	}
	return joinWithSpace(s)
}

type Field struct {
	*BaseExpr
	Object Expression
	As     string
	Alias  Expression
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
	if s, ok := f.Object.(String); ok {
		return s.Value()
	}
	if dt, ok := f.Object.(Datetime); ok {
		return dt.literal
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
	Value     Expression
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

type Case struct {
	*BaseExpr
	Case  string
	End   string
	Value Expression
	When  []Expression
	Else  Expression
}

func (c Case) String() string {
	s := []string{c.Case}
	if c.Value != nil {
		s = append(s, c.Value.String())
	}
	for _, v := range c.When {
		s = append(s, v.String())
	}
	if c.Else != nil {
		s = append(s, c.Else.String())
	}
	s = append(s, c.End)
	return joinWithSpace(s)
}

type CaseWhen struct {
	*BaseExpr
	When      string
	Then      string
	Condition Expression
	Result    Expression
}

func (cw CaseWhen) String() string {
	s := []string{cw.When, cw.Condition.String(), cw.Then, cw.Result.String()}
	return joinWithSpace(s)
}

type CaseElse struct {
	*BaseExpr
	Else   string
	Result Expression
}

func (ce CaseElse) String() string {
	s := []string{ce.Else, ce.Result.String()}
	return joinWithSpace(s)
}

type ListAgg struct {
	*BaseExpr
	ListAgg     string
	Distinct    Token
	Args        []Expression
	WithinGroup string
	OrderBy     Expression
}

func (e ListAgg) String() string {
	option := []string{}
	if !e.Distinct.IsEmpty() {
		option = append(option, e.Distinct.Literal)
	}
	option = append(option, listExpressions(e.Args))

	s := []string{e.ListAgg + "(" + joinWithSpace(option) + ")"}
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

func (e ListAgg) IsDistinct() bool {
	return !e.Distinct.IsEmpty()
}

type AnalyticFunction struct {
	*BaseExpr
	Name           string
	Distinct       Token
	Args           []Expression
	IgnoreNulls    bool
	IgnoreNullsLit string
	Over           string
	AnalyticClause AnalyticClause
}

func (e AnalyticFunction) String() string {
	option := []string{}
	if !e.Distinct.IsEmpty() {
		option = append(option, e.Distinct.Literal)
	}
	if e.Args != nil {
		option = append(option, listExpressions(e.Args))
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
	Partition     Expression
	OrderByClause Expression
}

func (e AnalyticClause) String() string {
	s := []string{}
	if e.Partition != nil {
		s = append(s, e.Partition.String())
	}
	if e.OrderByClause != nil {
		s = append(s, e.OrderByClause.String())
	}
	return joinWithSpace(s)
}

func (e AnalyticClause) PartitionValues() []Expression {
	if e.Partition == nil {
		return nil
	}
	return e.Partition.(Partition).Values
}

func (e AnalyticClause) OrderValues() []Expression {
	if e.OrderByClause == nil {
		return nil
	}

	items := e.OrderByClause.(OrderByClause).Items
	result := make([]Expression, len(items))
	for i, v := range items {
		result[i] = v.(OrderItem).Value
	}
	return result
}

type Partition struct {
	*BaseExpr
	PartitionBy string
	Values      []Expression
}

func (e Partition) String() string {
	s := []string{e.PartitionBy, listExpressions(e.Values)}
	return joinWithSpace(s)
}

type Variable struct {
	*BaseExpr
	Name string
}

func (v Variable) String() string {
	return v.Name
}

type VariableSubstitution struct {
	*BaseExpr
	Variable Variable
	Value    Expression
}

func (vs VariableSubstitution) String() string {
	return joinWithSpace([]string{vs.Variable.String(), SUBSTITUTION_OPERATOR, vs.Value.String()})
}

type VariableAssignment struct {
	*BaseExpr
	Variable Variable
	Value    Expression
}

func (va VariableAssignment) String() string {
	if va.Value == nil {
		return va.Variable.String()
	}
	return joinWithSpace([]string{va.Variable.String(), SUBSTITUTION_OPERATOR, va.Value.String()})
}

type VariableDeclaration struct {
	*BaseExpr
	Assignments []Expression
}

type DisposeVariable struct {
	*BaseExpr
	Variable Variable
}

type InsertQuery struct {
	*BaseExpr
	WithClause Expression
	Insert     string
	Into       string
	Table      Expression
	Fields     []Expression
	Values     string
	ValuesList []Expression
	Query      Expression
}

func (e InsertQuery) String() string {
	s := []string{}
	if e.WithClause != nil {
		s = append(s, e.WithClause.String())
	}
	s = append(s, e.Insert, e.Into, e.Table.String())
	if e.Fields != nil {
		s = append(s, putParentheses(listExpressions(e.Fields)))
	}
	if e.ValuesList != nil {
		s = append(s, e.Values)
		s = append(s, listExpressions(e.ValuesList))
	} else {
		s = append(s, e.Query.String())
	}
	return joinWithSpace(s)
}

type UpdateQuery struct {
	*BaseExpr
	WithClause  Expression
	Update      string
	Tables      []Expression
	Set         string
	SetList     []Expression
	FromClause  Expression
	WhereClause Expression
}

func (e UpdateQuery) String() string {
	s := []string{}
	if e.WithClause != nil {
		s = append(s, e.WithClause.String())
	}
	s = append(s, e.Update, listExpressions(e.Tables), e.Set, listExpressions(e.SetList))
	if e.FromClause != nil {
		s = append(s, e.FromClause.String())
	}
	if e.WhereClause != nil {
		s = append(s, e.WhereClause.String())
	}
	return joinWithSpace(s)
}

type UpdateSet struct {
	*BaseExpr
	Field Expression
	Value Expression
}

func (us UpdateSet) String() string {
	return joinWithSpace([]string{us.Field.String(), "=", us.Value.String()})
}

type DeleteQuery struct {
	*BaseExpr
	WithClause  Expression
	Delete      string
	Tables      []Expression
	FromClause  Expression
	WhereClause Expression
}

func (e DeleteQuery) String() string {
	s := []string{}
	if e.WithClause != nil {
		s = append(s, e.WithClause.String())
	}
	s = append(s, e.Delete)
	if e.Tables != nil {
		s = append(s, listExpressions(e.Tables))
	}
	s = append(s, e.FromClause.String())
	if e.WhereClause != nil {
		s = append(s, e.WhereClause.String())
	}
	return joinWithSpace(s)
}

type CreateTable struct {
	*BaseExpr
	CreateTable string
	Table       Identifier
	Fields      []Expression
}

func (e CreateTable) String() string {
	s := []string{
		e.CreateTable,
		e.Table.String(),
		putParentheses(listExpressions(e.Fields)),
	}
	return joinWithSpace(s)
}

type AddColumns struct {
	*BaseExpr
	AlterTable string
	Table      Expression
	Add        string
	Columns    []Expression
	Position   Expression
}

func (e AddColumns) String() string {
	s := []string{
		e.AlterTable,
		e.Table.String(),
		e.Add,
		putParentheses(listExpressions(e.Columns)),
	}
	if e.Position != nil {
		s = append(s, e.Position.String())
	}
	return joinWithSpace(s)
}

type ColumnDefault struct {
	*BaseExpr
	Column  Identifier
	Default string
	Value   Expression
}

func (e ColumnDefault) String() string {
	s := []string{e.Column.String()}
	if e.Value != nil {
		s = append(s, e.Default, e.Value.String())
	}
	return joinWithSpace(s)
}

type ColumnPosition struct {
	*BaseExpr
	Position Token
	Column   Expression
}

func (e ColumnPosition) String() string {
	s := []string{e.Position.Literal}
	if e.Column != nil {
		s = append(s, e.Column.String())
	}
	return joinWithSpace(s)
}

type DropColumns struct {
	*BaseExpr
	AlterTable string
	Table      Expression
	Drop       string
	Columns    []Expression
}

func (e DropColumns) String() string {
	s := []string{
		e.AlterTable,
		e.Table.String(),
		e.Drop,
		putParentheses(listExpressions(e.Columns)),
	}
	return joinWithSpace(s)
}

type RenameColumn struct {
	*BaseExpr
	AlterTable string
	Table      Expression
	Rename     string
	Old        Expression
	To         string
	New        Identifier
}

func (e RenameColumn) String() string {
	s := []string{
		e.AlterTable,
		e.Table.String(),
		e.Rename,
		e.Old.String(),
		e.To,
		e.New.String(),
	}
	return joinWithSpace(s)
}

type FunctionDeclaration struct {
	*BaseExpr
	Name       Identifier
	Parameters []Variable
	Statements []Statement
}

type AggregateDeclaration struct {
	*BaseExpr
	Name       Identifier
	Parameter  Identifier
	Statements []Statement
}

type Return struct {
	*BaseExpr
	Value Expression
}

type Print struct {
	*BaseExpr
	Value Expression
}

type Printf struct {
	*BaseExpr
	Format string
	Values []Expression
}

type Source struct {
	*BaseExpr
	FilePath Expression
}

type SetFlag struct {
	*BaseExpr
	Name  string
	Value Primary
}

type If struct {
	*BaseExpr
	Condition  Expression
	Statements []Statement
	ElseIf     []ProcExpr
	Else       ProcExpr
}

type ElseIf struct {
	*BaseExpr
	Condition  Expression
	Statements []Statement
}

type Else struct {
	*BaseExpr
	Statements []Statement
}

type While struct {
	*BaseExpr
	Condition  Expression
	Statements []Statement
}

type WhileInCursor struct {
	*BaseExpr
	Variables  []Variable
	Cursor     Identifier
	Statements []Statement
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
	Position  Expression
	Cursor    Identifier
	Variables []Variable
}

type FetchPosition struct {
	*BaseExpr
	Position Token
	Number   Expression
}

func (e FetchPosition) String() string {
	s := []string{e.Position.Literal}
	if e.Number != nil {
		s = append(s, e.Number.String())
	}
	return joinWithSpace(s)
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

type TableDeclaration struct {
	*BaseExpr
	Table  Identifier
	Fields []Expression
	Query  Expression
}

type DisposeTable struct {
	*BaseExpr
	Table Identifier
}

type TransactionControl struct {
	*BaseExpr
	Token int
}

type FlowControl struct {
	*BaseExpr
	Token int
}

func putParentheses(s string) string {
	return "(" + s + ")"
}

func joinWithSpace(s []string) string {
	return strings.Join(s, " ")
}

func listExpressions(exprs []Expression) string {
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
