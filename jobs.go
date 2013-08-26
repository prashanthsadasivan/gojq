package gojq

import (
    "fmt"
    "encoding/json"
    "reflect"
    "redis"
    "time"
)
type Job interface {
    Name() string
    Work(args map[string]string)
}

type JobQueue interface {
    Enqueue(j Job, args map[string] string)
}


func Register(j Job) {
    fmt.Printf("registering: %s\n",j.Name())
    jobTypeMap[j.Name()] = reflect.TypeOf(j)
}

var (
    jobschan = make(chan []byte)
    jobTypeMap = make(map[string]reflect.Type)
    redisClient *redis.Client
)

type DefaultJobQueue struct {
    name string
}

type marshalledJob struct {
    Name string
    Args map[string]string
}

func NewJobQueue(name string) JobQueue{
    if name != nil {
        //TODO: named queues
        panic("not yet supported!!")
    }
    dj := DefaultJobQueue {name: "gojq"}
    return dj
}

func (dj DefaultJobQueue) Enqueue(j Job, args map[string]string) {
    if jobTypeMap[j.Name()] == nil {
        fmt.Printf("NOT REGISTERED\n")
        Register(j)
    }
    mj := marshalledJob { Name: j.Name(), Args: args}
    b, err := json.Marshal(mj)
    if err != nil {
        panic(err.Error())
    } else {
        redisClient.Lpush(dj.name, b)
    }
}

func initRedisClient() {
    if redisClient == nil {
        redisClient, err = & redis.NewSyncClient()
        if err != nil {
            panic(err)
        }
    }
}


//only sleeps if it needs to sleep
func Dequeue() {
    for {
        //TODO: configureable sleep period
        msg, rediserr := redisClient.Rpop("gojq")
        if msg == nil  {
            fmt.Printf("msg was nil!")
            time.Sleep(5 * time.Second)
        }
        else if rediserr != nil {
            panic(rediserr.Error())
        }
        else if len(msg) == 0 {
            fmt.Printf("length was 0!")
            time.Sleep(5 * time.Second)
        } else {
            mj := marshalledJob {}
            jsonerr := json.Unmarshal(msg, &mj)
            if jsonerr != nil {
                panic(jsonerr.Error())
            } else {
                forkAndRunJob(&mj)
            }
        }
    }
}


