package agent1

import (
	"fmt"
	"github.com/behavioral-ai/caseofficer/assignment1"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/test"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/domain/guidance"
	"github.com/behavioral-ai/operative/agent1"
	"time"
)

var (
	shutdown   = messaging.NewControlMessage("", "", messaging.ShutdownEvent)
	dataChange = messaging.NewControlMessage("", "", messaging.DataChangeEvent)
)

func ExampleEmissary() {
	ch := make(chan struct{})
	traceDispatch := messaging.NewTraceDispatcher([]string{messaging.StartupEvent, messaging.ShutdownEvent}, "")
	agent := newAgent(common.Origin{Region: guidance.WestRegion}, test.NewAgent("agent-test").Notify, traceDispatch)
	dataChange.SetContent(guidance.ContentTypeCalendar, guidance.NewProcessingCalendar())

	go func() {
		go emissaryAttend(agent, assignment1.Entries, agent1.New, nil)
		agent.Message(dataChange)
		time.Sleep(time.Minute * 1)
		agent.Message(shutdown)

		fmt.Printf("test: emissaryAttend() -> [finalized:%v]\n", agent.IsFinalized())
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//test: emissaryAttend() -> [finalized:true]

}
