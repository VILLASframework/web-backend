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

package database

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"
)

func CheckScenarioPermissions(c *gin.Context, operation CRUD, scenarioIDsource string, scenarioIDbody int) (bool, Scenario) {

	var so Scenario

	err := ValidateRole(c, ModelScenario, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation of scenario failed): %v", err))
		return false, so
	}

	if operation == Create || (operation == Read && scenarioIDsource == "none") {
		return true, so
	}

	scenarioID, err := helper.GetIDOfElement(c, "scenarioID", scenarioIDsource, scenarioIDbody)
	if err != nil {
		return false, so
	}

	userID, _ := c.Get(UserIDCtx)

	db := GetDB()
	err = db.Find(&so, uint(scenarioID)).Error
	if DBError(c, err, so) {
		return false, so
	}

	u := User{}
	err = db.Find(&u, userID.(uint)).Error
	if err != nil {
		helper.UnprocessableEntityError(c, "Access denied (user has no access or scenario is locked).")
		return false, so
	}

	if u.Role == "Admin" {
		return true, so
	}

	scenarioUser := User{}
	err = db.Order("ID asc").Model(&so).Where("ID = ?", userID.(uint)).Related(&scenarioUser, "Users").Error
	if err != nil {
		helper.UnprocessableEntityError(c, "Access denied (user has no access or scenario is locked).")
		return false, so
	}

	if !scenarioUser.Active {
		helper.UnprocessableEntityError(c, "Access denied (user has no access or scenario is locked).")
		return false, so
	} else if so.IsLocked && operation != Read {
		helper.UnprocessableEntityError(c, "Access denied (user has no access or scenario is locked).")
		return false, so
	} else {
		return true, so
	}

}

func CheckComponentConfigPermissions(c *gin.Context, operation CRUD, configIDSource string, configIDBody int) (bool, ComponentConfiguration) {

	var m ComponentConfiguration

	err := ValidateRole(c, ModelComponentConfiguration, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation of Component Configuration failed): %v", err.Error()))
		return false, m
	}

	configID, err := helper.GetIDOfElement(c, "configID", configIDSource, configIDBody)
	if err != nil {
		return false, m
	}

	db := GetDB()
	err = db.Find(&m, uint(configID)).Error
	if DBError(c, err, m) {
		return false, m
	}

	ok, _ := CheckScenarioPermissions(c, operation, "body", int(m.ScenarioID))
	if !ok {
		return false, m
	}

	return true, m

}

func CheckSignalPermissions(c *gin.Context, operation CRUD) (bool, Signal) {

	var sig Signal

	err := ValidateRole(c, ModelSignal, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation of signal failed): %v", err.Error()))
		return false, sig
	}

	signalID, err := helper.GetIDOfElement(c, "signalID", "path", -1)
	if err != nil {
		return false, sig
	}

	db := GetDB()
	err = db.Find(&sig, uint(signalID)).Error
	if DBError(c, err, sig) {
		return false, sig
	}

	ok, _ := CheckComponentConfigPermissions(c, operation, "body", int(sig.ConfigID))
	if !ok {
		return false, sig
	}

	return true, sig

}

func CheckDashboardPermissions(c *gin.Context, operation CRUD, dabIDSource string, dabIDBody int) (bool, Dashboard) {

	var dab Dashboard

	err := ValidateRole(c, ModelDashboard, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err.Error()))
		return false, dab
	}

	dabID, err := helper.GetIDOfElement(c, "dashboardID", dabIDSource, dabIDBody)
	if err != nil {
		return false, dab
	}

	db := GetDB()
	err = db.Find(&dab, uint(dabID)).Error
	if DBError(c, err, dab) {
		return false, dab
	}

	ok, _ := CheckScenarioPermissions(c, operation, "body", int(dab.ScenarioID))
	if !ok {
		return false, dab
	}

	return true, dab

}

func CheckWidgetPermissions(c *gin.Context, operation CRUD, widgetIDBody int) (bool, Widget) {

	var w Widget
	var err error
	err = ValidateRole(c, ModelWidget, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation of widget failed): %v", err.Error()))
		return false, w
	}

	var widgetID int
	if widgetIDBody < 0 {
		widgetID, err = helper.GetIDOfElement(c, "widgetID", "path", -1)
		if err != nil {
			return false, w
		}
	} else {
		widgetID = widgetIDBody
	}

	db := GetDB()
	err = db.Find(&w, uint(widgetID)).Error
	if DBError(c, err, w) {
		return false, w
	}

	ok, _ := CheckDashboardPermissions(c, operation, "body", int(w.DashboardID))
	if !ok {
		return false, w
	}

	return true, w
}

func CheckFilePermissions(c *gin.Context, operation CRUD) (bool, File) {

	var f File

	err := ValidateRole(c, ModelFile, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation of file failed): %v", err.Error()))
		return false, f
	}

	fileID, err := helper.GetIDOfElement(c, "fileID", "path", -1)
	if err != nil {
		return false, f
	}

	db := GetDB()
	err = db.Find(&f, uint(fileID)).Error
	if DBError(c, err, f) {
		return false, f
	}

	if operation != Read {
		// check access to scenario only if operation is not Read (=download) of file
		ok, _ := CheckScenarioPermissions(c, operation, "body", int(f.ScenarioID))
		if !ok {
			return false, f
		}
	}

	return true, f
}

func CheckResultPermissions(c *gin.Context, operation CRUD, resultIDSource string, resultIDBody int) (bool, Result) {

	var result Result

	err := ValidateRole(c, ModelResult, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err.Error()))
		return false, result
	}

	resultID, err := helper.GetIDOfElement(c, "resultID", resultIDSource, resultIDBody)
	if err != nil {
		return false, result
	}

	db := GetDB()
	err = db.Find(&result, uint(resultID)).Error
	if DBError(c, err, result) {
		return false, result
	}

	ok, _ := CheckScenarioPermissions(c, operation, "body", int(result.ScenarioID))
	if !ok {
		return false, result
	}

	return true, result
}

func CheckICPermissions(c *gin.Context, modeltype ModelName, operation CRUD, hasID bool) (bool, InfrastructureComponent) {

	var s InfrastructureComponent

	err := ValidateRole(c, modeltype, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation of infrastructure component failed): %v", err.Error()))
		return false, s
	}

	if hasID {
		// Get the ID of the infrastructure component from the context
		ICID, err := helper.GetIDOfElement(c, "ICID", "path", -1)
		if err != nil {
			return false, s
		}
		db := GetDB()
		err = db.Find(&s, uint(ICID)).Error
		if DBError(c, err, s) {
			return false, s
		}
	}

	return true, s
}
