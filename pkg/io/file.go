package io

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"
)

func WriteStringTofile(s string, path string) error {
	return WriteBytesToFile([]byte(s), path)
}

func WriteBytesToFile(bytes []byte, path string) error {
	err := ioutil.WriteFile(path, bytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

func ReadFileToBytes(path string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

func DoesPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false
	}

	return true
}

func RenderTemplate(tpl string, data interface{}) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 100))
	tmpl, err := template.New("test").Parse(tpl)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func WriteTemplate(tpl string, data interface{}, path string) error {
	rendered, err := RenderTemplate(tpl, data)
	if err != nil {
		return err
	}
	return WriteStringTofile(rendered, path)
}
