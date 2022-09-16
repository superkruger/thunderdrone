package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// node represents data about a lightning node.
type node struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Port    string `json:"port"`
}

// channel represents data about a lightning channel between two nodes.
type channel struct {
	ID       string `json:"id"`
	Capacity int64  `json:"capacity"`
	Node1ID  string `json:"node-1-id"`
	Node2ID  string `json:"node-2-id"`
}

// nodes slice to seed node data.
var nodes = []node{
	{ID: "1", Address: "10.0.0.1", Port: "9735"},
	{ID: "2", Address: "10.0.0.2", Port: "9735"},
	{ID: "3", Address: "10.0.0.3", Port: "9735"},
}

// channels slice to seed channel data.
var channels = []channel{
	{ID: "1", Capacity: 5000000, Node1ID: "1", Node2ID: "2"},
	{ID: "2", Capacity: 2000000, Node1ID: "2", Node2ID: "3"},
}

func main() {
	fmt.Println("Thunderdrone starting...")
	router := gin.Default()
	applyCors(router)

	router.GET("/nodes", getNodes)
	router.GET("/channels", getChannels)

	fmt.Println("Listening on port 8080")

	router.Run(":8080")
}

func applyCors(r *gin.Engine) {
	corsConfig := cors.DefaultConfig()
	//hot reload CORS
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))
}

func getNodes(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, nodes)
}
func getChannels(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, channels)
}
