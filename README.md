# <img src="doc/pictures/villas_web.png" width=40 /> VILLASweb-backend-go

## Description
This is a rewrite in Go language of the backend for the VILLASweb
website. The term __frontend__ refers to this project, the server the
website (frontend) is talking to.  The backend serves the database
content to the website, handles authentication and persistent storage of
content. It does __not__ handle simulation data. For this have a look at
VILLASnode.

## Frameworks
The backend is build upon [gin-gonic](https://github.com/gin-gonic/gin)

## Quick start

### Database

Before running the application the user has to setup and configure
[PostgreSQL](https://www.postgresql.org/) on the local machine. 

To create a new database login to user `postgres` and start `psql`
```bash
$ su - postgres
$ psql
```
then
```sql
CREATE DATABASE villasdb ;
```

Some usefull commants for `psql`
```sql
\c somedb -- connect to a database 
\dt       -- list all tables of the database
\l        -- list all databases
```

The default `host` for postgres is `\tmp` and the ssl mode is disabled
in development. The user can change those setting in
`common/database.go`.

The current repo is tested on Fedora 29 and with PostgreSQL 11 and Go
1.11

## Copyright

2019, Institute for Automation of Complex Power Systems, EONERC  

## License

This project is released under the terms of the [GPL version 3](COPYING.md).

```
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
```

For other licensing options please consult [Prof. Antonello Monti](mailto:amonti@eonerc.rwth-aachen.de).

## Contact

[![EONERC ACS Logo](doc/pictures/eonerc_logo.png)](http://www.acs.eonerc.rwth-aachen.de)

 - Stefanos Mavros <stefanos.mavros@rwth-aachen.de>

[Institute for Automation of Complex Power Systems (ACS)](http://www.acs.eonerc.rwth-aachen.de)  
[EON Energy Research Center (EONERC)](http://www.eonerc.rwth-aachen.de)  
[RWTH University Aachen, Germany](http://www.rwth-aachen.de)  
