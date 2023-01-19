package main

import (
	"fmt"
	"net/http"
)

type Master struct {
	tasks []Task
}

type Task struct {
	ID      int
	user    string
	command string
}

func NewMaster() *Master {
	instance := new(Master)

	return instance
}

func (pInstance *Master) MasterTaskSubmission(w http.ResponseWriter, r *http.Request) {
	tQuery := r.URL.Query()
	if tQuery["command"] != nil {
		sCommand := tQuery["command"]
		fmt.Printf("Receive command: %v \n", sCommand)
	}
}
