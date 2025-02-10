package agent1

import (
	"github.com/behavioral-ai/core/aspect"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/domain/guidance"
	"time"
)

const (
	Class              = "case-officer"
	assignmentDuration = time.Second * 15
)

type caseOfficer struct {
	running bool
	agentId string
	origin  common.Origin

	ticker        *messaging.Ticker
	emissary      *messaging.Channel
	serviceAgents *messaging.Exchange
	handler       messaging.OpsAgent
	global        messaging.Dispatcher
	local         dispatcher
}

func AgentUri(origin common.Origin) string {
	return origin.Uri(Class)
}

// NewAgent - create a new case officer agent
func NewAgent(origin common.Origin, handler messaging.OpsAgent, global messaging.Dispatcher) messaging.OpsAgent {
	return newAgent(origin, handler, global, newDispatcher(false))
}

// newAgent - create a new case officer agent
func newAgent(origin common.Origin, handler messaging.OpsAgent, global messaging.Dispatcher, local dispatcher) *caseOfficer {
	c := new(caseOfficer)
	c.agentId = AgentUri(origin)
	c.origin = origin
	c.ticker = messaging.NewPrimaryTicker(assignmentDuration)
	c.emissary = messaging.NewEmissaryChannel(true)
	c.handler = handler
	c.serviceAgents = messaging.NewExchange()
	c.local = local
	c.global = global
	return c
}

// String - identity
func (c *caseOfficer) String() string { return c.Uri() }

// Uri - agent identifier
func (c *caseOfficer) Uri() string { return c.agentId }

// Message - message the agent
func (c *caseOfficer) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	c.emissary.C <- m
}

// Notify - notifier
func (c *caseOfficer) Notify(status *aspect.Status) *aspect.Status { return c.handler.Notify(status) }

// Trace - activity tracing
func (c *caseOfficer) Trace(agent messaging.Agent, channel, event, activity string) {
	c.handler.Trace(agent, channel, event, activity)
}

// Run - run the agent
func (c *caseOfficer) Run() {
	if c.running {
		return
	}
	c.running = true
	go emissaryAttend(c, guidance.Assign, nil, nil)
}

// Shutdown - shutdown the agent
func (c *caseOfficer) Shutdown() {
	if !c.running {
		return
	}
	c.running = false
	msg := messaging.NewControlMessage(c.Uri(), c.Uri(), messaging.ShutdownEvent)
	c.serviceAgents.Shutdown()
	c.emissary.C <- msg
}

func (c *caseOfficer) IsFinalized() bool {
	return c.emissary.IsFinalized() && c.ticker.IsFinalized() && c.serviceAgents.IsFinalized()
}

func (c *caseOfficer) startup() {
	c.ticker.Start(-1)
}

func (c *caseOfficer) finalize() {
	c.emissary.Close()
	c.ticker.Stop()
	c.serviceAgents.Shutdown()
}

func (c *caseOfficer) reviseTicker(newDuration time.Duration) {
	c.ticker.Start(newDuration)
}

func (c *caseOfficer) setup(event string) {
	if c.local != nil {
		c.local.setup(c, event)
	}
}

func (c *caseOfficer) dispatch(event string) {
	if c.global != nil {
		c.global.Dispatch(c, messaging.EmissaryChannel, event, "")
	}
	if c.local != nil {
		c.local.dispatch(c, event)
	}
}
