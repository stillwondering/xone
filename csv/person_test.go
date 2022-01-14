package csv

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stillwondering/xone"
)

func TestEncode(t *testing.T) {
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
			wantDst: "Harry,Potter,1980-07-31\n",
			wantErr: false,
		},
		{
			name: "Multiple persons",
			args: args{
				persons: []xone.Person{
					{FirstName: "Harry", LastName: "Potter", DateOfBirth: dateFromString(t, "1980-07-31")},
					{FirstName: "Ron", LastName: "Weasley", DateOfBirth: dateFromString(t, "1980-03-01")},
					{FirstName: "Hermione", LastName: "Granger", DateOfBirth: dateFromString(t, "1979-09-19")},
				},
			},
			wantDst: "Harry,Potter,1980-07-31\nRon,Weasley,1980-03-01\nHermione,Granger,1979-09-19\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &bytes.Buffer{}
			if err := Encode(dst, tt.args.persons); (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDst := dst.String(); gotDst != tt.wantDst {
				t.Errorf("Encode() = %v, want %v", gotDst, tt.wantDst)
			}
		})
	}
}

func TestDecode(t *testing.T) {
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
				src: strings.NewReader("Harry,Potter,1980-07-31\n"),
			},
			want: []xone.Person{
				{FirstName: "Harry", LastName: "Potter", DateOfBirth: dateFromString(t, "1980-07-31")},
			},
			wantErr: false,
		},
		{
			name: "Multiple records",
			args: args{
				src: strings.NewReader("Harry,Potter,1980-07-31\nRon,Weasley,1980-03-01\nHermione,Granger,1979-09-19\n"),
			},
			want: []xone.Person{
				{FirstName: "Harry", LastName: "Potter", DateOfBirth: dateFromString(t, "1980-07-31")},
				{FirstName: "Ron", LastName: "Weasley", DateOfBirth: dateFromString(t, "1980-03-01")},
				{FirstName: "Hermione", LastName: "Granger", DateOfBirth: dateFromString(t, "1979-09-19")},
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
			got, err := Decode(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
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
