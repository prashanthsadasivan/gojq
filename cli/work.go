package main

import (
    "github.com/prashanthsadasivan/gojq"
    "fmt"
    "sync"
)

type marshalledJob struct {
    Name string
    Args map[string]string
}

var cmdWork = &Command {
    UsageLine: "work [-w workers]",
    Short: "start gojq to work",
    Long: `
    Work starts gojq to work!  
    `,
}

func init() {
    cmdWork.Run = work
}

func work(args []string){
    fmt.Printf("len: %d\n", len(args))
    if len(args) != 0 {
        if len(args) != 2 {
            errorf("You done goofed.")
        }
    }
    if len(args) == 2 {
        errorf("Supplying the number of workers is not yet supported!")
    }
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        gojq.Dequeue(0)
        wg.Done()
    }()
    wg.Wait()
}

