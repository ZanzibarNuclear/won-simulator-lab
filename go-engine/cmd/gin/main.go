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
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"title":    "WoN Simulator",
			"template": "index",
		})
	})

	router.GET("/operator", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"title":    "Simulator Operator",
			"template": "operator",
		})
	})

	router.GET("/inspector", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"title":    "Component Inspector",
			"template": "inspector",
		})
	})

	// router.GET("/analytics", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "analytics.tmpl", gin.H{
	// 		"title": "Analytics Pie",
	// 	})
	// })

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

	router.POST("/api/sims", createSimulation)
	router.GET("/api/sims", getSimInfos)
	router.GET("/api/sims/:id", getSimInfo)
	router.GET("/api/sims/:id/status", getSimStatus)
	// router.PUT("/sims/:id", updateSimInfo)
	// router.DELETE("/sims/:id", deleteSimInfo)
	router.GET("/api/sims/:id/components", getComponents)

	router.Run(":8080")
}

func getSimInfos(c *gin.Context) {
	var simInfos []sim.SimInfo

	for _, simulation := range simCache {
		simInfos = append(simInfos, simulation.Info())
	}

	c.JSON(http.StatusOK, simInfos)
}

func getSimInfo(c *gin.Context) {
	simulationID := c.Param("id")
	simulation, exists := simCache[simulationID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Simulation not found"})
		return
	}

	c.JSON(http.StatusOK, simulation.Info())
}

func getComponents(c *gin.Context) {
	simulationID := c.Param("id")

	simulation, exists := simCache[simulationID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Simulation not found"})
		return
	}

	c.JSON(http.StatusOK, simulation.Components())
}

func createSimulation(c *gin.Context) {
	var simData struct {
		Name  string `json:"name" binding:"required"`
		Motto string `json:"motto" binding:"required"`
	}

	if err := c.ShouldBindJSON(&simData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newSim := sim.NewSimulation(simData.Name, simData.Motto)
	newSim.AddComponent(sim.NewBoiler("Boilerator 37"))
	newSim.AddComponent(sim.NewTurbine("Turbinator 42"))

	simCache[newSim.ID()] = newSim

	c.JSON(http.StatusCreated, newSim.Info())
}

func getSimStatus(c *gin.Context) {
	simulationID := c.Param("id")

	simulation, exists := simCache[simulationID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Simulation not found"})
		return
	}

	c.JSON(http.StatusOK, simulation.Status())
}

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
