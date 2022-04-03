package fsm

import "github.com/mohae/deepcopy"

// StateAction is a closure, which can take action and emit event to transfer State
// if there is no need to emit any event, just return nil
type StateAction func(ctx Context, event Event) Event

type StateCode string

type State interface {
	Code() StateCode
	GetStateData() interface{}
	setStateData(stateData interface{}) State
	withEntryAction(entryAction StateAction) State
	withExitAction(exitAction StateAction) State
	entryAction() StateAction
	exitAction() StateAction
}

type defaultState struct {
	stateCode  StateCode
	entryActor StateAction
	exitActor  StateAction
	stateData  interface{}
}

func (d *defaultState) Code() StateCode {
	return d.stateCode
}

// GetStateData return a snapshot(copy) of the state data
// NOTE: unexported field values are not copied
func (d *defaultState) GetStateData() interface{} {
	return deepcopy.Copy(d.stateData)
}

func (sc StateCode) instance() State {
	return &defaultState{
		stateCode: sc,
	}
}

func (d *defaultState) setStateData(stateData interface{}) State {
	d.stateData = stateData
	return d
}

func (d *defaultState) withEntryAction(entryAction StateAction) State {
	d.entryActor = entryAction
	return d
}

func (d *defaultState) withExitAction(exitAction StateAction) State {
	d.exitActor = exitAction
	return d
}

func (d *defaultState) entryAction() StateAction {
	return d.entryActor
}

func (d *defaultState) exitAction() StateAction {
	return d.exitActor
}
