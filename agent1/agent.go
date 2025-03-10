package agent1

import (
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/collective"
	"github.com/behavioral-ai/domain/common"
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
	running  bool
	uri      string
	origin   common.Origin
	duration time.Duration

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
	c := new(agentT)
	c.uri = agentUri(origin)
	c.origin = origin

	c.serviceAgents = messaging.NewExchange()

	c.duration = minDuration
	c.ticker = messaging.NewTicker(messaging.Emissary, c.duration)
	c.emissary = messaging.NewEmissaryChannel()
	if resolver == nil {
		c.resolver = collective.Resolver
	} else {
		c.resolver = resolver
	}
	c.dispatcher = dispatcher
	return c
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
	go emissaryAttend(a, timeseries1.Assignments, agent.New)
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

func (a *agentT) startup() {
	a.ticker.Start(-1)
}

func (a *agentT) finalize() {
	if !a.emissary.IsClosed() {
		a.emissary.Close()
	}
	if !a.ticker.IsStopped() {
		a.ticker.Stop()
	}
	a.serviceAgents.Shutdown()
}

func (a *agentT) reviseTicker(newDuration time.Duration) {
	a.ticker.Start(newDuration)
}
