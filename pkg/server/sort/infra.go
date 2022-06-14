package sort

import (
	"sort"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// BuildStartDesc sorts builds by start time descending (most recent first).
// For the opposite, use sort.Reverse.
type InfraStartDesc []*pb.Infra

func (s InfraStartDesc) Len() int      { return len(s) }
func (s InfraStartDesc) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s InfraStartDesc) Less(i, j int) bool {
	t1 := s[i].Status.StartTime.AsTime()
	t2 := s[j].Status.StartTime.AsTime()

	return t2.Before(t1)
}

var (
	_ sort.Interface = (InfraStartDesc)(nil)
)
