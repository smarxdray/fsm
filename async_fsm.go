package fsm

type AsyncFSM interface {
	RegisterEntryAction(state StateCode, entryAction StateAction) AsyncFSM
	RegisterExitAction(state StateCode, exitAction StateAction) AsyncFSM
	From(state StateCode) AsyncFSM
	Envoy() Envoy
}

type Envoy interface {
	Run()
	Stop()
	Error() error
	State() State
	EmitEvent(eventCode EventCode, eventData ...interface{})
}

type defaultAsyncFSM struct {
	base      FSM
	eventPipe chan Event
	stopSig   chan struct{}
}

// NewAsyncFSM pipeSize determines the max number of events emitted concurrently.
// pipeSize is 0 by default, which means emitting the event in blocking mode.
func NewAsyncFSM(stateGraph TransitionRule, pipeSize ...int) AsyncFSM {
	defaultPipeSize := 0
	if len(pipeSize) > 0 {
		defaultPipeSize = pipeSize[0]
	}
	defaultFSM := NewFSM(stateGraph)
	if defaultFSM == nil {
		return nil
	}
	m := &defaultAsyncFSM{
		base:      defaultFSM,
		eventPipe: make(chan Event, defaultPipeSize),
		stopSig:   make(chan struct{}),
	}
	return m
}

func (d *defaultAsyncFSM) From(state StateCode) AsyncFSM {
	d.base.From(state)
	return d
}

func (d *defaultAsyncFSM) RegisterEntryAction(state StateCode, entryAction StateAction) AsyncFSM {
	d.base.RegisterEntryAction(state, entryAction)
	return d
}

func (d *defaultAsyncFSM) RegisterExitAction(state StateCode, exitAction StateAction) AsyncFSM {
	d.base.RegisterExitAction(state, exitAction)
	return d
}

func (d *defaultAsyncFSM) Envoy() Envoy {
	return d
}

func (d *defaultAsyncFSM) EmitEvent(eventCode EventCode, eventData ...interface{}) {
	d.eventPipe <- NewEvent(eventCode, eventData)
}

func (d *defaultAsyncFSM) State() State {
	return d.base.State()
}

func (d *defaultAsyncFSM) Run() {
	for {
		select {
		case event, ok := <-d.eventPipe:
			if ok && event != nil {
				d.base.Transfer(event.Code(), event.GetEventData())
			}
		case <-d.stopSig:
			break
		}
	}
}

func (d *defaultAsyncFSM) Stop() {
	d.stopSig <- struct{}{}
}

func (d *defaultAsyncFSM) Error() error {
	return d.base.Error()
}
