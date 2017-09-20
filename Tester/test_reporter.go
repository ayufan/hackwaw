package main

import (
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/types"

	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type testReporter struct {
}

func (t *testReporter) Fail() {
	os.Exit(1)
}

func (t *testReporter) SpecSuiteWillBegin(config config.GinkgoConfigType, summary *types.SuiteSummary) {
	if *dashboardURL == "" {
		return
	}

	values := url.Values{}
	values.Add("type", "suite")
	values.Add("repository", *repository)
	values.Add("id", strconv.FormatInt(*suiteID, 10))
	values.Add("total", strconv.Itoa(summary.NumberOfSpecsThatWillBeRun))
	values.Add("state", "starting")

	resp, err := http.PostForm(*dashboardURL, values)
	if err != nil {
		log.Println("[SpecSuiteWillBegin]", values, err)
		return
	}
	defer resp.Body.Close()
}

func (t *testReporter) BeforeSuiteDidRun(setupSummary *types.SetupSummary) {

}

func (t *testReporter) SpecWillRun(specSummary *types.SpecSummary) {
	if *dashboardURL == "" {
		return
	}

	values := url.Values{}
	values.Add("type", "spec")
	values.Add("repository", *repository)
	values.Add("suiteId", strconv.FormatInt(*suiteID, 10))
	values.Add("name", formatComponentTexts(specSummary))
	values.Add("state", "starting")

	resp, err := http.PostForm(*dashboardURL, values)
	if err != nil {
		log.Println("[SpecWillRun]", values, err)
		return
	}
	defer resp.Body.Close()
}

func formatComponentTexts(specSummary *types.SpecSummary) string {
	texts := specSummary.ComponentTexts[1:]
	components := []string{
		strings.Join(texts[:len(texts)-1], " "),
		texts[len(texts)-1],
	}
	return strings.Join(components, " IT ")
}

func (t *testReporter) SpecDidComplete(specSummary *types.SpecSummary) {
	if *dashboardURL == "" {
		return
	}

	values := url.Values{}
	values.Add("type", "spec")
	values.Add("repository", *repository)
	values.Add("suiteId", strconv.FormatInt(*suiteID, 10))
	values.Add("name", formatComponentTexts(specSummary))
	values.Add("duration", fmt.Sprintf("%f", specSummary.RunTime.Seconds()))
	values.Add("message", specSummary.Failure.Message)
	switch specSummary.State {
	case types.SpecStateInvalid:
		values.Add("state", "invalid")
	case types.SpecStatePending:
		values.Add("state", "pending")
	case types.SpecStateSkipped:
		values.Add("state", "skipped")
	case types.SpecStatePassed:
		values.Add("state", "passed")
	case types.SpecStateFailed:
		values.Add("state", "failed")
	case types.SpecStatePanicked:
		values.Add("state", "panicked")
	case types.SpecStateTimedOut:
		values.Add("state", "timedout")
	}

	resp, err := http.PostForm(*dashboardURL, values)
	if err != nil {
		log.Println("[SpecDidComplete]", values, err)
		return
	}
	defer resp.Body.Close()
}

func (t *testReporter) AfterSuiteDidRun(setupSummary *types.SetupSummary) {
}

func (t *testReporter) SpecSuiteDidEnd(summary *types.SuiteSummary) {
	if *dashboardURL == "" {
		return
	}

	values := url.Values{}
	values.Add("type", "suite")
	values.Add("repository", *repository)
	values.Add("id", strconv.FormatInt(*suiteID, 10))
	values.Add("pending", strconv.Itoa(summary.NumberOfPendingSpecs))
	values.Add("skipped", strconv.Itoa(summary.NumberOfSkippedSpecs))
	values.Add("passed", strconv.Itoa(summary.NumberOfPassedSpecs))
	values.Add("failed", strconv.Itoa(summary.NumberOfFailedSpecs))
	values.Add("total", strconv.Itoa(summary.NumberOfTotalSpecs))
	values.Add("duration", fmt.Sprintf("%f", summary.RunTime.Seconds()))
	if summary.SuiteSucceeded {
		values.Add("state", "passed")
	} else {
		values.Add("state", "failed")
	}

	resp, err := http.PostForm(*dashboardURL, values)
	if err != nil {
		log.Println("[SpecSuiteWillBegin]", values, err)
		return
	}
	defer resp.Body.Close()
}
