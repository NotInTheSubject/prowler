package processing

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)


func DeepCopyHTTPRequest(req *http.Request) (*http.Request, error) {
	cloned := req.Clone(context.Background())
	if req.Body != nil {
		b, err := readBody(req)
		if err != nil {
			return nil, err
		}
		cloned.Body = ioutil.NopCloser(bytes.NewReader(b))
	}
	return cloned, nil
}

func readBody(req *http.Request) ([]byte, error) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to copy the http.Request.Body")
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(data))
	return data, nil
}
