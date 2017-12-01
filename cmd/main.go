package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type deployment struct {
	run  chan struct{}
	done chan struct{}
	err  chan error
	fd   string
	exe  string
}

const version = "1.1"

// Delay in seconds watching file or directory
const watchdelay = 5

// Number of allowed pending deployments
const pending = 1

// WATCHERR is a critical error from fsnotify and leads to program termination
var WATCHERR error

func banner() {
	fmt.Println(`
   _____  ____     _____ _____       _____ _____  
  / ____|/ __ \   / ____|_   _|     / ____|  __ \ 
 | |  __| |  | | | |      | |______| |    | |  | |
 | | |_ | |  | | | |      | |______| |    | |  | |
 | |__| | |__| | | |____ _| |_     | |____| |__| |
  \_____|\____/   \_____|_____|     \_____|_____/ 
                                                  
        The World's most basic CI/CD Tool
`)
}

func usage() {
	banner()
	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Usage of go-cicd:\n\n")
	flag.PrintDefaults()
	fmt.Println()
	os.Exit(0)
}

func watch(ctx context.Context, d *deployment) {
	tCh := time.NewTicker(time.Second * watchdelay).C
	last, err := os.Stat(d.fd)
	if err != nil {
		WATCHERR = fmt.Errorf("%v", err)
		d.err <- WATCHERR
		return
	}

	for {
		select {
		case <-tCh:
			mod, err := os.Stat(d.fd)
			if err != nil {
				WATCHERR = fmt.Errorf("%v", err)
				d.err <- WATCHERR
				return
			}

			if last.ModTime() != mod.ModTime() {
				log.Printf("info: %q changed\n", d.fd)
				last = mod
				select {
				case d.run <- struct{}{}:
				default:
					log.Println("info: deployment queue full, skipping...")
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func deploy(ctx context.Context, d *deployment) {
	for {
		select {
		case <-ctx.Done():
			close(d.done)
			return
		case <-d.run:
			log.Println("Starting deployment")
			cmd := exec.CommandContext(ctx, d.exe)
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					//fmt.Println("This is the case where it ran, but returned an exit code that wasn't 0.")
					d.err <- errors.New(exitError.Error())
				} else {
					//fmt.Println("This is the case where the command failed to run at all.  Message is: " + err.Error())
					d.err <- err
				}
				continue
			}
			//fmt.Println("This is the case where the Run() command's error was nil.  Should be good and have a return code of 0.")
			log.Println("Deployed successfully")
		}

	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetPrefix("[go-cicd] ")
	log.SetFlags(log.Ltime)

	d := &deployment{
		// allow n pending deployments
		run:  make(chan struct{}, pending),
		done: make(chan struct{}),
		err:  make(chan error, 1),
	}

	// TODO: Implement flag to select events, e.g. "c,w,d"
	// TODO: Implement option flags for cmd, if any
	flag.StringVar(&d.fd, "f", "main.go", "File or directory to watch")
	flag.StringVar(&d.exe, "e", "./deploy.sh", "Pipeline executable")
	flag.Usage = usage
	flag.Parse()

	// Print banner
	banner()

	// Check if d.exe exists
	_, err := exec.LookPath(d.exe)
	if err != nil {
		log.Fatalf("error: %s not found. Is it in $PATH? See -h for help", d.exe)
	}

	// Check if d.fd exists
	f, err := os.Stat(d.fd)
	if err != nil {
		log.Fatalf("error: %s not found, see -h for help", d.fd)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("Starting CI/CD process (executable: %q, watching: %q (isDir: %v))\n", d.exe, d.fd, f.IsDir())
	go watch(ctx, d)
	go deploy(ctx, d)

	// Catch signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Signal(syscall.SIGTERM))

	for {
		select {
		case s := <-signals:
			log.Printf("Interrupted (%s). Cleaning up and cancelling running deployments\n", s)
			cancel()
			for {
				select {
				// Read any errors during shutdown
				case err = <-d.err:
					log.Printf("error: %v\n", err)

				// Wait for pending deployment
				case <-d.done:
					log.Println("Shutting down")
					os.Exit(0)
				}
			}
		case err = <-d.err:
			if err == WATCHERR {
				log.Fatalf("error: %v", err)
			} else {
				log.Printf("error: %v\n", err)
			}
		}
	}

}
