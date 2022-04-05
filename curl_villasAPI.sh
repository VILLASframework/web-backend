#!/bin/bash
##########################################################################################
# This file is part of VILLASweb-backend-go
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
##########################################################################################

# Handy bash script using curl for testing the checkpoint of the
# VILLASweb frontend. See at the end of the file for usage cases

port='4000'
version='v2'
apiBase='localhost:'"$port"'/api/'"$version"

login () { 
    printf "> POST "$apiBase"/authenticate\n"
    curl "$apiBase"/authenticate -s \
        -H "Contet-Type: application/json" \
        -X POST \
        --data "$1" | jq -r '.token' > auth.jwt \
        && printf '\n' 
}

create_user () { 
    printf "> POST "$apiBase"/users to create newUser\n"
    curl "$apiBase"/users -s \
        -H "Contet-Type: application/json" \
        -H "Authorization: Bearer $(< auth.jwt)" \
        -X POST \
        --data "$1" | jq -r '.' && printf '\n'
}

read_users () {
    printf "> GET "$apiBase"/users\n"
    curl "$apiBase"/users -s \
        -H "Contet-Type: application/json" \
        -H "Authorization: Bearer $(< auth.jwt)" \
        -X GET | jq '.' && printf '\n'
}

read_infrastructure_components () {
    printf "> GET "$apiBase"/ic\n"
    curl "$apiBase"/ic -s \
        -H "Contet-Type: application/json" \
        -H "Authorization: Bearer $(< auth.jwt)" \
        -X GET | jq '.' && printf '\n'
}

read_user () {
    printf "> GET "$apiBase"/users/$1\n"
    curl "$apiBase"/users/$1 -s \
        -H "Contet-Type: application/json" \
        -H "Authorization: Bearer $(< auth.jwt)" \
        -X GET | jq '.' && printf '\n'
}

update_user () { 
    printf "> PUT "$apiBase"/users to update newUser\n"
    curl "$apiBase"/users/$1 -s \
        -H "Contet-Type: application/json" \
        -H "Authorization: Bearer $(< auth.jwt)" \
        -X PUT \
        --data "$2" | jq -r '.' && printf '\n'
}

delete_user () {
    printf "> DELETE "$apiBase"/users/$1\n"
    curl "$apiBase"/users/$1 -s \
        -H "Contet-Type: application/json" \
        -H "Authorization: Bearer $(< auth.jwt)" \
        -X DELETE | jq '.' && printf '\n'
}


admin='{ "Username" : "User_0", 
            "Password" : "xyz789", 
            "Role" : "Admin" }'
userA='{ "Username" : "User_A", 
            "Password" : "abc123",
            "Mail" : "m@i.l" }'
newUserW='{ "user": { "Username" : "User_W", 
            "Password" : "www747_tst", 
            "Role" : "User",
            "Mail" : "m@i.l" } }'
updUserW='{ "user": { "Mail" : "lalala", 
            "Role" : "Admin" } }'
#updUserW='{ "user": { "Username" : "User_Wupdated", 
            #"Password" : "pie314_test", 
            #"Role" : "Admin",
            #"Mail" : "NEW_m@i.l" } }'
userC='{ "user": { "Username" : "User_C", 
            "Password" : "abc123",
            "Mail" : "C_m@i.l",
            "Role" : "User"} }'

login "$admin"
#create_user "$userC"
#read_users
#read_user 1
#read_infrastructure_components
create_user "$newUserW"
#read_users
read_user 4
#login "$newUserW"
update_user 4 "$updUserW"
#login "$admin"
read_user 4
#login "$admin"
#read_user 4
#login "$updUserW"
#delete_user 2

