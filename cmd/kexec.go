// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package main

import (
	"fmt"
	"khalehla/kexec"
	"khalehla/kexec/config"
	"khalehla/kexec/exec"
	"khalehla/klog"
	"os"
	"strconv"
	"strings"
)

// This is effectively the SSP for the exec.
// We set up the exec, then boot it, and then watch over it until it stops,
// at which point we either restart it or invoke a dump (just in case) and terminate.

type context struct {
	configFileName *string
	jumpKeys       []bool
	helpFlagSet    bool
}

func showHelp() {
	fmt.Println("Usage:")
	fmt.Println("  kexec [-h] [-c config_file] [-jk jump_key,...]")
	fmt.Println("If config_file is not specified, system defaults will be used,")
	fmt.Println("  and disk packs and tape reels will be searched for in the current directory.")
	fmt.Println("jump_key values are integers from 1 to 36 - the recognized values are:")
	fmt.Printf("  jk1:  Force display of Modify Config message to allow facilities keyins prior to boot")
	fmt.Printf("  jk2:  Perform partial manual dump prior to booting the exec")
	fmt.Printf("  jk3:  Disable auto-recovery - when the exec stops, this application will terminate")
	fmt.Printf("  jk4:  Reload system libraries from system library tape")
	fmt.Printf("  jk6:  Performs full dump in conjunction with jk2")
	fmt.Printf("  jk7:  Solicits TIP initialization;")
	fmt.Printf("          solicits pack recovery/initialization during recovery")
	fmt.Printf("  jk9:  On initial boot, prevents recovery of backlog and print queues")
	fmt.Printf("  jk13: On initial boot, initializes mass storage - requires jk4")
}

func showUsage() {
	fmt.Println("Usage:")
	fmt.Println("  kexec [-h] [-c config_file] [-j jump_key,...]")
}

// processArgs reads the command line arguments and puts together a configuration and an array of jump keys
// to be set on exec startup... or an error
func processArgs(args []string) (*context, error) {
	context := context{}
	context.jumpKeys = make([]bool, 36)

	for ax := 0; ax < len(args); {
		sw := args[ax]
		ax++
		switch strings.ToLower(sw) {
		case "-c":
			if ax == len(args) {
				return nil, fmt.Errorf("no argument specified for -c switch")
			}
			context.configFileName = &args[ax]
			ax++

		case "-h":
			context.helpFlagSet = true

		case "-j":
			if ax == len(args) {
				return nil, fmt.Errorf("no argument specified for -j switch")
			}
			split := strings.Split(args[ax], ",")
			ax++
			for _, key := range split {
				jk, err := strconv.Atoi(key)
				if err != nil || jk < 1 || jk > 36 {
					return nil, fmt.Errorf("invalid jump key %v", key)
				}
				context.jumpKeys[jk-1] = true
			}
		}
	}

	return &context, nil
}

func main() {
	args := os.Args
	context, err := processArgs(args)
	if err != nil {
		fmt.Println(err)
		showUsage()
		return
	} else if context.helpFlagSet {
		showHelp()
		return
	}

	klog.ClearLoggers()
	klog.RegisterLogger(klog.NewTimestampedFileLogger(klog.LevelAll, "kexec"))
	klog.SetGlobalLevel(klog.LevelAll)
	defer klog.Close()

	cfg := config.NewConfiguration()
	if context.configFileName != nil {
		err := cfg.UpdateFromFile(*context.configFileName)
		if err != nil {
			fmt.Printf("Error in configuration file:%v", err)
			return
		}
	}

	e := exec.NewExec(cfg)
	err = e.Initialize()
	if err != nil {
		fmt.Println("::Cannot continue - error in exec initialization")
		fmt.Printf("::%v\n", err)
		e.Close()
		return
	}

	e.SetConfiguration(cfg)

	channel := make(chan kexec.StopCode)
	session := uint(0)

	for {
		if context.jumpKeys[kexec.JumpKey2Index] {
			fmt.Println("::Performing pre-boot system dump...")
			fileName, err := e.PerformDump(e.GetJumpKey(6))
			if err != nil {
				fmt.Printf("::Error producing dump file:%v\n", err)
			} else {
				fmt.Printf("::Dump written to file %v\n", fileName)
			}
		}

		fmt.Printf("::Starting KEXEC session %03v...\n", session)
		go e.Boot(session, context.jumpKeys, channel)

		stopCode := <-channel

		fmt.Printf("::System error %03v terminated session %03v\n", stopCode, session)
		if e.GetJumpKey(3) {
			fmt.Printf("::Auto-recovery inhibited - producing final post-mortem dump...\n")
			fileName, err := e.PerformDump(e.GetJumpKey(6))
			if err != nil {
				fmt.Printf("::Error producing dump file:%v\n", err)
			} else {
				fmt.Printf("::Dump written to file %v\n", fileName)
			}
			break
		} else {
			fmt.Println("::Recovering system...")
			session++
		}
	}

	close(channel)
	e.Close()
}
