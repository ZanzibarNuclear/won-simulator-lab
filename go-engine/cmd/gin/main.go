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

	// bootstrap starter simulations, something to work with
	for _, s := range starter {
		simulation := spawnSimulation(s["Name"], s["Motto"])
		simCache[simulation.ID()] = simulation
	}

	// routes

	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")

	// static assets

	router.Static("/static", "./web/assets")
	router.StaticFile("/favicon.ico", "./web/assets/favicon.ico")

	// page routes

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"active":   "home",
			"title":    "WoN Simulator",
			"template": "index",
		})
	})

	router.GET("/operator", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"active":   "operator",
			"title":    "Simulator Operator",
			"template": "operator",
		})
	})

	router.GET("/inspector", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"active":   "inspector",
			"title":    "Component Inspector",
			"template": "inspector",
		})
	})

	router.GET("/analysis", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", gin.H{
			"active":   "analysis",
			"title":    "Simulation Analysis",
			"template": "analysis",
		})
	})

	// API routes

	router.POST("/api/sims", createSimulation)
	router.GET("/api/sims", getSimInfos)
	router.GET("/api/sims/:id", getSimInfo)
	router.GET("/api/sims/:id/status", getSimStatus)
	router.GET("/api/sims/:id/components", getComponents)
	router.GET("/api/sims/:id/components/:name", getComponentStatus)
	router.PUT("/api/sims/:id/advance", advanceSim)
	router.PUT("/api/sims/:id/interrupt", interruptSim)

	router.Run(":8080")
}

func spawnSimulation(name, motto string) *sim.Simulation {
	simulation := sim.NewSimulation(name, motto)
	simulation.AddComponent(sim.NewPrimaryLoop("Primary Loop"))
	simulation.AddComponent(sim.NewSecondaryLoop("Secondary Loop"))
	simulation.AddComponent(sim.NewReactorCore("Reactor Core"))
	simulation.AddComponent(sim.NewPressurizer("Pressurizer"))
	simulation.AddComponent(sim.NewSteamGenerator("Steam Generator"))
	simulation.AddComponent(sim.NewSteamTurbine("Steam Turbine"))
	simulation.AddComponent(sim.NewCondenser("Condenser"))
	simulation.AddComponent(sim.NewGenerator("Power Generator"))
	return simulation
}

func advanceSim(c *gin.Context) {
	simulationID := c.Param("id")
	simulation, exists := simCache[simulationID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Simulation not found"})
		return
	}

	var advanceData struct {
		Steps int `json:"steps"`
	}

	if err := c.ShouldBindJSON(&advanceData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if advanceData.Steps == 0 {
		advanceData.Steps = 1
	}

	simulation.Advance(advanceData.Steps)

	c.JSON(http.StatusOK, simulation.Status())
}

func interruptSim(c *gin.Context) {
	simulationID := c.Param("id")
	simulation, exists := simCache[simulationID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Simulation not found"})
		return
	}

	simulation.Stop()

	c.JSON(http.StatusOK, simulation.Status())
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

	newSim := spawnSimulation(simData.Name, simData.Motto)
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

func getComponentStatus(c *gin.Context) {
	simulationID := c.Param("id")
	componentName := c.Param("name")

	simulation, exists := simCache[simulationID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Simulation not found"})
		return
	}

	var componentInfo map[string]interface{}

	switch componentName {
	case "PrimaryLoop":
		componentInfo = simulation.FindPrimaryLoop().Status()
	case "SecondaryLoop":
		componentInfo = simulation.FindSecondaryLoop().Status()
	case "ReactorCore":
		componentInfo = simulation.FindReactorCore().Status()
	case "Pressurizer":
		componentInfo = simulation.FindPressurizer().Status()
	case "SteamGenerator":
		componentInfo = simulation.FindSteamGenerator().Status()
	case "SteamTurbine":
		componentInfo = simulation.FindSteamTurbine().Status()
	case "Condenser":
		componentInfo = simulation.FindCondenser().Status()
	case "Generator":
		componentInfo = simulation.FindGenerator().Status()
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Component not found"})
		return
	}

	if componentInfo != nil {
		c.JSON(http.StatusOK, componentInfo)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Component not found"})
	}
}
