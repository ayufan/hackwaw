package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/goware/prefixer"
	. "github.com/onsi/ginkgo"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type App struct {
	Id   string
	Args []string
}

func dockerCommand(command string, args ...string) *exec.Cmd {
	cmd := exec.Command("docker", append([]string{command}, args...)...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	return cmd
}

func (a *App) request(method, path string, query url.Values, result interface{}) (err error) {
	url, err := a.URL(path, query)
	if err != nil {
		return
	}

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("Invalid status code: %d, %v", res.StatusCode, res.Status)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	log.Println("APP:", "REQUEST", method, url, string(data))
	err = json.Unmarshal(data, result)
	return
}

func (a *App) Fork(args ...string) (newApp *App, err error) {
	newApp = &App{
		Args: a.Args,
	}
	newApp.Args = append(newApp.Args, "--volumes-from", a.Id)
	newApp.Args = append(newApp.Args, args...)
	return newApp, newApp.Run()
}

func (a *App) Run() (err error) {
	if *appImage == "" {
		return errors.New("Missing -app.image")
	}

	log.Println("APP:", "Creating with", a.Args)
	cmd := dockerCommand("run", "-d")
	cmd.Args = append(cmd.Args, a.Args...)
	cmd.Args = append(cmd.Args, *appImage)
	id, err := cmd.Output()
	if err != nil {
		return
	}
	a.Id = strings.TrimSpace(string(id))

	if *containerLogs {
		go a.StartLogging()
	} else {
		td := CurrentGinkgoTestDescription()
		if len(td.ComponentTexts) > 1 {
			go a.StartLoggingToFile("logs/" + td.ComponentTexts[0] + ".txt")
		}
	}
	return
}

func (a *App) IP() (string, error) {
	if a.Id == "" {
		return "", errors.New("not started")
	}

	cmd := dockerCommand("inspect", "--format", "{{ .NetworkSettings.IPAddress }}", a.Id)
	ip, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(ip)), nil
}

func (a *App) URL(path string, query url.Values) (*url.URL, error) {
	ip, err := a.IP()
	if err != nil {
		return nil, err
	}
	return &url.URL{
		Scheme:   "http",
		Host:     ip + ":8080",
		Path:     path,
		RawQuery: query.Encode(),
	}, nil
}

func (a *App) Start() {
	if a.Id == "" {
		return
	}
	log.Println("APP:", "Starting", a.Id)
	dockerCommand("APP:", "start", a.Id).Run()
}

func (a *App) Remove() {
	if a.Id == "" {
		return
	}
	log.Println("APP:", "Removing", a.Id)
	dockerCommand("rm", "-v", "-f", a.Id).Run()
	a.Id = ""
}

func (a *App) Stop() {
	if a.Id == "" {
		return
	}
	log.Println("APP:", "Stopping", a.Id)
	dockerCommand("stop", a.Id).Run()
}

func (a *App) Recreate() (err error) {
	a.Stop()
	a.Remove()
	return a.Run()
}

func (a *App) Update(args ...string) (err error) {
	if a.Id == "" {
		return errors.New("not started")
	}
	log.Println("APP:", "Updating", a.Id)
	return dockerCommand("update", append(args, a.Id)...).Run()
}

func (a *App) Status() string {
	if a.Id == "" {
		return "not-found"
	}
	cmd := dockerCommand("inspect", "-f", "{{.State.Status}}", a.Id)
	status, err := cmd.Output()
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(string(status))
}

func (a *App) StartLogging() error {
	if a.Id == "" {
		return errors.New("not started")
	}
	r, w := io.Pipe()
	go func() {
		defer r.Close()
		prefixer.New(r, ANSI_BOLD_BLUE+"[container] "+ANSI_RESET).WriteTo(os.Stdout)
	}()
	defer w.Close()

	cmd := dockerCommand("logs", "--follow", a.Id)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}

func (a *App) StartLoggingToFile(name string) error {
	if a.Id == "" {
		return errors.New("not started")
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := dockerCommand("logs", "--follow", a.Id)
	cmd.Stdout = f
	cmd.Stderr = f
	return cmd.Run()
}

func (a *App) Logs() error {
	if a.Id == "" {
		return errors.New("not started")
	}

	var b bytes.Buffer
	log.Println("APP:", "Logs for", a.Id)
	cmd := dockerCommand("logs", a.Id)
	cmd.Stdout = &b
	cmd.Stderr = &b
	err := cmd.Run()
	if err != nil {
		return err
	}
	r := prefixer.New(&b, ANSI_BOLD_BLUE+"[container] "+ANSI_RESET)
	r.WriteTo(os.Stdout)
	return nil
}

func (a *App) Latest(n int) (messages []Tweet, err error) {
	query := url.Values{}
	query.Add("n", strconv.Itoa(n))
	err = a.request("GET", "/latest", query, &messages)
	return
}

func (a *App) LatestIds() (ids []int, err error) {
	tweets, err := a.Latest(0)
	ids = make([]int, len(tweets))
	for idx, tweet := range tweets {
		ids[idx] = tweet.UniqueID()
	}
	return
}

func (a *App) AppStatus() (string, error) {
	health, err := a.Health()
	return health.App, err
}

func (a *App) TwitterStatus() (string, error) {
	health, err := a.Health()
	return health.Twitter, err
}

func (a *App) SlackStatus() (string, error) {
	health, err := a.Health()
	return health.Slack, err
}

func (a *App) Health() (health Health, err error) {
	err = a.request("GET", "/health", nil, &health)
	return
}

func NewApp(args ...string) (a *App, err error) {
	a = &App{
		Args: args,
	}
	err = a.Run()
	return
}

func NewAppForServer(server *httptest.Server, args ...string) (a *App, err error) {
	args = append(args, "-e", "TWITTER_URL="+server.URL)
	args = append(args, "-e", "SLACK_URL="+server.URL)
	return NewApp(args...)
}
