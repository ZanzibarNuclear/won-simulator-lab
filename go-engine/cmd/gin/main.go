package main

import (
	"time"
	"github.com/gin-gonic/gin"
	"net/http"
)

type simInfo struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Motto string `json:"motto"`
	SpawnedAt time.Time `json:"spawned_at"`
}

var simInfos = []simInfo{
	{ID: "1", Name: "Simmy", Motto: "Make it hot. Make it go.", SpawnedAt: time.Now()},
	{ID: "2", Name: "Gloria", Motto: "Neutrons are my thing.", SpawnedAt: time.Now()},
	{ID: "3", Name: "Power Pete", Motto: "Meeting your energy demands, day by day.", SpawnedAt: time.Now()},
}

func main() {
	router := gin.Default()
	router.GET("/sims", getSimInfos)

	router.Run(":8080")
}

func getSimInfos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, simInfos)
}
