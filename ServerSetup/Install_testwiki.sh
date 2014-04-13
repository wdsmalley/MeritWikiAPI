#!/bin/bash
echo "Enter the version: "
read version_variable

set -x
number=$(sudo -u postgres psql -c "SELECT count(datname) FROM pg_database WHERE datname = 'wikidb' and datistemplate = false;" | (read; read; read count_tmp; read; echo "$count_tmp"))

if [ $number = "1" ]; then
	if [ $version_variable -eq 1 ]; then
		#always drop and intstall for version 1
		sudo -u postgres psql -c "drop database wikidb"
		sudo -u postgres createdb wikidb
		sudo -u postgres createuser -D -A -P postgres
		sudo -u postgres psql -d wikidb -c "create language plpgsql"
	else
		echo "Enter i to install or u to upgrade: "
		read install_variable
		if [ $install_variable  = "i" ]; then
			#this is an install, drop and recreate
			sudo -u postgres psql -c "drop database wikidb"
			sudo -u postgres createdb wikidb
			sudo -u postgres createuser -D -A -P postgres
			sudo -u postgres psql -d wikidb -c "create language plpgsql"
		fi
	fi
else
	#no database install
	sudo -u postgres createdb wikidb
	sudo -u postgres createuser -D -A -P postgres
	sudo -u postgres psql -d wikidb -c "create language plpgsql"
fi

	

if [ $version_variable -ge 1 ];then
	sudo -u postgres psql -d wikidb  -f /home/admin/SQLFiles/DBVersion1_DDL.sql
	sudo -u postgres psql -d wikidb  -f /home/admin/SQLFiles/SeedTablesForv1Tests.sql
fi	
if [ $version_variable -ge 2 ];then
	sudo -u postgres psql -d wikidb  -f /home/admin/SQLFiles/DBVersion2_DDL.sql
	sudo -u postgres psql -d wikidb  -f /home/admin/SQLFiles/SeedTablesForv2Tests.sql	
fi	


go run /usr/local/go/bin/src/Test/Tests.go

set +x
