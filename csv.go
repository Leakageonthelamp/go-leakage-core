package core

import (
	"bytes"
	"net/http"
	"os"

	csvs "encoding/csv"

	"github.com/gocarina/gocsv"
)

type ICSVOptions struct {
	FirstRowIsHeader bool
	Separator        string
}

type ICSV[T any] interface {
	//Reader
	ReadFromFile(data []byte, options *ICSVOptions) ([]T, error)
	ReadFromPath(path string, options *ICSVOptions) ([]T, error)
	ReadFromString(data string, options *ICSVOptions) ([]T, error)
	ReadFromURL(url string, options *ICSVOptions) ([]T, error)
	ReadFromFileMaps(data []byte, options *ICSVOptions) ([]map[string]interface{}, error)
}

type csv[T any] struct {
	ctx IContext
}

func NewCSV[T any](ctx IContext) ICSV[T] {
	return &csv[T]{
		ctx: ctx,
	}
}

func (c *csv[T]) ReadFromFile(data []byte, options *ICSVOptions) ([]T, error) {

	reader := bytes.NewReader(data)

	items := make([]T, 0)
	csvReader := csvs.NewReader(reader)
	csvReader.LazyQuotes = true
	if options != nil && !options.FirstRowIsHeader {
		err := gocsv.UnmarshalCSVWithoutHeaders(csvReader, &items)
		if err != nil {
			return nil, err
		}
	} else {
		err := gocsv.UnmarshalCSV(csvReader, &items)
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}

func (c *csv[T]) ReadFromPath(path string, options *ICSVOptions) ([]T, error) {
	clientsFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()

	items := make([]T, 0)
	csvReader := csvs.NewReader(clientsFile)
	csvReader.LazyQuotes = true
	if options != nil && !options.FirstRowIsHeader {
		err := gocsv.UnmarshalCSVWithoutHeaders(csvReader, &items)
		if err != nil {
			return nil, err
		}
	} else {
		err := gocsv.UnmarshalCSV(csvReader, &items)
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (c *csv[T]) ReadFromURL(url string, options *ICSVOptions) ([]T, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	items := make([]T, 0)
	csvReader := csvs.NewReader(resp.Body)
	csvReader.LazyQuotes = true
	if options != nil && !options.FirstRowIsHeader {
		err := gocsv.UnmarshalCSVWithoutHeaders(csvReader, &items)
		if err != nil {
			return nil, err
		}
		mapss := map[string]interface{}{}
		err = gocsv.UnmarshalCSVToMap(csvReader, mapss)
		if err != nil {
			return nil, err
		}

	} else {
		err := gocsv.UnmarshalCSV(csvReader, &items)
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}

func (c *csv[T]) ReadFromFileMaps(data []byte, options *ICSVOptions) ([]map[string]interface{}, error) {
	reader := bytes.NewReader(data)
	csvData, err := csvs.NewReader(reader).ReadAll()
	if err != nil {
		return nil, err
	}

	dataList := []map[string]interface{}{}
	var dataLineOne []string

	for i, line := range csvData {
		if i == 0 {
			dataLineOne = line
		} else {
			data := map[string]interface{}{}
			for i2, line2 := range dataLineOne {
				data[line2] = line[i2]
			}
			dataList = append(dataList, data)
		}
	}

	return dataList, nil
}

func (c *csv[T]) ReadFromString(data string, options *ICSVOptions) ([]T, error) {
	items := make([]T, 0)

	err := gocsv.UnmarshalString(data, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
