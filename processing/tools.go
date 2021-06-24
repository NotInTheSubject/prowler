package processing

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func SetRespBody(responseInterface interface{}, rc io.ReadCloser) {
	response := responseInterface.(*http.Response)
	response.Body = rc
}

func GetRespBody(responseInterface interface{}) ([]byte, error) {
	response := responseInterface.(*http.Response)
	return io.ReadAll(response.Body)
}

func SetReqBody(requestInterface interface{}, rc io.ReadCloser) {
	request := requestInterface.(*http.Request)
	request.Body = rc
}

func GetReqBody(requestInterface interface{}) ([]byte, error) {
	request := requestInterface.(*http.Request)
	return io.ReadAll(request.Body)
}


func DeepCopyHTTPResponse(resp *http.Response) (*http.Response, error) {
	cloned := resp
	if resp.Body != nil {
		b, err := readBody(resp, GetRespBody, SetRespBody)
		if err != nil {
			return nil, err
		}
		cloned.Body = ioutil.NopCloser(bytes.NewReader(b))
	}
	return cloned, nil
}

func DeepCopyHTTPRequest(req *http.Request) (*http.Request, error) {
	cloned := req.Clone(context.Background())
	if req.Body != nil {
		b, err := readBody(req, GetReqBody, SetReqBody)
		if err != nil {
			return nil, err
		}
		cloned.Body = ioutil.NopCloser(bytes.NewReader(b))
	}
	return cloned, nil
}

func readBody(src interface{}, reader func(interface{}) ([]byte, error), writer func(interface{}, io.ReadCloser)) ([]byte, error) {
	data, err := reader(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read the *.Body: %+v", err)
	}

	writer(src, ioutil.NopCloser(bytes.NewReader(data)))
	return data, nil
}
