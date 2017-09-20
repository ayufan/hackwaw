package main

import "flag"

var serverListen = flag.String("server-listen", "127.0.0.1:0", "if non-empty, httptest.NewServer serves on this address and blocks")
var apiServer = flag.String("remote-server", "https://hackwaw-proxy-app.herokuapp.com", "defines a server to which the requests are forwarded")
var appImage = flag.String("image", "", "application image to be used for tests")
var dashboardURL = flag.String("dashboard-url", "", "notification address for suite begin and end")
var suiteID = flag.Int64("suite-id", 0, "unique number of suite run")
var repository = flag.String("suite-repository", "group/project", "the path to repository")
var testerLogs = flag.Bool("tester-logs", false, "If to show tester debug logs")
var containerLogs = flag.Bool("container-logs", false, "If to show the container logs during runs")
var startupTimeout = flag.Float64("startup-timeout", 60.0, "Application startup time")
