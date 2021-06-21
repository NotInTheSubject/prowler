package processing

import (
	"net/http"

	"github.com/NotInTheSubject/prowler/mods"
)

type SequenceModifier interface {
	GetModifier(RequestID) mods.Modifier
	GetNextSequenceModifier() SequenceModifier
}

type defaultSeqMod struct {
	originOperator mods.Operator
	operators      map[RequestID]mods.Operator
	// probably rename
	isCompletedModded map[RequestID]bool
	// isNewModSequence  bool
}

func (sm defaultSeqMod) GetModifier(id RequestID) mods.Modifier {
	// TODO: think about wrapping a request before modifying
	var (
		mod         mods.Modifier
		identityMod mods.Modifier = func(r *http.Request) *http.Request { return r }
	)

	modsOp, found := sm.operators[id]
	if !found {
		modsOp = sm.originOperator.Copy()
		sm.operators[id] = modsOp
	}

	if _, isCompletedModded := sm.isCompletedModded[id]; !isCompletedModded {
		mod = modsOp.GetNewMod()
	} else {
		mod = modsOp.GetLastMod()
		if mod == nil {
			sm.isCompletedModded[id] = true
		}
	}

	if mod == nil {
		mod = identityMod
	}
	return mod
}

func (sm defaultSeqMod) GetNextSequenceModifier() SequenceModifier {
	return defaultSeqMod{}
}

func NewDefaultSequenceModifier(modOp mods.Operator) SequenceModifier {
	return defaultSeqMod{originOperator: modOp, operators: make(map[interface{}]mods.Operator)}
}
