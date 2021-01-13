/** Result package, methods.
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

package result

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"log"
)

type Result struct {
	database.Result
}

func (r *Result) save() error {
	db := database.GetDB()
	err := db.Create(r).Error
	return err
}

func (r *Result) ByID(id uint) error {
	db := database.GetDB()
	err := db.Find(r, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Result) addToScenario() error {
	db := database.GetDB()
	var sco scenario.Scenario
	err := sco.ByID(r.ScenarioID)
	if err != nil {
		return err
	}

	// save result to DB
	err = r.save()
	if err != nil {
		return err
	}

	// associate result with scenario
	err = db.Model(&sco).Association("Results").Append(r).Error

	return err
}

func (r *Result) update(modifiedResult Result) error {

	db := database.GetDB()

	err := db.Model(r).Updates(map[string]interface{}{
		"Description":     modifiedResult.Description,
		"ConfigSnapshots": modifiedResult.ConfigSnapshots,
		"ResultFileIDs":   modifiedResult.ResultFileIDs,
	}).Error

	return err
}

func (r *Result) delete() error {

	db := database.GetDB()
	var sco scenario.Scenario
	err := sco.ByID(r.ScenarioID)
	if err != nil {
		return err
	}

	// remove association between Result and Scenario
	err = db.Model(&sco).Association("Results").Delete(r).Error

	// Delete result files
	for _, fileid := range r.ResultFileIDs {
		var f file.File
		err := f.ByID(uint(fileid))
		if err != nil {
			log.Println("Unable to delete file with ID ", fileid, err)
			continue
		}
		err = f.Delete()
		if err != nil {
			return err
		}
	}

	// Delete result
	err = db.Delete(r).Error

	return err
}

func (r *Result) addResultFileID(fileID uint) error {

	oldResultFileIDs := r.ResultFileIDs
	newResultFileIDs := append(oldResultFileIDs, int64(fileID))

	db := database.GetDB()

	err := db.Model(r).Updates(map[string]interface{}{
		"ResultFileIDs": newResultFileIDs,
	}).Error

	return err

}

func (r *Result) removeResultFileID(fileID uint) error {
	oldResultFileIDs := r.ResultFileIDs
	var newResultFileIDs []int64

	for _, id := range oldResultFileIDs {
		if id != int64(fileID) {
			newResultFileIDs = append(newResultFileIDs, id)
		}
	}

	db := database.GetDB()

	err := db.Model(r).Updates(map[string]interface{}{
		"ResultFileIDs": newResultFileIDs,
	}).Error

	return err
}
