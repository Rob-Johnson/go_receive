package goreceive

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// PIDs of the Unicorn Master + oldMaster once respawned
var (
	masterPID = os.Getenv("UNICORN_PID")
	appDir    = os.Getenv("APP_DIR")
)

func RedeployCodebase(commitID string) error {
	err := checkEnv()
	if err != nil {
		return err
	}

	err = gitFetch()
	if err != nil {
		return err
	}

	err = gitReset(commitID)
	if err != nil {
		return err
	}

	err = gracefulRestart()
	if err != nil {
		return err
	}

	return err
}

func checkEnv() error {
	if appDir == "" {
		return errors.New("APP_DIR env missing")
	}
	if masterPID == "" {
		return errors.New("UNICORN_PID env missing")
	}
	return nil
}

func gitFetch() error {
	fetchCmd := exec.Command("git", "fetch", "origin", "master")
	fetchCmd.Dir = appDir
	err := fetchCmd.Run()
	return err

}

func gitReset(commitID string) error {
	resetCmd := exec.Command("git", "reset", "--hard", commitID)
	log.Print(resetCmd)
	resetCmd.Dir = appDir
	err := resetCmd.Run()
	return err
}

//send USR2 to the unicorn master
//let the app deal with killing of old process once it's alive
//http://unicorn.bogomips.org/SIGNALS.html
func gracefulRestart() error {
	pid, err := getPID(masterPID)
	if err != nil {
		return err
	}

	unicornPID, err := strconv.Atoi(pid)
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(unicornPID)
	err = proc.Signal(syscall.SIGUSR2)
	return err
}

func getPID(pidfile string) (string, error) {
	pidCmd, err := ioutil.ReadFile(pidfile)
	return strings.TrimSpace(fmt.Sprintf("%s", pidCmd)), err
}
