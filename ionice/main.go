package main

import "github.com/kormoc/ionice"
import "github.com/spf13/pflag"
import "log"

func main() {
    var classStr string
    var classdata uint32
    var pid int
    var pgid int
    var ignore bool
    var uid int

    pflag.StringVarP( &classStr,     "class", "c",    "", "name or number of scheduling class, 0: none, 1: realtime, 2: best-effort, 3: idle")
    pflag.Uint32VarP(&classdata, "classdata", "n",     0, "priority (0..7) in the specified scheduling class, only for the realtime and best-effort classes")
    pflag.IntVarP(         &pid,       "pid", "p",     0, "act on these already running processes")
    pflag.IntVarP(        &pgid,      "pgid", "P",     0, "act on already running processes in these groups")
    pflag.BoolVarP(     &ignore,    "ignore", "t", false, "ignore failures")
    pflag.IntVarP(         &uid,       "uid", "u",     0, "act on already running processes owned by these users")

    pflag.Parse()

    // Setup logger
    log.SetFlags(0)

    // Validate we have one target
    var whoCount = 0
    if pid != 0 {
        whoCount += 1
    }
    if pgid != 0 {
        whoCount += 1
    }
    if uid != 0 {
        whoCount += 1
    }
    if whoCount != 1 {
        log.Fatalln("We require only one of --pid, --pgid, or --uid")
    }

    // Figure out the which and who arguments

    var which int
    var who int

    if pid != 0 {
        which = ionice.IOPRIO_WHO_PROCESS
        who = pid
    }

    if pgid != 0 {
        which = ionice.IOPRIO_WHO_PGRP
        who = pgid
    }

    if uid != 0 {
        which = ionice.IOPRIO_WHO_USER
        who = uid
    }

    if classStr == "" {
        // Just viewing the current data
        prio, err := ionice.GetIOPriority(which, who)
        if err != nil {
            log.Fatalf("FATAL: %v\n", err)
        }
        class, classdata := ionice.PrioToClassAndClassdata(prio)

        if class == ionice.IOPRIO_CLASS_NONE || class == ionice.IOPRIO_CLASS_IDLE {
            log.Printf("%v\n", ionice.ClassToString[class])
        } else {
            log.Printf("%v: prio %v\n", ionice.ClassToString[class], classdata)
        }
        return
    } else {
        class, present := ionice.StringToClass[classStr]
        if !present {
            log.Fatalf("unknown scheduling class: '%s'\n", classStr)
        }

        prio := ionice.ClassAndClassdataToPrio(class, classdata)

        if err := ionice.SetIOPriority(which, who, prio); err != nil {
            log.Fatalf("FATAL: %v\n", err)
        }
        return
    }
}
