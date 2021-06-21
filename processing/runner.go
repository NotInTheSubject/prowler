package processing

import (
	"fmt"
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

func RunProwling(logger *logrus.Logger, es ExternalSystem, client http.Client, mods mods.Operator, stopCondition StopCondition) error {
	var (
		statistic    = FuzzingStatistic{RequestStatistic: make(map[interface{}]int)}
		commonLog    = logger.WithField("env", "RunProwling")
		lastResponse *http.Response
	)

	for !stopCondition(statistic) {
		sp, err := es.GetSequenceProducer()
		if err != nil {
			errReport := fmt.Sprintf("cannot get a sequence producer: %+v", err)
			commonLog.Error(errReport)
			return fmt.Errorf("RunProwling: " + errReport)
		}

		for !stopCondition(statistic) {
			var requestLog = commonLog.WithFields(logrus.Fields{})

			request, err := sp.GetRequest(lastResponse)
			if err != nil {
				errReport := fmt.Sprintf("cannot get a new request: %+v", err)
				requestLog.Error(errReport)
				return fmt.Errorf("RunProwling: " + errReport)
			}

			if request.Identifer == nil || request.Request == nil {
				break
			}

			requestLog = requestLog.WithField("request-id", request.Identifer)
			logWithRequest(requestLog, request.Request).Debug()
			
			// operate req by modifiers

			lastResponse, err = client.Do(request.Request)
			if err != nil {
				errReport := fmt.Sprintf("executing a request is failed: %+v", err)
				requestLog.Error(errReport)
				return fmt.Errorf("RunProwling: " + errReport)
			}
			statistic.RequestStatistic[request.Identifer]++
		}
		statistic.SequenceExecutedTimes++
	}

	return nil
}

func logWithRequest(logEntry *logrus.Entry, request *http.Request) *logrus.Entry {
	if request == nil {
		return logEntry
	}

	return logEntry.WithFields(logrus.Fields{
		"headers": request.Header,
		"method":  request.Method,
		"url":     request.URL,
	})
}
