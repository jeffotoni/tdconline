package main

import (
	"context"
	"fmt"
	"time"

	engdb_pb "github.com/jeffotoni/tdconline/grpc/proto"
	ggrpc "google.golang.org/grpc"
)

type Client struct {
	Host    string
	Timeout time.Duration
}

func (c *Client) Work(job *engdb_pb.Job) (res *engdb_pb.Result, err error) {
	conn, err := ggrpc.Dial(c.Host, ggrpc.WithInsecure())
	fmt.Printf("\nNew connection %v, timeout[%s]", c.Host, c.Timeout)
	if err != nil {
		fmt.Errorf("Error dial gRPC Worker client \n error[%v]", err)
		return nil, err
	}
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	w := engdb_pb.NewWorkerServiceClient(conn)
	res, err = w.Perform(ctx, job)
	if err != nil {
		fmt.Errorf("Error performing gRPC work error[%v] \n job[%v]", err, job)
		return res, err
	}
	return res, nil
}

func main() {
	var c = Client{"localhost:8001", time.Second * 30}
	fmt.Println(c.Work(createJob()))
}

func createJob() *engdb_pb.Job {
	job := &engdb_pb.Job{}
	job.Name = "Job Eng..2020..."
	job.Id = "1234696969624696xx"
	job.Params = []byte(`{"x": "42"}`)

	headers := make(map[string]string)
	headers["key1"] = "value1"
	headers["key2"] = "value2"
	job.Headers = headers

	configs := make(map[string]string)
	configs["etc1"] = "value1"
	configs["etc2"] = "value2"
	job.Configs = configs
	return job
}
