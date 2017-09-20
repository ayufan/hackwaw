package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"github.com/onsi/gomega/types"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"sort"
	"sync"
)

type sharedContext struct {
	ts                   *httptest.Server
	app                  *App
	appArgs              []string
	storageArgs          []string
	memoryArgs           []string
	tweetsHandler        http.HandlerFunc
	pushHandler          http.HandlerFunc
	tweetsForwardHandler http.HandlerFunc
	pushForwardHandler   http.HandlerFunc
	mountDir             string

	returnRecordedTweets bool
	recordedTweets       []Tweet

	lock            sync.RWMutex
	tweets          sort.IntSlice
	pushes          sort.IntSlice
	pushesErrored   sort.IntSlice
	tweetsRequested bool
	slackPushed     bool
	slackError      error
}

var previousServer *httptest.Server
var previousApp *App

var _ = AfterSuite(func() {
	if previousServer != nil {
		previousServer.Close()
	}
	if previousApp != nil {
		previousApp.Remove()
	}
})

func (s *sharedContext) ConfigureHandlers(once ...bool) {
	restoreHandlers := func() {
		s.tweetsHandler = validateTweetRequest(s.forwardTweets, s.handleTweet)
		s.pushHandler = validatePushRequest(s.forwardPushes, s.handlePush)
		s.slackError = nil
	}

	if len(once) > 0 && once[0] {
		BeforeGroup(restoreHandlers)
	} else {
		BeforeEach(restoreHandlers)
	}
}

func (s *sharedContext) TweetMatchers() (matchers []types.GomegaMatcher) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	matchers = make([]types.GomegaMatcher, len(s.tweets))
	for idx, tweetId := range s.tweets {
		matchers[idx] = ContainElement(tweetId)
	}
	return
}

func (s *sharedContext) ConfigureServer() {
	BeforeGroup(s.startServer)
	AfterGroup(s.finishServer)
}

func (s *sharedContext) ConfigureApp(once ...bool) {
	if len(once) > 0 && once[0] {
		JustBeforeGroup(s.startApp)
		JustBeforeGroup(func() {
			Eventually(s.app.Status).Should(Equal("running"))
		})
		AfterGroup(s.finishApp)
	} else {
		JustBeforeEach(s.startApp, 3.0)
		JustBeforeEach(func() {
			Eventually(s.app.Status).Should(Equal("running"))
		})
		AfterEach(s.finishApp)
	}
}

func (s *sharedContext) ConfigureBoot() {
	JustBeforeGroup(s.waitForBoot)
}

func (s *sharedContext) ConfigurePersistentStorage() {
	BeforeGroup(func() {
		var err error
		s.mountDir, err = ioutil.TempDir("", "mount")
		Expect(err).NotTo(HaveOccurred())
		Expect(s.mountDir).NotTo(BeEmpty())

		// Mount file to directory
		cmd := exec.Command("mount", "-t", "tmpfs", "-o", "size=1M", "tmpfs", s.mountDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		Expect(err).NotTo(HaveOccurred())

		// Update storage directory
		s.storageArgs = []string{"-v", s.mountDir + ":/storage"}
	})

	AfterGroup(func() {
		if s.mountDir == "" {
			return
		}
		cmd := exec.Command("umount", "-f", s.mountDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	})
}

func (s *sharedContext) forwardTweets(rw http.ResponseWriter, req *http.Request) {
	s.tweetsRequested = true
	if s.tweetsForwardHandler != nil {
		s.tweetsForwardHandler(rw, req)
	} else if s.returnRecordedTweets && len(s.recordedTweets) > 0 {
		jsonResponse(s.recordedTweets)
	} else {
		forwardTweets(rw, req)
	}
}

func (s *sharedContext) handleTweet(tweet Tweet) {
	s.lock.Lock()
	defer s.lock.Unlock()
	idx := s.tweets.Search(tweet.UniqueID())
	if idx < len(s.tweets) && s.tweets[idx] == tweet.UniqueID() {
		return
	}
	s.tweets = append(s.tweets, tweet.UniqueID())
	s.tweets.Sort()
	s.recordedTweets = append(s.recordedTweets, tweet)
}

func (s *sharedContext) forwardPushes(rw http.ResponseWriter, req *http.Request) {
	s.slackPushed = true
	if s.pushForwardHandler != nil {
		s.pushForwardHandler(rw, req)
	} else {
		forwardPushes(rw, req)
	}
}

func (s *sharedContext) handlePush(push Slack) error {
	s.lock.RLock()
	idx := sort.SearchInts(s.tweets, int(push.TweetId))
	s.lock.RUnlock()
	if idx < 0 {
		return errors.New("tweet not found")
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	if s.slackError != nil {
		s.pushesErrored = append(s.pushesErrored, int(push.TweetId))
		s.pushesErrored.Sort()
	} else {
		s.pushes = append(s.pushes, int(push.TweetId))
		s.pushes.Sort()
	}
	return s.slackError
}

func (s *sharedContext) startApp() {
	if previousApp != nil {
		previousApp.Remove()
		previousApp = nil
	}

	By("Starting application...")
	var err error
	var args []string
	args = append(args, s.appArgs...)
	args = append(args, s.storageArgs...)
	args = append(args, s.memoryArgs...)
	s.app, err = NewAppForServer(s.ts, args...)
	Expect(err).NotTo(HaveOccurred())
	previousApp = s.app
}

func (s *sharedContext) finishApp() {
	if s.app != nil {
		s.app.Remove()
	}
}

func (s *sharedContext) startServer() {
	By("Starting internal API server...")

	// TODO: This is needed, because gomega doesn't implement group contexes for shared-it's
	if previousServer != nil {
		previousServer.Close()
		previousServer = nil
	}

	s.ts = startServer(func(rw http.ResponseWriter, req *http.Request) {
		if s.tweetsHandler != nil {
			s.tweetsHandler(rw, req)
		} else {
			serverError(rw, req)
		}
	}, func(rw http.ResponseWriter, req *http.Request) {
		if s.pushHandler != nil {
			s.pushHandler(rw, req)
		} else {
			serverError(rw, req)
		}
	})
}

func (s *sharedContext) finishServer() {
	if s.ts != nil {
		s.ts.Close()
	}
}

func (s *sharedContext) waitForBoot() {
	By("Waiting for application to be up...")
	Eventually(s.app.AppStatus, *startupTimeout).Should(Equal(StatusOperational))
}

func NewSharedContext() *sharedContext {
	sc := &sharedContext{}
	sc.tweetsHandler = validateTweetRequest(sc.forwardTweets, sc.handleTweet)
	sc.pushHandler = validatePushRequest(sc.forwardPushes, sc.handlePush)
	sc.appArgs = []string{"--read-only", "--tmpfs=/tmp:size=10M"}
	sc.storageArgs = []string{"--tmpfs=/storage:size=10M"}
	sc.memoryArgs = []string{"--memory", "2G"}
	sc.returnRecordedTweets = false
	return sc
}
