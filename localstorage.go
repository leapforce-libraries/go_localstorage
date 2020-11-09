package localstorage

import (
	"encoding/json"
	"io/ioutil"

	fileio "github.com/leapforce-libraries/go_fileio"
)

const defaultFileName string = "__storage__"

type keyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LocalStorage struct {
	fileName string
	data     *[]keyValue
}

func NewLocalStorage(fileName *string) (*LocalStorage, error) {
	fileName_ := defaultFileName
	if fileName != nil {
		fileName_ = *fileName
	}
	ms := LocalStorage{fileName_, nil}
	err := ms.read()
	if err != nil {
		return nil, err
	}

	return &ms, nil
}

func (localStorage *LocalStorage) read() error {
	data := []keyValue{}

	if fileio.FileExists(localStorage.fileName) {
		b, err := ioutil.ReadFile(localStorage.fileName)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &data)
		if err != nil {
			return err
		}

	}

	localStorage.data = &data

	return nil
}

func (localStorage *LocalStorage) Get(key string) (*string, error) {
	if localStorage.data == nil {
		err := localStorage.read()
		if err != nil {
			return nil, err
		}
	}

	for _, m := range *localStorage.data {
		if m.Key == key {
			return &m.Value, nil
		}
	}

	//keyValue, _ := time.Parse("2006-01-02", "1800-01-01")
	return nil, nil
}

func (localStorage *LocalStorage) Set(key string, value string) error {
	found := false

	for i, m := range *localStorage.data {
		if m.Key == key {
			(*localStorage.data)[i].Value = value
			found = true
		}
	}

	if !found {
		data := append(*localStorage.data, keyValue{key, value})

		localStorage.data = &data
	}

	b, err := json.Marshal(*localStorage.data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(localStorage.fileName, b, 0644)
	if err != nil {
		return err
	}

	return nil
}
