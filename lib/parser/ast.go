package parser

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/ternary"
)

const TOKEN_UNDEFINED = 0
const DATETIME_FORMAT = "2006-01-02 15:04:05.999999999"

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

type ProcExpr interface{}

type Expression interface {
	String() string
}

type Primary interface {
	String() string
	Ternary() ternary.Value
}

type String struct {
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
	literal string
	value   time.Time
	format  string
}

func NewDatetimeFromString(s string) Datetime {
	t, _ := StrToTime(s)
	return Datetime{
		literal: s,
		value:   t,
		format:  DATETIME_FORMAT,
	}
}

func NewDatetime(t time.Time) Datetime {
	return Datetime{
		literal: t.Format(time.RFC3339Nano),
		value:   t,
		format:  DATETIME_FORMAT,
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

type Parentheses struct {
	Expr Expression
}

func (p Parentheses) String() string {
	return putParentheses(p.Expr.String())
}

type RowValue struct {
	Value Expression
}

func (e RowValue) String() string {
	return e.Value.String()
}

type ValueList struct {
	Values []Expression
}

func (e ValueList) String() string {
	return putParentheses(listExpressions(e.Values))
}

type RowValueList struct {
	RowValues []Expression
}

func (e RowValueList) String() string {
	return putParentheses(listExpressions(e.RowValues))
}

type SelectQuery struct {
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
	From   string
	Tables []Expression
}

func (f FromClause) String() string {
	s := []string{f.From, listExpressions(f.Tables)}
	return joinWithSpace(s)
}

type WhereClause struct {
	Where  string
	Filter Expression
}

func (w WhereClause) String() string {
	s := []string{w.Where, w.Filter.String()}
	return joinWithSpace(s)
}

type GroupByClause struct {
	GroupBy string
	Items   []Expression
}

func (gb GroupByClause) String() string {
	s := []string{gb.GroupBy, listExpressions(gb.Items)}
	return joinWithSpace(s)
}

type HavingClause struct {
	Having string
	Filter Expression
}

func (h HavingClause) String() string {
	s := []string{h.Having, h.Filter.String()}
	return joinWithSpace(s)
}

type OrderByClause struct {
	OrderBy string
	Items   []Expression
}

func (ob OrderByClause) String() string {
	s := []string{ob.OrderBy, listExpressions(ob.Items)}
	return joinWithSpace(s)
}

type LimitClause struct {
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
	With string
	Type Token
}

func (e LimitWith) String() string {
	s := []string{e.With, e.Type.Literal}
	return joinWithSpace(s)
}

type OffsetClause struct {
	Offset string
	Value  Expression
}

func (e OffsetClause) String() string {
	s := []string{e.Offset, e.Value.String()}
	return joinWithSpace(s)
}

type WithClause struct {
	With         string
	InlineTables []Expression
}

func (e WithClause) String() string {
	s := []string{e.With, listExpressions(e.InlineTables)}
	return joinWithSpace(s)
}

type InlineTable struct {
	Recursive Token
	Name      Identifier
	Columns   []Expression
	As        string
	Query     SelectQuery
}

func (e InlineTable) String() string {
	s := []string{}
	if !e.Recursive.IsEmpty() {
		s = append(s, e.Recursive.Literal)
	}
	s = append(s, e.Name.String())
	if e.Columns != nil {
		s = append(s, putParentheses(listExpressions(e.Columns)))
	}
	s = append(s, e.As, putParentheses(e.Query.String()))
	return joinWithSpace(s)
}

func (e InlineTable) IsRecursive() bool {
	return !e.Recursive.IsEmpty()
}

type Subquery struct {
	Query SelectQuery
}

func (sq Subquery) String() string {
	return putParentheses(sq.Query.String())
}

type Comparison struct {
	LHS      Expression
	Operator string
	RHS      Expression
}

func (c Comparison) String() string {
	s := []string{c.LHS.String(), c.Operator, c.RHS.String()}
	return joinWithSpace(s)
}

type Is struct {
	Is       string
	LHS      Expression
	RHS      Expression
	Negation Token
}

func (i *Is) IsNegated() bool {
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
	Between  string
	And      string
	LHS      Expression
	Low      Expression
	High     Expression
	Negation Token
}

func (b *Between) IsNegated() bool {
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
	In       string
	LHS      Expression
	Values   Expression
	Negation Token
}

func (i *In) IsNegated() bool {
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
	Like     string
	LHS      Expression
	Pattern  Expression
	Negation Token
}

func (l *Like) IsNegated() bool {
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
	Exists string
	Query  Subquery
}

func (e Exists) String() string {
	s := []string{e.Exists, e.Query.String()}
	return joinWithSpace(s)
}

type Arithmetic struct {
	LHS      Expression
	Operator int
	RHS      Expression
}

func (a Arithmetic) String() string {
	s := []string{a.LHS.String(), string(rune(a.Operator)), a.RHS.String()}
	return joinWithSpace(s)
}

type Logic struct {
	LHS      Expression
	Operator Token
	RHS      Expression
}

func (l Logic) String() string {
	s := []string{}
	if l.LHS != nil {
		s = append(s, l.LHS.String())
	}
	s = append(s, l.Operator.Literal, l.RHS.String())
	return joinWithSpace(s)
}

type Concat struct {
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
	Name string
	Args []Expression
}

func (e Function) String() string {
	return e.Name + "(" + listExpressions(e.Args) + ")"
}

type AggregateFunction struct {
	Name   string
	Option AggregateOption
}

func (e AggregateFunction) String() string {
	return e.Name + "(" + e.Option.String() + ")"
}

type AggregateOption struct {
	Distinct Token
	Args     []Expression
}

func (e AggregateOption) IsDistinct() bool {
	return !e.Distinct.IsEmpty()
}

func (e AggregateOption) String() string {
	s := []string{}
	if !e.Distinct.IsEmpty() {
		s = append(s, e.Distinct.Literal)
	}
	if e.Args != nil {
		s = append(s, listExpressions(e.Args))
	}
	return joinWithSpace(s)
}

type Table struct {
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

func (t *Table) Name() string {
	if t.Alias != nil {
		return t.Alias.(Identifier).Literal
	}

	if file, ok := t.Object.(Identifier); ok {
		return FormatTableName(file.Literal)
	}

	return t.Object.String()
}

type Join struct {
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

func (f *Field) Name() string {
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
}

func (ac AllColumns) String() string {
	return "*"
}

type Dual struct {
	Dual string
}

func (d Dual) String() string {
	return d.Dual
}

type Stdin struct {
	Stdin string
}

func (si Stdin) String() string {
	return si.Stdin
}

type OrderItem struct {
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
	Else   string
	Result Expression
}

func (ce CaseElse) String() string {
	s := []string{ce.Else, ce.Result.String()}
	return joinWithSpace(s)
}

type GroupConcat struct {
	GroupConcat  string
	Option       AggregateOption
	OrderBy      Expression
	SeparatorLit string
	Separator    string
}

func (gc GroupConcat) String() string {
	s := []string{}

	op := gc.Option.String()
	if 0 < len(op) {
		s = append(s, op)
	}
	if gc.OrderBy != nil {
		s = append(s, gc.OrderBy.String())
	}
	if 0 < len(gc.SeparatorLit) {
		s = append(s, gc.SeparatorLit)
	}
	if 0 < len(gc.Separator) {
		s = append(s, quoteString(gc.Separator))
	}
	return gc.GroupConcat + "(" + joinWithSpace(s) + ")"
}

type AnalyticFunction struct {
	Name           string
	Args           []Expression
	Over           string
	AnalyticClause AnalyticClause
}

func (e AnalyticFunction) String() string {
	s := []string{
		e.Name + "(" + listExpressions(e.Args) + ")",
		e.Over,
		"(" + e.AnalyticClause.String() + ")",
	}
	return joinWithSpace(s)
}

type AnalyticClause struct {
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
	PartitionBy string
	Values      []Expression
}

func (e Partition) String() string {
	s := []string{e.PartitionBy, listExpressions(e.Values)}
	return joinWithSpace(s)
}

type Variable struct {
	Name string
}

func (v Variable) String() string {
	return v.Name
}

type VariableSubstitution struct {
	Variable Variable
	Value    Expression
}

func (vs VariableSubstitution) String() string {
	return joinWithSpace([]string{vs.Variable.String(), ":=", vs.Value.String()})
}

type VariableAssignment struct {
	Name  string
	Value Expression
}

func (va VariableAssignment) String() string {
	if va.Value == nil {
		return va.Name
	}
	return joinWithSpace([]string{va.Name, "=", va.Value.String()})
}

type VariableDeclaration struct {
	Assignments []Expression
}

type InsertQuery struct {
	WithClause Expression
	Insert     string
	Into       string
	Table      Identifier
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
	Field FieldReference
	Value Expression
}

func (us UpdateSet) String() string {
	return joinWithSpace([]string{us.Field.String(), "=", us.Value.String()})
}

type DeleteQuery struct {
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
	AlterTable string
	Table      Identifier
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
	AlterTable string
	Table      Identifier
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
	AlterTable string
	Table      Identifier
	Rename     string
	Old        FieldReference
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
	Name       Identifier
	Parameters []Variable
	Statements []Statement
}

type Return struct {
	Value Expression
}

type Print struct {
	Value Expression
}

type Printf struct {
	Values []Expression
}

type Source struct {
	FilePath string
}

type SetFlag struct {
	Name  string
	Value Primary
}

type If struct {
	Condition  Expression
	Statements []Statement
	ElseIf     []ProcExpr
	Else       ProcExpr
}

type ElseIf struct {
	Condition  Expression
	Statements []Statement
}

type Else struct {
	Statements []Statement
}

type While struct {
	Condition  Expression
	Statements []Statement
}

type WhileInCursor struct {
	Variables  []Variable
	Cursor     Identifier
	Statements []Statement
}

type CursorDeclaration struct {
	Cursor Identifier
	Query  SelectQuery
}

type OpenCursor struct {
	Cursor Identifier
}

type CloseCursor struct {
	Cursor Identifier
}

type DisposeCursor struct {
	Cursor Identifier
}

type FetchCursor struct {
	Position  Expression
	Cursor    Identifier
	Variables []Variable
}

type FetchPosition struct {
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

type TableDeclaration struct {
	Table  Identifier
	Fields []Expression
	Query  Expression
}

type TransactionControl struct {
	Token int
}

type FlowControl struct {
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
