package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/coreos/pkg/flagutil"

	"flag"
	"io/ioutil"
	"log"
	"time"
	"os"
)

func main() {
	err := flagutil.SetFlagsFromEnv(flag.CommandLine, "APP")
	if err != nil {
		log.Fatalln(err)
	}
	flag.Parse()

	log.SetPrefix(ANSI_BOLD_YELLOW + "[tester] " + ANSI_RESET)
	log.SetFlags(log.Ltime)
	if !*testerLogs {
		log.SetOutput(ioutil.Discard)
	}

	SetDefaultEventuallyTimeout(time.Second * 10)
	SetDefaultEventuallyPollingInterval(time.Second)

	SetDefaultConsistentlyDuration(time.Second * 10)
	SetDefaultConsistentlyPollingInterval(time.Second)

	t := testReporter{}
	RegisterFailHandler(Fail)
	passed := RunSpecsWithDefaultAndCustomReporters(&t, "Suite", []Reporter{&t})
	if !passed {
		os.Exit(1)
	}
}
