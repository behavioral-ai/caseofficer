package agent1

import (
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/domain/content"
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
	resolver      content.Resolution
	dispatcher    messaging.Dispatcher
}

func agentUri(origin common.Origin) string {
	return fmt.Sprintf("%v%v#%v", NamespaceName, strconv.Itoa(version), origin)
}

// New - create a new case officer agent
func New(origin common.Origin, resolver content.Resolution, dispatcher messaging.Dispatcher) messaging.Agent {
	return newAgent(origin, resolver, dispatcher)
}

// newAgent - create a new case officer agent
func newAgent(origin common.Origin, resolver content.Resolution, dispatcher messaging.Dispatcher) *agentT {
	a := new(agentT)
	a.uri = agentUri(origin)
	a.origin = origin

	a.serviceAgents = messaging.NewExchange()

	a.ticker = messaging.NewTicker(messaging.Emissary, maxDuration)
	a.emissary = messaging.NewEmissaryChannel()
	if resolver == nil {
		a.resolver = content.Resolver
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
	go emissaryAttend(a, timeseries1.Assignments, operative.New, nil)
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

func (a *agentT) reviseTicker(s messaging.Spanner) {
	if s != nil {
		dur := s.Span()
		a.ticker.Start(dur)
		a.resolver.Notify(messaging.NewStatusMessage(http.StatusOK, fmt.Sprintf("revised ticker -> traffic: %v duration: %v", a.traffic, dur), a.uri))
		return
	}
	p, status := content.Resolve[metrics1.TrafficProfile](metrics1.ProfileName, 1, a.resolver)
	if !status.OK() {
		a.ticker.Start(maxDuration)
		if status.NotFound() {
			status.SetAgent(a.Uri())
		}
		a.resolver.Notify(status)
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
	a.resolver.Notify(messaging.NewStatusMessage(http.StatusOK, fmt.Sprintf("revised ticker -> traffic: %v duration: %v", a.traffic, dur), a.uri))
}

/*
func (a *agentT) startup() {
	a.ticker.Start(-1)
}
*/
