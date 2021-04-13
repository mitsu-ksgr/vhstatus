Notes
=====

vhstatus structured by following parts.

- web server
- log watcher
- data store


### web server
- src: `internal/web`
- a web server that provides the status of vhserver.
- this part get status of vhserver from the data store, and then deliever it.


### log watcher
- src: `internal/vhlogwatcher`
- monitor logs and collect data.
	- at first, read logs that already created (`vhserver-console-YYYY-MM-DD-hh:mm:ss.log`).
	- and then monitor the current log file (`vhserver-console.log`).


### data store
- src: `internal/vhstatus`
- store the status of the vhserver.


