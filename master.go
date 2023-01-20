package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tony2037/WooPoC/pb"
	"google.golang.org/grpc"
)

type Master struct {
	tasks      []Task
	add_client pb.AddServiceClient
}

type Task struct {
	ID      int
	user    string
	command string
}

func NewMaster() *Master {
	instance := new(Master)

	//	Connect to slave service
	//var target string = "add-service:3000"
	var target string = "localhost:3000"
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	instance.add_client = pb.NewAddServiceClient(conn)
	return instance
}

func (pInstance *Master) MasterTaskSubmission(w http.ResponseWriter, r *http.Request) {
	tQuery := r.URL.Query()
	if tQuery["command"] != nil {
		sCommand := tQuery["command"]
		fmt.Printf("Receive command: %v \n", sCommand)
	}
}

func (pInstance *Master) MasterAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	tQuery := r.URL.Query()
	a, err := strconv.ParseUint(tQuery.Get("a"), 10, 64)
	if err != nil {
		log.Fatal("Invalid parameter a")
		w.WriteHeader(http.StatusBadRequest)
	}
	b, err := strconv.ParseUint(tQuery.Get("b"), 10, 64)
	if err != nil {
		log.Fatal("Invalid parameter b")
		w.WriteHeader(http.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	req := &pb.AddRequest{A: a, B: b}
	if resp, err := pInstance.add_client.Compute(ctx, req); err == nil {
		msg := fmt.Sprintf("Summation is %d \n", resp.Result)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, msg)
	} else {
		msg := fmt.Sprintf("Internal server error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, msg)
	}
}
