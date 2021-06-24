package mods

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"
)

type modConstructor struct {
	modFunc func(*http.Request) bool
	about   string
}

func (m modConstructor) Modify(r *http.Request) bool {
	return m.modFunc(r)
}

func (m modConstructor) About() string {
	return m.about
}

func MakeModifier(modFunc func(*http.Request) bool, about string) Modifier {
	return modConstructor{modFunc: modFunc, about: about}
}

func DefaultModifiers() []Modifier {
	var result []Modifier

	for i := 0; i < 1000; i++ {
		var (
			pos        = rand.Int()
			changeData byte
		)
		for changeData != byte(0) {
			changeData = byte(rand.Int() % 10)
		}

		result = append(result, MakeModifier(func(r *http.Request) bool {
			log := logrus.
				WithField("func", "defaultModifier").
				WithField("pos", pos).
				WithField("changeData", changeData)

			b, err := httputil.DumpRequestOut(r, true)
			if err != nil {
				log.Error(err)
				return false
			}

			b[pos] += changeData
			_, err = http.ReadRequest(bufio.NewReader(bytes.NewReader(b)))
			if err != nil {
				log.Error(err)
				return false
			}

			return true
		}, fmt.Sprintf("defaultModifier: changing -- data[%v] += byte(%v)", pos, int(changeData))))
	}
	return result
}

type Modifier interface {
	Modify(*http.Request) bool
	About() string
}
