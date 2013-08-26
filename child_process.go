package gojq

import (
    "fmt"
    "reflect"
    "os"
)

func startRunningJob(mj *marshalledJob) {
    val := jobTypeMap[mj.Name]
    theJob := reflect.New(val).Elem().Interface().(Job)
    fmt.Printf("running Job: %s\n", theJob.Name())
    theJob.Work(mj.Args)
    fmt.Printf("after work call\n")
    fmt.Printf("done with child job\n")
    os.Exit(0)
}

