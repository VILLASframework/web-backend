package model

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

func ModelsRegister(r *gin.RouterGroup) {
	r.GET("/:SimulationID/models/", modelsReadEp)
	r.POST("/:SimulationID/models/", modelRegistrationEp)

	r.PUT("/:SimulationID/models/:ModelID", modelUpdateEp)
	r.GET("/:SimulationID/models/:ModelID", modelReadEp)
	r.DELETE("/:SimulationID/models/:ModelID", modelDeleteEp)

	// Files
	r.POST ("/:SimulationID/models/:ModelID/file", modelRegisterFileEp) // NEW in API
	r.GET("/:SimulationID/models/:ModelID/file", modelReadFileEp) // NEW in API
	r.PUT("/:SimulationID/models/:ModelID/file", modelUpdateFileEp) // NEW in API
	r.DELETE("/:SimulationID/models/:ModelID/file", modelDeleteFileEp) // NEW in API

	// Simulator
	r.PUT("/:SimulationID/models/:ModelID/simulator", modelUpdateSimulatorEp) // NEW in API
	r.GET("/:SimulationID/models/:ModelID/simulator", modelReadSimulatorEp) // NEW in API

	// Input and Output Samples
	r.POST("/:SimulationID/models/:ModelID/Samples/:Direction", modelRegisterSamplesEp) // NEW in API
	r.GET("/:SimulationID/models/:ModelID/Samples/:Direction", modelReadSamplesEp) // NEW in API
	r.PUT("/:SimulationID/models/:ModelID/Samples/:Direction", modelUpdateSamplesEp) // NEW in API
	r.DELETE("/:SimulationID/models/:ModelID/Samples/:Direction", modelDeleteSamplesEp) // NEW in API
}

func modelsReadEp(c *gin.Context) {
	allModels, _, _ := FindAllModels()
	serializer := ModelsSerializerNoAssoc{c, allModels}
	c.JSON(http.StatusOK, gin.H{
		"models": serializer.Response(),
	})
}

func modelRegistrationEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func modelUpdateEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func modelReadEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func modelDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}


func modelRegisterFileEp(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Save file locally and register file in DB, HTTP response is set by this method
	file.RegisterFile(c,-1, modelID, simulationID)

}

func modelReadFileEp(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Read file from disk and return in HTTP response, no change to DB
	file.ReadFile(c, -1, modelID, simulationID)
}

func modelUpdateFileEp(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Update file locally and update file entry in DB, HTTP response is set by this method
	file.UpdateFile(c,-1, modelID, simulationID)
}

func modelDeleteFileEp(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Delete file from disk and remove entry from DB, HTTP response is set by this method
	file.DeleteFile(c, -1, modelID, simulationID)


}


func GetModelID(c *gin.Context) (int, error) {

	modelID, err := strconv.Atoi(c.Param("ModelID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulation model ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, err
	} else {
		return modelID, err

	}
}

func getRequestParams(c *gin.Context) (int, int, error){
	simulationID, err := simulation.GetSimulationID(c)
	if err != nil{
		return -1, -1, err
	}

	modelID, err := GetModelID(c)
	if err != nil{
		return -1, -1, err
	}

	return simulationID, modelID, err
}