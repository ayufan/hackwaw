package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"errors"
	"fmt"
	"github.com/onsi/gomega/types"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var tweet Tweet = Tweet{
	Id:   100,
	Body: "long tweet",
	Link: "my link",
	Date: time.Now(),
}

var tweet2 Tweet = Tweet{
	Id:   101,
	Body: "long tweet",
	Link: "my link",
	Date: time.Now(),
}

var _ = Describe("Check startup time", func() {
	sc := NewSharedContext()
	sc.ConfigureHandlers(true)
	sc.ConfigureServer()
	sc.ConfigureApp(true)
	sc.ConfigureBoot()

	var startupTime float64

	JustBeforeGroup(func() {
		By("We are waiting for application to become start...")
		started := time.Now()
		Eventually(sc.app.AppStatus, 30.0).Should(Equal(StatusOperational))
		startupTime = time.Since(started).Seconds()
	})

	It("should start in 3 seconds", func() {
		Expect(startupTime).To(BeNumerically("<=", 3.0))
	})

	It("should start in 10 seconds", func() {
		Expect(startupTime).To(BeNumerically("<=", 10.0))
	})

	It("should start in 30 seconds", func() {
		Expect(startupTime).To(BeNumerically("<=", 30.0))
	})
})

var _ = Describe("Verify normal behavior", func() {
	sc := NewSharedContext()
	sc.ConfigureHandlers()
	sc.ConfigureServer()
	sc.ConfigureApp(true)
	sc.ConfigureBoot()

	It("should request tweets", func() {
		Eventually(func() bool {
			return sc.tweetsRequested
		}).Should(BeTrue())
	})

	It("should send slack notification for all received tweets", func() {
		Eventually(func() bool {
			return sc.slackPushed
		}).Should(BeTrue())
	})
})

var _ = Describe("Twitter", func() {
	sc := NewSharedContext()
	sc.ConfigureHandlers()
	sc.ConfigureServer()
	sc.ConfigureApp(true)
	sc.ConfigureBoot()

	JustBeforeEach(func() {
		By("Waiting for Twitter to be operational...")
		Eventually(sc.app.TwitterStatus, 20.0).Should(Equal(StatusOperational))
	})

	Context("is serving wrong response", func() {
		JustBeforeEach(func() {
			By("Twitter is serving invalid JSON...")
			sc.tweetsHandler = invalidJson
		})

		It("should report an ERROR when trying to access tweets", func() {
			Eventually(sc.app.TwitterStatus, 10.0).Should(Equal(StatusError))
		})
	})

	Context("is serving json without Content-Type", func() {
		JustBeforeEach(func() {
			By("Twitter is serving json without Content-Type...")
			sc.tweetsHandler = jsonNoContentType
		})

		It("should report an ERROR when trying to access tweets", func() {
			Eventually(sc.app.TwitterStatus, 10.0).Should(Equal(StatusError))
		})
	})

	Context("server is returning 500", func() {
		JustBeforeEach(func() {
			By("Twitter is returning 500...")
			sc.tweetsHandler = serverError
		})

		It("should report an ERROR when trying to access tweets", func() {
			Eventually(sc.app.TwitterStatus, 10.0).Should(Equal(StatusError))
		})
	})

	Context("connection to server is taking 10 second", func() {
		JustBeforeEach(func() {
			By("Twitter is slow...")
			sc.tweetsHandler = slowDown(time.Second*10, sc.tweetsHandler)
		})

		It("should report an SLOW when trying to access tweets", func() {
			Eventually(sc.app.TwitterStatus, 20.0).Should(Equal(StatusSlow))
		})
	})

	Context("connection is hanging", func() {
		JustBeforeEach(func() {
			By("Twitter is hanging on connection...")
			sc.tweetsHandler = slowDown(time.Minute, serverError)
		})

		It("should report an DOWN when trying to access tweets", func() {
			Eventually(sc.app.TwitterStatus, 30.0).Should(Equal(StatusDown))
		})
	})

	Context("is back to normal", func() {
		BeforeEach(func() {
			By("Twitter is back to normal...")
			sc.tweetsHandler = emptyJson
		})

		It("should report that twitter is OPERATIONAL", func() {
			Eventually(sc.app.TwitterStatus, 10.0).Should(Equal(StatusOperational))
		})

		It("should send all remaining tweets", func() {
			Eventually(func() sort.IntSlice {
				return sc.pushes
			}, 20.0).Should(Equal(sc.tweets))
		})
	})
})

var _ = Describe("Slack", func() {
	var failedPushes sort.IntSlice

	sc := NewSharedContext()
	sc.ConfigureHandlers()
	sc.ConfigureServer()
	sc.ConfigureApp(true)
	sc.ConfigureBoot()

	JustBeforeEach(func() {
		By("Waiting for Slack to be operational...")
		Eventually(sc.app.SlackStatus, 20.0).Should(Equal(StatusOperational))
	})

	AfterEach(func() {
		failedPushes = append(failedPushes, sc.pushesErrored...)
		failedPushes.Sort()
		sc.pushesErrored = sort.IntSlice{}
	})

	Context("server is returning 500", func() {
		JustBeforeEach(func() {
			By("Slack is returning 500...")
			sc.slackError = errors.New("server error")
		})

		It("should report an ERROR when trying to push notification", func() {
			Eventually(sc.app.SlackStatus, 10.0).Should(Equal(StatusError))
		})
	})

	Context("connection to server is taking 10 second", func() {
		JustBeforeEach(func() {
			By("Slack is taking long time to push notification...")
			sc.pushHandler = slowDown(time.Second*10, sc.pushHandler)
		})

		It("should report an SLOW when trying to push notification", func() {
			Eventually(sc.app.SlackStatus, 20.0).Should(Equal(StatusSlow))
		})
	})

	Context("connection is hanging", func() {
		JustBeforeEach(func() {
			By("Slack is hanging on connection...")
			sc.slackError = errors.New("timeout")
			sc.pushHandler = slowDown(time.Minute, serverError)
		})

		It("should report an DOWN when trying to push notification", func() {
			Eventually(sc.app.SlackStatus, 30.0).Should(Equal(StatusDown))
		})
	})

	Context("is back to normal", func() {
		It("should report that Slack is OPERATIONAL", func() {
			Eventually(sc.app.SlackStatus, 10.0).Should(Equal(StatusOperational))
		})

		It("should send all notifications from tweets", func() {
			Eventually(func() sort.IntSlice {
				return sc.pushes
			}, 20.0).Should(Equal(sc.tweets))
		})

		It("should retry all failed notifications", func() {
			matchers := make([]types.GomegaMatcher, len(failedPushes))
			for idx, failedPush := range failedPushes {
				matchers[idx] = ContainElement(failedPush)
			}

			Eventually(func() sort.IntSlice {
				sc.lock.RLock()
				defer sc.lock.RUnlock()
				return sc.pushes
			}, 20.0).Should(SatisfyAll(matchers...))
		})
	})
})

var _ = Describe("Verify performance", func() {
	sc := NewSharedContext()
	sc.ConfigureHandlers(true)
	sc.ConfigureServer()
	sc.ConfigureApp(true)
	sc.ConfigureBoot()

	It("should request tweets", func() {
		Eventually(func() bool {
			return sc.tweetsRequested
		}).Should(BeTrue())
	})

	It("should return non-empty list of tweets", func() {
		Eventually(func() ([]Tweet, error) {
			return sc.app.Latest(0)
		}).ShouldNot(BeEmpty())
	})

	Measure("it should process 100 requests", func(b Benchmarker) {
		time := b.Time("/latest request time", func() {
			messages, err := sc.app.Latest(0)
			Expect(err).NotTo(HaveOccurred())
			Expect(messages).NotTo(BeEmpty())
		})

		Expect(time.Seconds()).Should(BeNumerically("<", 1.0), "/latest shouldn't take longer then 1 second")
	}, 100)
})

var _ = Describe("For large response", func() {
	sc := NewSharedContext()
	sc.ConfigureHandlers(true)
	sc.ConfigureServer()
	sc.ConfigureApp(true)
	sc.ConfigureBoot()

	BeforeEach(func() {
		sc.tweetsHandler = func(rw http.ResponseWriter, req *http.Request) {
			data, err := json.Marshal(tweet)
			if err != nil {
				serverError(rw, req)
				return
			}
			data2, err := json.Marshal(tweet2)
			if err != nil {
				serverError(rw, req)
				return
			}

			rw.Header().Add("Content-Type", "application/json")
			fmt.Fprintln(rw, "[")
			fmt.Fprintln(rw, string(data), ",")
			fmt.Fprintln(rw, string(data), ",")
			fmt.Fprintln(rw, string(data), ",")
			for i := 0; i < 10000; i++ {
				fmt.Fprintln(rw, string(data2), ",")
			}
			fmt.Fprintln(rw, string(data2))
			fmt.Fprintln(rw, "]")
		}
		sc.pushForwardHandler = emptyText
	})

	Context("containing the same Tweet ID", func() {
		It("respond with last tweets", func() {
			Eventually(sc.app.LatestIds, 10).Should(ContainElement(tweet.UniqueID()))
		})

		It("should send only one notification", func() {
			pushes := func() []int {
				return sc.pushes
			}
			matchers := SatisfyAll(
				HaveLen(2),
				ContainElement(tweet.UniqueID()),
				ContainElement(tweet2.UniqueID()),
			)
			Eventually(pushes, 10).Should(matchers)
			Consistently(pushes, 3).Should(matchers)
		})
	})
})

var _ = Describe("For data persistence", func() {
	sc := NewSharedContext()
	sc.ConfigureHandlers(true)
	sc.ConfigureServer()
	sc.ConfigureApp(true)
	sc.ConfigureBoot()
	sc.ConfigurePersistentStorage()

	JustBeforeEach(func() {
		By("Wait for tweets to be received...")
		Eventually(func() []int {
			return sc.tweets
		}).ShouldNot(BeEmpty())
	})

	JustBeforeEach(func(done Done) {
		By("Wait for tweets to be received...")
		Eventually(sc.app.LatestIds).ShouldNot(BeEmpty())
		close(done)
	}, 20.0)

	JustBeforeEach(func() {
		By("Stopping Tweets feed...")
		sc.tweetsHandler = emptyJson
	})

	JustBeforeEach(func(done Done) {
		By("Recreating application...")
		Expect(sc.app.Recreate()).NotTo(HaveOccurred())
		close(done)
	}, 20.0)

	JustBeforeEach(func() {
		By("Waiting for application to start responding...")
		Eventually(sc.app.AppStatus, *startupTimeout).Should(Equal(StatusOperational))
	})

	It("should respond with previously fetched tweets", func() {
		By("Waiting for returning latest ids...")
		Eventually(sc.app.LatestIds, 5).Should(SatisfyAny(sc.TweetMatchers()...))
	})
})

var _ = Describe("Data sharing", func() {
	sc := NewSharedContext()
	sc.ConfigureHandlers(true)
	sc.ConfigureServer()
	sc.ConfigureApp(true)
	sc.ConfigureBoot()
	sc.ConfigurePersistentStorage()

	var (
		secondApp *App
	)

	JustBeforeEach(func() {
		var err error
		By("Forking application container...")
		secondApp, err = sc.app.Fork("-e", "TWITTER_URL=http://localhost/")
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		By("Waiting for second application to become alive...")
		Eventually(secondApp.AppStatus, *startupTimeout).Should(Equal(StatusOperational))
	})

	JustBeforeEach(func() {
		By("Waiting for tweets to be received...")
		Eventually(func() []int {
			return sc.tweets
		}).ShouldNot(BeEmpty())
	})

	It("second app should respond with tweets from the first one", func(done Done) {
		By("Waiting for second application to return a list of tweets...")
		Eventually(secondApp.LatestIds, 10).Should(SatisfyAny(sc.TweetMatchers()...))
		close(done)
	}, 90.0)
})

var _ = Describe("Verify memory limits", func() {
	sc := NewSharedContext()
	sc.ConfigureHandlers(true)
	sc.ConfigureServer()
	sc.ConfigureApp(true)
	sc.ConfigureBoot()

	Context("if starts on 50M", func() {
		BeforeGroup(func() {
			sc.memoryArgs = []string{"--memory", "50M"}
		})

		JustBeforeGroup(func() {
			By("Waiting for tweets to be received...")
			Eventually(func() []int {
				return sc.tweets
			}).ShouldNot(BeEmpty())
		})

		It("should return latest tweets", func() {
			Eventually(sc.app.LatestIds, 10).Should(SatisfyAny(sc.TweetMatchers()...))
		})
	})
})

var _ = Describe("Disk space", func() {
	sc := NewSharedContext()
	sc.ConfigureHandlers()
	sc.ConfigureServer()
	sc.ConfigureApp(true)
	sc.ConfigureBoot()
	sc.ConfigurePersistentStorage()

	temporaryFile := func() string {
		return filepath.Join(sc.mountDir, ".temporary")
	}

	fillDisk := func() {
		if sc.mountDir == "" {
			return
		}

		By("Creating temporary file that will fill the disk...")
		f, err := os.Create(temporaryFile())
		Expect(err).NotTo(HaveOccurred())
		defer f.Close()

		written := 0
		buf := make([]byte, 4096)

		for {
			var n int
			n, err = f.Write(buf)
			if err != nil {
				break
			}
			written += n
		}
		By(fmt.Sprintln("Wrote", written, "and failed with", err.Error()))
	}

	freeDisk := func() {
		if sc.mountDir == "" {
			return
		}
		By("Removing temporary file...")
		os.Remove(temporaryFile())
	}

	JustBeforeEach(func() {
		By("Waiting for tweets to be received...")
		Eventually(func() []int {
			return sc.tweets
		}, 10.0).ShouldNot(BeEmpty())
	})

	JustBeforeEach(func() {
		fillDisk()
	})

	Context("should respond with some tweets", func() {
		It("should respond with the newest received tweets", func() {
			By("Checking a list of tweets...")
			Eventually(sc.app.LatestIds, 10.0).Should(SatisfyAny(sc.TweetMatchers()...))
		})

		Context("when the disk space is freed", func() {
			JustBeforeEach(func() {
				sc.returnRecordedTweets = true
				freeDisk()
			})

			JustBeforeEach(func() {
				By("We give application 10 seconds to recover from low disk space scenario...")
				time.Sleep(time.Second * 10)
			})

			Context("and when application is restarted", func() {
				JustBeforeEach(func() {
					By("Stopping Tweets feed...")
					sc.tweetsHandler = emptyJson
				})

				JustBeforeEach(func() {
					By("Recreating application...")
					err := sc.app.Recreate()
					Expect(err).ToNot(HaveOccurred())
				})

				JustBeforeEach(func() {
					By("Waiting for application to start responding...")
					Eventually(sc.app.AppStatus, *startupTimeout).Should(Equal(StatusOperational))
				})

				It("should respond with the latealy received tweets", func() {
					By("Checking a list of tweets...")
					Eventually(sc.app.LatestIds, 10.0).Should(SatisfyAny(sc.TweetMatchers()...))
				})
			})
		})
	})
})
