package marathon

// All actions under command group

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	"github.com/marathonPac/marathonctl/check"
)

type GroupList struct {
	Clients *Client
	Formats Formatter
}

func (g GroupList) Apply(args []string) {
	switch len(args) {
	case 0:
		g.listGroups("")
	case 1:
		g.listGroups(args[0])
	default:
		check.Check(false, "expected 0 or 1 argument")
	}
}

func (g GroupList) listGroups(groupid string) {
	path := "/v2/groups"
	if groupid != "" {
		path += "/" + url.QueryEscape(groupid)
	}
	request := g.Clients.GET(path)
	response, e := g.Clients.Do(request)
	check.Check(e == nil, "failed to get response", e)
	defer response.Body.Close()
	fmt.Println(g.Formats.Format(response.Body, g.Humanize))
}

func (g GroupList) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var root Group
	e := dec.Decode(&root)
	check.Check(e == nil, "failed to unmarshal response", e)
	return columnizeGroup(&root)
}

func columnizeGroup(group *Group) string {
	title := "GROUPID VERSION GROUPS APPS\n"
	var b bytes.Buffer
	gatherGroup(group, &b)
	text := title + b.String()
	return Columnize(text)
}

func gatherGroup(g *Group, b *bytes.Buffer) {
	b.WriteString(g.GroupID)
	b.WriteString(" ")
	b.WriteString(g.Version)
	b.WriteString(" ")
	b.WriteString(strconv.Itoa(len(g.Groups)))
	b.WriteString(" ")
	b.WriteString(strconv.Itoa(len(g.Apps)))
	b.WriteString("\n")
	for _, group := range g.Groups {
		gatherGroup(group, b)
	}
}

type GroupCreate struct {
	Clients *Client
	Formats Formatter
}

func (g GroupCreate) Apply(args []string) {
	check.Check(len(args) == 1, "must supply 1 jsonfile")
	f, e := os.Open(args[0])
	check.Check(e == nil, "failed to open jsonfile", e)
	defer f.Close()
	request := g.Clients.POST("/v2/groups", f)
	response, e := g.Clients.Do(request)
	check.Check(e == nil, "failed to get response")
	defer response.Body.Close()
	check.Check(response.StatusCode != 409, "group already exists")
	fmt.Println(g.Formats.Format(response.Body, g.Humanize))
}

func (g GroupCreate) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var update Update
	e := dec.Decode(&update)
	check.Check(e == nil, "failed to decode response", e)
	title := "DEPLOYID VERSION\n"
	text := title + update.DeploymentID + " " + update.Version
	return Columnize(text)
}

type GroupDestroy struct {
	Clients *Client
	Formats Formatter
}

func (g GroupDestroy) Apply(args []string) {
	check.Check(len(args) == 1, "must specify groupid")
	groupid := url.QueryEscape(args[0])
	path := "/v2/groups/" + groupid
	request := g.Clients.DELETE(path)
	response, e := g.Clients.Do(request)
	check.Check(e == nil, "destroy group failed", e)
	defer response.Body.Close()
	c := response.StatusCode
	check.Check(c != 404, "unknown group")
	check.Check(c == 200, "destroy group bad status", c)

	fmt.Println(g.Formats.Format(response.Body, g.Humanize))
}

func (g GroupDestroy) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var versionmap map[string]string // ugh
	e := dec.Decode(&versionmap)
	check.Check(e == nil, "failed to decode response", e)
	v, ok := versionmap["version"]
	check.Check(ok, "version missing")
	return "VERSION\n" + v
}

type GroupUpdate struct {
	Clients *Client
	Formats Formatter
}

func (g GroupUpdate) Apply(args []string) {
	check.Check(len(args) == 2, "must specify groupid and jsonfile")
	groupid := url.QueryEscape(args[0])
	f, e := os.Open(args[1])
	check.Check(e == nil, "failed to open jsonfile", e)
	defer f.Close()
	request := g.Clients.PUT("/v2/groups/"+groupid, f)
	response, e := g.Clients.Do(request)
	check.Check(e == nil, "failed to get response", e)
	defer response.Body.Close()

	sc := response.StatusCode
	check.Check(sc == 200, "bad status code", sc)

	fmt.Println(g.Formats.Format(response.Body, g.Humanize))
}

func (g GroupUpdate) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var update Update
	e := dec.Decode(&update)
	check.Check(e == nil, "failed to decode response", e)
	title := "DEPLOYID VERSION\n"
	text := title + update.DeploymentID + " " + update.Version
	return Columnize(text)
}
