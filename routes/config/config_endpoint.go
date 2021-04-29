/** Healthz package, endpoints.
*
* @author Steffen Vogel <svogel2@eonerc.rwth-aachen.de>
* @copyright 2014-2021, Institute for Automation of Complex Power Systems, EONERC
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
package config

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"github.com/gin-gonic/gin"
)

func RegisterConfigEndpoint(r *gin.RouterGroup) {

	r.GET("", getConfig)
}

type ConfigAuthenticationExternal struct {
	Enabled      bool   `json:"enabled"`
	ProviderName string `json:"provider_name"`
	LoginURL     string `json:"login_url"`
}

type ConfigAuthentication struct {
	External  ConfigAuthenticationExternal `json:"external"`
	LogoutURL string                       `json:"logout_url"`
}

type ConfigContact struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
}

type Config struct {
	Title          string               `json:"title"`
	SubTitle       string               `json:"sub_title"`
	Mode           string               `json:"mode"`
	Contact        ConfigContact        `json:"contact"`
	Authentication ConfigAuthentication `json:"authentication"`
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

	c.JSON(200, resp)
}
