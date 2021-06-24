package processing

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/NotInTheSubject/prowler/mods"
	"github.com/sirupsen/logrus"
)

type RequestID = interface{}

type ExecutedTimes = int

type IdentifiedRequest struct {
	Request   *http.Request
	Identifer RequestID
}

type FuzzingStatistic struct {
	RequestStatistic      map[RequestID]ExecutedTimes
	SequenceExecutedTimes ExecutedTimes
}

type SequenceProducer interface {
	GetRequest(lastResponse *http.Response) (IdentifiedRequest, error)
}

type ExternalSystem interface {
	GetSequenceProducer() (SequenceProducer, error)
}

type StopCondition func(FuzzingStatistic) bool

func RunProwling(logger *logrus.Logger, es ExternalSystem, client http.Client, mods []mods.Modifier, stopCondition StopCondition) error {
	var (
		statistic        = FuzzingStatistic{RequestStatistic: make(map[interface{}]int)}
		commonLog        = logger.WithField("env", "RunProwling")
		sequenceModifier = NewSequenceModifier(mods)
	)

	for ; !stopCondition(statistic); sequenceModifier = sequenceModifier.GetNextSequenceModifier() {
		var (
			sequenceProducer func(*http.Response) (IdentifiedRequest, error)
			lastResponse *http.Response = nil
		)

		if sp, err := es.GetSequenceProducer(); err != nil {
			errReport := fmt.Sprintf("cannot get a sequence producer: %+v", err)
			commonLog.Error(errReport)
			return fmt.Errorf("RunProwling: " + errReport)
		} else {
			sequenceProducer = func(resp *http.Response) (IdentifiedRequest, error) {
				r, err := sp.GetRequest(resp)
				if err == nil {
					statistic.RequestStatistic[r.Identifer]++
				}
				if r.Request == nil || r.Identifer == nil {
					statistic.SequenceExecutedTimes++
				}
				return r, err
			}
		}

		for !stopCondition(statistic) {
			var (
				sequenceLog                 = commonLog
			)

			request, err := sequenceProducer(lastResponse)
			if err != nil {
				errReport := fmt.Sprintf("cannot get a new request: %+v", err)
				sequenceLog.Info(errReport)
				break
			}

			if request.Identifer == nil || request.Request == nil {
				sequenceLog.Info("the sequence is finished")
				break
			}

			sequenceLog = logWithRequest(sequenceLog.WithField("request-id", request.Identifer), request.Request)

			moddedRequest := sequenceModifier.Modify(request)

			lastResponse, err = client.Do(moddedRequest)
			if err != nil {
				errReport := fmt.Sprintf("executing a request is failed: %+v", err)
				sequenceLog.Error(errReport)
				break
			}
			sequenceLog = logWithResponse(sequenceLog, lastResponse)
			defer lastResponse.Body.Close()

			if lastResponse.StatusCode/100 == 4 {
				sequenceLog.Info("4** statusCode is received")
			}

			if lastResponse.StatusCode/100 == 5 {
				sequenceLog.Info("5** statusCode is received")
			}
		}
	}

	return nil
}

func logWithRequest(logEntry *logrus.Entry, request *http.Request) *logrus.Entry {
	if request == nil {
		return logEntry
	}

	return logEntry.WithFields(logrus.Fields{
		"req-headers": request.Header,
		"req-method":  request.Method,
		"req-url":     request.URL,
	})
}

func logWithResponse(logEntry *logrus.Entry, response *http.Response) *logrus.Entry {
	respCopy, err := DeepCopyHTTPResponse(response)
	if err != nil {
		logEntry.Errorf("cannot read response body: %+v", err)
	}
	bodyBytes, err := ioutil.ReadAll(respCopy.Body)
	if err != nil {
		logEntry.Errorf("cannot read response body: %+v", err)
	}
	bodyString := string(bodyBytes)

	return logEntry.WithFields(logrus.Fields{
		"resp-status-code": respCopy.StatusCode,
		"resp-headers":     respCopy.Header,
		"resp-body":        bodyString,
	})
}
