package csv

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stillwondering/xone"
)

func TestWrite(t *testing.T) {
	type args struct {
		persons []xone.Person
	}
	tests := []struct {
		name    string
		args    args
		wantDst string
		wantErr bool
	}{
		{
			name: "Empty person collection",
			args: args{
				persons: []xone.Person{},
			},
			wantDst: "",
			wantErr: false,
		},
		{
			name: "One person",
			args: args{
				persons: []xone.Person{
					{FirstName: "Harry", LastName: "Potter", DateOfBirth: dateFromString(t, "1980-07-31")},
				},
			},
			wantDst: "Harry,Potter,1980-07-31,Other\n",
			wantErr: false,
		},
		{
			name: "Multiple persons",
			args: args{
				persons: []xone.Person{
					{FirstName: "Harry", LastName: "Potter", DateOfBirth: dateFromString(t, "1980-07-31"), Gender: xone.Male},
					{FirstName: "Ron", LastName: "Weasley", DateOfBirth: dateFromString(t, "1980-03-01"), Gender: xone.Male},
					{FirstName: "Hermione", LastName: "Granger", DateOfBirth: dateFromString(t, "1979-09-19"), Gender: xone.Female},
				},
			},
			wantDst: "Harry,Potter,1980-07-31,Male\nRon,Weasley,1980-03-01,Male\nHermione,Granger,1979-09-19,Female\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &bytes.Buffer{}
			if err := Write(dst, tt.args.persons); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDst := dst.String(); gotDst != tt.wantDst {
				t.Errorf("Write() = %v, want %v", gotDst, tt.wantDst)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	dir := t.TempDir()

	type args struct {
		file    string
		persons []xone.Person
	}
	tests := []struct {
		name            string
		args            args
		compareWithFile string
		wantErr         bool
	}{
		{
			name: "EmptyFile",
			args: args{
				file:    path.Join(dir, "EmptyFile.csv"),
				persons: []xone.Person{},
			},
			compareWithFile: "testdata/EmptyFile.csv",
			wantErr:         false,
		},
		{
			name: "MultiplePeople",
			args: args{
				file: path.Join(dir, "MultiplePeople.csv"),
				persons: []xone.Person{
					{FirstName: "Harry", LastName: "Potter", DateOfBirth: dateFromString(t, "1980-07-31"), Gender: xone.Male},
					{FirstName: "Ron", LastName: "Weasley", DateOfBirth: dateFromString(t, "1980-03-01"), Gender: xone.Male},
					{FirstName: "Hermione", LastName: "Granger", DateOfBirth: dateFromString(t, "1979-09-19"), Gender: xone.Female},
				},
			},
			compareWithFile: "testdata/MultiplePeople.csv",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteFile(tt.args.file, tt.args.persons); (err != nil) != tt.wantErr {
				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !areContentsEqual(t, tt.compareWithFile, tt.args.file) {
				t.Errorf("WriteFile() does not match %s", tt.compareWithFile)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		src io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []xone.Person
		wantErr bool
	}{
		{
			name: "Empty string",
			args: args{
				src: strings.NewReader(""),
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "One record",
			args: args{
				src: strings.NewReader("Harry,Potter,1980-07-31,Other\n"),
			},
			want: []xone.Person{
				{FirstName: "Harry", LastName: "Potter", DateOfBirth: dateFromString(t, "1980-07-31"), Gender: xone.Other},
			},
			wantErr: false,
		},
		{
			name: "Multiple records",
			args: args{
				src: strings.NewReader("Harry,Potter,1980-07-31,Male\nRon,Weasley,1980-03-01,Male\nHermione,Granger,1979-09-19,Female\n"),
			},
			want: []xone.Person{
				{FirstName: "Harry", LastName: "Potter", DateOfBirth: dateFromString(t, "1980-07-31"), Gender: xone.Male},
				{FirstName: "Ron", LastName: "Weasley", DateOfBirth: dateFromString(t, "1980-03-01"), Gender: xone.Male},
				{FirstName: "Hermione", LastName: "Granger", DateOfBirth: dateFromString(t, "1979-09-19"), Gender: xone.Female},
			},
			wantErr: false,
		},
		{
			name: "Incorrect number of records",
			args: args{
				src: strings.NewReader("Harry,Potter,1980-07-31,incorrect field\n"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Incorrect date of birth format",
			args: args{
				src: strings.NewReader("Harry,Potter,31.07.1980\n"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    []xone.Person
		wantErr bool
	}{
		{
			name:    "IncorrectNumberOfColumns",
			args:    args{file: "testdata/IncorrectNumberOfColumns.csv"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "EmptyFile",
			args:    args{file: "testdata/EmptyFile.csv"},
			want:    nil,
			wantErr: false,
		},
		{
			name: "MultiplePeople",
			args: args{file: "testdata/MultiplePeople.csv"},
			want: []xone.Person{
				{FirstName: "Harry", LastName: "Potter", DateOfBirth: dateFromString(t, "1980-07-31"), Gender: xone.Male},
				{FirstName: "Ron", LastName: "Weasley", DateOfBirth: dateFromString(t, "1980-03-01"), Gender: xone.Male},
				{FirstName: "Hermione", LastName: "Granger", DateOfBirth: dateFromString(t, "1979-09-19"), Gender: xone.Female},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func dateFromString(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse(xone.FormatDateOfBirth, s)
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func areContentsEqual(t *testing.T, expected, actual string) bool {
	t.Helper()

	expectedContent, err := ioutil.ReadFile(expected)
	if err != nil {
		t.Fatal(err)
	}

	actualContent, err := ioutil.ReadFile(actual)
	if err != nil {
		t.Fatal(err)
	}

	return bytes.Equal(expectedContent, actualContent)
}

func TestNewPersonRepository(t *testing.T) {
	dir := t.TempDir()

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *PersonRepository
		wantErr bool
	}{
		{
			name:    "Directory",
			args:    args{filename: dir},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "NonexistentFile",
			args:    args{filename: path.Join(dir, "NonexistentFile.csv")},
			want:    &PersonRepository{file: path.Join(dir, "NonexistentFile.csv")},
			wantErr: false,
		},
		{
			name:    "ExistingFile",
			args:    args{filename: "testdata/EmptyFile.csv"},
			want:    &PersonRepository{file: "testdata/EmptyFile.csv"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPersonRepository(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPersonRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPersonRepository() = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				_, err = os.Stat(tt.args.filename)
				if err != nil {
					if os.IsNotExist(err) {
						t.Errorf("file %s should exist", tt.args.filename)
					} else {
						t.Fatal(err)
					}
				}
			}
		})
	}
}

func TestPersonRepository_GetAll(t *testing.T) {
	type fields struct {
		file string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []xone.Person
		wantErr bool
	}{
		{
			name:    "EmptyFile",
			fields:  fields{file: "testdata/EmptyFile.csv"},
			want:    nil,
			wantErr: false,
		},
		{
			name:   "MultiplePeople",
			fields: fields{file: "testdata/MultiplePeople.csv"},
			want: []xone.Person{
				{FirstName: "Harry", LastName: "Potter", DateOfBirth: dateFromString(t, "1980-07-31"), Gender: xone.Male},
				{FirstName: "Ron", LastName: "Weasley", DateOfBirth: dateFromString(t, "1980-03-01"), Gender: xone.Male},
				{FirstName: "Hermione", LastName: "Granger", DateOfBirth: dateFromString(t, "1979-09-19"), Gender: xone.Female},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PersonRepository{
				file: tt.fields.file,
			}
			got, err := r.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("PersonRepository.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PersonRepository.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
