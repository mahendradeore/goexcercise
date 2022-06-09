package main
import (
	"flag"
	"fmt"
	"os"
  "context"
  nw "github.com/mahendradeore/goexcercise/network-grpc/nw"
  "strconv"
  "strings"
  "time"
  "google.golang.org/grpc"
)

func CreateGrpcClient() *grpc.ClientConn {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		fmt.Errorf("Failed to start gRPC connection: %v", err)
	}
	return conn
}
func main() {

	var graph string
	var graphId string
	var src string
	var dest string
  //var msg string
  fmt.Printf("Client started ....")
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createCmd.StringVar(&graph, "graph", "", "graph edges example: 1-2,2-3 or 'null' for null graph")

	findCmd := flag.NewFlagSet("shortPath", flag.ExitOnError)
	findCmd.StringVar(&graphId, "id", "", "graph id")
	findCmd.StringVar(&src, "src", "", "source node")
	findCmd.StringVar(&dest, "dest", "", "destination node")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteCmd.StringVar(&graphId, "id", "", "graph id")

	if len(os.Args) < 2 {
		fmt.Println("expected operations 'create' or 'find' or 'delete'")
		os.Exit(1)
	}

	cmd := os.Args[1]
  fmt.Println(cmd)
	switch cmd {
	case "create":
    fmt.Println(os.Args[1],os.Args[2])
		createCmd.Parse(os.Args[2:])
		if createCmd.Parsed() {
      //fmt.Println("graph",graph)
			if graph == ""{
				createCmd.PrintDefaults()
        fmt.Println("graphId is missing",graph)
				os.Exit(1)
			}
      createGraph(graph)
		}
	case "shortPath":
		findCmd.Parse(os.Args[2:])
		if findCmd.Parsed() {
			if graphId == ""|| src == "" || dest == ""{
        fmt.Println("graphId is missing",graph)
				findCmd.PrintDefaults()
				os.Exit(1)
			}
			findPath(graphId, src, dest)
		}
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if deleteCmd.Parsed() {
			if graphId == ""{
				deleteCmd.PrintDefaults()
				os.Exit(1)
			}
			deleteGraph(graphId)
		}
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
func createGraph(graph string){
  //fmt.Printf("Starting createGraph\n")
	edges := make([]*nw.Edge, 0, 0)
	if graph != "null" {
		edgePairs := strings.Split(graph, ",")
		for i := 0; i < len(edgePairs); i++ {
			if !strings.Contains(edgePairs[i], "-") {
				fmt.Println("Invalid edge format", edgePairs[i])
				os.Exit(1)
			}
			edge := strings.Split(edgePairs[i], "-")
			if len(edge) != 2 {
				fmt.Println("Invalid edge format", edgePairs[i])
				os.Exit(1)
			}
      intVar1, _ := strconv.Atoi(edge[0])

      intVar2, _:= strconv.Atoi(edge[1])
			edges = append(edges, &nw.Edge{Source:int32(intVar1), Dest: int32(intVar2)})
		}
	}

//	request := &nw.CreateRequest{Edges:edges}
  //fmt.Printf("\nrequest:%v",edges)



	conn := CreateGrpcClient()
	defer conn.Close()
  //c := nw.NewChatServiceClient(conn)
	client := nw.NewGraphServiceClient(conn)
//context.Background()
	response, err := client.CreateGraph(context.Background(), &nw.CreateRequest{Edges:edges})
	if err != nil {
		fmt.Errorf("\nCreateGraph Failed [%v]", err)
    return
	}else {
    fmt.Println("\nCreateGraph Response : GraphId: ", response.GraphId)
  }


}


func findPath(graphId, src, dest string){
  intVar1, _ := strconv.Atoi(src)
  intVar2, _:= strconv.Atoi(dest)
	request := &nw.ShortestPathRequest{GraphId:graphId, Source: int32(intVar1), Destination: int32(intVar2)}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

  conn := CreateGrpcClient()
  defer conn.Close()

	client := nw.NewGraphServiceClient(conn)
	response, err := client.ShortestPath(ctx, request)
  //fmt.Errorf("\nShortestPath:Error %v", err)
  //fmt.Println("\nShortest Path:", response.ShortestPath)
	if err != nil {
		fmt.Errorf("\nShortestPath  : Error %v", err)
	}else {
		fmt.Println("\nShortestPath Response : ", response.ShortestPath)
	}
}

func deleteGraph(graphId string){

	request := &nw.DeleteRequest{GraphId:graphId}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

  conn := CreateGrpcClient()
	defer conn.Close()

	client := nw.NewGraphServiceClient(conn)
	response, err := client.DeleteGraph(ctx, request)
	if err != nil {
		fmt.Errorf("Error %v", err)
	}else {
		fmt.Println("\nDeleteGraph Response : ", response.Message)
	}
}
