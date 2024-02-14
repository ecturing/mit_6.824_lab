package mr

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/rpc"
	"os"
	"sort"
)

type Job struct{}

type DataHandle interface{
	//  序列化函数
	Serialization(data []KeyValue)
	// 反序列化函数
	DeSerialization() []KeyValue
}

func (j Job) Serialization(data []KeyValue) {
	intermediateFile,err:=os.Create("")
	if err != nil {
		return nil, err
	}
	writer:=bufio.NewWriter(intermediateFile)
	for _, v := range data {
		writer.WriteString(v.Key+":"+v.Value)
	}
	writer.Flush()

}

func (j Job) DeSerialization(file *os.File) []KeyValue{

}

// Map functions return a slice of KeyValue.
// Map函数返回一个KeyValue切片。
type KeyValue struct {
	Key   string
	Value string
}

type SortBy []KeyValue

func (a SortBy) Len() int           { return len(a) }
func (a SortBy) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortBy) Less(i, j int) bool { return a[i].Key < a[j].Key }

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
// 使用ihash(key) % NReduce选择每个KeyValue由Map发出的reduce任务编号。
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
		var J Job
	task := CallGetTask()
	switch task.TaskType {
	case MapTask:
		file,err:=os.Open(task.MapSourceAddr)
		if err != nil {
			return nil, err
		}
		fileBytes,_:=io.ReadAll(file)
		output := mapf(task.MapSourceAddr, string(fileBytes))
		J.Serialization(output)
		call("Coordinator.PostWorkerOut", &PostMapRes{MapOutput: output}, nil)
	case ReduceTask:
		file,err:=os.Open(task.ReduceSourceAddr)
		if err != nil {
			return nil, err
		}
		// 反序列化文件为[]keyvalue切片，对切片进行操作
		ReduceSource:=J.DeSerialization(file)
		sort.Sort(SortBy(ReduceSource))
		intermediate:=ReduceSource
		i := 0
		for i < len(intermediate) {
			j := i + 1
			for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
				j++
			}
			values := []string{}
			for k := i; k < j; k++ {
				values = append(values, intermediate[k].Value)
			}
			output := reducef(intermediate[i].Key, values)
			i = j
		}
	}
}

func CallGetTask() *TaskReply {
	args := TaskReply{}
	ok := call("Coordinator.GetTask", nil, &args)
	if ok {
		return &args
	} else {
		return nil
	}
}

// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
// 发送一个RPC请求到协调器，等待响应。通常返回true。如果出现问题，则返回false。
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
