// Package gracequit handle delete pid file at dead,
// set logging path and folder, setup http/pprof service and etc. common tasks for background service.
//
// Example:
// 	ctx := gracequit.New(
//		myPathPid,
//		myPathLog,
//		myProfHTTP,
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
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/giant-stone/go/gutil"
)

const (
	// DefaultQuitAfterSecs quit main process after n seconds
	DefaultQuitAfterSecs = 2
)

// GraceQuit internal struct holds settings
type GraceQuit struct {
	debug bool

	quitAfterSecs int
	ctx           context.Context

	pathPid string
	pathLog string

	httpProf string
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
			if it.debug {
				log.Printf("[info] signal %v captured, cancel context and exit after %v", sig, duration)
			}

			it.deletePid()

			cancel()
			<-time.After(duration)
			os.Exit(0)
		}
	}()

	it.setupLog()
	it.setupHttpProf()
	it.createPid()

	return ctx
}

// New initialize.
func New(pathPid string, httpProf string, pathLog string, quitAfterSecs int) *GraceQuit {
	return &GraceQuit{
		pathPid:       pathPid,
		quitAfterSecs: quitAfterSecs,
		httpProf:      httpProf,
		pathLog:       pathLog,
	}
}

// SetDebug set it true to enable logging.
func (it *GraceQuit) SetDebug(val bool) (rs *GraceQuit) {
	it.debug = val
	return it
}

func (it *GraceQuit) deletePid() {
	if it.pathPid == "" {
		return
	}
	gutil.CheckErr(os.Remove(it.pathPid))
}

func (it *GraceQuit) createPid() {
	if it.pathPid == "" {
		return
	}

	gutil.CheckErr(os.MkdirAll(path.Dir(it.pathPid), 0644))
	err := ioutil.WriteFile(it.pathPid, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	gutil.CheckErr(err)
}

// GetContext use context.Context share states.
func (it *GraceQuit) GetContext() context.Context {
	return it.ctx
}

func (it *GraceQuit) setupLog() {
	if it.pathLog == "" {
		return
	}

	log.SetOutput(&lumberjack.Logger{
		Filename:   it.pathLog,
		MaxSize:    100, // megabytes
		MaxBackups: 10,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	})
}

func (it *GraceQuit) setupHttpProf() {
	if it.httpProf == "" {
		return
	}
	go func() {
		if it.debug {
			log.Printf("[debug] net/http/pprof listen on %s", it.httpProf)
		}
		err := http.ListenAndServe(it.httpProf, nil)
		gutil.ExitOnErr(err)
	}()
}
