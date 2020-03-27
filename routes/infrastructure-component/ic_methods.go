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
)

type InfrastructureComponent struct {
	database.InfrastructureComponent
}

func (s *InfrastructureComponent) save() error {
	db := database.GetDB()
	err := db.Create(s).Error
	return err
}

func (s *InfrastructureComponent) ByID(id uint) error {
	db := database.GetDB()
	err := db.Find(s, id).Error
	return err
}

func (s *InfrastructureComponent) update(updatedIC InfrastructureComponent) error {

	db := database.GetDB()
	err := db.Model(s).Updates(updatedIC).Error

	return err
}

func (s *InfrastructureComponent) delete() error {
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