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

package config

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func RegisterConfigEndpoint(r *gin.RouterGroup) {

	r.GET("", getConfig)
}

type AuthenticationExternal struct {
	Enabled      bool   `json:"enabled"`
	ProviderName string `json:"provider_name"`
	LoginURL     string `json:"authorize_url"`
}

type Authentication struct {
	External  AuthenticationExternal `json:"external"`
	LogoutURL string                 `json:"logout_url"`
}

type Contact struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
}

type Kubernetes struct {
	RancherURL  string `json:"rancher_url"`
	ClusterName string `json:"cluster_name"`
}

type WebRTC struct {
	ICEServers []ICEServer `json:"ice_servers"`
}

type ICEServer struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	Title          string         `json:"title"`
	SubTitle       string         `json:"sub_title"`
	Mode           string         `json:"mode"`
	Contact        Contact        `json:"contact"`
	Authentication Authentication `json:"authentication"`
	Kubernetes     Kubernetes     `json:"kubernetes"`
	WebRTC         WebRTC         `json:"webrtc"`
}

// getHealth godoc
// @Summary Get config VILLASweb to be used by frontend
// @ID config
// @Produce json
// @Tags config
// @Success 200 {object} config.Config "The configuration"
// @Router /config [get]
func getConfig(c *gin.Context) {

	cfg := configuration.GlobalConfig

	resp := &Config{}

	resp.Mode, _ = cfg.String("mode")
	resp.Authentication.LogoutURL, _ = cfg.String("auth.logout-url")
	resp.Authentication.External.Enabled, _ = cfg.Bool("auth.external.enabled")
	resp.Authentication.External.LoginURL, _ = cfg.String("auth.external.login-url")
	resp.Authentication.External.ProviderName, _ = cfg.String("auth.external.provider-name")
	resp.Title, _ = cfg.String("title")
	resp.SubTitle, _ = cfg.String("sub-title")
	resp.Contact.Name, _ = cfg.String("contact.name")
	resp.Contact.Mail, _ = cfg.String("contact.mail")
	resp.Kubernetes.RancherURL, _ = cfg.String("k8s.rancher-url")
	resp.Kubernetes.ClusterName, _ = cfg.String("k8s.cluster-name")

	var servers []ICEServer

	// read config param webrtc.ice-urls and transform into data structure resp.WebRTC
	var ICEurlsRaw, err = cfg.String("webrtc.ice-urls")
	if err == nil {
		ICEurls := strings.Split(ICEurlsRaw, ",")

		for _, url := range ICEurls {
			elements := strings.Split(url, "@")
			if len(elements) == 1 {
				// anonymous server, no username and password
				server := ICEServer{URL: elements[0], Username: "", Password: ""}
				servers = append(servers, server)

			} else if len(elements) == 2 {
				// server requires username and password
				// split username and password
				credentials := strings.Split(elements[0], ":")
				if len(credentials) != 2 {
					// error
					log.Println("Error parsing WebRTC ICE URLs provided in config. Please check correct format of username and password of", url)
					break
				}
				server := ICEServer{URL: elements[1], Username: credentials[0], Password: credentials[1]}
				servers = append(servers, server)

			} else {
				// error
				log.Println("Error parsing WebRTC ICE URLs provided in config. Please check correct format of", url)
				break
			}
		}
		resp.WebRTC.ICEServers = servers
	} else {
		log.Println("Error parsing WeRTC ICE URLs:", err)
	}

	c.JSON(200, resp)
}
