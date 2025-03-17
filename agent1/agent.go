package agent1

import (
	"fmt"
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/domain/metrics1"
	"github.com/behavioral-ai/domain/timeseries1"
	operative "github.com/behavioral-ai/operative/agent"
	"net/http"
	"strconv"
	"time"
)

const (
	NamespaceName = "resiliency:agent/behavioral-ai/caseofficer"
	minDuration   = time.Second * 5
	maxDuration   = time.Second * 10
)

type agentT struct {
	running bool
	uri     string
	traffic string
	origin  common.Origin

	serviceAgents *messaging.Exchange
	ticker        *messaging.Ticker
	emissary      *messaging.Channel
	notifier      messaging.NotifyFunc
	activity      messaging.ActivityFunc
	dispatcher    messaging.Dispatcher
}

func agentUri(origin common.Origin) string {
	return fmt.Sprintf("%v%v#%v", NamespaceName, strconv.Itoa(version), origin)
}

// New - create a new case officer agent
func New(origin common.Origin, activity messaging.ActivityFunc, notifier messaging.NotifyFunc, dispatcher messaging.Dispatcher) messaging.Agent {
	return newAgent(origin, activity, notifier, dispatcher)
}

// newAgent - create a new case officer agent
func newAgent(origin common.Origin, activity messaging.ActivityFunc, notifier messaging.NotifyFunc, dispatcher messaging.Dispatcher) *agentT {
	a := new(agentT)
	a.uri = agentUri(origin)
	a.origin = origin

	a.serviceAgents = messaging.NewExchange()

	a.ticker = messaging.NewTicker(messaging.Emissary, maxDuration)
	a.emissary = messaging.NewEmissaryChannel()

	a.activity = activity
	a.notifier = notifier
	a.dispatcher = dispatcher
	return a
}

// String - identity
func (a *agentT) String() string { return a.Uri() }

// Uri - agent identifier
func (a *agentT) Uri() string { return a.uri }

// Name - agent class
func (a *agentT) Name() string { return NamespaceName }

// Message - message the agent
func (a *agentT) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	a.emissary.C <- m
}

// Run - run the agent
func (a *agentT) Run() {
	if a.running {
		return
	}
	a.running = true
	go emissaryAttend(a, timeseries1.Assignments, content.Resolver, operative.New, nil)
}

// Shutdown - shutdown the agent
func (a *agentT) Shutdown() {
	if !a.emissary.IsClosed() {
		a.emissary.C <- messaging.Shutdown
	}
}

func (a *agentT) notify(e messaging.NotifyItem) {
	if e == nil {
		return
	}
	if a.notifier != nil {
		a.notifier(e)
	} else {
		event.Agent.Message(messaging.NewNotifyMessage(e))
	}
}

func (a *agentT) addActivity(e messaging.ActivityItem) {
	if a.activity != nil {
		a.activity(e)
	} else {
		event.Agent.Message(messaging.NewActivityMessage(e))
	}
}
func (a *agentT) dispatch(channel any, event string) {
	messaging.Dispatch(a, a.dispatcher, channel, event)
}

func (a *agentT) finalize() {
	a.emissary.Close()
	a.ticker.Stop()
	a.serviceAgents.Shutdown()
}

func (a *agentT) reviseTicker(resolver *content.Resolution, s messaging.Spanner) {
	if s != nil {
		dur := s.Span()
		a.ticker.Start(dur)
		a.notify(messaging.NewStatusMessage(http.StatusOK, fmt.Sprintf("revised ticker -> traffic: %v duration: %v", a.traffic, dur), a.uri))
		return
	}
	p, status := content.Resolve[metrics1.TrafficProfile](metrics1.ProfileName, 1, resolver)
	if !status.OK() {
		a.ticker.Start(maxDuration)
		if status.NotFound() {
			status.SetAgent(a.Uri())
		}
		a.notify(status)
		return
	}
	traffic := p.Now()
	if p.IsMedium(traffic) || traffic == a.traffic {
		return
	}
	var dur time.Duration
	if p.IsLow(traffic) {
		dur = minDuration
	} else {
		dur = maxDuration
	}
	a.ticker.Start(dur)
	a.traffic = traffic
	a.notify(messaging.NewStatusMessage(http.StatusOK, fmt.Sprintf("revised ticker -> traffic: %v duration: %v", a.traffic, dur), a.uri))
}

/*
func (a *agentT) startup() {
	a.ticker.Start(-1)
}
*/
