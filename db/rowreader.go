package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// TimeFormat should be used to parse time strings
var TimeFormat = time.RFC3339Nano

// Types
var (
	typeMapJSON = reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf((*interface{})(nil)).Elem()) // type of map[string]interface{}

	typeByteSlice = reflect.SliceOf(reflect.TypeOf(byte(0))) // type of []byte

	typeTime = reflect.TypeOf(time.Time{})

	typeRawJSON = reflect.TypeOf(json.RawMessage{})
)

// Errors
var (
	ErrorNullValue   = errors.New("Null value encountered")
	ErrorWrongType   = errors.New("Unable to convert type")
	ErrorUnsupported = errors.New("Unsupported type")
)

type rowReader struct {
	rows      *sql.Rows
	columns   []string
	values    []interface{}
	valuePtrs []interface{}
	colIdxMap map[string]int
	lastError error
}

// RowReader is used to simplify reading sql.Rows object
type RowReader interface {
	ScanNext() bool
	Error() error

	ReadAllAsMap() map[string]interface{}
	ReadAllToStruct(p interface{})

	ReadByIdxBytes(columnIdx int) []byte
	ReadByNameBytes(columnName string) []byte

	ReadByIdxJSON(columnIdx int) map[string]interface{}
	ReadByNameJSON(columnName string) map[string]interface{}

	ReadByIdxString(columnIdx int) string
	ReadByNameString(columnName string) string

	ReadByIdxInt64(columnIdx int) int64
	ReadByNameInt64(columnName string) int64

	ReadByIdxFloat64(columnIdx int) float64
	ReadByNameFloat64(columnName string) float64

	ReadByIdxBool(columnIdx int) bool
	ReadByNameBool(columnName string) bool

	ReadByIdxTime(columnIdx int) time.Time
	ReadByNameTime(columnName string) time.Time
}

// GetRowReader returns RowReader interface that can be used to read data from sql.Rows
func GetRowReader(rows *sql.Rows) (RowReader, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	n := len(columns)

	rr := new(rowReader)
	rr.rows = rows
	rr.columns = columns
	rr.values = make([]interface{}, n)
	rr.valuePtrs = make([]interface{}, n)

	for i := 0; i < n; i++ {
		rr.valuePtrs[i] = &rr.values[i]
	}

	return rr, nil
}

func (rr *rowReader) getColumnIndex(name string) int {
	if rr.colIdxMap == nil {
		m := make(map[string]int)
		for idx, col := range rr.columns {
			m[col] = idx
		}
		rr.colIdxMap = m
		return m[name]
	}

	return rr.colIdxMap[name]
}

func (rr *rowReader) ScanNext() (hasMore bool) {
	if hasMore = rr.rows.Next(); hasMore {
		err := rr.rows.Scan(rr.valuePtrs...)
		rr.lastError = err
		if err != nil {
			hasMore = false
		}
	}

	return
}

func (rr *rowReader) Error() error {
	return rr.lastError
}

func (rr *rowReader) ReadAllAsMap() map[string]interface{} {
	m := make(map[string]interface{})
	for i, name := range rr.columns {
		switch v := rr.values[i].(type) {
		case []byte:
			m[name] = string(v)
		case bool:
			m[name] = v
		case string:
			m[name] = v
		case float64:
			m[name] = v
		case float32:
			m[name] = v
		case int64:
			m[name] = v
		case int32:
			m[name] = v
		case int16:
			m[name] = v
		case int8:
			m[name] = v
		case time.Time:
			m[name] = v
		case nil:
			m[name] = nil
		}
	}
	return m
}

func (rr *rowReader) ReadAllToStruct(p interface{}) {
	var value reflect.Value

	value = reflect.ValueOf(p)
	if value.Kind() != reflect.Ptr {
		return
	}

	value = reflect.Indirect(value)
	if value.Kind() != reflect.Struct {
		return
	}

	for columnIdx, columnName := range rr.columns {
		if rr.values[columnIdx] == nil {
			continue
		}

		column := caseInsenstiveFieldByName(value, columnName)
		if column == (reflect.Value{}) {
			continue
		}

		switch columnKind := column.Kind(); columnKind {
		case reflect.String:
			column.SetString(rr.ReadByIdxString(columnIdx))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			column.SetInt(rr.ReadByIdxInt64(columnIdx))
		case reflect.Float32, reflect.Float64:
			column.SetFloat(rr.ReadByIdxFloat64(columnIdx))
		case reflect.Bool:
			column.SetBool(rr.ReadByIdxBool(columnIdx))
		case reflect.Struct:
			switch columnType := column.Type(); columnType {
			case typeTime:
				column.Set(reflect.ValueOf(rr.ReadByIdxTime(columnIdx)))
			default:
				panic(ErrorUnsupported)
			}
		case reflect.Slice:
			if column.Type() != typeByteSlice && column.Type() != typeRawJSON {
				panic(ErrorUnsupported)
			}
			column.SetBytes(rr.ReadByIdxBytes(columnIdx)) // will panic if slice is not []byte
		case reflect.Map:
			if column.Type() != typeMapJSON {
				panic(ErrorUnsupported)
			}
			column.Set(reflect.ValueOf(rr.ReadByIdxJSON(columnIdx))) // will panic if map is not map[string]interface{}
		default:
			panic(ErrorUnsupported)
		}
	}
}

func caseInsenstiveFieldByName(v reflect.Value, name string) reflect.Value {
	name = strings.ToLower(name)
	return v.FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == name })
}

// Byte slice readers

func (rr *rowReader) ReadByIdxBytes(columnIdx int) []byte {
	switch v := rr.values[columnIdx].(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	case nil:
		panic(ErrorNullValue)
	default:
		panic(ErrorWrongType)
	}
}

func (rr *rowReader) ReadByNameBytes(columnName string) []byte {
	return rr.ReadByIdxBytes(rr.getColumnIndex(columnName))
}

// JSON readers

func (rr *rowReader) ReadByIdxJSON(columnIdx int) map[string]interface{} {
	var bytes []byte

	switch v := rr.values[columnIdx].(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	case nil:
		return nil
	default:
		panic(ErrorWrongType)
	}

	if j := make(map[string]interface{}); json.Unmarshal(bytes, &j) != nil {
		panic(ErrorWrongType)
	} else {
		return j
	}
}

func (rr *rowReader) ReadByNameJSON(columnName string) map[string]interface{} {
	return rr.ReadByIdxJSON(rr.getColumnIndex(columnName))
}

// String readers

func (rr *rowReader) ReadByIdxString(columnIdx int) string {
	switch v := rr.values[columnIdx].(type) {
	case []byte:
		return string(v)
	case string:
		return v
	case nil:
		panic(ErrorNullValue)
	default:
		panic(ErrorWrongType)
	}
}

func (rr *rowReader) ReadByNameString(columnName string) string {
	return rr.ReadByIdxString(rr.getColumnIndex(columnName))
}

// Int64 readers

func (rr *rowReader) ReadByIdxInt64(columnIdx int) int64 {
	switch v := rr.values[columnIdx].(type) {
	case int64:
		return v
	case []byte:
		s := string(v)
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(ErrorWrongType)
		}
		return i
	case nil:
		panic(ErrorNullValue)
	default:
		panic(ErrorWrongType)
	}
}

func (rr *rowReader) ReadByNameInt64(columnName string) int64 {
	return rr.ReadByIdxInt64(rr.getColumnIndex(columnName))
}

// Float64 readers

func (rr *rowReader) ReadByIdxFloat64(columnIdx int) float64 {
	switch v := rr.values[columnIdx].(type) {
	case float64:
		return v
	case int64:
		return float64(v)
	case []byte:
		s := string(v)
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(ErrorWrongType)
		}
		return f
	case nil:
		panic(ErrorNullValue)
	default:
		panic(ErrorWrongType)
	}
}

func (rr *rowReader) ReadByNameFloat64(columnName string) float64 {
	return rr.ReadByIdxFloat64(rr.getColumnIndex(columnName))
}

// Bool readers

func (rr *rowReader) ReadByIdxBool(columnIdx int) bool {
	switch v := rr.values[columnIdx].(type) {
	case bool:
		return v
	case []byte:
		s := string(v)
		b, err := strconv.ParseBool(s)
		if err != nil {
			panic(ErrorWrongType)
		}
		return b
	case int64:
		return v != 0
	case float64:
		return v != 0.0
	case nil:
		panic(ErrorNullValue)
	default:
		panic(ErrorWrongType)
	}
}

func (rr *rowReader) ReadByNameBool(columnName string) bool {
	return rr.ReadByIdxBool(rr.getColumnIndex(columnName))
}

// Time readers

func (rr *rowReader) ReadByIdxTime(columnIdx int) time.Time {
	switch v := rr.values[columnIdx].(type) {
	case time.Time:
		return v
	case []byte:
		q, err := time.Parse(TimeFormat, string(v))
		if err != nil {
			panic(ErrorWrongType)
		}
		return q
	case string:
		q, err := time.Parse(TimeFormat, v)
		if err != nil {
			panic(ErrorWrongType)
		}
		return q
	case nil:
		panic(ErrorNullValue)
	default:
		panic(ErrorWrongType)
	}
}

func (rr *rowReader) ReadByNameTime(columnName string) time.Time {
	return rr.ReadByIdxTime(rr.getColumnIndex(columnName))
}
