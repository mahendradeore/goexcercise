package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	//"graph/logger"
  nw "github.com/mahendradeore/goexcercise/network-grpc/nw"
	"log"
	"net"
	"os"
	"testing"
  "fmt"
)

const (
	bufSize = 1024 * 1024
)


func TestMain(m *testing.M) {
	//setup()
	retCode := m.Run()
	os.Exit(retCode)
}

var listener *bufconn.Listener

func init() {
	listener = bufconn.Listen(bufSize)
	server := grpc.NewServer()
	nw.RegisterGraphServiceServer(server, &graphserver{})
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}
func bufDialer(context.Context, string) (net.Conn, error) {
	return listener.Dial()
}

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "mockServer", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial mockServer: %v", err)
	}
	defer conn.Close()

	client := nw.NewGraphServiceClient(conn)
	resp, err := client.CreateGraph(ctx, &nw.CreateRequest{Edges:populateGraph()})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.GraphId)
	t.Log("GraphId: ", resp.GraphId)

	// Find Shortest Path
	sp, err := client.ShortestPath(ctx, &nw.ShortestPathRequest{GraphId:resp.GraphId, Source: 2, Destination: 1})
	assert.Nil(t, err)
	assert.NotNil(t, sp)
	assert.NotEmpty(t, sp.ShortestPath)
	t.Log("Shortest Path: ", sp.ShortestPath)
	assert.Equal(t,  "path = [2 1], dist = 0", sp.ShortestPath)

	// Invalid graph id - shortest path
	sp, err = client.ShortestPath(ctx, &nw.ShortestPathRequest{GraphId:"a123", Source: 4, Destination: 1})
  assert.Equal(t,  "Error! graphId[a123] not found", sp.ShortestPath)
	t.Log(sp.ShortestPath)

	// Delete graph
	dr, err := client.DeleteGraph(ctx, &nw.DeleteRequest{GraphId:resp.GraphId})
  msg := fmt.Sprintf("graphId[%v] Successfuly Deleted !!",resp.GraphId)
  assert.Equal(t,msg, dr.Message)
	assert.NotEmpty(t, dr.Message)

	// Invalid graph id - delete graph
	dr, err = client.DeleteGraph(ctx, &nw.DeleteRequest{GraphId:"a123"})
  assert.Equal(t,  "Error! graphId[a123] not found", dr.Message)
	t.Log(dr.Message)
}

func BenchmarkCreateGraph(b *testing.B) {
  ctx := context.Background()
  conn, err := grpc.DialContext(ctx, "mockServer", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
  if err != nil {
    log.Fatalf("Failed to dial mockServer: %v", err)
  }
  defer conn.Close()
  client := nw.NewGraphServiceClient(conn)
    for i := 0; i < b.N; i++ {
         client.CreateGraph(ctx, &nw.CreateRequest{Edges:populateGraph()})
    }
}
func populateGraph() []*nw.Edge {
	edges := make([]*nw.Edge, 0, 20)
	edges= append(edges, &nw.Edge{Source: 1, Dest: 2})
	edges= append(edges, &nw.Edge{Source: 2, Dest: 3})
	edges= append(edges, &nw.Edge{Source: 2, Dest: 4})
	return edges
}
