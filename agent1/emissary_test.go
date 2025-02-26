package agent1

import (
	"github.com/behavioral-ai/caseofficer/assignment1"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/operative/agent"
	"time"
)

var (
	shutdown = messaging.NewMessage(messaging.EmissaryChannel, messaging.ShutdownEvent)
	//dataChange = messaging.NewControlMessage("", "", messaging.DataChangeEvent)
)

func ExampleEmissary() {
	ch := make(chan struct{})
	traceDispatch := messaging.NewTraceDispatcher()
	officer := newAgent(nil, common.Origin{Region: common.WestRegion}, messaging.Notify, traceDispatch)
	//dataChange.SetContent(guidance.ContentTypeCalendar, guidance.NewProcessingCalendar())

	go func() {
		go emissaryAttend(officer, assignment1.Entries, agent.New)
		//agent.Message(dataChange)
		time.Sleep(time.Minute * 1)
		officer.Message(shutdown)

		//fmt.Printf("test: emissaryAttend() -> [finalized:%v]\n", officer.IsFinalized())
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//test: emissaryAttend() -> [finalized:true]

}
