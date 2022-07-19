// Package gracequit handle delete pid file at dead,
// set logging path and folder, setup http/pprof service and etc. common tasks for background service.
//
// Example:
// 	ctx := gracequit.New(
//		myPathPid,
//		2,
//	).Init()
//
// ...
// for {
//		select {
//		case <-ctx.Done():
//			{
//				return
//			}
//		case <-time.After(d):
//			{
//				// err := myOtherJob(...)
//				gutil.CheckErr(err)
//			}
//		}
//	}
package gracequit

import (
	"context"
	"fmt"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/giant-stone/go/glogging"
)

const (
	// DefaultQuitAfterSecs quit main process after n seconds
	DefaultQuitAfterSecs = 2
)

// GraceQuit internal struct holds settings
type GraceQuit struct {
	quitAfterSecs int
	ctx           context.Context

	pathPid string
}

func (it *GraceQuit) Init() (ctx context.Context) {
	if it.quitAfterSecs <= 0 {
		it.quitAfterSecs = DefaultQuitAfterSecs
	}
	duration := time.Second * time.Duration(it.quitAfterSecs)

	ctx, cancel := context.WithCancel(context.Background())
	it.ctx = ctx

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range sigs {
			glogging.Sugared.Debugf("signal %v captured, cancel context and exit after %v", sig, duration)

			it.deletePid()

			cancel()
			<-time.After(duration)
			os.Exit(0)
		}
	}()

	it.createPid()

	return ctx
}

// New initialize.
func New(pathPid string, quitAfterSecs int) *GraceQuit {
	return &GraceQuit{
		pathPid:       pathPid,
		quitAfterSecs: quitAfterSecs,
	}
}

func (it *GraceQuit) deletePid() {
	if it.pathPid == "" {
		return
	}
	err := os.Remove(it.pathPid)
	if err != nil {
		glogging.Sugared.Error("os.Remove", err, it.pathPid)
	}
}

func (it *GraceQuit) createPid() {
	if it.pathPid == "" {
		return
	}

	err := os.MkdirAll(path.Dir(it.pathPid), 0644)
	if err != nil {
		glogging.Sugared.Error("os.MkdirAll", err, it.pathPid)
	}

	err = ioutil.WriteFile(it.pathPid, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	if err != nil {
		glogging.Sugared.Error("ioutil.WriteFile", err, it.pathPid)
	}
}

// GetContext use context.Context share states.
func (it *GraceQuit) GetContext() context.Context {
	return it.ctx
}
