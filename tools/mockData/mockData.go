package main

import (
	"encoding/csv"
	"os"
	"path"
)

// User schema
type User struct {
	FirstName string
	LastName  string
	Email     string
	Gender    string
	IPAddress string
}

// ReadCsv accepts a file and returns its content as a multi-dimentional type
// with lines and each column. Only parses to string type.
func ReadCsv(filename string) ([][]string, error) {
	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return rows, nil
}

func getMockData() ([]*User, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	rows, err := ReadCsv(path.Join(dir, "tools", "mockData", "users.csv"))
	if err != nil {
		return nil, err
	}

	var data []*User
	for _, row := range rows[1:] {
		data = append(data, &User{
			FirstName: row[0],
			LastName:  row[1],
			Email:     row[2],
			Gender:    row[3],
			IPAddress: row[4],
		})
	}

	return data, nil
}
