package result

import (
	"github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewResult struct {
	Description     string           `form:"Description" validate:"omitempty"`
	ResultFileIDs   []int64          `form:"ResultFileIDs" validate:"omitempty"`
	ConfigSnapshots []postgres.Jsonb `form:"ConfigSnapshots" validate:"required"`
	ScenarioID      uint             `form:"ScenarioID" validate:"required"`
}

type validUpdatedResult struct {
	Description     string           `form:"Description" validate:"omitempty" json:"description"`
	ResultFileIDs   []int64          `form:"ResultFileIDs" validate:"omitempty" json:"resultFileIDs"`
	ConfigSnapshots []postgres.Jsonb `form:"ConfigSnapshots" validate:"omitempty" json:"configSnapshots"`
}

type addResultRequest struct {
	Result validNewResult `json:"result"`
}

type updateResultRequest struct {
	Result validUpdatedResult `json:"result"`
}

func (r *addResultRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *validUpdatedResult) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *addResultRequest) createResult() Result {
	var s Result

	s.Description = r.Result.Description
	s.ConfigSnapshots = r.Result.ConfigSnapshots
	s.ResultFileIDs = r.Result.ResultFileIDs
	s.ScenarioID = r.Result.ScenarioID

	return s
}

func (r *updateResultRequest) updatedResult(oldResult Result) Result {
	// Use the old Result as a basis for the updated Result `s`
	s := oldResult

	s.Result.Description = r.Result.Description
	s.ConfigSnapshots = r.Result.ConfigSnapshots
	s.ResultFileIDs = r.Result.ResultFileIDs

	return s
}
