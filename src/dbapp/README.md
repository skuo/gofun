# dbapp

Please reference https://astaxie.gitbooks.io/build-web-application-with-golang/en/05.2.html for the original information on Go DB interaction.

Go's database/sql provides very low level APIs that requires a new connection and statement for each db call.  A higher level API that support connection pooling and ORM is desirable.

## Good Reference on Go's database/sql package

http://go-database-sql.org/index.html

## build

From https://ieftimov.com/golang-package-multiple-binaries. 

If we see Go’s documentation on the go build command, we will find this segment:

When compiling a single main package, build writes the resulting executable to an output file named after the first source file (‘go build ed.go rx.go’ writes ‘ed’ or ‘ed.exe’) or the source code directory (‘go build unix/sam’ writes ‘sam’ or ‘sam.exe’). The ‘.exe’ suffix is added when writing a Windows executable.

Also, this:

When compiling multiple packages or a single non-main package, build compiles the packages but discards the resulting object, serving only as a check that the packages can be built.

```bash
# to build a command
cd cmd/mysql
go build # mysql executible is in cmd/mysql/

# go install can and does install all commands
go install ./...
```

## MySql Implementation

### Pull mysql docker image from registry
```bash
docker pull mysql/mysql-server:latest
```

### Start and Mount a Docker Volume for Persistence
```bash
# create a docker volume
docker volume create mysql-volume

# This sets root password, creates debdb, dbuser/dbpassword in the docker volume
docker run -d --name=mysql-server -e MYSQL_ROOT_PASSWORD=my-secret-pw -e MYSQL_DATABASE=devdb -e MYSQL_USER=dbuser -e MYSQL_PASSWORD=dbpassword --mount type=volume,src=mysql-volume,dst=/var/lib/mysql -p 3306:3306 mysql/mysql-server:latest

#Start mysql docker without changing the above values
docker run -d --name=mysql-server --mount type=volume,src=mysql-volume,dst=/var/lib/mysql -p 3306:3306 mysql/mysql-server:latest
```

### Use a Env file instead of Embedding Env Variables in Cmdline
```bash
# Use an docker.env file
docker run -d --name=mysql-server --env-file docker.env --mount type=volume,src=mysql-volume,dst=/var/lib/mysql -p 3306:3306 mysql/mysql-server:latest
```

## Sqlite Implementation

Sqlite is an embedded database. All codes needed to create a database, table and execute all DML operations are contained in cmd/sqlite/main.go.

## Postgres Implementation

### Start and Mount a Docker Volume for Persistence
```bash
# create a docker volume
docker volume create postgres-volume

# start postgres
docker run -d --name postgres-server -e POSTGRES_PASSWORD=my-secret-pw --mount type=volume,src=postgres-volume,dst=/var/lib/postgresql/data  -p 5432:5432 postgres:latest
```

### Create devuser and devdb
```bash
# attach to the running postgres docker proces
docker exec -it postgres-server psql -U postgres

# Create dev user and db
$ CREATE ROLE devuser WITH LOGIN PASSWORD 'devpassword' VALID UNTIL 'infinity';
$ CREATE DATABASE devdb WITH ENCODING='UTF8' OWNER=devuser CONNECTION LIMIT=-1;
```