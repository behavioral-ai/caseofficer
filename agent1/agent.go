package agent1

import (
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/collective"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/domain/metrics1"
	"github.com/behavioral-ai/domain/timeseries1"
	"github.com/behavioral-ai/operative/agent"
	"strconv"
	"time"
)

const (
	Name        = "resiliency:agent/caseofficer"
	minDuration = time.Second * 5
	maxDuration = time.Second * 10
)

type agentT struct {
	running bool
	uri     string
	traffic string
	origin  common.Origin

	serviceAgents *messaging.Exchange
	ticker        *messaging.Ticker
	emissary      *messaging.Channel
	resolver      collective.Resolution
	dispatcher    messaging.Dispatcher
}

func agentUri(origin common.Origin) string {
	return fmt.Sprintf("%v%v#%v", Name, strconv.Itoa(version), origin)
}

// New - create a new case officer agent
func New(origin common.Origin, resolver collective.Resolution, dispatcher messaging.Dispatcher) messaging.Agent {
	return newAgent(origin, resolver, dispatcher)
}

// newAgent - create a new case officer agent
func newAgent(origin common.Origin, resolver collective.Resolution, dispatcher messaging.Dispatcher) *agentT {
	a := new(agentT)
	a.uri = agentUri(origin)
	a.origin = origin

	a.serviceAgents = messaging.NewExchange()

	a.ticker = messaging.NewTicker(messaging.Emissary, maxDuration)
	a.emissary = messaging.NewEmissaryChannel()
	if resolver == nil {
		a.resolver = collective.Resolver
	} else {
		a.resolver = resolver
	}
	a.dispatcher = dispatcher
	return a
}

// String - identity
func (a *agentT) String() string { return a.Uri() }

// Uri - agent identifier
func (a *agentT) Uri() string { return a.uri }

// Name - agent class
func (a *agentT) Name() string { return Name }

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
	go emissaryAttend(a, timeseries1.Assignments, agent.New, 0)
}

// Shutdown - shutdown the agent
func (a *agentT) Shutdown() {
	if !a.emissary.IsClosed() {
		a.emissary.C <- messaging.Shutdown
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

func (a *agentT) reviseTicker(duration time.Duration) {
	if duration != 0 {
		a.ticker.Start(duration)
		return
	}
	p, status := collective.Resolve[metrics1.TrafficProfile](metrics1.ProfileName, 1, collective.Resolver)
	if !status.OK() {
		a.resolver.Notify(status)
		return
	}
	traffic := p.Now()
	if traffic == metrics1.TrafficMedium || traffic == a.traffic {
		return
	}
	if traffic == metrics1.TrafficLow {
		a.ticker.Start(minDuration)
	} else {
		a.ticker.Start(maxDuration)
	}
	a.traffic = traffic
}

/*
func (a *agentT) startup() {
	a.ticker.Start(-1)
}
*/
