package utils

import (
	"bytes"
	"encoding/csv"
)

func GenerateAddressCSV(fields []string, records []map[string]string) ([]byte, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	header := []string{}

	for _, field := range fields {
		header = append(header, AllowedAddressExportFields[field])
	}

	if err := writer.Write(header); err != nil {
		return nil, err
	}

	for _, record := range records {
		row := []string{}
		for _, field := range fields {
			row = append(row, record[field])
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buffer.Bytes(), nil
}