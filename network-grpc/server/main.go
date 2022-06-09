package main

import(
  "fmt"
  "net"
  "context"
  "google.golang.org/grpc"
  nw "github.com/mahendradeore/goexcercise/network-grpc/nw"
  "github.com/yourbasic/graph"
  "github.com/google/uuid"
  "errors"
  "sync"
)

type server struct {
	nw.UnimplementedChatServiceServer
}

type graphserver struct {
	nw.UnimplementedGraphServiceServer
  mu         sync.Mutex
}

var(
  store = make(map[string]*graph.Mutable)
)
// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *nw.Message) (*nw.Message, error) {
	//fmt.Printf("\nSayHello:Received: %v\n", in.Body)
	return &nw.Message{Body: "Hello " + in.Body}, nil
}

func transform(edge []*nw.Edge)([][]int32, int) {
	var strGraph [][]int32
  m := make(map[int32]bool)
	for i:=0; i < len(edge); i++ {
    if !m[edge[i].Source]{
      m[edge[i].Source]= true
    }
    if !m[edge[i].Dest]{
      m[edge[i].Dest]= true
    }
		strGraph = append(strGraph, []int32{edge[i].Source, edge[i].Dest})
	}
  //fmt.Printf("\n transform = %v",strGraph)
	return strGraph, len(m)
}
// SayHello implements helloworld.GreeterServer
func (s *graphserver) CreateGraph(ctx context.Context, in *nw.CreateRequest) (*nw.CreateResponse, error) {
	//fmt.Printf("\n CreateGraph:Received: %v", in.Edges)
  edges,count := transform(in.Edges)
  //  fmt.Printf("\n count: %v, edges = [%v]", count,edges)
  s.mu.Lock()
  g := graph.New(count*2)
  graphId := uuid.New().String()
  for  i := 0; i < len(edges); i++ {
      //fmt.Printf("\n edges[%v]) = %v, edges[i][0] = %v, edges[i][0] = %v", i,len(edges[i]),int(edges[i][0]),int(edges[i][1]))
  		if len(edges[i])== 2 {
        //fmt.Printf("\nAdding in graph")
       g.AddBoth(int(edges[i][0]), int(edges[i][1]))
  		}else{
  			return nil, errors.New(fmt.Sprintf("EdgeCountError = %d", i))
  		}
  }
  store[graphId] = g
  s.mu.Unlock()
	return &nw.CreateResponse{GraphId: graphId}, nil
}

func (s *graphserver)ShortestPath(ctx context.Context, in *nw.ShortestPathRequest) (*nw.ShortestPathResponse, error){
  //fmt.Printf("\n ShortestPath:Received id: %v, src = %v, dest = %v", in.GraphId,in.Source,in.Destination)
  s.mu.Lock()
  if mut,found := store[in.GraphId];found {
    path, dist := graph.ShortestPath(mut,int(in.Source), int(in.Destination))
      s.mu.Unlock()
	  return &nw.ShortestPathResponse{ShortestPath: fmt.Sprintf("path = %v, dist = %v", path, dist), Err : ""},nil
  } else {
  //  fmt.Println("\n ShortestPath:Graph not found")
  //  var err error
    s.mu.Unlock()
   return &nw.ShortestPathResponse{ShortestPath: fmt.Sprintf("Error! graphId[%v] not found", in.GraphId), Err : ""},nil
  }
}

func (s *graphserver)DeleteGraph(ctx context.Context, in *nw.DeleteRequest) (*nw.DeleteResponse, error){
  //fmt.Printf("\n DeleteGraph:Received id: %v", in.GraphId)
  s.mu.Lock()
  if _,found := store[in.GraphId];found {

    delete(store,in.GraphId)
    s.mu.Unlock()
    return &nw.DeleteResponse{Message: fmt.Sprintf("graphId[%v] Successfuly Deleted !!", in.GraphId,), Err : ""},nil
  } else {
    //  fmt.Println("\n DeleteGraph:Graph not found")
    s.mu.Unlock()
    return &nw.DeleteResponse{Message: fmt.Sprintf("Error! graphId[%v] not found", in.GraphId,), Err : ""},nil
  }
}

func main(){
  lis,err := net.Listen("tcp",":9000")
  if err!= nil {
    panic(err)
  }

  s := grpc.NewServer()
  nw.RegisterGraphServiceServer(s, &graphserver{})
	fmt.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}

}
