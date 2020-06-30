package types

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	String   = `[a-zA-Z0-9_]{2,20}`
	ReHandle = regexp.MustCompile(fmt.Sprintf(`^%s$`, String))
)

type Findable interface {
	ElAtIndex(index int) string
	Len() int
}

func (handler AllowedHandles) Find(id string) (handle string, found bool) {
	index := handler.find(id)
	if index == -1 {
		return handle, false
	}
	return handler.Handles[index], true
}

func (handlers AllowedHandles) find(id string) int {
	return FindUtil(handlers, id)
}

func FindUtil(group Findable, el string) int {
	if group.Len() == 0 {
		return -1
	}
	low := 0
	high := group.Len() - 1
	median := 0
	for low <= high {
		median = (low + high) / 2
		switch compare := strings.Compare(group.ElAtIndex(median), el); {
		case compare == 0:
			return median
		case compare == -1:
			low = median + 1
		default:
			high = median - 1
		}
	}
	return -1
}

func (handles AllowedHandles) ElAtIndex(index int) string { return handles.Handles[index] }
func (handles AllowedHandles) Len() int                   { return len(handles.Handles) }
func (handles AllowedHandles) Less(i, j int) bool {
	return strings.Compare(handles.Handles[i], handles.Handles[j]) == -1
}
func (handles AllowedHandles) Swap(i, j int) {
	handles.Handles[i], handles.Handles[j] = handles.Handles[j], handles.Handles[i]
}

var _ sort.Interface = AllowedHandles{}

func (handles AllowedHandles) Sort() AllowedHandles {
	sort.Sort(handles)
	return handles
}
