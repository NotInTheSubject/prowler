package processing

import (
	"net/http"

	"github.com/NotInTheSubject/prowler/mods"
)

type SequenceModifier interface {
	Modify(IdentifiedRequest) *http.Request
	GetNextSequenceModifier() SequenceModifier
}

func NewSequenceModifier(modifiers []mods.Modifier) SequenceModifier {
	return seqModifier{modifiers: modifiers, requestsProgress: make(map[RequestID]int), isNewCombination: false}
}

type seqModifier struct {
	modifiers        []mods.Modifier
	requestsProgress map[RequestID]int
	isNewCombination bool
}

func (sm seqModifier) Modify(r IdentifiedRequest) *http.Request {

	requestCopy, err := DeepCopyHTTPRequest(r.Request)
	if err != nil {
		return r.Request
	}

	if len(sm.modifiers) == 0 {
		return requestCopy
	}

	if _, progressFound := sm.requestsProgress[r.Identifer]; !progressFound {
		sm.requestsProgress[r.Identifer] = 0
	}

	if !sm.isNewCombination {
		sm.requestsProgress[r.Identifer]++
		if (sm.requestsProgress[r.Identifer] % len(sm.modifiers)) != 0 {
			sm.isNewCombination = true
		}
	}

	for i := sm.requestsProgress[r.Identifer]; i%len(sm.modifiers) != sm.requestsProgress[r.Identifer]%len(sm.modifiers); i++ {
		isChanged := sm.modifiers[i%len(sm.modifiers)].Modify(requestCopy)
		if isChanged {
			sm.requestsProgress[r.Identifer] = i
			break
		}
	}

	return requestCopy
}

func (sm seqModifier) GetNextSequenceModifier() SequenceModifier {
	var newRequestProgress = make(map[RequestID]int)
	for k, v := range sm.requestsProgress {
		newRequestProgress[k] = v
	}
	return seqModifier{modifiers: sm.modifiers, requestsProgress: newRequestProgress, isNewCombination: false}
}
