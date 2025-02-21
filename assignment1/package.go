package assignment1

import (
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
)

const (
	PkgPath = "github/behavioral-ai/caseofficer/assignment1"
)

func GetRegion(origin common.Origin) ([]HostEntry, *messaging.Status) {
	if origin.Region == common.WestRegion {
		return westData, messaging.StatusOK()
	}
	if origin.Region == common.CentralRegion {
		return centralData, messaging.StatusOK()
	}
	return []HostEntry{}, messaging.StatusOK()
}

// Assignments - assignments functions struct
type Assignments struct {
	All func(origin common.Origin) ([]HostEntry, *messaging.Status)
	New func(origin common.Origin) ([]HostEntry, *messaging.Status)
}

var Entries = func() *Assignments {
	return &Assignments{
		All: func(origin common.Origin) ([]HostEntry, *messaging.Status) {
			entry, status := GetRegion(origin)
			return entry, status
		},
		New: func(origin common.Origin) ([]HostEntry, *messaging.Status) {
			return nil, messaging.StatusNotFound()
		},
	}
}()
