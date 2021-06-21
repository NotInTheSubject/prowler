package prowler

import (
	"github.com/NotInTheSubject/prowler/mods"
	"github.com/NotInTheSubject/prowler/processing"
	"github.com/sirupsen/logrus"
)

type RequestID = processing.RequestID

type ExecutedTimes = processing.ExecutedTimes

type IdentifiedRequest = processing.IdentifiedRequest

type FuzzingStatistic = processing.FuzzingStatistic

type RequestModifier = mods.Operator

type StopCondition = processing.StopCondition

type SequenceProducer = processing.SequenceProducer

type ExternalSystem = processing.ExternalSystem

func RunProwling(log *logrus.Logger, es ExternalSystem, mods RequestModifier, stopCondition StopCondition) error {
	return processing.RunProwling(log, es, mods, stopCondition)
}
