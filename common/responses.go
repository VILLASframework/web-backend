package common

type FileResponse struct {
	Name              string `json:"name"`
	ID                uint   `json:"id"`
	Type              string `json:"type"`
	Size              uint   `json:"size"`
	ImageWidth        uint   `json:"imageHeight"`
	ImageHeight       uint   `json:"imageWidth"`
	Date              string `json:"date"`
	WidgetID          uint   `json:"widgetID"`
	SimulationModelID uint   `json:"simulationModelID"`
}

// Response messages

type ResponseMsg struct {
	Message string `json:"message"`
}

type ResponseMsgFiles struct {
	Files []FileResponse `json:"files"`
}

type ResponseMsgFile struct {
	File FileResponse `json:"file"`
}
