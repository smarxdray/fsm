package fsm

type TransitionRule interface {
	findSuccessor(state StateCode, event EventCode) StateCode
	allStates() map[StateCode]struct{}
	validState(state StateCode) bool
	validEvent(event EventCode) bool
}

type Transitions []Transition

type Transition struct {
	SrcState StateCode
	Event    EventCode
	DstState StateCode
}

type defaultTransitionRule struct {
	rule   map[StateCode]map[EventCode]StateCode
	states map[StateCode]struct{}
	events map[EventCode]struct{}
}

func NewTransitionRule(transitions Transitions) TransitionRule {
	tr := &defaultTransitionRule{
		rule:   make(map[StateCode]map[EventCode]StateCode),
		events: make(map[EventCode]struct{}),
		states: make(map[StateCode]struct{}),
	}
	return tr.register(transitions)
}

func (r *defaultTransitionRule) register(transitions Transitions) TransitionRule {
	for _, transition := range transitions {
		if r.rule[transition.SrcState] == nil {
			r.rule[transition.SrcState] = make(map[EventCode]StateCode)
		}
		r.rule[transition.SrcState][transition.Event] = transition.DstState
		r.states[transition.SrcState] = struct{}{}
		r.events[transition.Event] = struct{}{}
		r.states[transition.DstState] = struct{}{}
	}
	return r
}

func (r *defaultTransitionRule) findSuccessor(state StateCode, event EventCode) StateCode {
	if eventToSuccessor, ok := r.rule[state]; ok {
		if successor, ok := eventToSuccessor[event]; ok {
			return successor
		}
	}
	return state
}

func (r *defaultTransitionRule) validEvent(event EventCode) bool {
	if _, ok := r.events[event]; ok {
		return true
	}
	return false
}

func (r *defaultTransitionRule) validState(state StateCode) bool {
	if _, ok := r.states[state]; ok {
		return true
	}
	return false
}

func (r *defaultTransitionRule) allStates() map[StateCode]struct{} {
	return r.states
}
