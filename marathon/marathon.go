package marathon

// All actions under command marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/marathonPac/marathonctl/check"
)

// ping (todo ping all hosts)
type MarathonPing struct {
	Clients *Client
	Formats Formatter
}

func (p MarathonPing) Apply(args []string) {
	hosts := p.Clients.login.Hosts
	timings := make(map[string]time.Duration)
	for _, host := range hosts {
		request, e := http.NewRequest("/ping", host, nil)
		check.Check(e == nil, "could not create ping request")
		p.Clients.tweak(request)
		start := time.Now()
		_, err := p.Clients.client.Do(request)
		var elapsed time.Duration
		if err == nil {
			elapsed = time.Now().Sub(start)
		}
		timings[host] = elapsed
	}

	var b bytes.Buffer
	for host, duration := range timings {
		b.WriteString(host)
		b.WriteString(" ")
		if duration == 0 {
			b.WriteString("-")
		} else {
			b.WriteString(duration.String())
		}
		b.WriteString("\n")
	}
	fmt.Println(p.Formats.Format(strings.NewReader(b.String()), p.Humanize))
}

func (P MarathonPing) Humanize(body io.Reader) string {
	b, e := ioutil.ReadAll(body)
	check.Check(e == nil, "reading ping response failed", e)
	text := "HOST DURATION\n" + string(b)
	return Columnize(text)
}

// leader
type MarathonLeader struct {
	Clients *Client
	Formats Formatter
}

func (l MarathonLeader) Apply(args []string) {
	request := l.Clients.GET("/v2/leader")
	response, e := l.Clients.Do(request)
	check.Check(e == nil, "get leader failed", e)
	c := response.StatusCode
	check.Check(c == 200, "get leader bad status", c)
	defer response.Body.Close()
	fmt.Println(l.Formats.Format(response.Body, l.Humanize))
}

func (l MarathonLeader) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var which Which
	e := dec.Decode(&which)
	check.Check(e == nil, "failed to decode response", e)
	text := "LEADER\n" + which.Leader
	return Columnize(text)
}

// abdicate
type MarathonAbdicate struct {
	Clients *Client
	Formats Formatter
}

func (a MarathonAbdicate) Apply(args []string) {
	request := a.Clients.DELETE("/v2/leader")
	response, e := a.Clients.Do(request)
	check.Check(e == nil, "abdicate request failed", e)
	c := response.StatusCode
	check.Check(c == 200, "abdicate bad status", c)
	defer response.Body.Close()
	fmt.Println(a.Formats.Format(response.Body, a.Humanize))
}

func (a MarathonAbdicate) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var mess Message
	e := dec.Decode(&mess)
	check.Check(e == nil, "failed to decode response", e)
	return "MESSAGE\n" + mess.Message
}
