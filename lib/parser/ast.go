package parser

import (
	"errors"
	"fmt"
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

type Expression interface {
	String() string
}

type Primary interface {
	String() string
	Bool() bool
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

func (s String) Bool() bool {
	if b, err := strconv.ParseBool(s.Value()); err == nil {
		return b
	}
	return false
}

func (s String) Ternary() ternary.Value {
	return ternary.ParseBool(s.Bool())
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

func (i Integer) Bool() bool {
	if i.Value() == 1 {
		return true
	}
	return false
}

func (i Integer) Ternary() ternary.Value {
	return ternary.ParseBool(i.Bool())
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

func (f Float) Bool() bool {
	if f.Value() == 1 {
		return true
	}
	return false
}

func (f Float) Ternary() ternary.Value {
	return ternary.ParseBool(f.Bool())
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

func (b Boolean) Bool() bool {
	return b.value
}

func (b Boolean) Ternary() ternary.Value {
	return ternary.ParseBool(b.Bool())
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

func (t Ternary) Bool() bool {
	return t.Ternary().BoolValue()
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

func (dt Datetime) Bool() bool {
	return !dt.Value().IsZero()
}

func (dt Datetime) Ternary() ternary.Value {
	return ternary.ParseBool(dt.Bool())
}

func (dt Datetime) Format() string {
	return dt.value.Format(dt.format)
}

func (dt *Datetime) SetFormat(format string) {
	runes := []rune(format)
	gofmt := []rune{}

	escaped := false
	for _, r := range runes {
		if !escaped {
			switch r {
			case '%':
				escaped = true
			default:
				gofmt = append(gofmt, r)
			}
			continue
		}

		switch r {
		case 'a':
			gofmt = append(gofmt, []rune("Mon")...)
		case 'b':
			gofmt = append(gofmt, []rune("Jan")...)
		case 'c':
			gofmt = append(gofmt, []rune("1")...)
		case 'd':
			gofmt = append(gofmt, []rune("02")...)
		case 'E':
			gofmt = append(gofmt, []rune("_2")...)
		case 'e':
			gofmt = append(gofmt, []rune("2")...)
		case 'F':
			gofmt = append(gofmt, []rune(".999999")...)
		case 'f':
			gofmt = append(gofmt, []rune(".000000")...)
		case 'H':
			gofmt = append(gofmt, []rune("15")...)
		case 'h':
			gofmt = append(gofmt, []rune("03")...)
		case 'i':
			gofmt = append(gofmt, []rune("04")...)
		case 'l':
			gofmt = append(gofmt, []rune("3")...)
		case 'M':
			gofmt = append(gofmt, []rune("January")...)
		case 'm':
			gofmt = append(gofmt, []rune("01")...)
		case 'N':
			gofmt = append(gofmt, []rune(".999999999")...)
		case 'n':
			gofmt = append(gofmt, []rune(".000000000")...)
		case 'p':
			gofmt = append(gofmt, []rune("PM")...)
		case 'r':
			gofmt = append(gofmt, []rune("03:04:05 PM")...)
		case 's':
			gofmt = append(gofmt, []rune("05")...)
		case 'T':
			gofmt = append(gofmt, []rune("15:04:05")...)
		case 'W':
			gofmt = append(gofmt, []rune("Monday")...)
		case 'Y':
			gofmt = append(gofmt, []rune("2006")...)
		case 'y':
			gofmt = append(gofmt, []rune("06")...)
		case 'Z':
			gofmt = append(gofmt, []rune("Z07:00")...)
		case 'z':
			gofmt = append(gofmt, []rune("MST")...)
		default:
			gofmt = append(gofmt, r)
		}
		escaped = false
	}

	dt.format = string(gofmt)
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

func (n Null) Bool() bool {
	return false
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

func (i *Identifier) FieldRef() (string, string, error) {
	words := strings.Split(i.Literal, ".")
	ref := ""
	field := ""
	if 2 < len(words) {
		return "", "", errors.New(fmt.Sprintf("field identifier = %s, incorrect format", i.Literal))
	} else if len(words) < 2 {
		field = words[0]
	} else {
		ref = words[0]
		field = words[1]
	}
	return ref, field, nil
}

type Parentheses struct {
	Expr Expression
}

func (p Parentheses) String() string {
	return putParentheses(p.Expr.String())
}

type SelectQuery struct {
	SelectClause  Expression
	FromClause    Expression
	WhereClause   Expression
	GroupByClause Expression
	HavingClause  Expression
	OrderByClause Expression
	LimitClause   Expression
}

func (sq SelectQuery) String() string {
	s := []string{sq.SelectClause.String()}
	if sq.FromClause != nil {
		s = append(s, sq.FromClause.String())
	}
	if sq.WhereClause != nil {
		s = append(s, sq.WhereClause.String())
	}
	if sq.GroupByClause != nil {
		s = append(s, sq.GroupByClause.String())
	}
	if sq.HavingClause != nil {
		s = append(s, sq.HavingClause.String())
	}
	if sq.OrderByClause != nil {
		s = append(s, sq.OrderByClause.String())
	}
	if sq.LimitClause != nil {
		s = append(s, sq.LimitClause.String())
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
	Limit  string
	Number int64
}

func (l LimitClause) String() string {
	s := []string{l.Limit, strconv.FormatInt(l.Number, 10)}
	return joinWithSpace(s)
}

type Subquery struct {
	Query SelectQuery
}

func (sq Subquery) String() string {
	return putParentheses(sq.Query.String())
}

type Comparison struct {
	LHS      Expression
	Operator Token
	RHS      Expression
}

func (c Comparison) String() string {
	s := []string{c.LHS.String(), c.Operator.Literal, c.RHS.String()}
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
	List     []Expression
	Query    Subquery
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
	s = append(s, i.In)
	if i.List != nil {
		s = append(s, putParentheses(listExpressions(i.List)))
	} else {
		s = append(s, i.Query.String())
	}
	return joinWithSpace(s)
}

type All struct {
	All      string
	LHS      Expression
	Operator Token
	Query    Subquery
}

func (a All) String() string {
	s := []string{a.LHS.String(), a.Operator.Literal, a.All, a.Query.String()}
	return joinWithSpace(s)
}

type Any struct {
	Any      string
	LHS      Expression
	Operator Token
	Query    Subquery
}

func (a Any) String() string {
	s := []string{a.LHS.String(), a.Operator.Literal, a.Any, a.Query.String()}
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

type Option struct {
	Distinct Token
	Args     []Expression
}

func (o Option) IsDistinct() bool {
	return !o.Distinct.IsEmpty()
}

func (o Option) String() string {
	s := []string{}
	if !o.Distinct.IsEmpty() {
		s = append(s, o.Distinct.Literal)
	}
	if o.Args != nil {
		s = append(s, listExpressions(o.Args))
	}
	return joinWithSpace(s)
}

type Function struct {
	Name   string
	Option Option
}

func (f Function) String() string {
	return f.Name + "(" + f.Option.String() + ")"
}

type Table struct {
	Object Expression
	As     Token
	Alias  Expression
}

func (t Table) String() string {
	s := []string{t.Object.String()}
	if !t.As.IsEmpty() {
		s = append(s, t.As.Literal)
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
	As     Token
	Alias  Expression
}

func (f Field) String() string {
	s := []string{f.Object.String()}
	if !f.As.IsEmpty() {
		s = append(s, f.As.Literal)
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
	Item      Expression
	Direction Token
}

func (oi OrderItem) String() string {
	s := []string{oi.Item.String()}
	if !oi.Direction.IsEmpty() {
		s = append(s, oi.Direction.Literal)
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
	Option       Option
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
	Var         string
	Assignments []Expression
}

func (vd VariableDeclaration) String() string {
	return joinWithSpace([]string{vd.Var, listExpressions(vd.Assignments)})
}

type InsertQuery struct {
	Insert     string
	Into       string
	Table      Identifier
	Fields     []Expression
	Values     string
	ValuesList []Expression
	Query      Expression
}

func (iq InsertQuery) String() string {
	s := []string{iq.Insert, iq.Into, iq.Table.String()}
	if iq.Fields != nil {
		s = append(s, putParentheses(listExpressions(iq.Fields)))
	}
	if iq.ValuesList != nil {
		s = append(s, iq.Values)
		s = append(s, listExpressions(iq.ValuesList))
	} else {
		s = append(s, iq.Query.String())
	}
	return joinWithSpace(s)
}

type InsertValues struct {
	Values []Expression
}

func (iv InsertValues) String() string {
	return putParentheses(listExpressions(iv.Values))
}

type UpdateQuery struct {
	Update      string
	Tables      []Expression
	Set         string
	SetList     []Expression
	FromClause  Expression
	WhereClause Expression
}

func (uq UpdateQuery) String() string {
	s := []string{uq.Update, listExpressions(uq.Tables), uq.Set, listExpressions(uq.SetList)}
	if uq.FromClause != nil {
		s = append(s, uq.FromClause.String())
	}
	if uq.WhereClause != nil {
		s = append(s, uq.WhereClause.String())
	}
	return joinWithSpace(s)
}

type UpdateSet struct {
	Field Identifier
	Value Expression
}

func (us UpdateSet) String() string {
	return joinWithSpace([]string{us.Field.String(), "=", us.Value.String()})
}

type DeleteQuery struct {
	Delete      string
	Tables      []Expression
	FromClause  Expression
	WhereClause Expression
}

func (dq DeleteQuery) String() string {
	s := []string{dq.Delete}
	if dq.Tables != nil {
		s = append(s, listExpressions(dq.Tables))
	}
	s = append(s, dq.FromClause.String())
	if dq.WhereClause != nil {
		s = append(s, dq.WhereClause.String())
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
	Old        Identifier
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

type Print struct {
	Print string
	Value Expression
}

func (e Print) String() string {
	s := []string{e.Print, e.Value.String()}
	return joinWithSpace(s)
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
