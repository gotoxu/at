package at

import (
	"fmt"
	"testing"
	"time"

	"github.com/gotoxu/assert"
)

func TestAddJob(t *testing.T) {
	at := New()
	at.AddFunc(time.Now().Add(5*time.Second), func() {
		fmt.Println("Hello world")
	})

	assert.DeepEqual(t, at.entries.Len(), 1)
}

func TestAddJobAfterStarted(t *testing.T) {
	at := New()
	at.Start()
	assert.True(t, at.running)

	at.AddFunc(time.Now().Add(5*time.Second), func() {
		fmt.Println("Hello world")
	})

	assert.DeepEqual(t, at.entries.Len(), 1)
}

func TestRun(t *testing.T) {
	at := New()
	at.AddFunc(time.Now().Add(5*time.Second), func() {
		fmt.Println("Hello world")
	})

	at.Start()
	time.Sleep(6 * time.Second)
	assert.DeepEqual(t, at.entries.Len(), 0)
}
