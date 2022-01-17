package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/stillwondering/xone"
)

func Write(dst io.Writer, persons []xone.Person) error {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	for _, person := range persons {
		err := writer.Write([]string{
			person.FirstName,
			person.LastName,
			person.DateOfBirth.Format(xone.FormatDateOfBirth),
		})

		if err != nil {
			return err
		}
	}

	writer.Flush()

	_, err := io.Copy(dst, &buf)

	return err
}

func Parse(src io.Reader) ([]xone.Person, error) {
	reader := getDefaultReader(src)

	var persons []xone.Person

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		dob, err := time.Parse(xone.FormatDateOfBirth, record[2])
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid date of birth", record[2])
		}

		persons = append(persons, xone.Person{
			FirstName:   record[0],
			LastName:    record[1],
			DateOfBirth: dob,
		})
	}

	return persons, nil
}

func getDefaultReader(src io.Reader) *csv.Reader {
	reader := csv.NewReader(src)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = 3

	return reader
}
