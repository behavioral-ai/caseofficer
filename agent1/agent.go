package agent1

import (
	"fmt"
	"github.com/behavioral-ai/caseofficer/assignment1"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/collective"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/operative/agent"
	"strconv"
	"time"
)

const (
	Name               = "resiliency:agent/caseofficer/agent"
	assignmentDuration = time.Second * 15
)

type caseOfficer struct {
	running       bool
	uri           string
	origin        common.Origin
	serviceAgents *messaging.Exchange

	ticker     *messaging.Ticker
	emissary   *messaging.Channel
	notifier   messaging.NotifyFunc
	dispatcher messaging.Dispatcher
}

func AgentUri(origin common.Origin) string {
	return fmt.Sprintf("%v%v#%v", Name, strconv.Itoa(version), origin)
}

// New - create a new case officer agent
func New(origin common.Origin, notifier messaging.NotifyFunc, dispatcher messaging.Dispatcher) messaging.Agent {
	return newAgent(origin, notifier, dispatcher)
}

// newAgent - create a new case officer agent
func newAgent(origin common.Origin, notifier messaging.NotifyFunc, dispatcher messaging.Dispatcher) *caseOfficer {
	c := new(caseOfficer)
	c.uri = AgentUri(origin)
	c.origin = origin
	c.serviceAgents = messaging.NewExchange()

	c.ticker = messaging.NewPrimaryTicker(assignmentDuration)
	c.emissary = messaging.NewEmissaryChannel(true)
	c.notifier = notifier
	if c.notifier == nil {
		c.notifier = collective.Resolver.Notify
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
	go emissaryAttend(c, assignment1.Entries, agent.New)
}

// Shutdown - shutdown the agent
func (c *caseOfficer) Shutdown() {
	if !c.running {
		return
	}
	c.running = false
	c.serviceAgents.Shutdown()
	c.emissary.C <- messaging.Shutdown
}

func (c *caseOfficer) notify(e messaging.Event) {
	c.notifier(e)
}

func (c *caseOfficer) dispatch(channel any, event string) {
	messaging.Dispatch(c, c.dispatcher, channel, event)
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
