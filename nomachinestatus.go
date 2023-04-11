package main

import (
	"errors"
	"os"

	"github.com/mitchellh/go-ps"
)

type nomachineStatus struct {
	HostName         string `msgpack:",omitempty"`
	NoMachineRunning bool   `msgpack:",omitempty"`
	ClientAttached   bool   `msgpack:",omitempty"`
}

func getFirstProcessByName(name string) (int, error) {
	processes, err := ps.Processes()
	if err != nil {
		return -1, err
	}

	for _, process := range processes {
		if process.Executable() == name {
			return process.Pid(), nil
		}
	}

	return -1, errors.New("Could not find process " + name)
}

// getStatus returns a structure containing information on the status of NoMachine
func getStatus() (nomachineStatus, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return nomachineStatus{}, err
	}
	status := nomachineStatus{HostName: hostName}

	noMachinePid, _ := getFirstProcessByName("nxserver.bin")
	noMachineClientPid, _ := getFirstProcessByName("nxexec")

	if noMachinePid >= 0 {
		status.NoMachineRunning = true
	}
	if noMachineClientPid >= 0 {
		status.ClientAttached = true
	}
	return status, nil
}
