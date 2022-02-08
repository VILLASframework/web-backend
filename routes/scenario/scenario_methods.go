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
	"log"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"github.com/jinzhu/gorm"
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
		if err != nil {
			return err
		}
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

func (s *Scenario) delete() []error {
	db := database.GetDB()

	var errs []error

	// delete all files of the scenario
	var files []database.File
	err := db.Order("ID asc").Model(s).Related(&files, "Files").Error
	if err != nil {
		errs = append(errs, err)
	}

	for _, f := range files {
		// delete file from s3 bucket
		if f.Key != "" {
			// TODO we do not delete the file from s3 object storage
			// to ensure that no data is lost if multiple File objects reference the same S3 data object
			// This behavior should be replaced by a different file handling in the future
			//err = f.deleteS3()
			//if err != nil {
			//	return err
			//}
			//log.Println("Deleted file in S3 object storage")
			log.Printf("Did NOT delete file with key %v in S3 object storage!\n", f.Key)
		}

		log.Println("DELETE file ", f.ID, "(name="+f.Name+")")
		err = db.Delete(&f).Error
		if err != nil {
			errs = append(errs, err)
		}
	}

	// delete all results of the scenario
	var results []database.Result
	err = db.Order("ID asc").Model(s).Related(&results, "Results").Error
	if err != nil {
		errs = append(errs, err)
	}

	for _, r := range results {
		log.Println("DELETE result ", r.ID, "(desc="+r.Description+")")
		err = db.Delete(&r).Error
		if err != nil {
			errs = append(errs, err)
		}
	}

	// delete all dashboards (and widgets) of the scenario
	var dab []database.Dashboard
	err = db.Order("ID asc").Model(s).Related(&dab, "Dashboards").Error
	if err != nil {
		errs = append(errs, err)
	}

	for _, d := range dab {
		// get all widgets of the dashboard
		var widgets []database.Widget
		err = db.Order("ID asc").Model(&d).Related(&widgets, "Widgets").Error
		if err != nil {
			errs = append(errs, err)
		}

		// Delete widgets
		for _, widget := range widgets {
			log.Println("DELETE widget ", widget.ID, "(name="+widget.Name+")")
			err = db.Delete(&widget).Error
			if err != nil {
				errs = append(errs, err)
			}
		}

		// Delete dashboard
		log.Println("DELETE dashboard ", d.ID, "(name="+d.Name+")")
		err = db.Delete(&d).Error
		if err != nil {
			errs = append(errs, err)
		}
	}

	// delete all component configs (and signals) of the scenario
	var configs []database.ComponentConfiguration
	err = db.Order("ID asc").Model(s).Related(&configs, "ComponentConfigurations").Error
	if err != nil {
		errs = append(errs, err)
	}

	for _, config := range configs {

		// Get Signals of InputMapping and delete them
		var InputMappingSignals []database.Signal
		err = db.Model(&config).Related(&InputMappingSignals, "InputMapping").Error
		if err != nil {
			errs = append(errs, err)
		}
		for _, sig := range InputMappingSignals {
			log.Println("DELETE signal ", sig.ID, "(name="+sig.Name+")")
			err = db.Delete(&sig).Error
			if err != nil {
				errs = append(errs, err)
			}
		}

		// Get Signals of OutputMapping and delete them
		var OutputMappingSignals []database.Signal
		err = db.Model(&config).Related(&OutputMappingSignals, "OutputMapping").Error
		if err != nil {
			errs = append(errs, err)
		}
		for _, sig := range OutputMappingSignals {
			log.Println("DELETE signal ", sig.ID, "(name="+sig.Name+")")
			err = db.Delete(&sig).Error
			if err != nil {
				errs = append(errs, err)
			}
		}

		var ic database.InfrastructureComponent
		err = db.Find(&ic, config.ICID).Error
		if err == nil {
			// remove association between Infrastructure component and config
			log.Println("DELETE ASSOCIATION to IC ", ic.ID, "(name="+ic.Name+")")
			err = db.Model(&ic).Association("ComponentConfigurations").Delete(&config).Error
			if err != nil {
				errs = append(errs, err)
			}

			// if IC has state gone and there is no component configuration associated with it: delete IC
			no_configs := db.Model(&ic).Association("ComponentConfigurations").Count()
			if no_configs == 0 && ic.State == "gone" {
				log.Println("DELETE IC with state gone, last component config deleted", ic.UUID)
				err = db.Delete(&ic).Error
				if err != nil {
					errs = append(errs, err)
				}

			}
		} else {
			if err == gorm.ErrRecordNotFound {
				log.Printf("SKIPPING IC association removal, IC with id=%v not found\n", config.ICID)
			} else {
				errs = append(errs, err)
			}
		}

		// delete component configuration
		log.Println("DELETE component config ", config.ID, "(name="+config.Name+")")
		err = db.Delete(&config).Error
		if err != nil {
			errs = append(errs, err)
		}
	}

	// delete scenario from all users and vice versa

	users, no_users, err := s.getUsers()
	if err != nil {
		errs = append(errs, err)
	}
	if no_users > 0 {
		for _, u := range users {
			// remove user from scenario
			log.Println("DELETE ASSOCIATION to user", u.ID, "(name="+u.Username+")")
			err = db.Model(s).Association("Users").Delete(&u).Error
			if err != nil {
				errs = append(errs, err)
			}
			// remove scenario from user
			err = db.Model(&u).Association("Scenarios").Delete(s).Error
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	// Delete scenario
	err = db.Delete(s).Error
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}
