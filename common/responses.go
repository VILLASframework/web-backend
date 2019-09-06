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

type SignalResponse struct {
	Name              string `json:"name"`
	Unit              string `json:"unit"`
	Index             uint   `json:"index"`
	Direction         string `json:"direction"`
	SimulationModelID uint   `json:"simulationModelID"`
}

// Response messages

type ResponseMsg struct {
	Message string `json:"message"`
}

type ResponseMsgSignals struct {
	Signals []SignalResponse `json:"signals"`
}

type ResponseMsgSignal struct {
	Signal SignalResponse `json:"signal"`
}

type ResponseMsgFiles struct {
	Files []FileResponse `json:"files"`
}

type ResponseMsgFile struct {
	File FileResponse `json:"file"`
}
