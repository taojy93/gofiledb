package gofiledb

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type DB struct {
	sync.Mutex
	Name   string
	Path   string
	Tables map[string]*Table
}

func NewDBClient(name string) (*DB, error) {
	dbPath := filepath.Join(".", name)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if err := os.Mkdir(dbPath, 0755); err != nil {
			return nil, err
		}
	}

	db := &DB{
		Name:   name,
		Path:   dbPath,
		Tables: make(map[string]*Table),
	}

	files, err := os.ReadDir(dbPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			tableName := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
			table := NewTable(tableName)
			if err := table.Load(filepath.Join(dbPath, file.Name())); err != nil {
				return nil, err
			}
			db.Tables[tableName] = table
		}
	}

	return db, nil
}

func (db *DB) SaveTable(tableName string) error {
	db.Lock()
	defer db.Unlock()

	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	filePath := filepath.Join(db.Path, tableName+".json")
	return table.Save(filePath)
}

func (db *DB) CreateTable(tableName string) error {
	db.Lock()
	defer db.Unlock()

	if _, exists := db.Tables[tableName]; exists {
		return fmt.Errorf("table %s already exists", tableName)
	}

	table := NewTable(tableName)
	db.Tables[tableName] = table

	db.Unlock() // Unlock DB before saving the table
	err := db.SaveTable(tableName)
	db.Lock() // Lock DB again after saving the table

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DropTable(tableName string) error {
	db.Lock()
	defer db.Unlock()

	if _, exists := db.Tables[tableName]; !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	filePath := filepath.Join(db.Path, tableName+".json")
	if err := os.Remove(filePath); err != nil {
		return err
	}

	delete(db.Tables, tableName)
	return nil
}

func (db *DB) AddRecord(tableName string, record Record) error {
	db.Lock()
	defer db.Unlock()

	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	table.AddRecord(record)

	db.Unlock() // Unlock DB before saving the table
	err := db.SaveTable(tableName)
	db.Lock() // Lock DB again after saving the table

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UpdateRecord(tableName string, id int, data interface{}) error {
	db.Lock()
	defer db.Unlock()

	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	if err := table.UpdateRecord(id, data); err != nil {
		return err
	}

	db.Unlock() // Unlock DB before saving the table
	err := db.SaveTable(tableName)
	db.Lock() // Lock DB again after saving the table

	if err != nil {
		return err
	}
	return nil

}

func (db *DB) DeleteRecord(tableName string, id int) error {
	db.Lock()
	defer db.Unlock()

	table, exists := db.Tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	if err := table.DeleteRecord(id); err != nil {
		return err
	}

	db.Unlock() // Unlock DB before saving the table
	err := db.SaveTable(tableName)
	db.Lock() // Lock DB again after saving the table

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetRecord(tableName string, id int) (*Record, error) {
	db.Lock()
	defer db.Unlock()

	table, exists := db.Tables[tableName]
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	return table.GetRecord(id)
}
