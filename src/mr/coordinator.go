package mr

import (
	"bufio"
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
	// Your definitions here.
	// 你的定义在这里
}

// Your code here -- RPC handlers for the worker to call.
//你的代码写在这 -- worker调用的RPC处理程序。

// worker节点获取任务
func (c *Coordinator)GetTask()  {
	
}


// an example RPC handler.
//一个示例RPC处理程序。
// the RPC argument and reply types are defined in rpc.go.
//RPC参数和回复类型在rpc.go中定义。
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//开启一个线程，监听来自worker.go的RPC
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

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//main/mrcoordinator.go定期调用Done()，以找出整个作业是否已经完成。就是判断是否所有的任务都已经完成了
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.


	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//创建一个Coordinator。main/mrcoordinator.go调用这个函数。nReduce是要使用的reduce任务的数量。
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	// Your code here.

	for _, v := range files {
		file,err:=os.Open(v)
		if err != nil {
			log.Fatal("cannot open %v",v)
		}
		content,err:=io.ReadAll(file)
	}

	c.server()
	return &c
}


func ReadFile(f []string)  {
	for _, v := range f {
		file,err:=os.Open(v)
		if err != nil {
			log.Fatal("cannot open %v",v)
		}
		content,err:=io.ReadAll(file)
	}
}