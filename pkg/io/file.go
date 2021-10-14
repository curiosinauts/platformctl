package io

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

// WriteStringTofile writes given string to file
func WriteStringTofile(s string, path string) error {
	return WriteBytesToFile([]byte(s), path)
}

// WriteBytesToFile writes given bytes to file
func WriteBytesToFile(bytes []byte, path string) error {
	err := ioutil.WriteFile(path, bytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

// ReadFileToBytes reads files and returns the content as bytes
func ReadFileToBytes(path string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

// DoesPathExists checks for existence of file
func DoesPathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false
	}

	return true
}

// RenderTemplate renders template
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

// WriteTemplate executes go template and writes the content to file
func WriteTemplate(tpl string, data interface{}, path string) error {
	rendered, err := RenderTemplate(tpl, data)
	if err != nil {
		return err
	}
	return WriteStringTofile(rendered, path)
}

// MustOpen ensures file can be opened
func MustOpen(f string) (*os.File, error) {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		log.Println("PWD: ", pwd)
		return nil, err
	}
	return r, nil
}
