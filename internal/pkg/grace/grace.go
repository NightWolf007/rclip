// Package grace provides methods to implement graceful shutdown.
package grace

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Grace represents a configuration for graceful shutdown.
// It implements 3 exit points:
// - Done - graceful shutdown, all shutdown tasks are done
// - ForceQuit - when process received second termination signal,
//               Run function completes without waiting for tasks graceful shutdown.
// - Timeout - when it takes too long to graceful shutdown tasks,
//             Run function completes without waiting for tasks graceful shutdown.
type Grace struct {
	// Shutdown function must provide shutdown routine
	// Ex. cancelling context, closing sockets etc.
	// It can be both synchronous or asynchronous (then you should provide Wait function)
	Shutdown func()
	// Wait function should provide routine to wait for graceful shutdown.
	// It is only needed when Shutdown is asynchronous.
	Wait func()
	// Timeout is a maximum time to wait for graceful shutdown before exit.
	// To disable timeout function, leave this value as zero.
	Timeout time.Duration

	// Hooks can be used to print some info to the user or logging.

	// OnShutdown hook will be called right before shutdown call.
	OnShutdown func()
	// OnDone hook runs on graceful shutdown.
	OnDone func()
	// OnForceQuit hook runs when the second termination signal appears.
	OnForceQuit func()
	// OnTimeout hook runs when timeout exceeded.
	// It never runs if Timeout option is zero.
	OnTimeout func()

	// This functions are only used for mocking purposes.
	SignalNotify func(chan<- os.Signal, ...os.Signal)
	TimeAfter    func(time.Duration) <-chan time.Time
}

// Run blocks current thread and waits for SIGINT or SIGTERM signals to start shutdown.
func (g *Grace) Run() {
	termCh := make(chan os.Signal, 2) // nolint:gomnd // Intercept only two termination signals
	g.signalNotify(termCh, syscall.SIGINT, syscall.SIGTERM)

	<-termCh

	doneCh := make(chan struct{})

	go func() {
		g.onShutdown()
		g.shutdown()
		g.wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		g.onDone()

		return
	case <-termCh:
		g.onForceQuit()

		return
	case <-g.timeAfter(g.Timeout):
		g.onTimeout()

		return
	}
}

func (g *Grace) shutdown() {
	if g.Shutdown != nil {
		g.Shutdown()
	}
}

func (g *Grace) wait() {
	if g.Wait != nil {
		g.Wait()
	}
}

func (g *Grace) onShutdown() {
	if g.OnShutdown != nil {
		g.OnShutdown()
	}
}

func (g *Grace) onDone() {
	if g.OnDone != nil {
		g.OnDone()
	}
}

func (g *Grace) onForceQuit() {
	if g.OnForceQuit != nil {
		g.OnForceQuit()
	}
}

func (g *Grace) onTimeout() {
	if g.OnTimeout != nil {
		g.OnTimeout()
	}
}

func (g *Grace) signalNotify(c chan<- os.Signal, sig ...os.Signal) {
	if g.SignalNotify != nil {
		g.SignalNotify(c, sig...)

		return
	}

	signal.Notify(c, sig...)
}

func (g *Grace) timeAfter(d time.Duration) <-chan time.Time {
	if g.Timeout == 0 {
		return make(chan time.Time)
	}

	if g.TimeAfter != nil {
		return g.TimeAfter(d)
	}

	return time.After(d)
}
