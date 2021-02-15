package terminal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func HandleSignals(cancel context.CancelFunc, cb func() error) {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		syscall.SIGUSR1,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt,
	)

	closeDone := make(chan struct{}, 1)
	go func() {
		for {
			sig := <-signals
			switch sig {
			case syscall.SIGUSR1:
				log.Println("caught USR1 signal")

			case syscall.SIGHUP:
				log.Println("caught HUP signal")

			// Graceful shut-down on SIGINT/SIGTERM
			case os.Interrupt, syscall.SIGTERM:

				log.Infof("Got signal [%v] to exit.\n", sig)
				cancel()

				select {
				case <-signals:
					// send signal again, return directly
					log.Infof("\nGot signal [%v] again to exit.\n", sig)
					os.Exit(1)

				case <-time.After(10 * time.Second):
					log.Infof("\nWait 10s for closed, force exit\n")
					os.Exit(1)

				case <-closeDone:
					log.Info("Gracefully exited")
					os.Exit(0)
					return
				}

			default:
				log.Printf("caught unhandled signal %+v\n", sig)
			}
		}
	}()

	if err := cb(); err != nil {
		log.Warn("Application failed: ", err)
	}

	// close handlers
	cancel()
	closeDone <- struct{}{}
}
