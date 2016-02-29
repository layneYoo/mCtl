package check

import (
	"fmt"
	"os"
)

const Help = `marathonctl <flags...> [action] <args...>
 Actions
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

    image
       build  [buildpath] [registrypath] [gitURL] [deployJson]
                              - build docker images in buildpath and named registrypath:gitversion
       upload [buildpath] [registry]
                              - upload docker images to docker registry

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

func Check(b bool, args ...interface{}) {
	if !b {
		fmt.Fprintln(os.Stderr, args...)
		os.Exit(1)
	}
}
