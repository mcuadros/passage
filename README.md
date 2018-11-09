# Passage - SSH tunnels on steroids [![Latest Stable Version](http://img.shields.io/github/release/mcuadros/passage.svg?style=flat)](https://github.com/mcuadros/passage/releases) [![Build Status](http://img.shields.io/travis/mcuadros/passage.svg?style=flat)](https://travis-ci.org/mcuadros/passage)

<img src="https://i.imgsafe.org/c6c2d16.png" align="right" width="415" height="279px" vspace="20" />


**Passage** is a modern SSH tunneling tool, build on Go. 

Passage server allows you to have multiple SSH tunnels over the same SSH connection or several SSH connections, also allows to make easy connections to _Docker_ containers without exported ports, along other features.

#### ok but what for is a SSH tunnel ...

A [SSH tunnel](https://en.wikipedia.org/wiki/Tunneling_protocol#Secure_Shell_tunneling) forward a specified local port to a port on the remote machine with the only requirement of having a working SSH connection. 

This can be done easily using a basic command on any *nix machine like `ssh example.com -L 80:localhost:80 -N`

Installation
------------

```
wget https://github.com/mcuadros/passage/releases/download/v0.1.0/passage_v0.1.0_linux_amd64.tar.gz
tar -xvzf passage_v0.1.0_linux_amd64.tar.gz
cp passage_v0.1.0_linux_amd64/passage /usr/local/bin/
```

browse the [`releases`](https://github.com/mcuadros/passage/releases) section to see other archs and versions

Usage
-----

## Running the server

Passage runs as daemon, the command to run it is: 
```
passage server
```
It loads by default the config file from `$HOME/.passage.yaml` ([config file reference](#config)). 

All the connections to the SSH server are done it lazy mode, this means that until you open a connection to a `passage` the connection to the SSH server is close.

## Quering the local address of a passage

Using the command `passage get <passage-name>` you can retrieve the local address for this passage. 

As an example if we have a passage server running with the following configuration: 
```yaml
servers:
  example-server:
    address: your-ssh-server.com:22
    passages:
      nginx:
        address: 127.0.0.1:80
``` 

You can wget the nginx server using the command:
```sh
wget $(passage get nginx)
``` 

Config <a name="config" />
------

The format of the config file is `yaml` and the structure is as follows:

```yaml
servers:                     # [multiple] SSH servers you can have as many as you want
  <server-name>:             # [mandatory] name of the server to connect
    address: <host:port>     # [mandatory] address and port of the SSH server
    user: <username>         # [optional] the SSH username, by default $USER is used
    retries: <int>           # [optional] number of reconnect retries if a connection fails.
    passages:                # [multiple] you can many different passage over the same SSH connection
      <passage-name>:        # [mandatory] name of the passage, the name provided to the `get` 
        address: <host:port> # [mandatory]address and port of the local service, the address can be 
                             # a localhost server or a remote one, remember this is an internal
                             # connection, so you don't need a port rechable from outside
        local: <host:port>   # [optional] address where the passage will be listening (eg `:8080`)
                             # if empty a random port will be assigned
```


License
-------

MIT, see [LICENSE](LICENSE)
