package main

import (
	"net/http"

	"won/sim-lab/go-engine/internal/sim"

	"github.com/gin-gonic/gin"
)

// type simmer struct {
// 	about sim.SimInfo
// 	sim   *sim.Simulation
// }

// var simulators = []sim.SimInfo{
// 	{ID: "1", Name: "Simmy", Motto: "Make it hot. Make it go.", SpawnedAt: time.Now()},
// 	{ID: "2", Name: "Gloria", Motto: "Neutrons are my thing.", SpawnedAt: time.Now()},
// 	{ID: "3", Name: "Power Pete", Motto: "Meeting your energy demands, day by day.", SpawnedAt: time.Now()},
// }

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")

	router.StaticFile("/favicon.ico", "./web/assets/favicon.ico")

	
	simulation := sim.NewSimulation()
	simulation.AddComponent(sim.NewBoiler("Main Boiler"))
	simulation.AddComponent(sim.NewTurbine("Steam Turbine"))

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

	router.GET("/api/component/:name", func(c *gin.Context) {
		componentName := c.Param("name")
		// Send back the component name as received
		// c.JSON(http.StatusOK, gin.H{"name": componentName})

		var componentInfo map[string]interface{}

		switch componentName {
		case "boiler":
			// TODO: move finders to simulation so that it can supply the component slice
			if boiler := sim.FindBoiler(simulation.Components()); boiler != nil {
				componentInfo = boiler.Status()
			}
		case "turbine":
			if turbine := sim.FindTurbine(simulation.Components()); turbine != nil {
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

	// router.GET("/sims", getSimInfos)
	// router.POST("/sims", postSimInfo)
	// router.GET("/sims/:id", getSimInfoByID)
	// router.PUT("/sims/:id", updateSimInfo)
	// router.DELETE("/sims/:id", deleteSimInfo)

	router.Run(":8080")
}

// func getSimInfos(c *gin.Context) {
// 	c.IndentedJSON(http.StatusOK, simInfos)
// }

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
