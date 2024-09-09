package main

import (
	"net/http"

	"won/sim-lab/go-engine/internal/sim"

	"github.com/gin-gonic/gin"
)

var starter = []map[string]string{
	{"Name": "Simmy", "Motto": "Make it hot. Make it go."},
	{"Name": "Gloria", "Motto": "Neutrons are my thing."},
	{"Name": "Power Pete", "Motto": "Meeting your energy demands, day by day."},
}

var simCache = make(map[string]*sim.Simulation)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")

	router.StaticFile("/favicon.ico", "./web/assets/favicon.ico")

	for _, s := range starter {
		simulation := sim.NewSimulation(s["Name"], s["Motto"])
		simulation.AddComponent(sim.NewBoiler("Billy Boiler"))
		simulation.AddComponent(sim.NewTurbine("Tilly Turner"))
		simCache[simulation.ID()] = simulation
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.tmpl", gin.H{
			"title":    "WoN Simulator",
			"template": "index",
		})
	})

	router.GET("/inspector", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.tmpl", gin.H{
			"title":    "Component Inspector",
			"template": "inspector",
		})
	})

	router.GET("/api/sims/:id/components/:name", func(c *gin.Context) {
		simulationID := c.Param("id")
		componentName := c.Param("name")

		simulation, exists := simCache[simulationID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Simulation not found"})
			return
		}

		var componentInfo map[string]interface{}

		switch componentName {
		case "boiler":
			if boiler := simulation.FindBoiler(); boiler != nil {
				componentInfo = boiler.Status()
			}
		case "turbine":
			if turbine := simulation.FindTurbine(); turbine != nil {
				componentInfo = turbine.Status()
			}
		default:
			c.JSON(http.StatusNotFound, gin.H{"error": "Component not found"})
			return
		}

		if componentInfo != nil {
			c.JSON(http.StatusOK, componentInfo)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Component not found"})
		}
	})

	// router.GET("/analytics", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "analytics.tmpl", gin.H{
	// 		"title": "Analytics Pie",
	// 	})
	// })

	router.GET("/api/sims", getSimInfos)
	// router.POST("/sims", postSimInfo)
	// router.GET("/sims/:id", getSimInfoByID)
	// router.PUT("/sims/:id", updateSimInfo)
	// router.DELETE("/sims/:id", deleteSimInfo)

	router.Run(":8080")
}

func getSimInfos(c *gin.Context) {
	var simInfos []sim.SimInfo

	for _, simulation := range simCache {
		simInfos = append(simInfos, simulation.Info())
	}

	c.JSON(http.StatusOK, simInfos)
}

// func postSimInfo(c *gin.Context) {
// 	var newSimInfo simInfo

// 	if err := c.BindJSON(&newSimInfo); err != nil {
// 		return
// 	}

// 	newSimInfo.SpawnedAt = time.Now()

// 	simInfos = append(simInfos, newSimInfo)
// 	c.IndentedJSON(http.StatusCreated, newSimInfo)
// }

// func getSimInfoByID(c *gin.Context) {
// 	id := c.Param("id")
// 	for _, a := range simInfos {
// 		if a.ID == id {
// 			c.IndentedJSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sim not found"})
// }

// func updateSimInfo(c *gin.Context) {
// 	id := c.Param("id")
// 	var updatedSim simInfo

// 	// Bind the JSON body to the updatedSim variable
// 	if err := c.BindJSON(&updatedSim); err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
// 		return
// 	}

// 	for i, sim := range simInfos {
// 		if sim.ID == id {
// 			// Update the sim info, preserving the ID and SpawnedAt
// 			updatedSim.ID = sim.ID
// 			updatedSim.SpawnedAt = sim.SpawnedAt
// 			simInfos[i] = updatedSim
// 			c.IndentedJSON(http.StatusOK, updatedSim)
// 			return
// 		}
// 	}

// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sim not found"})
// }

// func deleteSimInfo(c *gin.Context) {
// 	id := c.Param("id")
// 	for index, a := range simInfos {
// 		if a.ID == id {
// 			simInfos = append(simInfos[:index], simInfos[index+1:]...)
// 			c.IndentedJSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "sim not found"})
// }
