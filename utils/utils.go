package utils

import (
	"encoding/csv"
	"os"
)

func ReadCSV(filepath string) ([][]string, error) {

	csvfile, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	fields, err := reader.ReadAll()

	return fields, nil
}
