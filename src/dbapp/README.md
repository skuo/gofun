# dbapp

Please reference https://astaxie.gitbooks.io/build-web-application-with-golang/en/05.2.html for the original information on Go DB interaction.

Go's database/sql provides very low level APIs that requires a new connection and statement for each db call.  A higher level API that support connection pooling and ORM is desirable.

## build

```bash
# to build all cmds
> go build ./...
```
