package serverz_test

import (
	"testing"

	"errors"

	"github.com/Sirupsen/logrus"
	logrus_test "github.com/Sirupsen/logrus/hooks/test"
	"github.com/sagikazarmark/serverz"
)

func CreateShutdownFunc() (serverz.ShutdownHandler, *bool) {
	var called bool

	f := serverz.ShutdownFunc(func() {
		called = true
	})

	return f, &called
}

func CreateOrderedShutdownFuncs(num int) ([]serverz.ShutdownHandler, *[]int) {
	funcs := make([]serverz.ShutdownHandler, num)
	called := []int{}

	for index := 0; index < num; index++ {
		funcs[index] = func(index int) serverz.ShutdownHandler {
			return serverz.ShutdownFunc(func() {
				called = append(called, index+1)
			})
		}(index)
	}

	return funcs, &called
}

func TestShutdownFunc_CallsUnderlyingFunc(t *testing.T) {
	f, called := CreateShutdownFunc()

	nilf := func() error {
		return nil
	}

	if got, want := f(), nilf(); got != want {
		t.Fatalf("wrapped functions are expected to return nil, error received: %v", got)
	}

	if got, want := *called, true; got != want {
		t.Fatal("shutdown: the wrapped function is expected to be called")
	}
}

func TestShutdownRegister(t *testing.T) {
	f, called := CreateShutdownFunc()

	logger, _ := logrus_test.NewNullLogger()
	shutdown := serverz.NewShutdown(logger)

	shutdown.Register(f)
	shutdown.Handle()

	if got, want := *called, true; got != want {
		t.Fatal("shutdown: the shutdown handler is expected to be called")
	}
}

func TestShutdownRegister_ExecutedInOrder(t *testing.T) {
	funcs, called := CreateOrderedShutdownFuncs(2)

	logger, _ := logrus_test.NewNullLogger()
	shutdown := serverz.NewShutdown(logger)

	shutdown.Register(funcs[0], funcs[1])
	shutdown.Handle()

	if got, want := (*called)[0], 1; got != want {
		t.Fatal("shutdown: the first shutdown handler is expected to be called first")
	}

	if got, want := (*called)[1], 2; got != want {
		t.Fatal("shutdown: the second shutdown handler is expected to be called second")
	}
}

func TestShutdownRegisterAsFirst(t *testing.T) {
	funcs, called := CreateOrderedShutdownFuncs(2)

	logger, _ := logrus_test.NewNullLogger()
	shutdown := serverz.NewShutdown(logger)

	shutdown.Register(funcs[1])
	shutdown.RegisterAsFirst(funcs[0])
	shutdown.Handle()

	if got, want := (*called)[0], 1; got != want {
		t.Fatal("shutdown: the first shutdown handler is expected to be called first")
	}

	if got, want := (*called)[1], 2; got != want {
		t.Fatal("shutdown: the second shutdown handler is expected to be called second")
	}
}

func TestShutdownHandle(t *testing.T) {
	logger, hook := logrus_test.NewNullLogger()
	shutdown := serverz.NewShutdown(logger)

	shutdown.Handle()

	if got, want, gotLevel, wantLevel := hook.Entries[0].Message, "Shutting down", hook.Entries[0].Level, logrus.InfoLevel; got != want || gotLevel != wantLevel {
		t.Fatal("shutdown: shutting down should be logged as info")
	}
}

func TestShutdownHandle_LogErrors(t *testing.T) {
	logger, hook := logrus_test.NewNullLogger()
	shutdown := serverz.NewShutdown(logger)

	shutdown.Register(func() error {
		return errors.New("error")
	})

	shutdown.Handle()

	if got, want, gotLevel, wantLevel := hook.LastEntry().Message, "error", hook.LastEntry().Level, logrus.ErrorLevel; got != want || gotLevel != wantLevel {
		t.Fatal("shutdown: errors ocurred during shutdown should be logged")
	}
}

func TestShutdownHandle_RecoverFromPanic(t *testing.T) {
	logger, hook := logrus_test.NewNullLogger()
	shutdown := serverz.NewShutdown(logger)

	func() {
		defer shutdown.Handle()

		func() {
			panic(errors.New("error"))
		}()
	}()

	if got, want, gotLevel, wantLevel := hook.Entries[0].Message, "error", hook.Entries[0].Level, logrus.ErrorLevel; got != want || gotLevel != wantLevel {
		t.Fatal("shutdown: errors ocurred during shutdown should be logged")
	}
}
