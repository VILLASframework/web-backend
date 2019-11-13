package widget

import (
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsf/jsondiff"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewWidget struct {
	Name             string         `form:"name" validate:"required"`
	Type             string         `form:"type" validate:"required"`
	Width            uint           `form:"width" validate:"required"`
	Height           uint           `form:"height" validate:"required"`
	MinWidth         uint           `form:"minWidth" validate:"omitempty"`
	MinHeight        uint           `form:"minHeight" validate:"omitempty"`
	X                int            `form:"x" validate:"omitempty"`
	Y                int            `form:"y" validate:"omitempty"`
	Z                int            `form:"z" validate:"omitempty"`
	DashboardID      uint           `form:"dashboardID" validate:"required"`
	IsLocked         bool           `form:"isLocked" validate:"omitempty"`
	CustomProperties postgres.Jsonb `form:"customProperties" validate:"omitempty"`
}

type validUpdatedWidget struct {
	Name             string         `form:"name" validate:"omitempty"`
	Type             string         `form:"type" validate:"omitempty"`
	Width            uint           `form:"width" validate:"omitempty"`
	Height           uint           `form:"height" validate:"omitempty"`
	MinWidth         uint           `form:"minWidth" validate:"omitempty"`
	MinHeight        uint           `form:"minHeight" validate:"omitempty"`
	X                int            `form:"x" validate:"omitempty"`
	Y                int            `form:"y" validate:"omitempty"`
	Z                int            `form:"z" validate:"omitempty"`
	IsLocked         bool           `form:"isLocked" validate:"omitempty"`
	CustomProperties postgres.Jsonb `form:"customProperties" validate:"omitempty"`
}

type addWidgetRequest struct {
	Widget validNewWidget `json:"widget"`
}

type updateWidgetRequest struct {
	Widget validUpdatedWidget `json:"widget"`
}

func (r *addWidgetRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *validUpdatedWidget) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *addWidgetRequest) createWidget() Widget {
	var s Widget

	s.Name = r.Widget.Name
	s.Type = r.Widget.Type
	s.Width = r.Widget.Width
	s.Height = r.Widget.Height
	s.MinWidth = r.Widget.MinWidth
	s.MinHeight = r.Widget.MinHeight
	s.X = r.Widget.X
	s.Y = r.Widget.Y
	s.Z = r.Widget.Z
	s.IsLocked = r.Widget.IsLocked
	s.CustomProperties = r.Widget.CustomProperties
	s.DashboardID = r.Widget.DashboardID
	return s
}

func (r *updateWidgetRequest) updatedWidget(oldWidget Widget) Widget {
	// Use the old Widget as a basis for the updated Widget `s`
	s := oldWidget

	if r.Widget.Name != "" {
		s.Name = r.Widget.Name
	}

	s.Type = r.Widget.Type
	s.Width = r.Widget.Width
	s.Height = r.Widget.Height
	s.MinWidth = r.Widget.MinWidth
	s.MinHeight = r.Widget.MinHeight
	s.X = r.Widget.X
	s.Y = r.Widget.Y
	s.Z = r.Widget.Z
	s.IsLocked = r.Widget.IsLocked

	// only update custom props if not empty
	var emptyJson postgres.Jsonb
	// Serialize empty json and params
	emptyJson_ser, _ := json.Marshal(emptyJson)
	customprops_ser, _ := json.Marshal(r.Widget.CustomProperties)
	opts := jsondiff.DefaultConsoleOptions()
	diff, _ := jsondiff.Compare(emptyJson_ser, customprops_ser, &opts)
	if diff.String() != "FullMatch" {
		s.CustomProperties = r.Widget.CustomProperties
	}

	return s
}
