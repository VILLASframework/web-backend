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

package user

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func IsAlreadyDuplicated(sc *database.Scenario, u *database.User) bool {
	duplicateName := fmt.Sprintf("%s %s", sc.Name, u.Username)
	db := database.GetDB()
	var scenarios []database.Scenario
	db.Find(&scenarios, "name = ?", duplicateName)

	return (len(scenarios) > 0)
}

// check if access of U to SC is exclusively granted by UG
func IsExclusiveAccess(sc *database.Scenario, u *database.User, ug *database.UserGroup) bool {
	db := database.GetDB()
	var ugs []database.UserGroup
	db.Model(u).Association("UserGroups").Find(&ugs)
	for _, asc_ug := range ugs {
		if ug.ID == asc_ug.ID {
			continue
		}
		var sms []database.ScenarioMapping
		db.Model(&asc_ug).Association("ScenarioMappings").Find(&sms)
		for _, sm := range sms {
			if sm.ScenarioID == sc.ID {
				return false
			}
		}
	}
	return true
}

func DuplicateScenarioForUser(s database.Scenario, user *database.User, uuidstr string) {
	if IsAlreadyDuplicated(&s, user) {
		return
	}
	// get all component configs of the scenario
	db := database.GetDB()
	var configs []database.ComponentConfiguration
	err := db.Order("ID asc").Model(s).Related(&configs, "ComponentConfigurations").Error
	if err != nil {
		log.Printf("Warning: scenario to duplicate (id=%d) has no component configurations", s.ID)
	}

	icIdmap := make(map[uint]uint)             // key: original IC id, value: duplicated IC id
	duplicatedICuuids := make(map[uint]string) // key: original icID; value: UUID of duplicate

	// iterate over component configs to check for ICs to duplicate
	for _, config := range configs {
		icID := config.ICID

		if _, ok := duplicatedICuuids[icID]; ok { // this IC was already added
			log.Println("IC already added while ranging configs")
			continue
		}

		var ic database.InfrastructureComponent
		err = db.Find(&ic, icID).Error

		if err != nil {
			log.Printf("Cannot find IC with id %d in DB, will not duplicate for User %s: %s", icID, user.Username, err)
			continue
		}

		alreadyDuplicated, dupID := isICalreadyDuplicated(ic, user.Username)
		// create new kubernetes simulator OR use existing IC
		if ic.Category == "simulator" && ic.Type == "kubernetes" && !alreadyDuplicated {
			duplicateUUID, err := duplicateIC(ic, user.Username, uuidstr)
			if err != nil {
				log.Printf("duplication of IC (id=%d) unsuccessful, err: %s", icID, err)
				continue
			}

			duplicatedICuuids[ic.ID] = duplicateUUID
		} else if alreadyDuplicated {
			duplicatedICuuids[ic.ID] = "alreadyduplicated"
			icIdmap[ic.ID] = dupID
			err = nil
		} else { // use existing IC
			duplicatedICuuids[ic.ID] = "useoriginal"
			icIdmap[ic.ID] = ic.ID
			err = nil
		}
	}

	// copy scenario after all new external ICs are in DB
	icsToWaitFor := len(duplicatedICuuids)
	var timeout = 20 // seconds

	for i := 0; i < timeout; i++ {
		// duplicate scenario after all duplicated ICs have been found in the DB
		log.Printf("Number of ICs to wait for: %d", icsToWaitFor)
		if icsToWaitFor <= 0 {
			err := duplicateScenario(s, icIdmap, user)
			if err != nil {
				log.Printf("duplicate scenario %v fails with error %v", s.Name, err.Error())
			}

			return
		} else {
			time.Sleep(1 * time.Second)
		}

		// check for new ICs with previously created UUIDs
		for icId, uuid_r := range duplicatedICuuids {
			if uuid_r == "alreadyduplicated" || uuid_r == "useoriginal" {
				icsToWaitFor--
				continue
			}
			log.Printf("Looking for duplicated IC with UUID %s", uuid_r)

			// a new IC is searched by UUID, otherwise by string comparison of IC name
			// this makes sense in testing and also in reality: if an external IC was requested
			// it should be checked for, checking shouldn't be omitted due to old duplicated IC
			var duplicatedIC database.InfrastructureComponent
			err = db.Find(&duplicatedIC, "UUID = ?", uuid_r).Error
			if err != nil {
				log.Printf("Didn't find duplicated IC: %s, keep waiting for it..", err)
			} else {
				log.Printf("Found duplicated IC! Original IC id: %d, duplicated IC id: %d", icId, duplicatedIC.ID)
				icsToWaitFor--
				icIdmap[icId] = duplicatedIC.ID
				uuid_r = ""
			}
		}
	}
}

func RemoveDuplicate(sc *database.Scenario, u *database.User) error {
	db := database.GetDB()

	var nsc database.Scenario
	duplicateName := fmt.Sprintf("%s %s", sc.Name, u.Username)
	err := db.Find(&nsc, "Name = ?", duplicateName).Error
	if err != nil {
		return err
	}

	err = db.Delete(&nsc).Error
	return err
}

func RemoveAccess(sc *database.Scenario, u *database.User, ug *database.UserGroup) error {
	if !IsExclusiveAccess(sc, u, ug) {
		return nil
	}
	db := database.GetDB()
	err := db.Model(&sc).Association("Users").Delete(&u).Error
	if err != nil {
		return err
	}
	err = db.Model(&u).Association("Scenarios").Delete(&sc).Error
	return err
}

func duplicateScenario(s database.Scenario, icIds map[uint]uint, user *database.User) error {

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
	fileidmap := make(map[uint]uint) // key: original file id, value: duplicated file id
	err = db.Order("ID asc").Model(s).Related(&files, "Files").Error
	if err != nil {
		log.Printf("error getting files for scenario %d", s.ID)
	}
	for _, f := range files {
		duplicateFileID, err := duplicateFile(f, duplicateSo.ID)
		if err != nil {
			log.Printf("error creating duplicate file %d: %s", f.ID, err)
			continue
		} else {
			fileidmap[f.ID] = duplicateFileID
		}
	}

	var configs []database.ComponentConfiguration
	configidmap := make(map[uint]uint) // key: original config id, value: duplicated config id
	// map existing signal IDs to duplicated signal IDs for widget duplication
	signalMap := make(map[uint]uint)
	err = db.Order("ID asc").Model(s).Related(&configs, "ComponentConfigurations").Error
	if err == nil {
		for _, c := range configs {
			duplicatConfigID, err := duplicateComponentConfig(c, duplicateSo.ID, icIds, &signalMap)
			if err != nil {
				log.Printf("Error duplicating component config %d: %s", c.ID, err)
				continue
			} else {
				configidmap[c.ID] = duplicatConfigID
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
		err = duplicateDashboard(dab, duplicateSo.ID, signalMap, configidmap, fileidmap, icIds)
		if err != nil {
			log.Printf("Error duplicating dashboard %d: %s", dab.ID, err)
			continue
		}
	}

	return err
}

func duplicateFile(f database.File, scenarioID uint) (uint, error) {

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
		return 0, err
	}

	// Create association of duplicate file to scenario ID of duplicate file

	var so database.Scenario
	err = db.Find(&so, scenarioID).Error
	if err != nil {
		return 0, err
	}

	err = db.Model(&so).Association("Files").Append(&dup).Error

	return dup.ID, err
}

func duplicateComponentConfig(m database.ComponentConfiguration, scenarioID uint, icIds map[uint]uint, signalMap *map[uint]uint) (uint, error) {

	db := database.GetDB()

	var dup database.ComponentConfiguration
	dup.Name = m.Name
	dup.StartParameters = m.StartParameters
	dup.ScenarioID = scenarioID

	if val, ok := icIds[m.ICID]; ok {
		var duplicatedIC database.InfrastructureComponent
		err := db.Find(&duplicatedIC, "ID = ?", val).Error
		if err != nil {
			log.Print(err)
			return 0, err
		}
		dup.ICID = duplicatedIC.ID
	} else {
		dup.ICID = m.ICID
	}

	// save duplicate to DB and create associations with IC and scenario
	var so database.Scenario
	err := db.Find(&so, scenarioID).Error
	if err != nil {
		return 0, err
	}

	// save component configuration to DB
	err = db.Create(&dup).Error
	if err != nil {
		return 0, err
	}

	// associate IC with component configuration
	var ic database.InfrastructureComponent
	err = db.Find(&ic, dup.ICID).Error
	if err != nil {
		return 0, err
	}
	err = db.Model(&ic).Association("ComponentConfigurations").Append(&dup).Error
	if err != nil {
		return 0, err
	}

	// associate component configuration with scenario
	err = db.Model(&so).Association("ComponentConfigurations").Append(&dup).Error
	if err != nil {
		return 0, err
	}

	// duplication of signals
	var sigs []database.Signal
	err = db.Order("ID asc").Model(&m).Related(&sigs, "OutputMapping").Error
	if err != nil {
		return 0, err
	}
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
			return 0, err
		}

		// associate signal with component configuration in correct direction
		if sigDup.Direction == "in" {
			err = db.Model(&dup).Association("InputMapping").Append(&sigDup).Error
		} else {
			err = db.Model(&dup).Association("OutputMapping").Append(&sigDup).Error
		}

		if err != nil {
			return 0, err
		}

		smap[s.ID] = sigDup.ID
	}

	return dup.ID, nil
}

func duplicateDashboard(d database.Dashboard, scenarioID uint, signalMap map[uint]uint,
	configIDmap map[uint]uint, fileIDmap map[uint]uint, icIds map[uint]uint) error {

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

		err = duplicateWidget(w, duplicateD.ID, signalMap, configIDmap, fileIDmap, icIds)
		if err != nil {
			log.Printf("error creating duplicate for widget %d: %s", w.ID, err)
			continue
		}
	}

	return nil
}

func duplicateWidget(w database.Widget, dashboardID uint, signalMap map[uint]uint,
	configIDmap map[uint]uint, fileIDmap map[uint]uint, icIds map[uint]uint) error {

	var duplicateW database.Widget
	duplicateW.DashboardID = dashboardID
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

	if w.Type == "ICstatus" {
		duplicateW.CustomProperties = duplicateICStatusCustomProps(w.CustomProperties, icIds)
	} else if w.Type == "Player" {
		duplicateW.CustomProperties = duplicatePlayerCustomProps(w.CustomProperties, configIDmap)
	} else if w.Type == "Image" {
		duplicateW.CustomProperties = duplicateImageCustomProps(w.CustomProperties, fileIDmap)
	} else {
		duplicateW.CustomProperties = w.CustomProperties
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

func duplicateICStatusCustomProps(customProps postgres.Jsonb, icIds map[uint]uint) postgres.Jsonb {
	type ICstatusCustomProps struct {
		CheckedIDs []uint `json:"checkedIDs"`
	}

	var props ICstatusCustomProps
	err := json.Unmarshal(customProps.RawMessage, &props)
	if err != nil {
		log.Printf("ICstatus duplication, unmarshalling failed: err: %s", err)
		return customProps
	}

	var IDs []string
	for _, id := range props.CheckedIDs {
		IDs = append(IDs, strconv.FormatUint(uint64(icIds[id]), 10))
	}

	customProperties := fmt.Sprintf(`{"checkedIDs": [%s]}`, strings.Join(IDs, ","))

	return postgres.Jsonb{RawMessage: json.RawMessage(customProperties)}
}

func duplicatePlayerCustomProps(customProps postgres.Jsonb, configIDmap map[uint]uint) postgres.Jsonb {
	type PlayerCustomProps struct {
		ConfigID      string `json:"configID"`
		ConfigIDs     []int  `json:"configIDs"`
		UploadResults bool   `json:"uploadResults"`
	}

	var props PlayerCustomProps
	err := json.Unmarshal(customProps.RawMessage, &props)
	if err != nil {
		log.Printf("Player duplication, unmarshalling failed: err: %v", err)
		return customProps
	}

	// get configID of original config
	u, err := strconv.ParseUint(props.ConfigID, 10, 64)
	if err != nil {
		log.Printf("Player duplication, parsing file ID failed: err: %v", err)
		return customProps
	}

	// get configID of duplicated config, save it in PlayerCustomProps struct
	props.ConfigID = strconv.FormatUint(uint64(configIDmap[uint(u)]), 10)
	customProperties, err := json.Marshal(props)
	if err != nil {
		log.Printf("Player duplication, marshalling failed: err: %v", err)
		return customProps
	}

	return postgres.Jsonb{RawMessage: customProperties}
}

func duplicateImageCustomProps(customProps postgres.Jsonb, fileIDmap map[uint]uint) postgres.Jsonb {
	type ImageCustomProps struct {
		File   string `json:"file"`
		Update bool   `json:"update"`
		Lock   bool   `json:"lockAspect"`
	}

	var props ImageCustomProps
	err := json.Unmarshal(customProps.RawMessage, &props)

	if err != nil {
		log.Printf("Image duplication, unmarshalling failed: err: %v", err)
		return customProps
	}

	// get fileID of original file
	u, err := strconv.ParseUint(props.File, 10, 64)
	if err != nil {
		log.Printf("Image duplication, parsing file ID failed: err: %v", err)
		return customProps
	}
	// get fileID of duplicated file, save it in ImageCustomProps struct
	props.File = strconv.FormatUint(uint64(fileIDmap[uint(u)]), 10)
	customProperties, err := json.Marshal(props)
	if err != nil {
		log.Printf("Image duplication, marshalling failed: err: %v", err)
		return customProps
	}

	return postgres.Jsonb{RawMessage: customProperties}
}

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
	Active   int         `json:"activeDeadlineSeconds"`
	Template JobTemplate `json:"template"`
}

type JobMetaData struct {
	JobName string `json:"name"`
}

type KubernetesJob struct {
	Spec     JobSpec     `json:"spec"`
	MetaData JobMetaData `json:"metadata"`
}

type ICPropertiesKubernetesJob struct {
	Job         KubernetesJob `json:"job"`
	UUID        string        `json:"uuid"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Location    string        `json:"location"`
	Owner       string        `json:"owner"`
	Category    string        `json:"category"`
	Type        string        `json:"type"`
}

type ICStatus struct {
	State     string  `json:"state"`
	Version   string  `json:"version"`
	Uptime    float64 `json:"uptime"`
	Result    string  `json:"result"`
	Error     string  `json:"error"`
	ManagedBy string  `json:"managed_by"`
}

type ICUpdateKubernetesJob struct {
	Properties ICPropertiesKubernetesJob `json:"properties"`
	Status     ICStatus                  `json:"status"`
	Schema     json.RawMessage           `json:"schema"`
}

func duplicateIC(ic database.InfrastructureComponent, userName string, uuidstr string) (string, error) {

	//WARNING: this function only works with the kubernetes-simple manager of VILLAScontroller
	if ic.Category != "simulator" || ic.Type != "kubernetes" {
		return "", fmt.Errorf("IC to duplicate is not a kubernetes simulator (%s %s)", ic.Type, ic.Category)
	}

	newUUID := uuid.New().String()

	var lastUpdate ICUpdateKubernetesJob
	err := json.Unmarshal(ic.StatusUpdateRaw.RawMessage, &lastUpdate)
	if err != nil {
		return newUUID, err
	}

	if uuidstr != "" {
		newUUID = "4854af30-325f-44a5-ad59-b67b2597de68"
	}

	msg := `{"name": "` + lastUpdate.Properties.Name + ` ` + userName + `",` +
		`"location": "` + lastUpdate.Properties.Location + `",` +
		`"category": "` + lastUpdate.Properties.Category + `",` +
		`"type": "` + lastUpdate.Properties.Type + `",` +
		`"uuid": "` + newUUID + `",` +
		`"jobname": "` + lastUpdate.Properties.Job.MetaData.JobName + `-` + userName + `",` +
		`"activeDeadlineSeconds": "` + strconv.Itoa(lastUpdate.Properties.Job.Spec.Active) + `",` +
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

	if session != nil {
		if session.IsReady {
			err = session.Send(payload, ic.Manager)
			return newUUID, err
		} else {
			return "", fmt.Errorf("could not send IC create action, AMQP session is not ready")
		}
	} else {
		return "", fmt.Errorf("could not send IC create action, AMQP session is nil")
	}

}

func isICalreadyDuplicated(ic database.InfrastructureComponent, username string) (bool, uint) {
	db := database.GetDB()
	var duplicateICs []database.InfrastructureComponent
	duplicateName := fmt.Sprintf("%s %s", ic.Name, username)
	err := db.Find(&duplicateICs, "Name = ?", duplicateName).Error
	if err != nil {
		log.Printf("Error looking for duplicated ICs: %s", err)
	}

	// return the first duplicated IC that was recently updated
	for _, value := range duplicateICs {
		lastUpdate := value.UpdatedAt
		diff := time.Since(lastUpdate)
		if diff.Seconds() < 60 {
			return true, value.ID
		}
	}

	return false, 0
}
