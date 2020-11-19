package result

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
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

	// TODO delete Result + files (if any)

	return err
}
