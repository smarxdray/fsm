package main

import (
	"fmt"
	"github.com/smarxdray/fsm"
)

func main() {
	/*
	   Create an FSM instance and actions can be registered here
	*/
	sm := fsm.NewFSM(
		// fsm acts according to the transition rule
		fsm.NewTransitionRule(fsm.Transitions{
			{"state0", "event01", "state1"},
			{"state0", "event02", "state2"},
			{"state0", "event03", "state3"},
			{"state2", "event24", "state4"},
			{"state2", "event25", "state5"},
			{"state5", "event56", "state6"},
			{"state5", "event57", "state7"},
		}),
	).
		// register exitAction for state0
		RegisterExitAction(
			"state0",
			func(ctx fsm.Context, event fsm.Event) fsm.Event {
				// within the action, you can:

				// 1. get extra data from the event triggering the current action
				triggerEventData := event.GetEventData()
				fmt.Printf("state0.exitAction: triggerEvent: %+v\n", triggerEventData)

				// 2. record info into the current state (state0 here)
				ctx.WriteIntoCurrentState(map[string]string{
					"triggerEvent": fmt.Sprintf("%+v", triggerEventData),
				})
				// 3. create a new event & set detail info
				emittedEvent := fsm.NewEvent("event02",
					map[string]string{
						"eventName":    "event02",
						"eventEmitter": "state0",
					})
				fmt.Printf("state0.exitAction: emittedEvent: %+v\n", emittedEvent.GetEventData())
				// 4. emit the new event
				return emittedEvent
			},
		).RegisterEntryAction(
		"state2",
		func(ctx fsm.Context, event fsm.Event) fsm.Event {
			// emit event without extra data
			return fsm.NewEvent("event25")
		},
	).RegisterEntryAction(
		"state5",
		func(ctx fsm.Context, event fsm.Event) fsm.Event {
			return fsm.NewEvent("event56")
		},
	).RegisterExitAction(
		"state5",
		func(ctx fsm.Context, event fsm.Event) fsm.Event {
			// emit nothing, which will not trigger any state transition
			return nil
		},
	)
	/*
	   Choose a state to start and trigger the state transition
	*/
	fmt.Printf("sm starts from: state0\n")
	state := sm.
		From("state0").
		Transfer(
			"event01",
			map[string]string{
				"eventName":    "event01",
				"eventEmitter": "external",
			},
		).
		State()
	fmt.Printf("sm stops at: %s\n", state)
}
