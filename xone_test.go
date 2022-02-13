package xone

import (
	"testing"
	"time"
)

func TestPerson_Age(t *testing.T) {
	type fields struct {
		ID          int
		FirstName   string
		LastName    string
		DateOfBirth time.Time
	}
	type args struct {
		today time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Today before date of birth",
			fields: fields{
				DateOfBirth: time.Date(1980, time.July, 31, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				today: time.Date(1980, time.July, 30, 0, 0, 0, 0, time.UTC),
			},
			want: 0,
		},
		{
			name: "Birthday is today",
			fields: fields{
				DateOfBirth: time.Date(1980, time.July, 31, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				today: time.Date(2000, time.July, 31, 0, 0, 0, 0, time.UTC),
			},
			want: 20,
		},
		{
			name: "Today is after a birthday",
			fields: fields{
				DateOfBirth: time.Date(1980, time.July, 31, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				today: time.Date(2000, time.August, 1, 0, 0, 0, 0, time.UTC),
			},
			want: 20,
		},
		{
			name: "Date of birth is missing",
			fields: fields{
				DateOfBirth: time.Time{},
			},
			args: args{
				today: time.Now(),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Person{
				ID:          tt.fields.ID,
				FirstName:   tt.fields.FirstName,
				LastName:    tt.fields.LastName,
				DateOfBirth: tt.fields.DateOfBirth,
			}
			if got := p.Age(tt.args.today); got != tt.want {
				t.Errorf("Person.Age() = %v, want %v", got, tt.want)
			}
		})
	}
}
