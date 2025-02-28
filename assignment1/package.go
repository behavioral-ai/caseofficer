package assignment1

import (
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
)

const (
	PkgPath = "github/behavioral-ai/caseofficer/assignment1"
)

// Entry - host
type Entry struct {
	//CreatedTS time.Time     `json:"created-ts"`
	Origin common.Origin `json:"origin"`
}

type SelectAssignments func(origin common.Origin) ([]Entry, *messaging.Status)

// Assignments - assignments functions struct
type Assignments struct {
	All func(origin common.Origin) ([]Entry, *messaging.Status)
	New func(origin common.Origin) ([]Entry, *messaging.Status)
}

var Entries = func() *Assignments {
	return &Assignments{
		All: func(origin common.Origin) ([]Entry, *messaging.Status) {
			return get(origin)
		},
		New: func(origin common.Origin) ([]Entry, *messaging.Status) {
			return nil, messaging.StatusNotFound()
		},
	}
}()
