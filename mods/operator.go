package mods

import (
	"bufio"
	"bytes"
	"math/rand"
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"
)

type Operator interface {
	GetLastMod() Modifier
	GetNewMod() Modifier
	Copy() Operator
}

func defaultModifier(r *http.Request) *http.Request {
	b, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		logrus.WithField("func", "defaultModifier").Error(err)
		return r
	}
	for i := 0; i < 100; i++ {
		b[rand.Int()%len(b)] += byte(rand.Int())
		resultReq, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(b)))
		if err == nil {
			return resultReq
		}
	}
	logrus.WithField("func", "defaultModifier").Error(err)
	return r
}

// type operatorImpl struct {
// 	currentPosition int
// 	mods            []Modifier
// }

type defaultOperator struct {}

func (defaultOperator) GetLastMod() Modifier {
	return defaultModifier
}

func (defaultOperator) GetNewMod() Modifier {
	return defaultModifier
}

func (defaultOperator) Copy() Operator {
	return defaultOperator{}
}

func GetDefaultOperator() Operator {
	return defaultOperator{}
}

// func OperatorProducer(typicalReq *http.Request, mods []Modifier) Operator {

// }

type Modifier func(*http.Request) *http.Request
