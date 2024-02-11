package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//RPC 定义，记得将所有名称大写。

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//示例显示如何声明RPC的参数和回复。

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.

type WorkerAliveReply struct{
	Done bool
}

type TaskReply struct{
	TaskSource bytes.Buffer
	TaskName string
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
//准备一个在/var/tmp中唯一的UNIX域套接字名称，用于协调器。不能使用当前目录，因为Athena AFS不支持UNIX域套接字。
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
