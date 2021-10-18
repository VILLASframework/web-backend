/** Scenario package, methods.
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
package scenario

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	component_configuration "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/component-configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/file"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type Scenario struct {
	database.Scenario
}

func (s *Scenario) ByID(id uint) error {
	db := database.GetDB()
	err := db.Find(s, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Scenario) getUsers() ([]database.User, int, error) {
	db := database.GetDB()
	var users []database.User
	err := db.Order("ID asc").Model(s).Where("Active = ?", true).Related(&users, "Users").Error
	return users, len(users), err
}

func (s *Scenario) save() error {
	db := database.GetDB()
	err := db.Create(s).Error
	return err
}

func (s *Scenario) update(updatedScenario Scenario) error {

	// TODO: if the field is empty member shouldn't be updated
	s.Name = updatedScenario.Name
	s.IsLocked = updatedScenario.IsLocked
	s.StartParameters = updatedScenario.StartParameters

	db := database.GetDB()
	err := db.Model(s).Update(updatedScenario).Error
	if err != nil {
		return err
	}

	// extra update for bool IsLocked since it is ignored if false
	err = db.Model(s).Updates(map[string]interface{}{"IsLocked": updatedScenario.IsLocked}).Error
	return err
}

func (s *Scenario) addUser(u *database.User) error {

	db := database.GetDB()
	err := db.Model(s).Association("Users").Append(u).Error
	return err
}

func (s *Scenario) deleteUser(username string) error {
	db := database.GetDB()

	var deletedUser database.User
	err := db.Find(&deletedUser, "Username = ?", username).Error
	if err != nil {
		return err
	}

	no_users := db.Model(s).Association("Users").Count()

	if no_users > 1 {
		// remove user from scenario
		err = db.Model(s).Association("Users").Delete(&deletedUser).Error
		if err != nil {
			return err
		}
		// remove scenario from user
		err = db.Model(&deletedUser).Association("Scenarios").Delete(s).Error
		if err != nil {
			return err
		}
	} else {
		// There is only one associated user
		var remainingUser database.User
		err = db.Model(s).Related(&remainingUser, "Users").Error
		if remainingUser.Username == username {
			// if the remaining user is the one to be deleted
			return fmt.Errorf("cannot delete last user from scenario without deleting scenario itself, doing nothing")
		} else {
			// the remaining user is NOT the one to be deleted
			// that means the user to be deleted is not associated with the scenario
			return gorm.ErrRecordNotFound
		}
	}

	return nil
}

func (s *Scenario) delete() error {
	db := database.GetDB()

	// delete scenario from all users and vice versa

	users, no_users, err := s.getUsers()
	if err != nil {
		return err
	}

	if no_users > 0 {
		for _, u := range users {
			// remove user from scenario
			err = db.Model(s).Association("Users").Delete(&u).Error
			if err != nil {
				return err
			}
			// remove scenario from user
			err = db.Model(&u).Association("Scenarios").Delete(s).Error
			if err != nil {
				return err
			}
		}
	}

	// Delete scenario
	err = db.Delete(s).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *Scenario) DuplicateScenarioForUser(user *database.User) <-chan error {
	errs := make(chan error, 1)

	go func() {

		// get all component configs of the scenario
		db := database.GetDB()
		var configs []database.ComponentConfiguration
		err := db.Order("ID asc").Model(s).Related(&configs, "ComponentConfigurations").Error
		if err != nil {
			log.Printf("Warning: scenario to duplicate (id=%d) has no component configurations", s.ID)
		}

		// iterate over component configs to check for ICs to duplicate
		duplicatedICuuids := make(map[uint]string) // key: original icID; value: UUID of duplicate
		var externalUUIDs []string                 // external ICs to wait for
		for _, config := range configs {
			icID := config.ICID
			if duplicatedICuuids[icID] != "" { // this IC was already added
				continue
			}

			var ic infrastructure_component.InfrastructureComponent
			err = ic.ByID(icID)

			if err != nil {
				errs <- fmt.Errorf("Cannot find IC with id %d in DB, will not duplicate for User %s: %s", icID, user.Username, err)
				continue
			}

			// create new kubernetes simulator OR use existing IC
			if ic.Category == "simulator" && ic.Type == "kubernetes" {
				duplicateUUID, err := ic.RequestICcreateAMQPsimpleManager(user.Username)
				if err != nil {
					errs <- fmt.Errorf("Duplication of IC (id=%d) unsuccessful, err: %s", icID, err)
					continue
				}

				duplicatedICuuids[ic.ID] = duplicateUUID
				externalUUIDs = append(externalUUIDs, duplicateUUID)
			} else { // use existing IC
				duplicatedICuuids[ic.ID] = ""
				err = nil
			}
		}

		// copy scenario after all new external ICs are in DB
		icsToWaitFor := len(externalUUIDs)
		//var duplicatedScenario database.Scenario
		var timeout = 20 // seconds

		for i := 0; i < timeout; i++ {
			// duplicate scenario after all duplicated ICs have been found in the DB
			if icsToWaitFor == 0 {
				err := s.duplicateScenario(duplicatedICuuids, user)
				if err != nil {
					errs <- fmt.Errorf("duplicate scenario %v fails with error %v", s.Name, err.Error())
				}

				close(errs)
				return
			} else {
				time.Sleep(1 * time.Second)
			}

			// check for new ICs with previously created UUIDs
			for _, uuid := range externalUUIDs {
				if uuid == "" {
					continue
				}
				log.Printf("Looking for duplicated IC with UUID %s", uuid)
				var duplicatedIC database.InfrastructureComponent
				err = db.Find(&duplicatedIC, "UUID = ?", uuid).Error
				if err != nil {
					errs <- fmt.Errorf("Error looking up duplicated IC: %s", err)
				} else {
					icsToWaitFor--
					uuid = ""
				}
			}
		}

		errs <- fmt.Errorf("ALERT! Timed out while waiting for IC duplication, scenario not properly duplicated")
		close(errs)

	}()

	return errs
}

func (s *Scenario) duplicateScenario(icIds map[uint]string, user *database.User) error {

	db := database.GetDB()

	var duplicateSo Scenario
	duplicateSo.Name = s.Name + ` ` + user.Username
	duplicateSo.StartParameters.RawMessage = s.StartParameters.RawMessage
	err := duplicateSo.save()
	if err != nil {
		log.Printf("Could not create duplicate of scenario %d", s.ID)
		return err
	}

	// associate user to new scenario
	err = duplicateSo.addUser(user)
	if err != nil {
		log.Printf("Could not associate User %s to scenario %d", user.Username, duplicateSo.ID)
	}
	log.Println("Associated user to duplicated scenario")

	// duplicate files
	var files []file.File
	err = db.Order("ID asc").Model(s).Related(&files, "Files").Error
	if err != nil {
		log.Printf("error getting files for scenario %d", s.ID)
	}
	for _, f := range files {
		err = f.Duplicate(duplicateSo.ID)
		if err != nil {
			log.Print("error creating duplicate file %d: %s", f.ID, err)
			continue
		}
	}

	var configs []component_configuration.ComponentConfiguration
	// map existing signal IDs to duplicated signal IDs for widget duplication
	signalMap := make(map[uint]uint)
	err = db.Order("ID asc").Model(s).Related(&configs, "ComponentConfigurations").Error
	if err == nil {
		for _, c := range configs {
			err = c.Duplicate(duplicateSo.ID, icIds, &signalMap)
			//err = duplicateComponentConfig(&c, duplicateSo, icIds, userName, &signalMap)
			if err != nil {
				log.Printf("Error duplicating component config %d: %s", c.ID, err)
				continue
			}
		}
	} else {
		return err
	}

	var dabs []dashboard.Dashboard
	err = db.Order("ID asc").Model(s).Related(&dabs, "Dashboards").Error
	if err != nil {
		log.Printf("Error getting dashboards for scenario %d: %s", s.ID, err)
	}

	for _, dab := range dabs {
		err = dab.Duplicate(duplicateSo.ID, signalMap)
		if err != nil {
			log.Printf("Error duplicating dashboard %d: %s", dab.ID, err)
			continue
		}
	}

	return err
}
