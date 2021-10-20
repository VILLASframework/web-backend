/** User package, methods.
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

import "fmt"

type UsernameAlreadyTaken struct {
	Username string
}

func (e *UsernameAlreadyTaken) Error() string {
	return fmt.Sprintf("username is already taken: %s", e.Username)
}

type ForbiddenError struct {
	Reason string
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("permission denied: %s", e.Reason)
}
