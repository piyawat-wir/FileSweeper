# File Sweeper

A Script written in Go for searching Victim's directories and files using keywords.

Keywords can be edited which are listed in the following files

- Directory: `keyword/Directory.txt`
- File: `keyword/File.txt`
- Extensions: `keyword/Extension.txt`

# How to use

## Parameters

```console
filesweep -host <hostname> [-port <port> -depth <2> -verbose <true>]
```

## Build and Run

From inside the source code directory, run these following commands:

```console
go build .
./filesweep -host http://localhost
```

## Run without build

From inside the source code directory, run these following commands:

```console
go run . -host http://localhost
```
