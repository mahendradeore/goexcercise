# 1. Introduction
  gRPC service for undirected non-weighted graph.
    Service creates a graph from the specified payload and returns an id for the graph
    Finds the shortest path for the given graph id
    Deletes the graph for the given graph id
    server  has concurrent support.
    use go pkg github.com/yourbasic/graph for graph operation.
# 2 How to build
  **# 2.1 Compile server**
         $cd server
         ** For Windows **
          $ go build -o server.exe main.go
          
         **For linux **
         $ go build -o server main.go
         
  * #2.2 Compile client**
       $cd client/
        ** For Windows **
       $ go build -o client.exe main.go
        ** For Linux **
        $ go build -o client main.go
 # 3 Start Server
       $cd server
        ** For Windows **
        $./server.exe 
         server listening at [::]:9000           
        **For linux **
         $./server
          server listening at [::]:9000   
          
          
 # 4 Start Client
  ** $ cd client/**
   
** 4.1 Create Graph**
**$./client.exe create -graph=1-2,2-1,3-1,4-5,5-6,6-1**
Client started ....create
create -graph=1-2,2-1,3-1,4-5,5-6,6-1

**CreateGraph Response : GraphId:  8a052cc7-8ddb-4055-b294-c50ac6ea4b59**

**4.2 get shortest path**
./client.exe shortPath -id=8a052cc7-8ddb-4055-b294-c50ac6ea4b59 -src=3 -dest=1
Client started ....shortPath

**ShortestPath Response :  path = [3 1], dist = 0**
                                                                                                                    
**4.3 Delete Graph**
./client.exe delete  -id=8a052cc7-8ddb-4055-b294-c50ac6ea4b59
Client started ....delete

**DeleteGraph Response :  graphId[8a052cc7-8ddb-4055-b294-c50ac6ea4b59] Successfuly Deleted !!**



**# 5 Go Test for testing functional and performance test**
$ cd server
$go test -bench .

goos: windows
goarch: amd64
pkg: github.com/mahendradeore/goexcercise/network-grpc/server
BenchmarkCreateGraph-8             21182             51186 ns/op
PASS
ok      github.com/mahendradeore/goexcercise/network-grpc/server        2.009s








