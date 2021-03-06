package stringset

import (
	"sync"
)

type Interface interface {
	Add(s string) bool
	Remove(s string) bool
	RemoveAll()
	Contains(s string) bool
	Size() int
	Elements() []string
	Difference(other Interface) []string
}

type setType struct {
	set map[string]bool
}

type synchronized struct {
	*setType
	syncMutex sync.Mutex
}

func New() Interface {
	return newSetType()
}

func NewSynchronized() Interface {
	return &synchronized{setType: newSetType()}
}

func newSetType() *setType {
	return &setType{set: make(map[string]bool)}
}

func (set *setType) Add(s string) bool {
	_, found := set.set[s]
	set.set[s] = true

	return !found
}

func (set *setType) Contains(s string) bool {
	_, found := set.set[s]
	return found
}

func (set *setType) Size() int {
	return len(set.set)
}

func (set *setType) Remove(s string) bool {
	_, found := set.set[s]
	delete(set.set, s)

	return found
}

func (set *setType) RemoveAll() {
	for v := range set.set {
		delete(set.set, v)
	}
}

func (set *setType) Elements() []string {
	elements := make([]string, len(set.set))
	i := 0

	for v := range set.set {
		elements[i] = v
		i++
	}

	return elements
}

func (set *setType) Difference(other Interface) []string {
	if otherSet, ok := other.(*setType); ok {
		return diff(otherSet, set.Contains)
	}

	return diff2(other.Elements(), set.Contains)
}

func diff(set *setType, contains func(string) bool) []string {
	notFound := []string{}

	for item := range set.set {
		if !contains(item) {
			notFound = append(notFound, item)
		}
	}

	return notFound
}

func diff2(s []string, contains func(string) bool) []string {
	notFound := []string{}

	for _, item := range s {
		if !contains(item) {
			notFound = append(notFound, item)
		}
	}

	return notFound
}

func (set *synchronized) Add(s string) bool {
	set.syncMutex.Lock()
	defer set.syncMutex.Unlock()

	return set.setType.Add(s)
}

func (set *synchronized) Contains(s string) bool {
	set.syncMutex.Lock()
	defer set.syncMutex.Unlock()

	return set.containsSafe(s)
}

func (set *synchronized) containsSafe(s string) bool {
	return set.setType.Contains(s)
}

func (set *synchronized) Size() int {
	set.syncMutex.Lock()
	defer set.syncMutex.Unlock()

	return set.setType.Size()
}

func (set *synchronized) Remove(s string) bool {
	set.syncMutex.Lock()
	defer set.syncMutex.Unlock()

	return set.setType.Remove(s)
}

func (set *synchronized) RemoveAll() {
	set.syncMutex.Lock()
	defer set.syncMutex.Unlock()

	set.setType.RemoveAll()
}

func (set *synchronized) Elements() []string {
	set.syncMutex.Lock()
	defer set.syncMutex.Unlock()

	return set.setType.Elements()
}

func (set *synchronized) Difference(other Interface) []string {
	set.syncMutex.Lock()
	defer set.syncMutex.Unlock()

	if otherSet, ok := other.(*synchronized); ok {
		otherSet.syncMutex.Lock()
		defer otherSet.syncMutex.Unlock()

		return diff(otherSet.setType, set.containsSafe)
	}

	return diff2(other.Elements(), set.containsSafe)
}
