# A Finite State Machine implementation : )

## Example
![image](https://user-images.githubusercontent.com/24579932/161443875-8e2de082-5ff5-40a4-9975-d37f2571e2b9.png)
### 1. Basic FSM
```go

fsm.NewFSM(
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
From("state0"). // start from state0
Transfer("event01").
State() // stop at state1
```
### 2. FSM with Action
State can have Actions, which can emit Event (which may carry some data) that may trigger a follow-up transition. Because of these internal events, an external event may actually result in mutiple transitions.
```go

fsm.NewFSM(
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
RegisterExitAction( // register ExitAction for state0
    "state0",
    func(ctx fsm.Context, event fsm.Event) fsm.Event {
        // within the action, you can:
        
        // 1. get extra data from the event triggering the current action
        triggerEventData := event.GetEventData()
        
        // 2. record info into the current state (state0 here)
        ctx.WriteIntoCurrentState(map[string]string){
            "triggerEvent": fmt.Sprintf("%+v", triggerEventData),
        })
        
        // 3. create a new event, which can carry some extra data
        emittedEvent := fsm.NewEvent("event02",
            map[string]string{
                "eventName":        "event02",
                "eventEmitter":     "state0",
            })
        
        // 4. emit the new event, which may trigger another state transition
        return emittedEvent
    }, 
).
RegisterEntryAction( // register EntryAction for state2
    "state2",
    func(ctx fsm.Context, event fsm.Event) fsm.Event {
        // emit event without extra data
        return fsm.NewEvent("event25")
    },
).
RegisterEntryAction( // register EntryAction for state5
    "state5",
    func(ctx fsm.Context, event fsm.Event) fsm.Event {
        return fsm.NewEvent("event56")
    },
).
RegisterExitAction( // register ExitAction for state5
    "state5",
    func(ctx fsm.Context, event fsm.Event) fsm.Event {
        // emit nothing, which will not trigger any state transition
        return nil
    },
).
From("state0"). // start from state0 
Transfer(
      "event01",
      map[string]string{ // carry some extra data
          "eventName":    "event01",
          "eventEmitter": "external",
      }),
).
State() // stop at state6
```
### 3. Asynchronous FSM
AsyncFSM runs continuously in a standalone goroutine, which can accept external events in a nonblocking way and provides the same abilities as the synchronous FSM does.
```go
envoy := fsm.NewAsyncFSM( // take an envoy to communicate with the FSM
    fsm.NewTransitionRule(fsm.Transitions{
        {"state0", "event01", "state1"},
        {"state0", "event02", "state2"},
        {"state0", "event03", "state3"},
        {"state2", "event24", "state4"},
        {"state2", "event25", "state5"},
        {"state5", "event56", "state6"},
        {"state5", "event57", "state7"},
    }), 3). // event emitting buffer of size 3
    .RegisterExitAction("state0", func(ctx fsm.Context, event fsm.Event) fsm.Event {
        return NewEvent("event02")
    }).RegisterEntryAction("state2", func(ctx fsm.Context, event fsm.Event) fsm.Event {
        return NewEvent("event25")
    }).RegisterExitAction("state2", func(ctx fsm.Context, event fsm.Event) fsm.Event {
        return nil
    }).RegisterEntryAction("state5", func(ctx fsm.Context, event fsm.Event) fsm.Event {
        return NewEvent("event56")
    }).RegisterExitAction("state5", func(ctx fsm.Context, event fsm.Event) fsm.Event {
        return nil
    }).From("state0"). // start from state0
       Envoy()
    // let the FSM run in a standalone goroutine
    go envoy.Run()
    // emit event asynchronously
    envoy.EmitEvent("event01")
    // wait for a moment for the result of state transition
    time.Sleep(time.Millisecond)
    // stop at state6
    got := envoy.State().Code()
    // stop the FSM running behind the scenes
    envoy.Stop()
```
