package io

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

// FormDataFileUpload upload file using html file upload
func FormDataFileUpload(remoteURL string, filepath string) error {
	client := http.DefaultClient
	//prepare the reader instances to encode
	file, err := MustOpen(filepath)
	if err != nil {
		return err
	}
	values := map[string]io.Reader{
		"file": file, // lets assume its this file
	}
	err = formDataUpload(client, remoteURL, values)
	if err != nil {
		return err
	}
	return nil
}

func formDataUpload(client *http.Client, url string, values map[string]io.Reader) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

// ByteStreamFileUpload uploads file with given byte stream
// CREDIT: https://gist.github.com/ebraminio/576fdfdff425bf3335b51a191a65dbdb
func ByteStreamFileUpload(remoteURL, remoteFolder, filename string, datastream io.Reader) (string, error) {
	req, err := http.NewRequest("POST", remoteURL, datastream)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "binary/octet-stream")
	req.Header.Set("X-Filename", filename)
	req.Header.Set("X-Folder", remoteFolder)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code is %d", res.StatusCode)
	}

	message, _ := ioutil.ReadAll(res.Body)

	return string(message), nil
}
