package mr

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

// Coordinator is the master of a MapReduce job.
//协调器是MapReduce作业的主节点。

type Coordinator struct {
	TaskList   chan *TaskReply
	WorkerList map[string]bool
}

// worker节点获取任务
func (c *Coordinator) GetTask() *TaskReply {
	select {
	case r := <-c.TaskList:
		c.WorkerList[r.TaskName] = true
		return r
	}
}

func (c *Coordinator) PostWorkerOut(output *PostMapRes) {
	c.TaskList <- &TaskReply{
		ReduceSource: output.MapOutput,
		TaskName:     ,
	}
}

// start a thread that listens for RPCs from worker.go
// 开启一个线程，监听来自worker.go的RPC
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
// main/mrcoordinator.go定期调用Done()，以找出整个作业是否已经完成。就是判断是否所有的任务都已经完成了
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
// 创建一个Coordinator。main/mrcoordinator.go调用这个函数。nReduce是要使用的reduce任务的数量。
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	c.TaskList = make(chan *TaskReply)
	c.WorkerList = make(map[string]bool)
	// Your code here.

	for _, v := range files {
		file, err := os.Open(v)
		if err != nil {
			log.Fatal("cannot open %v", v)
		}
		content, err := io.ReadAll(file)
		res := bytes.NewBuffer(content)
		c.TaskList <- &TaskReply{
			MapSource: *res,
			TaskName:  v,
			TaskType:  MapTask,
		}
	}
	c.server()
	return &c
}
