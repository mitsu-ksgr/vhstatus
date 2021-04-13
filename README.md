vhstatus
========

This is a golang implementation of [cklokmose/vhstatus](https://github.com/cklokmose/vhstatus).

vhstatus provides a status page for a [valheim](https://www.valheimgame.com/) dedicated server deployed by [vhserver](https://linuxgsm.com/lgsm/vhserver/).


### Features.
- Status page (easy to customize)
	- server state
	- number of currently active player
	- the list of online and offline players
	- server information
- API Endpoint


### How to install
TODO


### Quick Test with docker-compose

```sh
$ docker-compose build
$ docker-compose up -d
```

and open http://localhost:8002/


### Development
#### Run tests
```sh
$ ./scripts/run_test.sh
```


