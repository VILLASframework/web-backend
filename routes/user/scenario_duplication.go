/** User package, authentication endpoint.
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

package user

import (
	"encoding/json"
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"github.com/google/uuid"
	"log"
	"time"
)

func duplicateScenarioForUser(s database.Scenario, user *database.User) <-chan error {
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

			var ic database.InfrastructureComponent
			err = db.Find(&ic, icID).Error

			if err != nil {
				log.Printf("Cannot find IC with id %d in DB, will not duplicate for User %s: %s", icID, user.Username, err)
				continue
			}

			// create new kubernetes simulator OR use existing IC
			if ic.Category == "simulator" && ic.Type == "kubernetes" {
				duplicateUUID, err := duplicateIC(ic, user.Username)
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
				err := duplicateScenario(s, duplicatedICuuids, user)
				if err != nil {
					errs <- fmt.Errorf("duplicate scenario %v fails with error %v", s.Name, err.Error())
				}

				close(errs)
				return
			} else {
				time.Sleep(1 * time.Second)
			}

			// check for new ICs with previously created UUIDs
			for _, uuid_r := range externalUUIDs {
				if uuid_r == "" {
					continue
				}
				log.Printf("Looking for duplicated IC with UUID %s", uuid_r)
				var duplicatedIC database.InfrastructureComponent
				err = db.Find(&duplicatedIC, "UUID = ?", uuid_r).Error
				if err != nil {
					errs <- fmt.Errorf("Error looking up duplicated IC: %s", err)
				} else {
					icsToWaitFor--
					uuid_r = ""
				}
			}
		}

		errs <- fmt.Errorf("ALERT! Timed out while waiting for IC duplication, scenario not properly duplicated")
		close(errs)

	}()

	return errs
}

func duplicateScenario(s database.Scenario, icIds map[uint]string, user *database.User) error {

	db := database.GetDB()

	var duplicateSo database.Scenario
	duplicateSo.Name = s.Name + ` ` + user.Username
	duplicateSo.StartParameters.RawMessage = s.StartParameters.RawMessage

	err := db.Create(&duplicateSo).Error
	if err != nil {
		log.Printf("Could not create duplicate of scenario %d", s.ID)
		return err
	}

	// associate user to new scenario
	err = db.Model(&duplicateSo).Association("Users").Append(user).Error
	if err != nil {
		log.Printf("Could not associate User %s to scenario %d", user.Username, duplicateSo.ID)
	}
	log.Println("Associated user to duplicated scenario")

	// duplicate files
	var files []database.File
	err = db.Order("ID asc").Model(s).Related(&files, "Files").Error
	if err != nil {
		log.Printf("error getting files for scenario %d", s.ID)
	}
	for _, f := range files {
		err = duplicateFile(f, duplicateSo.ID)
		if err != nil {
			log.Printf("error creating duplicate file %d: %s", f.ID, err)
			continue
		}
	}

	var configs []database.ComponentConfiguration
	// map existing signal IDs to duplicated signal IDs for widget duplication
	signalMap := make(map[uint]uint)
	err = db.Order("ID asc").Model(s).Related(&configs, "ComponentConfigurations").Error
	if err == nil {
		for _, c := range configs {
			err = duplicateComponentConfig(c, duplicateSo.ID, icIds, &signalMap)
			if err != nil {
				log.Printf("Error duplicating component config %d: %s", c.ID, err)
				continue
			}
		}
	} else {
		return err
	}

	var dabs []database.Dashboard
	err = db.Order("ID asc").Model(s).Related(&dabs, "Dashboards").Error
	if err != nil {
		log.Printf("Error getting dashboards for scenario %d: %s", s.ID, err)
	}

	for _, dab := range dabs {
		err = duplicateDashboard(dab, duplicateSo.ID, signalMap)
		if err != nil {
			log.Printf("Error duplicating dashboard %d: %s", dab.ID, err)
			continue
		}
	}

	return err
}

func duplicateFile(f database.File, scenarioID uint) error {

	var dup database.File
	dup.Name = f.Name
	dup.Key = f.Key
	dup.Type = f.Type
	dup.Size = f.Size
	dup.Date = f.Date
	dup.ScenarioID = scenarioID
	dup.FileData = f.FileData
	dup.ImageHeight = f.ImageHeight
	dup.ImageWidth = f.ImageWidth

	// file duplicate will point to the same data blob in the DB (SQL or postgres)

	// Add duplicate File object with parameters to DB
	db := database.GetDB()
	err := db.Create(&dup).Error
	if err != nil {
		return err
	}

	// Create association of duplicate file to scenario ID of duplicate file

	var so database.Scenario
	err = db.Find(&so, scenarioID).Error
	if err != nil {
		return err
	}

	err = db.Model(&so).Association("Files").Append(&dup).Error

	return err
}

func duplicateComponentConfig(m database.ComponentConfiguration, scenarioID uint, icIds map[uint]string, signalMap *map[uint]uint) error {

	db := database.GetDB()

	var dup database.ComponentConfiguration
	dup.Name = m.Name
	dup.StartParameters = m.StartParameters
	dup.ScenarioID = scenarioID

	if icIds[m.ICID] == "" {
		dup.ICID = m.ICID
	} else {
		var duplicatedIC database.InfrastructureComponent
		err := db.Find(&duplicatedIC, "UUID = ?", icIds[m.ICID]).Error
		if err != nil {
			log.Print(err)
			return err
		}
		dup.ICID = duplicatedIC.ID
	}

	// save duplicate to DB and create associations with IC and scenario
	var so database.Scenario
	err := db.Find(&so, m.ScenarioID).Error
	if err != nil {
		return err
	}

	// save component configuration to DB
	err = db.Create(&dup).Error
	if err != nil {
		return err
	}

	// associate IC with component configuration
	var ic database.InfrastructureComponent
	err = db.Find(&ic, dup.ICID).Error
	if err != nil {
		return err
	}
	err = db.Model(&ic).Association("ComponentConfigurations").Append(&dup).Error
	if err != nil {
		return err
	}

	// associate component configuration with scenario
	err = db.Model(&so).Association("ComponentConfigurations").Append(&dup).Error
	if err != nil {
		return err
	}

	// duplication of signals
	var sigs []database.Signal
	err = db.Order("ID asc").Model(&m).Related(&sigs, "OutputMapping").Error
	smap := *signalMap
	for _, s := range sigs {
		var sigDup database.Signal
		sigDup.Direction = s.Direction
		sigDup.Index = s.Index
		sigDup.Name = s.Name // + ` ` + userName
		sigDup.ScalingFactor = s.ScalingFactor
		sigDup.Unit = s.Unit
		sigDup.ConfigID = dup.ID

		// save signal to DB
		err = db.Create(&sigDup).Error
		if err != nil {
			return err
		}

		// associate signal with component configuration in correct direction
		if s.Direction == "in" {
			err = db.Model(&dup).Association("InputMapping").Append(&s).Error
		} else {
			err = db.Model(&dup).Association("OutputMapping").Append(&s).Error
		}

		if err != nil {
			return err
		}

		smap[s.ID] = sigDup.ID
	}

	return nil
}

func duplicateDashboard(d database.Dashboard, scenarioID uint, signalMap map[uint]uint) error {

	var duplicateD database.Dashboard
	duplicateD.Grid = d.Grid
	duplicateD.Name = d.Name
	duplicateD.ScenarioID = scenarioID
	duplicateD.Height = d.Height

	db := database.GetDB()
	var so database.Scenario
	err := db.Find(&so, duplicateD.ScenarioID).Error
	if err != nil {
		return err
	}

	// save dashboard to DB
	err = db.Create(&duplicateD).Error
	if err != nil {
		return err
	}

	// associate dashboard with scenario
	err = db.Model(&so).Association("Dashboards").Append(&duplicateD).Error

	if err != nil {
		return err
	}

	// add widgets to duplicated dashboard
	var widgets []database.Widget
	err = db.Order("ID asc").Model(d).Related(&widgets, "Widgets").Error
	if err != nil {
		log.Printf("Error getting widgets for dashboard %d: %s", d.ID, err)
	}
	for _, w := range widgets {

		err = duplicateWidget(w, duplicateD.ID, signalMap)
		if err != nil {
			log.Printf("error creating duplicate for widget %d: %s", w.ID, err)
			continue
		}
	}

	return nil
}

func duplicateWidget(w database.Widget, dashboardID uint, signalMap map[uint]uint) error {
	var duplicateW database.Widget
	duplicateW.DashboardID = dashboardID
	duplicateW.CustomProperties = w.CustomProperties
	duplicateW.Height = w.Height
	duplicateW.Width = w.Width
	duplicateW.MinHeight = w.MinHeight
	duplicateW.MinWidth = w.MinWidth
	duplicateW.Name = w.Name
	duplicateW.Type = w.Type
	duplicateW.X = w.X
	duplicateW.Y = w.Y
	duplicateW.Z = w.Z

	duplicateW.SignalIDs = []int64{}
	for _, id := range w.SignalIDs {
		duplicateW.SignalIDs = append(duplicateW.SignalIDs, int64(signalMap[uint(id)]))
	}

	db := database.GetDB()
	var dab database.Dashboard
	err := db.Find(&dab, duplicateW.DashboardID).Error
	if err != nil {
		return err
	}

	// save widget to DB
	err = db.Create(&duplicateW).Error
	if err != nil {
		return err
	}

	// associate widget with dashboard
	err = db.Model(&dab).Association("Widgets").Append(&duplicateW).Error
	return err
}

func duplicateIC(ic database.InfrastructureComponent, userName string) (string, error) {

	//WARNING: this function only works with the kubernetes-simple manager of VILLAScontroller
	if ic.Category != "simulator" || ic.Type == "kubernetes" {
		return "", nil
	}

	newUUID := uuid.New().String()
	log.Printf("New IC UUID: %s", newUUID)

	type Container struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	type TemplateSpec struct {
		Containers []Container `json:"containers"`
	}

	type JobTemplate struct {
		Spec TemplateSpec `json:"spec"`
	}

	type JobSpec struct {
		Active   string      `json:"activeDeadlineSeconds"`
		Template JobTemplate `json:"template"`
	}

	type JobMetaData struct {
		JobName string `json:"name"`
	}

	type KubernetesJob struct {
		Spec     JobSpec     `json:"spec"`
		MetaData JobMetaData `json:"metadata"`
	}

	type ICPropertiesToCopy struct {
		Job         KubernetesJob `json:"job"`
		UUID        string        `json:"uuid"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Location    string        `json:"location"`
		Owner       string        `json:"owner"`
		Category    string        `json:"category"`
		Type        string        `json:"type"`
	}

	type ICUpdateToCopy struct {
		Properties ICPropertiesToCopy `json:"properties"`
		Status     json.RawMessage    `json:"status"`
		Schema     json.RawMessage    `json:"schema"`
	}

	var lastUpdate ICUpdateToCopy
	log.Println(ic.StatusUpdateRaw.RawMessage)
	err := json.Unmarshal(ic.StatusUpdateRaw.RawMessage, &lastUpdate)
	if err != nil {
		return newUUID, err
	}

	msg := `{"name": "` + lastUpdate.Properties.Name + ` ` + userName + `",` +
		`"location": "` + lastUpdate.Properties.Location + `",` +
		`"category": "` + lastUpdate.Properties.Category + `",` +
		`"type": "` + lastUpdate.Properties.Type + `",` +
		`"uuid": "` + newUUID + `",` +
		`"jobname": "` + lastUpdate.Properties.Job.MetaData.JobName + `-` + userName + `",` +
		`"activeDeadlineSeconds": "` + lastUpdate.Properties.Job.Spec.Active + `",` +
		`"containername": "` + lastUpdate.Properties.Job.Spec.Template.Spec.Containers[0].Name + `-` + userName + `",` +
		`"image": "` + lastUpdate.Properties.Job.Spec.Template.Spec.Containers[0].Image + `",` +
		`"uuid": "` + newUUID + `"}`

	type Action struct {
		Act        string          `json:"action"`
		When       int64           `json:"when"`
		Parameters json.RawMessage `json:"parameters,omitempty"`
		Model      json.RawMessage `json:"model,omitempty"`
		Results    json.RawMessage `json:"results,omitempty"`
	}

	actionCreate := Action{
		Act:        "create",
		When:       time.Now().Unix(),
		Parameters: json.RawMessage(msg),
	}

	payload, err := json.Marshal(actionCreate)
	if err != nil {
		return "", err
	}

	if session.IsReady {
		err = session.Send(payload, ic.Manager)
		return newUUID, err
	} else {
		return "", fmt.Errorf("could not send IC create action, AMQP session is not ready")
	}

}
