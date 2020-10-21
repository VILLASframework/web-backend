/** InfrastructureComponent package, methods.
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
package infrastructure_component

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"log"
	"time"
)

type InfrastructureComponent struct {
	database.InfrastructureComponent
}

func (s *InfrastructureComponent) Save() error {
	db := database.GetDB()
	err := db.Create(s).Error
	return err
}

func (s *InfrastructureComponent) ByID(id uint) error {
	db := database.GetDB()
	err := db.Find(s, id).Error
	return err
}

func (s *InfrastructureComponent) ByUUID(uuid string) error {
	db := database.GetDB()
	err := db.Find(s, "UUID = ?", uuid).Error
	return err
}

func (s *InfrastructureComponent) Update(updatedIC InfrastructureComponent) error {

	db := database.GetDB()
	err := db.Model(s).Updates(updatedIC).Error

	return err
}

func (s *InfrastructureComponent) delete(receivedViaAMQP bool) error {
	if s.ManagedExternally && !receivedViaAMQP {
		var action Action
		action.Act = "delete"
		action.When = time.Now().Unix()
		action.Properties.UUID = new(string)
		*action.Properties.UUID = s.UUID

		err := sendActionAMQP(action)
		return err
	}

	db := database.GetDB()

	no_configs := db.Model(s).Association("ComponentConfigurations").Count()

	if no_configs > 0 {
		return fmt.Errorf("Infrastructure Component cannot be deleted as it is still used in configurations (active or dangling)")
	}

	// delete InfrastructureComponent from DB (does NOT remain as dangling)
	err := db.Delete(s).Error
	return err
}

func (s *InfrastructureComponent) getConfigs() ([]database.ComponentConfiguration, int, error) {
	db := database.GetDB()
	var configs []database.ComponentConfiguration
	err := db.Order("ID asc").Model(s).Related(&configs, "ComponentConfigurations").Error
	return configs, len(configs), err
}

func createNewICviaAMQP(payload ICUpdate) error {

	var newICReq AddICRequest
	newICReq.InfrastructureComponent.UUID = payload.Properties.UUID
	if payload.Properties.Name == nil ||
		payload.Properties.Category == nil ||
		payload.Properties.Type == nil {
		// cannot create new IC because required information (name, type, and/or category missing)
		return fmt.Errorf("AMQP: Cannot create new IC, required field(s) is/are missing: name, type, category")
	}
	newICReq.InfrastructureComponent.Name = *payload.Properties.Name
	newICReq.InfrastructureComponent.Category = *payload.Properties.Category
	newICReq.InfrastructureComponent.Type = *payload.Properties.Type

	// add optional params
	if payload.State != nil {
		newICReq.InfrastructureComponent.State = *payload.State
	} else {
		newICReq.InfrastructureComponent.State = "unknown"
	}
	// TODO check if state is "gone" and abort creation of IC in this case

	if payload.Properties.WS_url != nil {
		newICReq.InfrastructureComponent.WebsocketURL = *payload.Properties.WS_url
	}
	if payload.Properties.API_url != nil {
		newICReq.InfrastructureComponent.APIURL = *payload.Properties.API_url
	}
	if payload.Properties.Location != nil {
		newICReq.InfrastructureComponent.Location = *payload.Properties.Location
	}
	if payload.Properties.Description != nil {
		newICReq.InfrastructureComponent.Description = *payload.Properties.Description
	}
	// TODO add JSON start parameter scheme

	// set managed externally to true because this IC is created via AMQP
	newICReq.InfrastructureComponent.ManagedExternally = newTrue()

	// Validate the new IC
	err := newICReq.validate()
	if err != nil {
		return fmt.Errorf("AMQP: Validation of new IC failed: %v", err)
	}

	// Create the new IC
	newIC, err := newICReq.createIC(true)
	if err != nil {
		return fmt.Errorf("AMQP: Creating new IC failed: %v", err)
	}

	// save IC
	err = newIC.Save()
	if err != nil {
		return fmt.Errorf("AMQP: Saving new IC to DB failed: %v", err)
	}

	return nil
}

func (s *InfrastructureComponent) updateICviaAMQP(payload ICUpdate) error {
	var updatedICReq UpdateICRequest
	if payload.State != nil {
		updatedICReq.InfrastructureComponent.State = *payload.State

		if *payload.State == "gone" {
			// remove IC from DB
			err := s.delete(true)
			if err != nil {
				// if component could not be deleted there are still configurations using it in the DB
				// continue with the update to save the new state of the component and get back to the deletion later
				log.Println("Could not delete IC because there is a config using it, deletion postponed")
			}

		}
	}
	if payload.Properties.Type != nil {
		updatedICReq.InfrastructureComponent.Type = *payload.Properties.Type
	}
	if payload.Properties.Category != nil {
		updatedICReq.InfrastructureComponent.Category = *payload.Properties.Category
	}
	if payload.Properties.Name != nil {
		updatedICReq.InfrastructureComponent.Name = *payload.Properties.Name
	}
	if payload.Properties.WS_url != nil {
		updatedICReq.InfrastructureComponent.WebsocketURL = *payload.Properties.WS_url
	}
	if payload.Properties.API_url != nil {
		updatedICReq.InfrastructureComponent.APIURL = *payload.Properties.API_url
	}
	if payload.Properties.Location != nil {
		//postgres.Jsonb{json.RawMessage(`{"location" : " ` + *payload.Properties.Location + `"}`)}
		updatedICReq.InfrastructureComponent.Location = *payload.Properties.Location
	}
	if payload.Properties.Description != nil {
		updatedICReq.InfrastructureComponent.Description = *payload.Properties.Description
	}
	// TODO add JSON start parameter scheme

	// set managed externally to true because this IC is updated via AMQP
	updatedICReq.InfrastructureComponent.ManagedExternally = newTrue()

	// Validate the updated IC
	err := updatedICReq.validate()
	if err != nil {
		return fmt.Errorf("AMQP: Validation of updated IC failed: %v", err)
	}

	// Create the updated IC from old IC
	updatedIC, err := updatedICReq.updatedIC(*s, true)
	if err != nil {
		return fmt.Errorf("AMQP: Unable to update IC %v : %v", s.Name, err)
	}

	// Finally update the IC in the DB
	err = s.Update(updatedIC)
	if err != nil {
		return fmt.Errorf("AMQP: Unable to update IC %v in DB: %v", s.Name, err)
	}

	return err
}

func newTrue() *bool {
	b := true
	return &b
}
