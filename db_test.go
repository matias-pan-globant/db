package db

import (
	"testing"
)

func TestCreate(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		data    map[string]value
		args    args
		want    map[string]value
		wantErr bool
	}{
		{name: "key exists", data: map[string]value{"hi": value{data: "value"}}, args: args{key: "hi", value: "nope"}, wantErr: true},
		{name: "key does not exist", data: map[string]value{"bye": value{data: "nope"}}, args: args{key: "hi", value: "nope"}, wantErr: false},
	}
	for _, tt := range tests {
		db := &FileDB{
			data: tt.data,
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := db.Create(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("db.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]value
		key     string
		want    string
		wantErr bool
	}{
		{name: "key does not exist", data: map[string]value{"key": value{data: "value"}}, key: "nope", want: "", wantErr: true},
		{name: "key exists", data: map[string]value{"key": value{data: "value"}}, key: "key", want: "value", wantErr: false},
	}
	for _, tt := range tests {
		db := &FileDB{
			data: tt.data,
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.Read(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("db.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		data    map[string]value
		args    args
		wantErr bool
	}{
		{name: "update existing key", data: map[string]value{"key": value{data: "value"}}, args: args{key: "key", value: "val"}, wantErr: false},
	}
	for _, tt := range tests {
		db := &FileDB{
			data: tt.data,
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := db.Update(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("db.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]value
		key     string
		want    string
		wantErr bool
	}{
		{name: "key does not exist", data: map[string]value{"key": value{data: "value"}}, key: "nope", want: "", wantErr: true},
		{name: "key does exist", data: map[string]value{"key": value{data: "value"}}, key: "key", want: "value", wantErr: false},
	}
	for _, tt := range tests {
		db := &FileDB{
			data: tt.data,
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.Delete(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("db.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseData(t *testing.T) {
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
