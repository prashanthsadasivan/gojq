package gojq

import (
    "fmt"
    "encoding/json"
    "reflect"
)
type Job interface {
    Name() string
    Work(args map[string]string)
}

type marshalledjob struct {
    Name string
    Args map[string]string
}

func Register(j Job) {
    fmt.Printf("registering: %s\n",j.Name())
    jobsmap[j.Name()] = reflect.TypeOf(j)
}

var (
    jobschan = make(chan []byte)
    jobsmap = make(map[string]reflect.Type)

)

func Enqueue(j Job, args map[string]string) {
    mj := marshalledjob { Name: j.Name(), Args: args}
    b, err := json.Marshal(mj)
    if err != nil {
        panic(err.Error())
    } else {
        jobschan <- b
    }
}

func dequeue(id int) {
    for {
        msg := <-jobschan
        mj := marshalledjob {}
        err := json.Unmarshal(msg, &mj)
        if err != nil {
            panic(err.Error())
        } else {
            val := jobsmap[mj.Name]
            theJob := reflect.New(val).Elem().Interface().(Job)
            go theJob.Work(mj.Args)
        }
    }
}

func init() {
    go dequeue(0)
}


