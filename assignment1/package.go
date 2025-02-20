package assignment1

import (
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
)

const (
	PkgPath    = "github/behavioral-ai/caseofficer/assignment1"
	WestRegion = "us-west1"
	WestZoneA  = "w-a"
	WestZoneB  = "w-b"

	CentralRegion = "us-central1"
	CentralZoneA  = "c-a"
	CentralZoneB  = "c-b"

	EastRegion = "us-east1"
	EastZoneA  = "e-a"
	EastZoneB  = "e-b"
)

func GetRegion(origin common.Origin) ([]HostEntry, *messaging.Status) {
	if origin.Region == WestRegion {
		return westData, messaging.StatusOK()
	}
	if origin.Region == CentralRegion {
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
