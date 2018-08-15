package db

import (
	"os"
	"testing"
)

func TestNewFileDBErrors(t *testing.T) {
	if _, err := NewFileDB("/var/nopermtestHASH"); err == nil {
		t.Fatalf("expected error when opening file we don't own")
	}
	if _, err := NewFileDB("testdata/wrongdata.data"); err == nil {
		t.Fatalf("expected error when data of file is corrupted")
	}
}

func TestFilePersistence(t *testing.T) {
	f, err := os.Create("testdata/testdata.data")
	if err != nil {
		t.Fatalf("err opening file: %s", err)
	}
	db := &FileDB{
		file: f,
		data: map[string]string{"key1": "value1", "key2": "value2"},
	}
	if err = db.Close(); err != nil {
		t.Fatalf("failed to close DB: %s", err)
	}

	db, err = NewFileDB("testdata/testdata.data")
	if err != nil {
		t.Fatalf("failed to open DB: %s", err)
	}
	if _, ok := db.data["key1"]; !ok {
		t.Errorf("expected key1 to be in file")
	}
	if _, ok := db.data["key2"]; !ok {
		t.Errorf("expected key1 to be in file")
	}
}

func TestClosedDB(t *testing.T) {
	db, err := NewFileDB("testdata/testdata.data")
	if err != nil {
		t.Fatalf("err when opening file: %s", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("err when closing db: %s", err)
	}
	if err := db.Create("asda", "sdasd"); err == nil {
		t.Errorf("Create() on closed DB should fail")
	}
	if err := db.Update("asda", "sdasd"); err == nil {
		t.Errorf("Updatk() on closed DB should fail")
	}
	if _, err := db.Read("asda"); err == nil {
		t.Errorf("Read() on closed DB should fail")
	}
	if _, err := db.Delete("asda"); err == nil {
		t.Errorf("Delete() on closed DB should fail")
	}
	if err := db.Close(); err == nil {
		t.Errorf("Close() on closed DB should fail")
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		key   string
		value string
	}
	cases := []struct {
		name    string
		data    map[string]string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{name: "key exists", data: map[string]string{"hi": "value"}, args: args{key: "hi", value: "nope"}, wantErr: true},
		{name: "key does not exist", data: map[string]string{"bye": "nope"}, args: args{key: "hi", value: "nope"}, wantErr: false},
		{name: "invalid key format", args: args{key: "asd$asd", value: "asda"}, wantErr: true},
	}
	for _, c := range cases {
		db := &FileDB{
			data: c.data,
		}
		t.Run(c.name, func(t *testing.T) {
			if err := db.Create(c.args.key, c.args.value); (err != nil) != c.wantErr {
				t.Errorf("db.Create() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}

func TestRead(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		data    map[string]string
		key     string
		want    string
		wantErr bool
	}{
		{name: "key does not exist", data: map[string]string{"key": "value"}, key: "nope", want: "", wantErr: true},
		{name: "key exists", data: map[string]string{"key": "value"}, key: "key", want: "value", wantErr: false},
	}
	for _, c := range cases {
		db := &FileDB{
			data: c.data,
		}
		t.Run(c.name, func(t *testing.T) {
			got, err := db.Read(c.key)
			if (err != nil) && !c.wantErr {
				t.Errorf("db.Read() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if got != c.want {
				t.Errorf("db.Read() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		key   string
		value string
	}
	cases := []struct {
		name    string
		data    map[string]string
		args    args
		wantErr bool
	}{
		{name: "update existing key", data: map[string]string{"key": "value"}, args: args{key: "key", value: "val"}, wantErr: false},
		{name: "key does not exist", data: map[string]string{"key": "value"}, args: args{key: "asdas", value: "asda"}, wantErr: true},
	}
	for _, c := range cases {
		db := &FileDB{
			data: c.data,
		}
		t.Run(c.name, func(t *testing.T) {
			if err := db.Update(c.args.key, c.args.value); (err != nil) != c.wantErr {
				t.Errorf("db.Update() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		data    map[string]string
		key     string
		want    string
		wantErr bool
	}{
		{name: "key does not exist", data: map[string]string{"key": "value"}, key: "nope", want: "", wantErr: true},
		{name: "key does exist", data: map[string]string{"key": "value"}, key: "key", want: "value", wantErr: false},
	}
	for _, c := range cases {
		db := &FileDB{
			data: c.data,
		}
		t.Run(c.name, func(t *testing.T) {
			got, err := db.Delete(c.key)
			if (err != nil) != c.wantErr {
				t.Errorf("db.Delete() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if got != c.want {
				t.Errorf("db.Delete() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestParseData(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{name: "valid one-line string", data: `key:"{"key": "value"}"`, wantErr: false},
		{name: "valid multi-line string", data: `key:"{"key":"value"}"
otherKey:"{"nope":"val"}"`, wantErr: false},
		{name: "invalid one-line string", data: `sdfa$sdf$:value`, wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			_, err := parseData(c.data)
			if err != nil && !c.wantErr {
				t.Errorf("parseData error = %s, wantErr = %v", err, c.wantErr)
			}
		})
	}
}
