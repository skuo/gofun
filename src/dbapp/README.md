# dbapp

Please reference https://astaxie.gitbooks.io/build-web-application-with-golang/en/05.2.html for the original information on Go DB interaction.

Go's database/sql provides very low level APIs that requires a new connection and statement for each db call.  A higher level API that support connection pooling and ORM is desirable.

## build

```bash
# to build all cmds
go build ./...
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
