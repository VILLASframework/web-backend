package endpoints

import (
	"github.com/gin-gonic/gin"
)

func SimulationsRegister(r *gin.RouterGroup) {

	r.GET("/", simulationReadAllEp)
	r.POST("/", simulationRegistrationEp)
	r.POST("/:SimulationID", simulationCloneEp)
	r.PUT("/:SimulationID", simulationUpdateEp)
	r.GET("/:SimulationID", simulationReadEp)
	r.DELETE("/:SimulationID", simulationDeleteEp)

	// Users
	r.GET("/:SimulationID/users", userReadAllSimEp)
	r.PUT("/:SimulationID/user/:username", userUpdateSimEp)
	r.DELETE("/:SimulationID/user/:username", userDeleteSimEp)

	// Models
	r.GET("/:SimulationID/models/", modelReadAllEp)
	r.POST("/:SimulationID/models/", modelRegistrationEp)
	r.POST("/:SimulationID/models/:ModelID", modelCloneEp)
	r.PUT("/:SimulationID/models/:ModelID", modelUpdateEp)
	r.GET("/:SimulationID/models/:ModelID", modelReadEp)
	r.DELETE("/:SimulationID/models/:ModelID", modelDeleteEp)

	// Simulators
	r.PUT("/:SimulationID/models/:ModelID/simulator", simulatorUpdateModelEp) // NEW in API
	r.GET("/:SimulationID/models/:ModelID/simulator", simulatorReadModelEp) // NEW in API

	// Input and Output Signals
	r.POST("/:SimulationID/models/:ModelID/signals/:Direction", signalRegistrationEp) // NEW in API
	r.GET("/:SimulationID/models/:ModelID/signals/:Direction", signalReadAllEp) // NEW in API
	r.PUT("/:SimulationID/models/:ModelID/signals/:Direction", signalUpdateEp) // NEW in API
	r.DELETE("/:SimulationID/models/:ModelID/signals/:Direction", signalDeleteEp) // NEW in API

	// Visualizations
	r.GET("/:SimulationID/visualizations", visualizationReadAllEp)
	r.POST("/:SimulationID/visualization", visualizationRegistrationEp)
	r.POST("/:SimulationID/visualization/:visualizationID", visualizationCloneEp)
	r.PUT("/:SimulationID/visualization/:visualizationID", visualizationUpdateEp)
	r.GET("/:SimulationID/visualization/:visualizationID", visualizationReadEp)
	r.DELETE("/:SimulationID/visualization/:visualizationID", visualizationDeleteEp)

	// Widgets
	r.GET("/:SimulationID/visualization/:visualizationID/widgets", widgetReadAllEp)
	r.POST("/:SimulationID/visualization/:visualizationID/widget", widgetRegistrationEp)
	r.POST("/:SimulationID/visualization/:visualizationID/widget:widgetID", widgetCloneEp)
	r.PUT("/:SimulationID/visualization/:visualizationID/widget/:widgetID", widgetUpdateEp)
	r.GET("/:SimulationID/visualization/:visualizationID/widget/:widgetID", widgetReadEp)
	r.DELETE("/:SimulationID/visualization/:visualizationID/widget/:widgetID", widgetDeleteEp)

	// Files
	// Files of Models
	r.GET("/:SimulationID/models/:ModelID/files", fileMReadAllEp) // NEW in API
	r.POST ("/:SimulationID/models/:ModelID/file", fileMRegistrationEp) // NEW in API
	//r.POST ("/:SimulationID/models/:ModelID/file", fileMCloneEp) // NEW in API
	r.GET("/:SimulationID/models/:ModelID/file", fileMReadEp) // NEW in API
	r.PUT("/:SimulationID/models/:ModelID/file", fileMUpdateEp) // NEW in API
	r.DELETE("/:SimulationID/models/:ModelID/file", fileMDeleteEp) // NEW in API

	// Files of Widgets
	r.GET("/:SimulationID/visualizations/:VisID/widgets/:WidgetID/files", fileWReadAllEp) // NEW in API
	r.POST ("/:SimulationID/visualizations/:VisID/widgets/:WidgetID/file", fileWRegistrationEp) // NEW in API
	//r.POST ("/:SimulationID/visualizations/:VisID/widgets/:WidgetID/file", fileWCloneEp) // NEW in API
	r.GET("/:SimulationID/visualizations/:VisID/widgets/:WidgetID/file", fileWReadEp) // NEW in API
	r.PUT("/:SimulationID/visualizations/:VisID/widgets/:WidgetID/file", fileWUpdateEp) // NEW in API
	r.DELETE("/:SimulationID/visualizations/:VisID/widgets/:WidgetID/file", fileWDeleteEp) // NEW in API

}

func UsersRegister(r *gin.RouterGroup) {
	r.GET("/", userReadAllEp)
	r.POST("/", userRegistrationEp)
	r.PUT("/:UserID", userUpdateEp)
	r.GET("/:UserID", userReadEp)
	r.DELETE("/:UserID", userDeleteEp)
	//r.GET("/me", userSelfEp) // TODO redirect to users/:UserID
}


func SimulatorsRegister(r *gin.RouterGroup) {
	r.GET("/", simulatorReadAllEp)
	r.POST("/", simulatorRegistrationEp)
	r.PUT("/:SimulatorID", simulatorUpdateEp)
	r.GET("/:SimulatorID", simulatorReadEp)
	r.DELETE("/:SimulatorID", simulatorDeleteEp)
	r.POST("/:SimulatorID", simulatorSendActionEp)
}