package query

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

var Functions = map[string]func(parser.Function, []parser.Primary) (parser.Primary, error){
	"COALESCE":         Coalesce,
	"IF":               If,
	"IFNULL":           Ifnull,
	"NULLIF":           Nullif,
	"CEIL":             Ceil,
	"FLOOR":            Floor,
	"ROUND":            Round,
	"ABS":              Abs,
	"ACOS":             Acos,
	"ASIN":             Asin,
	"ATAN":             Atan,
	"ATAN2":            Atan2,
	"COS":              Cos,
	"SIN":              Sin,
	"TAN":              Tan,
	"EXP":              Exp,
	"EXP2":             Exp2,
	"EXPM1":            Expm1,
	"LOG":              Log,
	"LOG10":            Log10,
	"LOG2":             Log2,
	"LOG1P":            Log1p,
	"SQRT":             Sqrt,
	"POW":              Pow,
	"BIN_TO_DEC":       BinToDec,
	"OCT_TO_DEC":       OctToDec,
	"HEX_TO_DEC":       HexToDec,
	"ENOTATION_TO_DEC": EnotationToDec,
	"BIN":              Bin,
	"OCT":              Oct,
	"HEX":              Hex,
	"ENOTATION":        Enotation,
	"RAND":             Rand,
	"TRIM":             Trim,
	"LTRIM":            Ltrim,
	"RTRIM":            Rtrim,
	"UPPER":            Upper,
	"LOWER":            Lower,
	"BASE64_ENCODE":    Base64Encode,
	"BASE64_DECODE":    Base64Decode,
	"HEX_ENCODE":       HexEncode,
	"HEX_DECODE":       HexDecode,
	"LEN":              Len,
	"BYTE_LEN":         ByteLen,
	"LPAD":             Lpad,
	"RPAD":             Rpad,
	"SUBSTR":           Substr,
	"REPLACE":          Replace,
	"FORMAT":           Format,
	"MD5":              Md5,
	"SHA1":             Sha1,
	"SHA256":           Sha256,
	"SHA512":           Sha512,
	"MD5_HMAC":         Md5Hmac,
	"SHA1_HMAC":        Sha1Hmac,
	"SHA256_HMAC":      Sha256Hmac,
	"SHA512_HMAC":      Sha512Hmac,
	"NOW":              Now,
	"DATETIME_FORMAT":  DatetimeFormat,
	"YEAR":             Year,
	"MONTH":            Month,
	"DAY":              Day,
	"HOUR":             Hour,
	"MINUTE":           Minute,
	"SECOND":           Second,
	"MILLISECOND":      Millisecond,
	"MICROSECOND":      Microsecond,
	"NANOSECOND":       Nanosecond,
	"WEEKDAY":          Weekday,
	"UNIX_TIME":        UnixTime,
	"UNIX_NANO_TIME":   UnixNanoTime,
	"DAY_OF_YEAR":      DayOfYear,
	"WEEK_OF_YEAR":     WeekOfYear,
	"ADD_YEAR":         AddYear,
	"ADD_MONTH":        AddMonth,
	"ADD_DAY":          AddDay,
	"ADD_HOUR":         AddHour,
	"ADD_MINUTE":       AddMinute,
	"ADD_SECOND":       AddSecond,
	"ADD_MILLI":        AddMilli,
	"ADD_MICRO":        AddMicro,
	"ADD_NANO":         AddNano,
	"DATE_DIFF":        DateDiff,
	"TIME_DIFF":        TimeDiff,
	"TIME_NANO_DIFF":   TimeNanoDiff,
	"STRING":           String,
	"INTEGER":          Integer,
	"FLOAT":            Float,
	"BOOLEAN":          Boolean,
	"TERNARY":          Ternary,
	"DATETIME":         Datetime,
	"CALL":             Call,
}

func Coalesce(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) < 1 {
		return nil, NewFunctionArgumentLengthErrorWithCustomArgs(fn, fn.Name, "at least 1 argument")
	}

	for _, arg := range args {
		if !parser.IsNull(arg) {
			return arg, nil
		}
	}
	return parser.NewNull(), nil
}

func If(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 3 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{3})
	}

	if args[0].Ternary() == ternary.TRUE {
		return args[1], nil
	}
	return args[2], nil
}

func Ifnull(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 2 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{2})
	}

	if parser.IsNull(args[0]) {
		return args[1], nil
	}
	return args[0], nil
}

func Nullif(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 2 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{2})
	}

	if EqualTo(args[0], args[1]) == ternary.TRUE {
		return parser.NewNull(), nil
	}
	return args[0], nil
}

func roundParams(args []parser.Primary) (number float64, place float64, isnull bool, argsErr bool) {
	if len(args) < 1 || 2 < len(args) {
		argsErr = true
		return
	}

	f := parser.PrimaryToFloat(args[0])
	if parser.IsNull(f) {
		isnull = true
		return
	}
	number = f.(parser.Float).Value()

	if len(args) == 2 {
		f := parser.PrimaryToInteger(args[1])
		if parser.IsNull(f) {
			isnull = true
			return
		}
		place = float64(f.(parser.Integer).Value())
	}
	return
}

func Ceil(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	number, place, isnull, argsErr := roundParams(args)
	if argsErr {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1, 2})
	}
	if isnull {
		return parser.NewNull(), nil
	}

	pow := math.Pow(10, place)
	r := math.Ceil(pow*number) / pow
	return parser.Float64ToPrimary(r), nil
}

func Floor(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	number, place, isnull, argsErr := roundParams(args)
	if argsErr {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1, 2})
	}
	if isnull {
		return parser.NewNull(), nil
	}

	pow := math.Pow(10, place)
	r := math.Floor(pow*number) / pow
	return parser.Float64ToPrimary(r), nil
}

func round(f float64, place float64) float64 {
	pow := math.Pow(10, place)
	var r float64
	if f < 0 {
		r = math.Ceil(pow*f-0.5) / pow
	} else {
		r = math.Floor(pow*f+0.5) / pow
	}
	return r
}

func Round(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	number, place, isnull, argsErr := roundParams(args)
	if argsErr {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1, 2})
	}
	if isnull {
		return parser.NewNull(), nil
	}

	return parser.Float64ToPrimary(round(number, place)), nil
}

func execMath1Arg(fn parser.Function, args []parser.Primary, mathf func(float64) float64) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	f := parser.PrimaryToFloat(args[0])
	if parser.IsNull(f) {
		return parser.NewNull(), nil
	}

	result := mathf(f.(parser.Float).Value())
	if math.IsInf(result, 0) || math.IsNaN(result) {
		return parser.NewNull(), nil
	}
	return parser.Float64ToPrimary(result), nil
}

func execMath2Args(fn parser.Function, args []parser.Primary, mathf func(float64, float64) float64) (parser.Primary, error) {
	if len(args) != 2 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{2})
	}

	f1 := parser.PrimaryToFloat(args[0])
	if parser.IsNull(f1) {
		return parser.NewNull(), nil
	}

	f2 := parser.PrimaryToFloat(args[1])
	if parser.IsNull(f2) {
		return parser.NewNull(), nil
	}

	result := mathf(f1.(parser.Float).Value(), f2.(parser.Float).Value())
	if math.IsInf(result, 0) || math.IsNaN(result) {
		return parser.NewNull(), nil
	}
	return parser.Float64ToPrimary(result), nil
}

func Abs(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Abs)
}

func Acos(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Acos)
}

func Asin(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Asin)
}

func Atan(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Atan)
}

func Atan2(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath2Args(fn, args, math.Atan2)
}

func Cos(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Cos)
}

func Sin(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Sin)
}

func Tan(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Tan)
}

func Exp(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Exp)
}

func Exp2(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Exp2)
}

func Expm1(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Expm1)
}

func Log(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Log)
}

func Log10(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Log10)
}

func Log2(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Log2)
}

func Log1p(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Log1p)
}

func Sqrt(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath1Arg(fn, args, math.Sqrt)
}

func Pow(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execMath2Args(fn, args, math.Pow)
}

func execParseInt(fn parser.Function, args []parser.Primary, base int) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	p := parser.PrimaryToString(args[0])
	if parser.IsNull(p) {
		return parser.NewNull(), nil
	}

	s := p.(parser.String).Value()
	if base == 16 {
		s = ltrim(s, "0x")
	}

	i, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		return parser.NewNull(), nil
	}

	return parser.NewInteger(i), nil
}

func execFormatInt(fn parser.Function, args []parser.Primary, base int) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	p := parser.PrimaryToInteger(args[0])
	if parser.IsNull(p) {
		return parser.NewNull(), nil
	}

	s := strconv.FormatInt(p.(parser.Integer).Value(), base)
	return parser.NewString(s), nil
}

func BinToDec(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execParseInt(fn, args, 2)
}

func OctToDec(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execParseInt(fn, args, 8)
}

func HexToDec(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execParseInt(fn, args, 16)
}

func EnotationToDec(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	p := parser.PrimaryToString(args[0])
	if parser.IsNull(p) {
		return parser.NewNull(), nil
	}

	s := p.(parser.String).Value()

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return parser.NewNull(), nil
	}

	return parser.Float64ToPrimary(f), nil
}

func Bin(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execFormatInt(fn, args, 2)
}

func Oct(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execFormatInt(fn, args, 8)
}

func Hex(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execFormatInt(fn, args, 16)
}

func Enotation(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	p := parser.PrimaryToFloat(args[0])
	if parser.IsNull(p) {
		return parser.NewNull(), nil
	}

	s := strconv.FormatFloat(p.(parser.Float).Value(), 'e', -1, 64)
	return parser.NewString(s), nil
}

func Rand(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if 0 < len(args) && len(args) != 2 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{0, 2})
	}

	r := cmd.GetRand()

	if len(args) == 0 {
		return parser.NewFloat(r.Float64()), nil
	}

	p1 := parser.PrimaryToInteger(args[0])
	if parser.IsNull(p1) {
		return nil, NewFunctionInvalidArgumentError(fn, fn.Name, "the first argument must be an integer")
	}
	p2 := parser.PrimaryToInteger(args[1])
	if parser.IsNull(p2) {
		return nil, NewFunctionInvalidArgumentError(fn, fn.Name, "the second argument must be an integer")
	}

	low := p1.(parser.Integer).Value()
	high := p2.(parser.Integer).Value()
	if high <= low {
		return nil, NewFunctionInvalidArgumentError(fn, fn.Name, "the second argument must be greater than the first argument")
	}
	delta := high - low + 1
	return parser.NewInteger(r.Int63n(delta) + low), nil
}

func execStrings1Arg(fn parser.Function, args []parser.Primary, stringsf func(string) string) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	s := parser.PrimaryToString(args[0])
	if parser.IsNull(s) {
		return parser.NewNull(), nil
	}

	result := stringsf(s.(parser.String).Value())
	return parser.NewString(result), nil
}

func execStringsTrim(fn parser.Function, args []parser.Primary, stringsf func(string, string) string) (parser.Primary, error) {
	if len(args) < 1 || 2 < len(args) {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1, 2})
	}

	s := parser.PrimaryToString(args[0])
	if parser.IsNull(s) {
		return parser.NewNull(), nil
	}

	cutset := ""
	if 2 == len(args) {
		cs := parser.PrimaryToString(args[1])
		if parser.IsNull(cs) {
			return parser.NewNull(), nil
		}
		cutset = cs.(parser.String).Value()
	}

	result := stringsf(s.(parser.String).Value(), cutset)
	return parser.NewString(result), nil
}

func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func base64Decode(s string) string {
	bytes, _ := base64.StdEncoding.DecodeString(s)
	return string(bytes)
}

func hexEncode(s string) string {
	return hex.EncodeToString([]byte(s))
}

func hexDecode(s string) string {
	bytes, _ := hex.DecodeString(s)
	return string(bytes)
}

func trim(s string, cutset string) string {
	if len(cutset) < 1 {
		return strings.TrimSpace(s)
	}
	return strings.Trim(s, cutset)
}

func ltrim(s string, cutset string) string {
	if len(cutset) < 1 {
		return strings.TrimLeftFunc(s, unicode.IsSpace)
	}
	return strings.TrimLeft(s, cutset)
}

func rtrim(s string, cutset string) string {
	if len(cutset) < 1 {
		return strings.TrimRightFunc(s, unicode.IsSpace)
	}
	return strings.TrimRight(s, cutset)
}

func execStringsLen(fn parser.Function, args []parser.Primary, stringsf func(string) int) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	s := parser.PrimaryToString(args[0])
	if parser.IsNull(s) {
		return parser.NewNull(), nil
	}

	result := stringsf(s.(parser.String).Value())
	return parser.NewInteger(int64(result)), nil
}

func execStringsPadding(fn parser.Function, args []parser.Primary, direction rune) (parser.Primary, error) {
	if len(args) != 3 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{3})
	}

	s := parser.PrimaryToString(args[0])
	if parser.IsNull(s) {
		return parser.NewNull(), nil
	}
	str := s.(parser.String).Value()

	l := parser.PrimaryToInteger(args[1])
	if parser.IsNull(l) {
		return parser.NewNull(), nil
	}
	length := int(l.(parser.Integer).Value())

	p := parser.PrimaryToString(args[2])
	if parser.IsNull(p) {
		return parser.NewNull(), nil
	}
	padstr := p.(parser.String).Value()

	strLen := utf8.RuneCountInString(str)
	padstrLen := utf8.RuneCountInString(padstr)

	if length <= strLen {
		return args[0], nil
	}

	padLen := length - strLen
	repeat := int(math.Ceil(float64(padLen) / float64(padstrLen)))
	padding := strings.Repeat(padstr, repeat)
	padding = string([]rune(padding)[:padLen])

	if direction == 'r' {
		str = str + padding
	} else {
		str = padding + str
	}

	return parser.NewString(str), nil
}

func execCrypto(fn parser.Function, args []parser.Primary, cryptof func() hash.Hash) (parser.Primary, error) {
	if 1 != len(args) {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	s := parser.PrimaryToString(args[0])
	if parser.IsNull(s) {
		return parser.NewNull(), nil
	}

	h := cryptof()
	h.Write([]byte(s.(parser.String).Value()))
	r := hex.EncodeToString(h.Sum(nil))
	return parser.NewString(r), nil

}

func execCryptoHMAC(fn parser.Function, args []parser.Primary, cryptof func() hash.Hash) (parser.Primary, error) {
	if 2 != len(args) {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{2})
	}

	s := parser.PrimaryToString(args[0])
	if parser.IsNull(s) {
		return parser.NewNull(), nil
	}

	key := parser.PrimaryToString(args[1])
	if parser.IsNull(key) {
		return parser.NewNull(), nil
	}

	h := hmac.New(cryptof, []byte(key.(parser.String).Value()))
	h.Write([]byte(s.(parser.String).Value()))
	r := hex.EncodeToString(h.Sum(nil))
	return parser.NewString(r), nil
}

func byteLen(s string) int {
	return len(s)
}

func Trim(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStringsTrim(fn, args, trim)
}

func Ltrim(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStringsTrim(fn, args, ltrim)
}

func Rtrim(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStringsTrim(fn, args, rtrim)
}

func Upper(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStrings1Arg(fn, args, strings.ToUpper)
}

func Lower(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStrings1Arg(fn, args, strings.ToLower)
}

func Base64Encode(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStrings1Arg(fn, args, base64Encode)
}

func Base64Decode(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStrings1Arg(fn, args, base64Decode)
}

func HexEncode(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStrings1Arg(fn, args, hexEncode)
}

func HexDecode(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStrings1Arg(fn, args, hexDecode)
}

func Len(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStringsLen(fn, args, utf8.RuneCountInString)
}

func ByteLen(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStringsLen(fn, args, byteLen)
}

func Lpad(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStringsPadding(fn, args, 'l')
}

func Rpad(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execStringsPadding(fn, args, 'r')
}

func Substr(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) < 2 || 3 < len(args) {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{2, 3})
	}

	s := parser.PrimaryToString(args[0])
	if parser.IsNull(s) {
		return parser.NewNull(), nil
	}

	runes := []rune(s.(parser.String).Value())
	strlen := len(runes)
	start := 0
	end := strlen

	i := parser.PrimaryToInteger(args[1])
	if parser.IsNull(i) {
		return parser.NewNull(), nil
	}
	start = int(i.(parser.Integer).Value())
	if start < 0 {
		start = strlen + start
	}
	if start < 0 || strlen <= start {
		return parser.NewNull(), nil
	}

	if 3 == len(args) {
		i := parser.PrimaryToInteger(args[2])
		if parser.IsNull(i) {
			return parser.NewNull(), nil
		}
		sublen := int(i.(parser.Integer).Value())
		if sublen < 0 {
			return parser.NewNull(), nil
		}
		end = start + sublen
		if strlen < end {
			end = strlen
		}
	}

	return parser.NewString(string(runes[start:end])), nil
}

func Replace(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if 3 != len(args) {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{3})
	}

	s := parser.PrimaryToString(args[0])
	if parser.IsNull(s) {
		return parser.NewNull(), nil
	}

	oldstr := parser.PrimaryToString(args[1])
	if parser.IsNull(oldstr) {
		return parser.NewNull(), nil
	}

	newstr := parser.PrimaryToString(args[2])
	if parser.IsNull(newstr) {
		return parser.NewNull(), nil
	}

	r := strings.Replace(s.(parser.String).Value(), oldstr.(parser.String).Value(), newstr.(parser.String).Value(), -1)
	return parser.NewString(r), nil
}

func Format(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) < 1 {
		return nil, NewFunctionArgumentLengthErrorWithCustomArgs(fn, fn.Name, "at least 1 argument")
	}

	format := parser.PrimaryToString(args[0])
	if parser.IsNull(format) {
		return nil, NewFunctionInvalidArgumentError(fn, fn.Name, "the first argument must be a string")
	}

	str, err := FormatString(format.(parser.String).Value(), args[1:])
	if err != nil {
		return nil, NewFunctionInvalidArgumentError(fn, fn.Name, err.(AppError).ErrorMessage())
	}
	return parser.NewString(str), nil
}

func Md5(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execCrypto(fn, args, md5.New)
}

func Sha1(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execCrypto(fn, args, sha1.New)
}

func Sha256(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execCrypto(fn, args, sha256.New)
}

func Sha512(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execCrypto(fn, args, sha512.New)
}

func Md5Hmac(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execCryptoHMAC(fn, args, md5.New)
}

func Sha1Hmac(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execCryptoHMAC(fn, args, sha1.New)
}

func Sha256Hmac(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execCryptoHMAC(fn, args, sha256.New)
}

func Sha512Hmac(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execCryptoHMAC(fn, args, sha512.New)
}

func Now(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if 0 < len(args) {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{0})
	}

	return parser.NewDatetime(cmd.Now()), nil
}

func DatetimeFormat(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 2 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{2})
	}

	p := parser.PrimaryToDatetime(args[0])
	if parser.IsNull(p) {
		return parser.NewNull(), nil
	}
	format := parser.PrimaryToString(args[1])
	if parser.IsNull(format) {
		return parser.NewNull(), nil
	}

	dt := p.(parser.Datetime)
	return parser.NewString(dt.Format(parser.DatetimeFormats.Get(format.(parser.String).Value()))), nil
}

func execDatetimeToInt(fn parser.Function, args []parser.Primary, timef func(time.Time) int64) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	dt := parser.PrimaryToDatetime(args[0])
	if parser.IsNull(dt) {
		return parser.NewNull(), nil
	}

	result := timef(dt.(parser.Datetime).Value())
	return parser.NewInteger(result), nil
}

func execDatetimeAdd(fn parser.Function, args []parser.Primary, timef func(time.Time, int) time.Time) (parser.Primary, error) {
	if len(args) != 2 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{2})
	}

	p1 := parser.PrimaryToDatetime(args[0])
	if parser.IsNull(p1) {
		return parser.NewNull(), nil
	}
	p2 := parser.PrimaryToInteger(args[1])
	if parser.IsNull(p2) {
		return parser.NewNull(), nil
	}

	dt := p1.(parser.Datetime).Value()
	i := int(p2.(parser.Integer).Value())
	return parser.NewDatetime(timef(dt, i)), nil
}

func year(t time.Time) int64 {
	return int64(t.Year())
}

func month(t time.Time) int64 {
	return int64(t.Month())
}

func day(t time.Time) int64 {
	return int64(t.Day())
}

func hour(t time.Time) int64 {
	return int64(t.Hour())
}

func minute(t time.Time) int64 {
	return int64(t.Minute())
}

func second(t time.Time) int64 {
	return int64(t.Second())
}

func millisecond(t time.Time) int64 {
	return int64(round(float64(t.Nanosecond())/float64(1000000), 0))
}

func microsecond(t time.Time) int64 {
	return int64(round(float64(t.Nanosecond())/float64(1000), 0))
}

func nanosecond(t time.Time) int64 {
	return int64(t.Nanosecond())
}

func weekday(t time.Time) int64 {
	return int64(t.Weekday())
}

func unixTime(t time.Time) int64 {
	return t.Unix()
}

func unixNanoTime(t time.Time) int64 {
	return t.UnixNano()
}

func dayOfYear(t time.Time) int64 {
	return int64(t.YearDay())
}

func weekOfYear(t time.Time) int64 {
	_, w := t.ISOWeek()
	return int64(w)
}

func addYear(t time.Time, duration int) time.Time {
	return t.AddDate(duration, 0, 0)
}

func addMonth(t time.Time, duration int) time.Time {
	return t.AddDate(0, duration, 0)
}

func addDay(t time.Time, duration int) time.Time {
	return t.AddDate(0, 0, duration)
}

func addHour(t time.Time, duration int) time.Time {
	dur := time.Duration(duration)
	return t.Add(dur * time.Hour)
}

func addMinute(t time.Time, duration int) time.Time {
	dur := time.Duration(duration)
	return t.Add(dur * time.Minute)
}

func addSecond(t time.Time, duration int) time.Time {
	dur := time.Duration(duration)
	return t.Add(dur * time.Second)
}

func addMilli(t time.Time, duration int) time.Time {
	dur := time.Duration(duration)
	return t.Add(dur * time.Millisecond)
}

func addMicro(t time.Time, duration int) time.Time {
	dur := time.Duration(duration)
	return t.Add(dur * time.Microsecond)
}

func addNano(t time.Time, duration int) time.Time {
	dur := time.Duration(duration)
	return t.Add(dur * time.Nanosecond)
}

func Year(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, year)
}

func Month(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, month)
}

func Day(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, day)
}

func Hour(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, hour)
}

func Minute(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, minute)
}

func Second(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, second)
}

func Millisecond(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, millisecond)
}

func Microsecond(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, microsecond)
}

func Nanosecond(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, nanosecond)
}

func Weekday(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, weekday)
}

func UnixTime(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, unixTime)
}

func UnixNanoTime(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, unixNanoTime)
}

func DayOfYear(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, dayOfYear)
}

func WeekOfYear(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeToInt(fn, args, weekOfYear)
}

func AddYear(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeAdd(fn, args, addYear)
}

func AddMonth(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeAdd(fn, args, addMonth)
}

func AddDay(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeAdd(fn, args, addDay)
}

func AddHour(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeAdd(fn, args, addHour)
}

func AddMinute(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeAdd(fn, args, addMinute)
}

func AddSecond(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeAdd(fn, args, addSecond)
}

func AddMilli(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeAdd(fn, args, addMilli)
}

func AddMicro(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeAdd(fn, args, addMicro)
}

func AddNano(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return execDatetimeAdd(fn, args, addNano)
}

func DateDiff(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 2 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{2})
	}

	p1 := parser.PrimaryToDatetime(args[0])
	if parser.IsNull(p1) {
		return parser.NewNull(), nil
	}
	p2 := parser.PrimaryToDatetime(args[1])
	if parser.IsNull(p2) {
		return parser.NewNull(), nil
	}

	dt1 := p1.(parser.Datetime).Value()
	dt2 := p2.(parser.Datetime).Value()

	subdt1 := time.Date(dt1.Year(), dt1.Month(), dt1.Day(), 0, 0, 0, 0, cmd.GetLocation())
	subdt2 := time.Date(dt2.Year(), dt2.Month(), dt2.Day(), 0, 0, 0, 0, cmd.GetLocation())
	dur := subdt1.Sub(subdt2)

	return parser.NewInteger(int64(dur.Hours() / 24)), nil
}

func timeDiff(fn parser.Function, args []parser.Primary, durf func(time.Duration) parser.Primary) (parser.Primary, error) {
	if len(args) != 2 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{2})
	}

	p1 := parser.PrimaryToDatetime(args[0])
	if parser.IsNull(p1) {
		return parser.NewNull(), nil
	}
	p2 := parser.PrimaryToDatetime(args[1])
	if parser.IsNull(p2) {
		return parser.NewNull(), nil
	}

	dt1 := p1.(parser.Datetime).Value()
	dt2 := p2.(parser.Datetime).Value()

	dur := dt1.Sub(dt2)
	return durf(dur), nil
}

func durationSeconds(dur time.Duration) parser.Primary {
	return parser.Float64ToPrimary(dur.Seconds())
}

func durationNanoseconds(dur time.Duration) parser.Primary {
	return parser.NewInteger(dur.Nanoseconds())
}

func TimeDiff(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return timeDiff(fn, args, durationSeconds)
}

func TimeNanoDiff(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	return timeDiff(fn, args, durationNanoseconds)
}

func String(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	switch args[0].(type) {
	case parser.Boolean:
		return parser.NewString(strconv.FormatBool(args[0].(parser.Boolean).Value())), nil
	case parser.Ternary:
		return parser.NewString(args[0].(parser.Ternary).Ternary().String()), nil
	case parser.Datetime:
		return parser.NewString(args[0].(parser.Datetime).Format(time.RFC3339Nano)), nil
	default:
		return parser.PrimaryToString(args[0]), nil
	}
}

func Integer(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	switch args[0].(type) {
	case parser.Integer:
		return args[0], nil
	case parser.Float:
		return parser.NewInteger(int64(round(args[0].(parser.Float).Value(), 0))), nil
	case parser.String:
		s := strings.TrimSpace(args[0].(parser.String).Value())
		if i, e := strconv.ParseInt(s, 10, 64); e == nil {
			return parser.NewInteger(i), nil
		}
		if f, e := strconv.ParseFloat(s, 64); e == nil {
			return parser.NewInteger(int64(round(f, 0))), nil
		}
	case parser.Datetime:
		return parser.NewInteger(args[0].(parser.Datetime).Value().Unix()), nil
	}
	return parser.NewNull(), nil
}

func Float(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	switch args[0].(type) {
	case parser.Datetime:
		t := args[0].(parser.Datetime).Value()
		f := float64(t.Unix())
		if t.Nanosecond() > 0 {
			f = f + float64(t.Nanosecond())/float64(1000000000)
		}
		return parser.NewFloat(f), nil
	default:
		return parser.PrimaryToFloat(args[0]), nil
	}
}

func Boolean(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	return parser.PrimaryToBoolean(args[0]), nil
}

func Ternary(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	return parser.NewTernary(args[0].Ternary()), nil
}

func Datetime(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) != 1 {
		return nil, NewFunctionArgumentLengthError(fn, fn.Name, []int{1})
	}

	return parser.PrimaryToDatetime(args[0]), nil
}

func Call(fn parser.Function, args []parser.Primary) (parser.Primary, error) {
	if len(args) < 1 {
		return nil, NewFunctionArgumentLengthErrorWithCustomArgs(fn, fn.Name, "at least 1 argument")
	}

	cmdargs := make([]string, len(args))
	for i, v := range args {
		var s string
		switch v.(type) {
		case parser.String:
			s = v.(parser.String).Value()
		case parser.Integer:
			s = v.(parser.Integer).String()
		case parser.Float:
			s = v.(parser.Float).String()
		case parser.Boolean:
			s = v.(parser.Boolean).String()
		case parser.Ternary:
			s = v.(parser.Ternary).String()
		case parser.Datetime:
			s = v.(parser.Datetime).Format(time.RFC3339Nano)
		case parser.Null:
			s = ""
		}
		cmdargs[i] = s
	}

	buf, err := exec.Command(cmdargs[0], cmdargs[1:]...).Output()
	if err != nil {
		return nil, err
	}
	return parser.NewString(string(buf)), nil
}
