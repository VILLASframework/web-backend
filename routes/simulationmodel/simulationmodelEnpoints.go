package simulationmodel

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

func SimulationModelsRegister(r *gin.RouterGroup) {
	r.GET("/:SimulationID/models/", simulationmodelsReadEp)
	r.POST("/:SimulationID/models/", simulationmodelRegistrationEp)

	r.PUT("/:SimulationID/models/:SimulationModelID", simulationmodelUpdateEp)
	r.GET("/:SimulationID/models/:SimulationModelID", simulationmodelReadEp)
	r.DELETE("/:SimulationID/models/:SimulationModelID", simulationmodelDeleteEp)

	// Files
	r.POST ("/:SimulationID/models/:SimulationModelID/file", simulationmodelRegisterFileEp) // NEW in API
	r.GET("/:SimulationID/models/:SimulationModelID/file", simulationmodelReadFileEp) // NEW in API
	r.PUT("/:SimulationID/models/:SimulationModelID/file", simulationmodelUpdateFileEp) // NEW in API
	r.DELETE("/:SimulationID/models/:SimulationModelID/file", simulationmodelDeleteFileEp) // NEW in API


}

func simulationmodelsReadEp(c *gin.Context) {
	allSimulationModels, _, _ := FindAllSimulationModels()
	serializer := SimulationModelsSerializerNoAssoc{c, allSimulationModels}
	c.JSON(http.StatusOK, gin.H{
		"simulationmodels": serializer.Response(),
	})
}

func simulationmodelRegistrationEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationmodelUpdateEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationmodelReadEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationmodelDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}


func simulationmodelRegisterFileEp(c *gin.Context) {

	simulationID, simulationmodelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Save file locally and register file in DB, HTTP response is set by this method
	file.RegisterFile(c,-1, simulationmodelID, simulationID)

}

func simulationmodelReadFileEp(c *gin.Context) {

	simulationID, simulationmodelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Read file from disk and return in HTTP response, no change to DB
	file.ReadFile(c, -1, simulationmodelID, simulationID)
}

func simulationmodelUpdateFileEp(c *gin.Context) {

	simulationID, simulationmodelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Update file locally and update file entry in DB, HTTP response is set by this method
	file.UpdateFile(c,-1, simulationmodelID, simulationID)
}

func simulationmodelDeleteFileEp(c *gin.Context) {

	simulationID, simulationmodelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Delete file from disk and remove entry from DB, HTTP response is set by this method
	file.DeleteFile(c, -1, simulationmodelID, simulationID)


}


func GetSimulationmodelID(c *gin.Context) (int, error) {

	simulationmodelID, err := strconv.Atoi(c.Param("SimulationModelID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulation model ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, err
	} else {
		return simulationmodelID, err

	}
}

func getRequestParams(c *gin.Context) (int, int, error){
	simulationID, err := simulation.GetSimulationID(c)
	if err != nil{
		return -1, -1, err
	}

	simulationmodelID, err := GetSimulationmodelID(c)
	if err != nil{
		return -1, -1, err
	}

	return simulationID, simulationmodelID, err
}