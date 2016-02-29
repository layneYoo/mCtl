package marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/marathonPac/marathonctl/check"
)

// deploy [actions]

type DeployList struct {
	Clients *Client
	Formats Formatter
}

func (d DeployList) Apply(args []string) {
	check.Check(len(args) == 0, "no arguments")
	request := d.Clients.GET("/v2/deployments")
	response, e := d.Clients.Do(request)
	check.Check(e == nil, "failed to get response")

	defer response.Body.Close()
	fmt.Println(d.Formats.Format(response.Body, d.Humanize))
}

func (d DeployList) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var deploys Deploys
	e := dec.Decode(&deploys)
	check.Check(e == nil, "failed to unmarshal response", e)
	title := "DEPLOYID VERSION PROGRESS APPS\n"
	var b bytes.Buffer
	for _, deploy := range deploys {
		b.WriteString(deploy.DeployID)
		b.WriteString(" ")
		b.WriteString(deploy.Version)
		b.WriteString(" ")
		b.WriteString(strconv.Itoa(deploy.CurrentStep))
		b.WriteString("/")
		b.WriteString(strconv.Itoa(deploy.TotalSteps))
		b.WriteString(" ")
		for _, app := range deploy.AffectedApps {
			b.WriteString(app)
		}
		b.UnreadRune()
		b.WriteString("\n")
	}
	text := title + b.String()
	return Columnize(text)
}

type DeployCancel struct {
	Clients *Client
	Formats Formatter
}

func (d DeployCancel) Apply(args []string) {
	check.Check(len(args) == 1, "must supply deployid")
	deployid := url.QueryEscape(args[0])
	path := "/v2/deployments/" + deployid
	request := d.Clients.DELETE(path)
	response, e := d.Clients.Do(request)
	check.Check(e == nil, "failed to cancel deploy", e)
	defer response.Body.Close()
}

func (d DeployCancel) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var rollback Update
	e := dec.Decode(&rollback)
	check.Check(e == nil, "failed to decode response", e)
	title := "DEPLOYID VERSION\n"
	text := title + rollback.DeploymentID + " " + rollback.Version
	return Columnize(text)
}
