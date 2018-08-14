package db

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	// ErrDuplicatedKey indicates an error when using a key
	// that exists.
	ErrDuplicatedKey = errors.New("key already exists")
	// ErrWrongFormat indicates the content doesn't follow the format.
	ErrWrongFormat = errors.New("format is not correct")
	// ErrKeyNotFound indicates the key is not in the DB.
	ErrKeyNotFound = errors.New("key is not present in the DB")
	// ErrOpeningFile happens when we are unable to open file.
	ErrOpeningFile = errors.New("failed to open file")
	// ErrSavingToFile happens when writes to the file fail.
	ErrSavingToFile = errors.New("failed to write to file")
)

var (
	keyFormat  = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)
	lineFormat = regexp.MustCompile(`^[a-zA-Z0-9_-]*:.*$`)
)

const (
	keyValueSep = ":"
)

// DB is a database with the basic CRUD operations.
type DB interface {
	Create(key, value string) error
	Read(key string) (string, error)
	Update(key, value string) error
	Delete(key string) (string, error)
}

type value struct {
	saved bool
	data  string
}

// FileDB is a DB holding data in-memory and making
// persistence to a file.
type FileDB struct {
	mu   sync.Mutex
	data map[string]value
	file *os.File
}

// NewFileDB returns a DB with the data of the
// file loaded.
func NewFileDB(filename string) (*FileDB, error) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, ErrOpeningFile
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, ErrOpeningFile
	}
	data, err := parseData(string(b))
	if err != nil {
		return nil, err
	}
	return &FileDB{
		data: data,
		file: f,
	}, nil
}

func parseData(data string) (map[string]value, error) {
	d := make(map[string]value)
	s := bufio.NewScanner(strings.NewReader(data))
	for s.Scan() {
		line := strings.TrimSuffix(s.Text(), "\n")
		if !lineFormat.MatchString(line) {
			return map[string]value{}, ErrWrongFormat
		}
		key, v := line[:strings.Index(line, keyValueSep)], line[strings.Index(line, keyValueSep)+1:]
		fmt.Printf("%s=%s\n", key, v)
		d[key] = value{saved: true, data: v}
	}
	return d, nil
}

// Close dumps all the data into the file.
func (db *FileDB) Close() error {
	for k, v := range db.data {
		if v.saved {
			continue
		}
		b := append([]byte(k), []byte(v.data)...)
		if _, err := db.file.Write(append(b, []byte("\n")...)); err != nil {
			return ErrSavingToFile
		}
	}
	return db.file.Close()
}

// Create implements the create method of DB.
// If the key already exists it returns ErrDuplicatedKey.
// If the  value doesn't follow the basic format it returns
// ErrWrongFormat.
func (db *FileDB) Create(key, val string) error {
	if !keyFormat.MatchString(key) {
		return ErrWrongFormat
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.data[key]; ok {
		return ErrDuplicatedKey
	}
	db.data[key] = value{data: val}
	return nil
}

// Read retrieves the value from the database, if it not exists
// it returns ErrKeyNotFound.
func (db *FileDB) Read(key string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	v, ok := db.data[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return v.data, nil
}

// Update updates the `key` with `value`.
// If the key already exists it returns ErrDuplicatedKey.
// If the  value doesn't follow the basic format it returns
// ErrWrongFormat.
func (db *FileDB) Update(key, val string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.data[key]; !ok {
		return ErrKeyNotFound
	}
	db.data[key] = value{data: val}
	return nil
}

// Delete retrieves the value from the database and deletes it.
// If it not exists it returns ErrKeyNotFound.
func (db *FileDB) Delete(key string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	v, ok := db.data[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	delete(db.data, key)
	return v.data, nil
}
