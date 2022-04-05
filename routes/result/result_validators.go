/**
* This file is part of VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

package result

import (
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsf/jsondiff"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewResult struct {
	Description     string         `form:"Description" validate:"omitempty"`
	ResultFileIDs   []int64        `form:"ResultFileIDs" validate:"omitempty"`
	ConfigSnapshots postgres.Jsonb `form:"ConfigSnapshots" validate:"required"`
	ScenarioID      uint           `form:"ScenarioID" validate:"required"`
}

type validUpdatedResult struct {
	Description     string         `form:"Description" validate:"omitempty" json:"description"`
	ResultFileIDs   []int64        `form:"ResultFileIDs" validate:"omitempty" json:"resultFileIDs"`
	ConfigSnapshots postgres.Jsonb `form:"ConfigSnapshots" validate:"omitempty" json:"configSnapshots"`
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
	if r.Result.ResultFileIDs == nil {
		s.ResultFileIDs = []int64{}
	} else {
		s.ResultFileIDs = r.Result.ResultFileIDs
	}

	s.ScenarioID = r.Result.ScenarioID

	return s
}

func (r *updateResultRequest) updatedResult(oldResult Result) Result {
	// Use the old Result as a basis for the updated Result `s`
	s := oldResult

	s.Result.Description = r.Result.Description
	if r.Result.ResultFileIDs == nil {
		s.ResultFileIDs = []int64{}
	} else {
		s.ResultFileIDs = r.Result.ResultFileIDs
	}

	// only update snapshots if not empty
	var emptyJson postgres.Jsonb
	// Serialize empty json and params
	emptyJson_ser, _ := json.Marshal(emptyJson)
	configSnapshots_ser, _ := json.Marshal(r.Result.ConfigSnapshots)
	opts := jsondiff.DefaultConsoleOptions()
	diff, _ := jsondiff.Compare(emptyJson_ser, configSnapshots_ser, &opts)
	if diff.String() != "FullMatch" {
		s.ConfigSnapshots = r.Result.ConfigSnapshots
	}

	return s
}
