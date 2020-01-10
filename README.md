# s3mon - simple service status monitor

**KISS service status monitoring in your terminal.**

This tool has been developed due to a lack of fitting solution to run a basic monitoring 
dashboard on my Raspberry Pi. It is inspired by the amazing
[sampler](https://github.com/sqshq/sampler) by Alexander Lukyanchikov, but is more
lightweight with only minimalistic UI.

_**DISCLAIMER:** This tool is still under wild development!_

![s3mon demo](doc/res/s3mon.gif)

## Download as Go dependency

`go get github.com/netrixone/s3mon`

## Build

`make`

## Usage

```bash
usage: s3mon [-h|--help] [-v|--version] [-V|--verbose] [-c|--config "<value>"]

             s3mon - simple service status monitor v1.1 by stuchl4n3k

Arguments:

  -h  --help     Print help information
  -v  --version  Print version and exit
  -V  --verbose  Be more verbose
  -c  --config   Config file

```
