package fsm

import (
	"testing"
	"time"
)

func Test_defaultAsyncFSM_Run(t *testing.T) {
	testTransitionRule := NewTransitionRule(Transitions{
		{"state0", "event01", "state1"},
		{"state0", "event02", "state2"},
		{"state0", "event03", "state3"},
		{"state2", "event24", "state4"},
		{"state2", "event25", "state5"},
		{"state5", "event56", "state6"},
		{"state5", "event57", "state7"},
	})
	type args struct {
		initState    StateCode
		triggerEvent EventCode
		s2Entry      EventCode
		s2Exit       EventCode
		s5Entry      EventCode
		s5Exit       EventCode
		s0Exit       EventCode
	}
	type wants struct {
		state StateCode
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "from s0 to s1",
			args: args{
				initState:    "state0",
				triggerEvent: "event01",
			},
			wants: wants{
				state: "state1",
			},
		},
		{
			name: "from s0 to s1 with action",
			args: args{
				initState:    "state0",
				triggerEvent: "event01",
				s0Exit:       "event01",
				s2Entry:      "event24",
				s5Entry:      "event57",
			},
			wants: wants{
				state: "state1",
			},
		},
		{
			name: "from s0 to s3",
			args: args{
				initState:    "state0",
				triggerEvent: "event01",
				s0Exit:       "event03",
				s2Entry:      "event24",
				s5Entry:      "event57",
			},
			wants: wants{
				state: "state3",
			},
		},
		{
			name: "from s0 to s4",
			args: args{
				initState:    "state0",
				triggerEvent: "event01",
				s0Exit:       "event02",
				s2Entry:      "event24",
				s5Entry:      "event57",
			},
			wants: wants{
				state: "state4",
			},
		},
		{
			name: "from s0 to s6",
			args: args{
				initState:    "state0",
				triggerEvent: "event01",
				s0Exit:       "event02",
				s2Entry:      "event25",
				s5Entry:      "event56",
			},
			wants: wants{
				state: "state6",
			},
		},
		{
			name: "from s0 to s7",
			args: args{
				initState:    "state0",
				triggerEvent: "event01",
				s0Exit:       "event02",
				s2Entry:      "event25",
				s5Entry:      "event57",
			},
			wants: wants{
				state: "state7",
			},
		},
		{
			name: "from s2 to s5",
			args: args{
				initState:    "state2",
				triggerEvent: "event25",
			},
			wants: wants{
				state: "state5",
			},
		},
		{
			name: "from s2 to s5 by exitActor",
			args: args{
				initState:    "state2",
				triggerEvent: "event24",
				s2Exit:       "event25",
			},
			wants: wants{
				state: "state5",
			},
		},
		{
			name: "from s2 to s6",
			args: args{
				initState:    "state2",
				triggerEvent: "event25",
				s5Entry:      "event56",
			},
			wants: wants{
				state: "state6",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envoy := NewAsyncFSM(testTransitionRule, 3).
				RegisterEntryAction("state2", func(ctx Context, event Event) Event {
					return NewEvent(tt.args.s2Entry)
				}).RegisterExitAction("state2", func(ctx Context, event Event) Event {
				return NewEvent(tt.args.s2Exit)
			}).RegisterEntryAction("state5", func(ctx Context, event Event) Event {
				return NewEvent(tt.args.s5Entry)
			}).RegisterExitAction("state5", func(ctx Context, event Event) Event {
				return NewEvent(tt.args.s5Exit)
			}).RegisterExitAction("state0", func(ctx Context, event Event) Event {
				return NewEvent(tt.args.s0Exit)
			}).From(tt.args.initState).
				From(tt.args.initState).Envoy()
			go envoy.Run()
			envoy.EmitEvent(tt.args.triggerEvent)
			time.Sleep(time.Millisecond)
			got := envoy.State().Code()
			envoy.Stop()
			if got != tt.wants.state {
				t.Errorf("Run() = %v, want %v", got, tt.wants.state)
			}
		})
	}
}
