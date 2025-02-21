package agent1

import (
	"github.com/behavioral-ai/caseofficer/assignment1"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/operative/agent"
	"time"
)

const (
	Class              = "case-officer"
	assignmentDuration = time.Second * 15
)

type caseOfficer struct {
	running bool
	uri     string
	origin  common.Origin

	ticker        *messaging.Ticker
	emissary      *messaging.Channel
	serviceAgents *messaging.Exchange
	notifier      messaging.NotifyFunc
	dispatcher    messaging.Dispatcher
}

func AgentUri(origin common.Origin) string {
	return origin.Uri(Class)
}

// NewAgent - create a new case officer agent
func NewAgent(origin common.Origin, notifier messaging.NotifyFunc, dispatcher messaging.Dispatcher) messaging.Agent {
	return newAgent(origin, notifier, dispatcher)
}

// newAgent - create a new case officer agent
func newAgent(origin common.Origin, notifier messaging.NotifyFunc, dispatcher messaging.Dispatcher) *caseOfficer {
	c := new(caseOfficer)
	c.uri = AgentUri(origin)
	c.origin = origin
	c.ticker = messaging.NewPrimaryTicker(assignmentDuration)
	c.emissary = messaging.NewEmissaryChannel(true)
	c.notifier = notifier
	c.serviceAgents = messaging.NewExchange()

	c.dispatcher = dispatcher
	return c
}

// String - identity
func (c *caseOfficer) String() string { return c.Uri() }

// Uri - agent identifier
func (c *caseOfficer) Uri() string { return c.uri }

// Name - agent class
func (c *caseOfficer) Name() string { return c.uri }

// Message - message the agent
func (c *caseOfficer) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	c.emissary.C <- m
}

// Notify - notifier
func (c *caseOfficer) Notify(status *messaging.Status) {
	if c.notifier != nil {
		c.notifier(status)
	}
}

// Run - run the agent
func (c *caseOfficer) Run() {
	if c.running {
		return
	}
	c.running = true
	go emissaryAttend(c, assignment1.Entries, agent.New)
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

func (c *caseOfficer) dispatch(channel any, event string) {
	if c.dispatcher == nil || channel == nil {
		return
	}
	if ch, ok := channel.(*messaging.Channel); ok {
		c.dispatcher.Dispatch(c, ch.Name(), event)
		return
	}
	if t, ok := channel.(*messaging.Ticker); ok {
		c.dispatcher.Dispatch(c, t.Name(), event)
	}
}

/*
func (c *caseOfficer) dispatch(event string) {
	if c.global != nil {
		c.global.Dispatch(c, messaging.EmissaryChannel, event, "")
	}
	if c.local != nil {
		c.local.dispatch(c, event)
	}
}

*/

/*
func (c *caseOfficer) setup(event string) {
	if c.local != nil {
		c.local.setup(c, event)
	}
}


*/
