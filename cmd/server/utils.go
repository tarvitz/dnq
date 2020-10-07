package server

import (
	"io"
	"io/ioutil"
	"net/http"
)

func Close(target io.Closer) {
	_ = target.Close()
}

func body(request *http.Request) string {
	var (
		content []byte
		err error
	)

	defer Close(request.Body)

	err = request.ParseForm()
	if err != nil {
		return err.Error()
	}
	content, err = ioutil.ReadAll(request.Body)
	if err == nil {
		return string(content)
	}
	return err.Error()
}
