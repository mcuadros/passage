# Passage - SSH tunnels on steroids [![Latest Stable Version](http://img.shields.io/github/release/mcuadros/passage.svg?style=flat)](https://github.com/mcuadros/passage/releases) [![Build Status](http://img.shields.io/travis/mcuadros/passage.svg?style=flat)](https://travis-ci.org/mcuadros/passage)

<img src="https://i.imgsafe.org/c6c2d16.png" align="right" width="415" height="279px" vspace="20" />


**Passage** is a moderm SSH tunneling tool, build on Go. 

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

browse the [`releases`](https://github.com/tyba/beanstool/releases) section to see other archs and versions

License
-------

MIT, see [LICENSE](LICENSE)

