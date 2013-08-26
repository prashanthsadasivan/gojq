package main
import (
    "flag"
    "fmt"
    "strings"
    "os"
    "io"
    "text/template"
)

type Command struct {
    Run func(args []string)
    UsageLine, Short, Long string
}

func (cmd *Command) Name() string {
    name := cmd.UsageLine
    i := strings.Index(name, " ")
    if i >= 0 {
        name = name[:i]
    }
    return name
}

var commands = []*Command {
    cmdWork,
}

func main() {
    fmt.Fprintf(os.Stdout, header)
    flag.Usage = func() { usage(1) }
    flag.Parse()
    args := flag.Args()
    if len(args) < 1 || args[0] == "help" {
        if len(args) == 1 {
            usage(0)
        }
        if len(args) > 1 {
            for _, cmd := range commands {
                if cmd.Name() == args[1] {
                    tmpl(os.Stdout, helpTemplate, cmd)
                    return
                }
            }
        }
        usage(2)
    }

    for _, cmd := range commands {
        if cmd.Name() == args[0] {
            cmd.Run(args[1:])
            return
        }
    }
    errorf("unknown command%q\nRun gojq help for usage help\n", args[0])
}

func errorf(format string, args ...interface{}) {
    if !strings.HasSuffix(format, "\n") {
        format += "\n"
    }
    fmt.Fprintf(os.Stderr, format, args...)
    os.Exit(0)
}


const header = `
gojq: because things take time

`

const usageTemplate = `
usage: gojq command [arguments]
The commands are: 
{{range .}} 
{{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "gojq help [command]" for more information.
`

var helpTemplate = `
usage: gojq {{.UsageLine}} {{.Long}}

`

func usage(exitCode int) {
    tmpl(os.Stderr, usageTemplate, commands)
    os.Exit(exitCode)
}

func tmpl(w io.Writer,text string, data interface{}) {
    t := template.New("top")
    template.Must(t.Parse(text))
    if err := t.Execute(w, data); err != nil {
        panic(err)
    }
}


