// Author Seth Hoenig 2015

// Command marathonctl is a CLI tool for Marathon
package main

import (
	"flag"
	"fmt"
	//"os"

	"github.com/marathonPac/marathonctl/check"
	mctl "github.com/marathonPac/marathonctl/marathon"
)

/*
const Help = `marathonctl <flags...> [mctl.Action] <args...>
 mctl.Actions
    app
       list                      - list all apps
       versions [id]             - list all versions of apps of id
       show [id]                 - show config and status of app of id (latest version)
       show [id] [version]       - show config and status of app of id and version
       create [jsonfile]         - deploy application defined in jsonfile
       update [jsonfile]         - update application as defined in jsonfile
       update [id] [jsonfile]    - update application id as defined in jsonfile
       update cpu [id] [cpu%]    - update application id to have cpu% of cpu share
       update memory [id] [MB]   - update application id to have MB of memory
       update instances [id] [N] - update application id to have N instances
       restart [id]              - restart app of id
       destroy [id]              - destroy and remove all instances of id

    task
       list               - list all tasks
       list [id]          - list tasks of app of id
       kill [id]          - kill all tasks of app id
       kill [id] [taskid] - kill task taskid of app id
       queue              - list all queued tasks

    group
       list                        - list all groups
       list [groupid]              - list groups in groupid
       create [jsonfile]           - create a group defined in jsonfile
       update [groupid] [jsonfile] - update group groupid as defined in jsonfile
       destroy [groupid]           - destroy group of groupid

    deploy
       list               - list all active deploys
       destroy [deployid] - cancel deployment of [deployid]

    marathon
       leader   - get the current Marathon leader
       abdicate - force the current leader to relinquish control
       ping     - ping Marathon master host[s]

    artifact
       upload [path] [file]   - upload artifact to artifacts store
       get [path]             - get artifact from store
       delete [path]          - delete artifact from store

 Flags
  -c [config file]
  -h [host]
  -u [user:password] (separated by colon)
  -f [format]
       human  (simplified columns, default)
       json   (json on one line)
       jsonpp (json pretty printed)
       raw    (the exact response from Marathon)
`

func Usage() {
	fmt.Fprintln(os.Stderr, Help)
	os.Exit(1)
}
*/

func main() {
	host, login, format, e := mctl.Config()

	if e != nil {
		fmt.Printf("config error: %s\n\n", e)
		check.Usage()
	}

	f := mctl.NewFormatter(format)
	l := mctl.NewLogin(host, login)
	c := mctl.NewClient(l)
	app := &mctl.Category{
		Actions: map[string]mctl.Action{
			"list":     mctl.AppList{c, f},
			"versions": mctl.AppVersions{c, f},
			"show":     mctl.AppShow{c, f},
			"create":   mctl.AppCreate{c, f},
			"update":   mctl.AppUpdate{c, f},
			"restart":  mctl.AppRestart{c, f},
			"destroy":  mctl.AppDestroy{c, f},
		},
	}
	task := &mctl.Category{
		Actions: map[string]mctl.Action{
			"list":  mctl.TaskList{c, f},
			"kill":  mctl.TaskKill{c, f},
			"queue": mctl.TaskQueue{c, f},
		},
	}
	group := &mctl.Category{
		Actions: map[string]mctl.Action{
			"list":    mctl.GroupList{c, f},
			"create":  mctl.GroupCreate{c, f},
			"update":  mctl.GroupUpdate{c, f},
			"destroy": mctl.GroupDestroy{c, f},
		},
	}
	deploy := &mctl.Category{
		Actions: map[string]mctl.Action{
			"list":   mctl.DeployList{c, f},
			"cancel": mctl.DeployCancel{c, f},
		},
	}
	marathon := &mctl.Category{
		Actions: map[string]mctl.Action{
			"leader":   mctl.MarathonLeader{c, f},
			"abdicate": mctl.MarathonAbdicate{c, f},
			"ping":     mctl.MarathonPing{c, f},
		},
	}
	artifact := &mctl.Category{
		Actions: map[string]mctl.Action{
			"upload": mctl.ArtifactUpload{c, f},
			"get":    mctl.ArtifactGet{c, f},
			"delete": mctl.ArtifactDelete{c, f},
		},
	}
	t := &mctl.Tool{
		Selections: map[string]mctl.Selector{
			"app":      app,
			"task":     task,
			"group":    group,
			"deploy":   deploy,
			"marathon": marathon,
			"artifact": artifact,
		},
	}

	t.Start(flag.Args())
}

/*
func Check(b bool, args ...interface{}) {
	if !b {
		fmt.Fprintln(os.Stderr, args...)
		os.Exit(1)
	}
}
*/
