<p align="center">
  <img src="/html/static/wadoh.png" width="200" align="center"/>
</p>
<h2 align="center">Wadoh</h2>

**Wadoh** is WhatsApp Web API powered by [whatsmeow](https://github.com/tulir/whatsmeow). It provides API from
connected device that can be integrated with other app or system. This repo is just the Wadoh web itself, the
devices are managed in [wadoh-be](https://github.com/9d4/wadoh) through gRPC.

This supposed to be like:
- Multi user
- Multi devices (each user can have many devices)
- Human like action (typing and online before send message)
- Providing API for devices
- Providing webhook for received messages

**Disclaimer:** This is not official WhatsApp API, it's not guaranteed you will not be blocked by using this.
### Running
Frontend is using parcel to build css and js, it's located in [html/](/html). During development we may run
```sh
$ npm run watch
```

to build and run the web server, add `dev` tag so it will use os FS instead of embed FS.
```sh
$ go run -tags dev .
```
#### Configuration
Wadoh supports yaml, json, and environment variable. See yaml [example](./wadoh.yml). Environment variable should be prefixed
with `WADOH_` and use `__` for children.
```yml
log_level: -1 # see zerolog.Level for reference: https://github.com/rs/zerolog/?tab=readme-ov-file#leveled-logging
http:
  address: 0.0.0.0:8080 # or WADOH_HTTP__ADDRESS=127.0.0.1:8989
  jwt_secret: abc
storage:
  provider: mysql
  dsn: root:@tcp(localhost:3306)/wadoh?parseTime=true
```
Custom config path is also supported, use `-c` or `--config` to use custom config path.
