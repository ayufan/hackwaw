package main

import (
	. "github.com/onsi/ginkgo"
	"sync"
)

func BeforeGroup(body func()) bool {
	once := &sync.Once{}
	return BeforeEach(func() {
		once.Do(body)
	})
}

func JustBeforeGroup(body func()) bool {
	once := &sync.Once{}
	return JustBeforeEach(func() {
		once.Do(body)
	})
}

func AfterGroup(body func()) bool {
	// TODO: Unfortunatelly gomega doesn't support group contexes
	// We use it to speed-up testing of application
	// The time required to start the application is significant
	// We use Groups in context of Describe to define that we Describe single scenario,
	// that shares application
	return true
}
