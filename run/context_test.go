package run_test

import (
	"sync"
	"testing"
	"time"

	"github.com/anchore/go-make/lang"
	"github.com/anchore/go-make/require"
	"github.com/anchore/go-make/run"
)

func Test_Cancel(t *testing.T) {
	panicked := false
	startTime := time.Now()
	sleepDone := sync.WaitGroup{}
	sleepDone.Add(1)
	sleepReady := sync.WaitGroup{}
	sleepReady.Add(1)
	go func() {
		defer sleepDone.Done()
		panicked = nil != lang.Catch(func() {
			sleepReady.Done()
			run.Command("sleep", run.Args("60")) // sleep for 5 seconds
		})
	}()
	go func() {
		sleepReady.Wait()
		// let the sleep command run
		time.Sleep(100 * time.Millisecond)
		run.Cancel()
	}()
	sleepDone.Wait() // will wait for 1 minute if not canceled

	elapsed := time.Since(startTime)
	t.Log(elapsed)

	require.True(t, panicked)
	require.True(t, elapsed < 1*time.Second)

	goOutput := run.Command("go", run.Args("help"))
	require.True(t, goOutput != "")
}
