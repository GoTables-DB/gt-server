package table

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

type Table struct {
	columns []Column
	rows    []map[string]any
}

func (t *Table) GetColumns() []Column {
	return t.columns
}

func (t *Table) GetRows() []map[string]any {
	return t.rows
}

func (t *Table) GetColumn(name string) (Column, error) {
	index, err := getColumnIndex(name, t.columns)
	if err != nil {
		return Column{}, err
	}
	column := t.columns[index]
	return column, nil
}

func (t *Table) SetColumnName(name string, value string) error {
	_, notFound := getColumnIndex(value, t.columns)
	if notFound == nil {
		return errors.New("column with name " + value + " does already exist")
	}
	index, err := getColumnIndex(name, t.columns)
	if err != nil {
		return err
	}
	t.columns[index].Name = value
	return nil
}

func (t *Table) SetColumnDefault(name string, value string) error {
	index, err := getColumnIndex(name, t.columns)
	if err != nil {
		return err
	}
	// t.columns[index].Default = t.columns[index].Default
	return errors.New("cannot set default of column " + name + " at index " + strconv.Itoa(index) + " to " + value)
}

func (t *Table) GetRow(index int) (map[string]any, error) {
	if index < 0 && index >= len(t.rows) {
		return nil, errors.New("index " + strconv.Itoa(index) + " is out of range")
	}
	return t.rows[index], nil
}

func (t *Table) SetRow(index int, value map[string]any) error {
	if index < 0 && index >= len(t.rows) {
		return errors.New("index " + strconv.Itoa(index) + " is out of range")
	}
	err := checkRow(value, t.columns)
	if err != nil {
		return err
	}
	row := cleanRow(value, t.columns)
	t.rows[index] = row
	return nil
}

func (t *Table) AddColumn(value Column) error {
	_, notFound := getColumnIndex(value.Name, t.columns)
	if notFound == nil {
		return errors.New("column with name " + value.Name + " does already exist")
	}
	if findDatatype(value.Type) == nil {
		return errors.New("invalid datatype")
	}
	t.columns = append(t.columns, value)
	def := defaultValue(value.Type)
	for i := 0; i < len(t.rows); i++ {
		t.rows[i][value.Name] = def
	}
	return nil
}

func (t *Table) AddRow(value map[string]any) error {
	err := checkRow(value, t.columns)
	if err != nil {
		return err
	}
	row := cleanRow(value, t.columns)
	t.rows = append(t.rows, row)
	return nil
}

func (t *Table) DeleteColumn(name string) error {
	index, err := getColumnIndex(name, t.columns)
	if err != nil {
		return err
	}
	t.columns = append(t.columns[:index], t.columns[index+1:]...)
	for i := 0; i < len(t.rows); i++ {
		t.rows[i] = cleanRow(t.rows[i], t.columns)
	}
	return nil
}

func (t *Table) DeleteRow(index int) error {
	if index < 0 && index >= len(t.rows) {
		return errors.New("index " + strconv.Itoa(index) + " is out of range")
	}
	t.rows = append(t.rows[:index], t.rows[index+1:]...)
	return nil
}

/// Type conversions

// ToU - table to tableU
func (t *Table) ToU() TableU {
	out := TableU{}
	out.Rows = t.rows
	out.Columns = t.columns
	return out
}

/// Helper functions

func getColumnIndex(name string, columns []Column) (int, error) {
	for i := 0; i < len(columns); i++ {
		if name == columns[i].Name {
			return i, nil
		}
	}
	return 0, errors.New("column " + name + " not in table")
}

func findDatatype(datatype string) reflect.Type {
	var ret reflect.Type
	switch datatype {
	// String
	case "str":
		ret = reflect.TypeOf("")
	// Integer
	case "int":
		ret = reflect.TypeOf(0)
	// Float
	case "flt":
		ret = reflect.TypeOf(0.0)
	// Boolean
	case "bol":
		ret = reflect.TypeOf(false)
	// Date
	case "dat":
		ret = reflect.TypeOf(time.Time{})
	// Table
	case "tbl":
		ret = reflect.TypeOf(Table{})
	default:
		ret = nil
	}
	return ret
}

func defaultValue(datatype string) any {
	switch datatype {
	// String
	case "str":
		var ret string
		return ret
	// Integer
	case "int":
		var ret int
		return ret
	// Float
	case "flt":
		var ret float64
		return ret
	// Boolean
	case "bol":
		var ret bool
		return ret
	// Date
	case "dat":
		var ret time.Time
		return ret
	// Table
	case "tbl":
		var ret Table
		return ret
	default:
		return nil
	}
}

func correctDatatype(data any, datatype string) bool {
	var ret bool
	switch data.(type) {
	// String
	case string:
		if datatype == "str" {
			ret = true
		}
	// Integer
	case int:
		if datatype == "int" {
			ret = true
		}
	// Float
	case float32, float64:
		if datatype == "flt" {
			ret = true
		}
	// Boolean
	case bool:
		if datatype == "bol" {
			ret = true
		}
	// Date
	case time.Time:
		if datatype == "dat" {
			ret = true
		}
	// Table
	case Table, TableU:
		if datatype == "tbl" {
			ret = true
		}
	default:
		ret = false
	}
	return ret
}

// Ensure that the row does not contain entries that do not exist and that it contains an entry for all columns
func cleanRow(row map[string]any, columns []Column) map[string]any {
	out := map[string]any{}
	for i := 0; i < len(columns); i++ {
		item := row[columns[i].Name]
		out[columns[i].Name] = item
	}
	return out
}

func checkRow(row map[string]any, columns []Column) error {
	for i := 0; i < len(columns); i++ {
		item := row[columns[i].Name]
		if !correctDatatype(item, columns[i].Type) {
			return errors.New("item in column " + columns[i].Name + " has the wrong datatype")
		}
	}
	return nil
}
