package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func MkdirIfNotExist(p string) error {
	if f, err := os.Stat(p); err == nil { // 如果已经存在了
		if !f.IsDir() { // 如果不是目录，报错
			return errors.New(fmt.Sprintf("%s exist and is a file", p))
		}
	} else if os.IsNotExist(err) { // 如果不存在
		if err := os.MkdirAll(p, 0744); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func WriteToFile(filename string, data []byte) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0744)
	if err != nil {
		log.Printf("[error]: AppendOrCreate open file %s error, %+v\n", filename, err)
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		log.Printf("[error]: AppendOrCreate write file %s error, %+v\n", filename, err)
	}
	return err
}

func ReadFromFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("[error]: ReadFromFile read file %s error, %+v\n", filename, err)
	}
	return data, err
}
