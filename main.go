package main

import (
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"time"

	"os"
	"os/signal"

	"github.com/urfave/cli"
)

type Params struct {
	address   string
	token     string
	directory string
	verbose   bool
}

var (
	version           string
	fate              Fate
	param             Params
	watcher           Watcher
	readDir           = ioutil.ReadDir
	readFile          = ioutil.ReadFile
	signalNotify      = signal.Notify
	entrypoint        = run
	delayReadDuration = time.Duration(1500) * time.Millisecond
)

var helpTmpl = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[options...]{{end}}

EXAMPLE:
   {{.HelpName}} --vault-address '0.0.0.0:8200' --vault-token root --dir $PWD/your-configs-dir
{{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}
{{- if .Commands}}
OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
{{- end}}
{{- if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
{{end}}
{{- if .Version}}
VERSION:
   {{.Version}}
{{end}}
`

func main() {
	app := cli.NewApp()

	app.Name = "Dr. Vault"
	app.Usage = "Vault folder monitoring"
	app.Version = version

	app.Author = "Ninh Pham #ReeganExE -> ninh.js.org"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "vault-address, a",
			EnvVar:      "VAULT_DEV_LISTEN_ADDRESS",
			Destination: &param.address,
			Value:       "0.0.0.0:8200",
			Usage:       "Vault address that dr-vault will connect to",
		},
		cli.StringFlag{
			Name:        "vault-token, t",
			EnvVar:      "VAULT_DEV_ROOT_TOKEN_ID",
			Value:       "root",
			Usage:       "A writable token",
			Destination: &param.token,
		},
		cli.StringFlag{
			Name:        "dir, d",
			EnvVar:      "MONITOR_DIR",
			Value:       "/var/source",
			Usage:       "Specify a directory to monitor.",
			Destination: &param.directory,
		},
		cli.BoolFlag{
			Name:        "verbose, p",
			EnvVar:      "VERBOSE",
			Destination: &param.verbose,
		},
	}

	cli.AppHelpTemplate = helpTmpl

	app.Action = entrypoint

	e := app.Run(os.Args)

	if e != nil {
		log.Fatal(e)
	}
}

func run(c *cli.Context) {
	fate.client = NewVaultClient(param.address, param.token)

	dir := param.directory

	fsWatcher, e := fsnotify.NewWatcher()

	if e != nil {
		panic(e)
	}

	watcher = &FsWatcher{fsWatcher}

	scanDir(dir)
	watchDir(dir, watcher)
}
