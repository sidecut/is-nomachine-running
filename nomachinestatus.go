package main

import (
	"errors"
	"os"

	"github.com/mitchellh/go-ps"
)

const (
	ProcessNotFound   = -1
	serverProcessName = "nxserver.bin"
	clientProcessName = "nxexec"
)

// nomachineStatus represents the current state of NoMachine server and client
type nomachineStatus struct {
	HostName         string `json:"host_name,omitempty"`
	NoMachineRunning bool   `json:"no_machine_running,omitempty"`
	ClientAttached   bool   `json:"client_attached,omitempty"`
}

// IsActive returns true if NoMachine is running and has a client attached
func (s nomachineStatus) IsActive() bool {
	return s.NoMachineRunning && s.ClientAttached
}

func getFirstProcessByName(name string) (int, error) {
	if name == "" {
		return ProcessNotFound, errors.New("process name cannot be empty")
	}

	processes, err := ps.Processes()
	if err != nil {
		return ProcessNotFound, err
	}

	for _, process := range processes {
		if process.Executable() == name {
			return process.Pid(), nil
		}
	}

	return ProcessNotFound, errors.New("could not find process " + name)
}

// getStatus returns information about the NoMachine server and client status
func getStatus() (nomachineStatus, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return nomachineStatus{}, err
	}

	status := nomachineStatus{HostName: hostName}

	noMachinePid, err := getFirstProcessByName(serverProcessName)
	if err != nil {
		// Consider logging the error
	}

	noMachineClientPid, err := getFirstProcessByName(clientProcessName)
	if err != nil {
		// Consider logging the error
	}

	if noMachinePid >= 0 {
		status.NoMachineRunning = true
	}
	if noMachineClientPid >= 0 {
		status.ClientAttached = true
	}

	return status, nil
}
