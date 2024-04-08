package table

import "encoding/json"

type TableU struct {
	Columns []Column         `json:"columns"`
	Rows    []map[string]any `json:"rows"`
}

/// Type conversions

// ToJ tableU to json
func (t *TableU) ToJ() ([]byte, error) {
	out, err := json.Marshal(t)
	return out, err
}

// FromJ json to tableU
func (t *TableU) FromJ(in []byte) error {
	err := json.Unmarshal(in, t)
	return err
}

// ToT tableU to table
func (t *TableU) ToT() (Table, error) {
	out := Table{}
	for i := 0; i < len(t.Columns); i++ {
		err := out.AddColumn(t.Columns[i])
		if err != nil {
			return out, err
		}
	}
	for i := 0; i < len(t.Rows); i++ {
		err := out.AddRow(t.Rows[i])
		if err != nil {
			return out, err
		}
	}
	return out, nil
}
