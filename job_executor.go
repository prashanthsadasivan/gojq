package gojq

import (
    "fmt"
    "syscall"
    "os"
    "sync"
)

type jobDescriptor struct {
    mj *marshalledJob
    pid uintptr
}

var (
    jobPidMap = make(map[uintptr]*jobDescriptor)
)


func forkAndRunJob(mj *marshalledJob) {
    //Fork to run the job
    pid, failflags, err := syscall.Syscall(syscall.SYS_FORK, 0,0,0)

    if err != 0 {
        panic("error!\n")
    }

    if failflags < 0 {
        panic("failflags!\n")
    }

    if pid > 0 {
        //parent proc
        //keep track of the job/pid for the future, incase we
        //want to cancel this job
        var wg sync.WaitGroup
        wg.Add(1)
        jd := &jobDescriptor { mj : mj, pid:pid}
        jobPidMap[pid] = jd
        fmt.Printf("parent: childpid: %d\n", pid)
        //Wait for child to finish, and when it does  
        //remove it from the process bookkeeping
        go func(childpid uintptr) {
            proc, err := os.FindProcess(int(childpid))
            if err != nil {
                panic(err.Error)
            }
            procstate, err := proc.Wait()
            fmt.Printf("parent: done with child")
            if err != nil {
                fmt.Printf("error procstate: %s\n", err.Error())
            }
            fmt.Printf("procstate: %s\n", procstate)
            if procstate.Exited() {
                fmt.Printf("parent: deleting from bookkeeping\n")
                delete(jobPidMap, childpid)
            }
            wg.Done()
        }(pid)
        wg.Wait()
        return
    } else {
        //child proc
        startRunningJob(mj)
        fmt.Printf("child: after job\n")
        return
    }
}

