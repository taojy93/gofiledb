package gofiledb

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Table struct {
	sync.Mutex
	Name    string   `json:"name"`
	Records []Record `json:"records"`
}

func NewTable(name string) *Table {
	return &Table{
		Name:    name,
		Records: make([]Record, 0),
	}
}

func (t *Table) Load(filePath string) error {
	t.Lock()
	defer t.Unlock()

	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(file, t)
}

func (t *Table) Save(filePath string) error {
	t.Lock()
	defer t.Unlock()

	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

func (t *Table) AddRecord(record Record) {
	t.Lock()
	defer t.Unlock()

	record.ID = len(t.Records) + 1
	t.Records = append(t.Records, record)
}

func (t *Table) UpdateRecord(id int, data interface{}) error {
	t.Lock()
	defer t.Unlock()

	for i, record := range t.Records {
		if record.ID == id {
			t.Records[i].Data = data
			return nil
		}
	}
	return fmt.Errorf("record with ID %d not found", id)
}

func (t *Table) DeleteRecord(id int) error {
	t.Lock()
	defer t.Unlock()

	for i, record := range t.Records {
		if record.ID == id {
			t.Records = append(t.Records[:i], t.Records[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("record with ID %d not found", id)
}

func (t *Table) GetRecord(id int) (*Record, error) {
	t.Lock()
	defer t.Unlock()

	for _, record := range t.Records {
		if record.ID == id {
			return &record, nil
		}
	}
	return nil, fmt.Errorf("record with ID %d not found", id)
}
