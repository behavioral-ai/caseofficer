package agent1

import (
	"fmt"
	"github.com/behavioral-ai/caseofficer/assignment1"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/operative/agent"
	"time"
)

var (
	shutdown = messaging.NewControlMessage("", "", messaging.ShutdownEvent)
	//dataChange = messaging.NewControlMessage("", "", messaging.DataChangeEvent)
)

func ExampleEmissary() {
	ch := make(chan struct{})
	traceDispatch := messaging.NewTraceDispatcher([]string{messaging.StartupEvent, messaging.ShutdownEvent}, "")
	officer := newAgent(common.Origin{Region: common.WestRegion}, func(status *messaging.Status) {
		fmt.Printf("test: Agent() -> [status:%v]\n", status)
	}, traceDispatch)
	//dataChange.SetContent(guidance.ContentTypeCalendar, guidance.NewProcessingCalendar())

	go func() {
		go emissaryAttend(officer, assignment1.Entries, agent.New)
		//agent.Message(dataChange)
		time.Sleep(time.Minute * 1)
		officer.Message(shutdown)

		fmt.Printf("test: emissaryAttend() -> [finalized:%v]\n", officer.IsFinalized())
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//test: emissaryAttend() -> [finalized:true]

}
