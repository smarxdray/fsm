package fsm

import "errors"

type FSM interface {
	RegisterEntryAction(state StateCode, entryAction StateAction) FSM
	RegisterExitAction(state StateCode, exitAction StateAction) FSM
	From(state StateCode) FSM
	Transfer(event EventCode, eventData ...interface{}) FSM
	State() State
	Error() error
}

type Context interface {
	WriteIntoCurrentState(stateData interface{}) FSM
	ReadFromCurrentState() interface{}
}

type defaultFSM struct {
	stateInstances map[StateCode]State
	transitionRule TransitionRule
	currentState   State
	err            error
}

func NewFSM(stateGraph TransitionRule) FSM {
	if stateGraph == nil {
		return nil
	}
	m := &defaultFSM{
		stateInstances: make(map[StateCode]State),
		transitionRule: stateGraph,
	}
	for code := range m.transitionRule.allStates() {
		m.stateInstances[code] = code.instance()
	}
	return m
}

func (m *defaultFSM) From(state StateCode) FSM {
	if s, ok := m.stateInstances[state]; ok {
		m.currentState = s
		return m
	}
	m.err = errors.New("invalid starting State")
	return m
}

// Transfer take an event with optional data to trigger transition
func (m *defaultFSM) Transfer(eventCode EventCode, eventData ...interface{}) FSM {
	if !m.transitionRule.validEvent(eventCode) {
		m.err = errors.New("invalid event")
		return m
	}
	if m.currentState == nil {
		m.err = errors.New("current State is nil")
		return m
	}
	event := NewEvent(eventCode, eventData)
	for m.successorFound(event) {
		// trigger the exitAction of the current state before transfer to the succeeded state
		if exitAction := m.currentState.exitAction(); exitAction != nil {
			// emit new event, if any
			if exitEvent := exitAction(m, event); exitEvent != nil {
				// overwrite event and re-find the succeeded state
				if event = exitEvent; !m.successorFound(event) {
					break
				}
			}
		}
		// transfer to the succeeded state
		nextStateCode := m.transitionRule.findSuccessor(m.currentState.Code(), event.Code())
		m.currentState = m.stateInstances[nextStateCode]
		// trigger the entryAction of the succeeded state first
		if entryAction := m.currentState.entryAction(); entryAction != nil {
			// emit entryEvent and overwrite event directly
			// if entryEvent exists, state transition continues in the next loop
			event = entryAction(m, event)
		} else { // or vanish the already used event
			event = nil
		}
	}
	return m
}

func (m *defaultFSM) successorFound(event Event) bool {
	if event == nil {
		return false
	}
	currStateCode := m.currentState.Code()
	nextStateCode := m.transitionRule.findSuccessor(currStateCode, event.Code())
	return nextStateCode != currStateCode
}

func (m *defaultFSM) RegisterEntryAction(state StateCode, entryAction StateAction) FSM {
	if _, ok := m.stateInstances[state]; !ok {
		return m
	}
	m.stateInstances[state].withEntryAction(entryAction)
	return m
}

func (m *defaultFSM) RegisterExitAction(state StateCode, exitAction StateAction) FSM {
	if _, ok := m.stateInstances[state]; !ok {
		return m
	}
	m.stateInstances[state].withExitAction(exitAction)
	return m
}

func (m *defaultFSM) State() State {
	return m.currentState
}

func (m *defaultFSM) Error() error {
	return m.err
}

func (m *defaultFSM) WriteIntoCurrentState(stateData interface{}) FSM {
	if m.currentState == nil {
		m.err = errors.New("current state is nil")
		return m
	}
	m.currentState.setStateData(stateData)
	return m
}

func (m *defaultFSM) ReadFromCurrentState() interface{} {
	if m.currentState == nil {
		m.err = errors.New("current state is nil")
		return m
	}
	return m.currentState.GetStateData()
}
