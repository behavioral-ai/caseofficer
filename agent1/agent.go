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
	Name            = "resiliency:agent/caseofficer"
	defaultDuration = time.Second * 15
)

type caseOfficer struct {
	running bool
	uri     string
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
func newAgent(origin common.Origin, resolver collective.Resolution, dispatcher messaging.Dispatcher) *caseOfficer {
	c := new(caseOfficer)
	c.uri = agentUri(origin)
	c.origin = origin

	c.serviceAgents = messaging.NewExchange()

	c.ticker = messaging.NewTicker(messaging.Emissary, defaultDuration)
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
func (c *caseOfficer) String() string { return c.Uri() }

// Uri - agent identifier
func (c *caseOfficer) Uri() string { return c.uri }

// Name - agent class
func (c *caseOfficer) Name() string { return Name }

// Message - message the agent
func (c *caseOfficer) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	c.emissary.C <- m
}

// Run - run the agent
func (c *caseOfficer) Run() {
	if c.running {
		return
	}
	c.running = true
	go emissaryAttend(c, timeseries1.Assignments, agent.New)
}

// Shutdown - shutdown the agent
func (c *caseOfficer) Shutdown() {
	if !c.emissary.IsClosed() {
		c.emissary.C <- messaging.Shutdown
	}
}

func (c *caseOfficer) dispatch(channel any, event string) {
	messaging.Dispatch(c, c.dispatcher, channel, event)
}

func (c *caseOfficer) startup() {
	c.ticker.Start(-1)
}

func (c *caseOfficer) finalize() {
	if !c.emissary.IsClosed() {
		c.emissary.Close()
	}
	if !c.ticker.IsStopped() {
		c.ticker.Stop()
	}
	c.serviceAgents.Shutdown()
}

func (c *caseOfficer) reviseTicker(newDuration time.Duration) {
	c.ticker.Start(newDuration)
}
