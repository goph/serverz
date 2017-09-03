package serverz

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

// Daemon is a long-running process in the background.
type Daemon interface {
	// Run has a loop inside and quits when it gets a signal from the quit channel.
	Run(quit <-chan struct{}) error
}

// DaemonServer is a server without network, running a daemon.
type DaemonServer struct {
	Daemon Daemon

	mu       sync.Mutex
	doneChan chan struct{}
	quitChan chan struct{}
}

func (s *DaemonServer) getQuitChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getQuitChanLocked()
}

func (s *DaemonServer) getQuitChanLocked() chan struct{} {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

func (s *DaemonServer) closeQuitChanLocked() {
	ch := s.getDoneChanLocked()

	select {
	case <-ch:
		// Already closed. Don't close again.

	default:
		// Safe to close here. We're the only closer,
		// guarded by s.mu.
		close(ch)
	}
}

func (s *DaemonServer) getDoneChanLocked() chan struct{} {
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

func (s *DaemonServer) closeDoneChanLocked() {
	ch := s.getDoneChanLocked()

	select {
	case <-ch:
		// Already closed. Don't close again.

	default:
		// Safe to close here. We're the only closer,
		// guarded by s.mu.
		close(ch)
	}
}

// Serve starts the daemon.
func (s *DaemonServer) Serve(l net.Listener) error {
	if s.Daemon == nil {
		return errors.New("No daemon is specified")
	}

	err := s.Daemon.Run(s.getQuitChan())

	s.mu.Lock()
	s.closeDoneChanLocked()
	s.mu.Unlock()

	return err
}

// Shutdown gracefully stops the daemon by sending a quit signal to it and waiting for it to be done.
func (s *DaemonServer) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	s.closeQuitChanLocked()
	s.mu.Unlock()

	select {
	case <-s.doneChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close sends a quit signal to the daemon.
func (s *DaemonServer) Close() error {
	s.mu.Lock()
	s.closeQuitChanLocked()
	s.mu.Unlock()

	return nil
}

// CronJob is ran scheduled by the CronDaemon.
type CronJob interface {
	// Run runs the job.
	Run() error
}

// CronDaemon is a daemon with an internal scheduler for a CronJob.
type CronDaemon struct {
	Job    CronJob
	Ticker *time.Ticker
}

// Run implements the Daemon interface.
func (d *CronDaemon) Run(quit <-chan struct{}) error {
	if d.Job == nil {
		return errors.New("No job is specified")
	}

	if d.Ticker == nil {
		return d.Job.Run()
	}

	errChan := make(chan error)

	for {
		select {
		case <-quit:
			return nil

		case err := <-errChan:
			return err

		case <-d.Ticker.C:
			go func() {
				err := d.Job.Run()

				if err != nil {
					errChan <- err
				}
			}()
		}
	}
}
