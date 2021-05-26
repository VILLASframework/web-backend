/** infrastructure-component package, API queries.
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
	"encoding/json"
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"strconv"
	"strings"
	"time"
)

func QueryICAPIs(d time.Duration) {

	client := resty.New()
	//client.SetDebug(true)

	go func() {

		for _ = range time.Tick(d) {
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

				if ic.ManagedExternally {
					continue
				}

				if ic.APIURL == "" || (!strings.HasPrefix(ic.APIURL, "http://") && !strings.HasPrefix(ic.APIURL, "https://")) {
					continue
				}

				if ic.Category == "gateway" && ic.Type == "villas-node" {

					log.Println("External API: checking for villas-node gateway", ic.Name)
					statusResponse, err := client.R().SetHeader("Accept", "application/json").Get(ic.APIURL + "/status")
					if err != nil {
						log.Println("Error querying status of", ic.Name, err)
						continue
					}
					var status map[string]interface{}
					err = json.Unmarshal(statusResponse.Body(), &status)
					if err != nil {
						log.Println("Error unmarshalling status of", ic.Name, err)
						continue
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
					timeNow, myerr := strconv.ParseFloat(fmt.Sprintf("%v", status["time_now"]), 64)
					if myerr != nil {
						log.Println("Error parsing time_now to float", myerr.Error())
						continue
					}
					timeStarted, myerr := strconv.ParseFloat(fmt.Sprintf("%v", status["time_started"]), 64)
					if myerr != nil {
						log.Println("Error parsing time_started to float", myerr.Error())
						continue
					}
					uptime := timeNow - timeStarted
					updatedIC.InfrastructureComponent.Uptime = uptime

					// validate the update
					err = updatedIC.validate()
					if err != nil {
						log.Println("Error validating updated villas-node gateway", ic.Name, ic.UUID, err.Error())
						continue
					}

					// create the update and update IC in DB
					var x InfrastructureComponent
					err = x.byID(ic.ID)
					if err != nil {
						log.Println("Error getting villas-node gateway by ID", ic.Name, err)
						continue
					}
					u := updatedIC.updatedIC(x)
					err = x.update(u)
					if err != nil {
						log.Println("Error updating villas-node gateway", ic.Name, ic.UUID, err.Error())
						continue
					}

				} else if ic.Category == "manager" && ic.Type == "villas-relay" {

					log.Println("External API: checking for villas-relay manager", ic.Name)
					statusResponse, err := client.R().SetHeader("Accept", "application/json").Get(ic.APIURL)
					if err != nil {
						log.Println("Error querying API of", ic.Name, err)
						continue
					}
					var status map[string]interface{}
					err = json.Unmarshal(statusResponse.Body(), &status)
					if err != nil {
						log.Println("Error unmarshalling status villas-relay manager", ic.Name, err)
						continue
					}

					var updatedIC UpdateICRequest
					statusRaw, _ := json.Marshal(status)
					updatedIC.InfrastructureComponent.StatusUpdateRaw = postgres.Jsonb{RawMessage: statusRaw}
					updatedIC.InfrastructureComponent.UUID = fmt.Sprintf("%v", status["uuid"])

					// validate the update
					err = updatedIC.validate()
					if err != nil {
						log.Println("Error validating updated villas-relay manager", ic.Name, ic.UUID, err.Error())
						continue
					}

					// create the update and update IC in DB
					var x InfrastructureComponent
					err = x.byID(ic.ID)
					if err != nil {
						log.Println("Error getting villas-relay manager by ID", ic.Name, err)
						continue
					}
					u := updatedIC.updatedIC(x)
					err = x.update(u)
					if err != nil {
						log.Println("Error updating villas-relay manager", ic.Name, ic.UUID, err.Error())
						continue
					}

				} else if ic.Category == "gateway" && ic.Type == "villas-relay" {

					// TODO add code here once API for VILLASrelay sessions is available

				}
			}
		}
	}()
}
