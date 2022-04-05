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

package infrastructure_component

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func QueryICAPIs(d time.Duration) {

	go func() {

		for range time.Tick(d) {
			//log.Println("Querying IC APIs at time:", x)
			var err error

			db := database.GetDB()
			var ics []database.InfrastructureComponent
			err = db.Order("ID asc").Find(&ics).Error
			if err != nil {
				log.Println("Error getting ICs from DB:", err.Error())
				continue
			}

			// iterate over ICs in DB
			for _, ic := range ics {
				err := queryIC(&ic)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
}

func queryIC(ic *database.InfrastructureComponent) error {
	if ic.ManagedExternally || ic.APIURL == "" || (!strings.HasPrefix(ic.APIURL, "http://") && !strings.HasPrefix(ic.APIURL, "https://")) {
		return nil
	}

	if ic.Category == "gateway" {
		if ic.Type == "villas-node" {
			err := queryVillasNodeGateway(ic)
			if err != nil {
				return err
			}
		} else if ic.Type == "villas-relay" {
			err := queryVillasRelayGateway(ic)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func queryVillasNodeGateway(ic *database.InfrastructureComponent) error {
	client := resty.New()

	log.Println("External API: checking for villas-node gateway", ic.Name)
	statusResponse, err := client.R().SetHeader("Accept", "application/json").Get(ic.APIURL + "/status")
	if err != nil {
		return fmt.Errorf("failed to query the status of %s: %w", ic.Name, err)
	}
	var status map[string]interface{}
	err = json.Unmarshal(statusResponse.Body(), &status)
	if err != nil {
		return fmt.Errorf("failed to unmarshal status of %s: %w", ic.Name, err)
	}

	parts := strings.Split(ic.WebsocketURL, "/")
	if len(parts) > 0 && parts[len(parts)-1] != "" {

		configResponse, _ := client.R().SetHeader("Accept", "application/json").Get(ic.APIURL + "/node/" + parts[len(parts)-1])
		statsResponse, _ := client.R().SetHeader("Accept", "application/json").Get(ic.APIURL + "/node/" + parts[len(parts)-1] + "/stats")

		var config map[string]interface{}
		err = json.Unmarshal(configResponse.Body(), &config)
		if err == nil {
			status["config"] = config
		}
		var stats map[string]interface{}
		err = json.Unmarshal(statsResponse.Body(), &stats)
		if err == nil {
			status["statistics"] = stats
		}
	}

	var updatedIC UpdateICRequest
	statusRaw, _ := json.Marshal(status)
	updatedIC.InfrastructureComponent.StatusUpdateRaw = postgres.Jsonb{RawMessage: statusRaw}
	updatedIC.InfrastructureComponent.State = fmt.Sprintf("%v", status["state"])
	updatedIC.InfrastructureComponent.UUID = fmt.Sprintf("%v", status["uuid"])
	timeNow, err := strconv.ParseFloat(fmt.Sprintf("%v", status["time_now"]), 64)
	if err != nil {
		return fmt.Errorf("failed to parse time_now to float: %w", err)
	}
	timeStarted, err := strconv.ParseFloat(fmt.Sprintf("%v", status["time_started"]), 64)
	if err != nil {
		return fmt.Errorf("failed to parse time_started to float: %w", err)
	}
	uptime := timeNow - timeStarted
	updatedIC.InfrastructureComponent.Uptime = uptime

	// validate the update
	err = updatedIC.validate()
	if err != nil {
		return fmt.Errorf("failed to validate updated villas-node gateway: %s (%s): %w", ic.Name, ic.UUID, err)
	}

	// create the update and update IC in DB
	var x InfrastructureComponent
	err = x.ByID(ic.ID)
	if err != nil {
		return fmt.Errorf("failed to get villas-node gateway by ID %s (%s): %w", ic.Name, ic.UUID, err)
	}
	u := updatedIC.updatedIC(x)
	err = x.update(u)
	if err != nil {
		return fmt.Errorf("failed to update villas-node gateway %s (%s): %w", ic.Name, ic.UUID, err)
	}

	return nil
}

func queryVillasRelayGateway(ic *database.InfrastructureComponent) error {
	client := resty.New()

	log.Println("External API: checking for villas-relay manager", ic.Name)
	statusResponse, err := client.R().SetHeader("Accept", "application/json").Get(ic.APIURL)
	if err != nil {
		return fmt.Errorf("failed querying API of %s (%s): %w", ic.Name, ic.UUID, err)
	}

	var status map[string]interface{}
	err = json.Unmarshal(statusResponse.Body(), &status)
	if err != nil {
		return fmt.Errorf("failed to unmarshal status villas-relay manager %s (%s): %w", ic.Name, ic.UUID, err)
	}

	var updatedIC UpdateICRequest
	statusRaw, _ := json.Marshal(status)
	updatedIC.InfrastructureComponent.StatusUpdateRaw = postgres.Jsonb{RawMessage: statusRaw}
	updatedIC.InfrastructureComponent.UUID = fmt.Sprintf("%v", status["uuid"])

	// validate the update
	err = updatedIC.validate()
	if err != nil {
		return fmt.Errorf("failed to validate updated villas-relay manager %s (%s): %w", ic.Name, ic.UUID, err)
	}

	// create the update and update IC in DB
	var x InfrastructureComponent
	err = x.ByID(ic.ID)
	if err != nil {
		return fmt.Errorf("failed to get villas-relay manager by ID %s (%s): %w", ic.Name, ic.UUID, err)
	}
	u := updatedIC.updatedIC(x)
	err = x.update(u)
	if err != nil {
		return fmt.Errorf("failed to update villas-relay manager %s (%s): %w", ic.Name, ic.UUID, err)
	}

	return nil
}
