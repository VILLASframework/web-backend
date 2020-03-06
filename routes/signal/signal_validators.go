/** Signal package, validators.
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
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
package signal

import (
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

type validNewSignal struct {
	Name      string `form:"Name" validate:"required"`
	Unit      string `form:"unit" validate:"omitempty"`
	Index     uint   `form:"index" validate:"required"`
	Direction string `form:"direction" validate:"required,oneof=in out"`
	ConfigID  uint   `form:"configID" validate:"required"`
}

type validUpdatedSignal struct {
	Name  string `form:"Name" validate:"omitempty"`
	Unit  string `form:"unit" validate:"omitempty"`
	Index uint   `form:"index" validate:"omitempty"`
}

type addSignalRequest struct {
	Signal validNewSignal `json:"signal"`
}

type updateSignalRequest struct {
	Signal validUpdatedSignal `json:"signal"`
}

func (r *addSignalRequest) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *validUpdatedSignal) validate() error {
	validate = validator.New()
	errs := validate.Struct(r)
	return errs
}

func (r *addSignalRequest) createSignal() Signal {
	var s Signal

	s.Name = r.Signal.Name
	s.Unit = r.Signal.Unit
	s.Index = r.Signal.Index
	s.Direction = r.Signal.Direction
	s.ConfigID = r.Signal.ConfigID

	return s
}

func (r *updateSignalRequest) updatedSignal(oldSignal Signal) Signal {
	// Use the old Signal as a basis for the updated Signal `s`
	s := oldSignal

	if r.Signal.Name != "" {
		s.Name = r.Signal.Name
	}

	if r.Signal.Index != 0 {
		// TODO this implies that we start indexing at 1
		s.Index = r.Signal.Index
	}

	if r.Signal.Unit != "" {
		s.Unit = r.Signal.Unit
	}

	return s
}
