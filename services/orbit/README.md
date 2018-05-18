# Orbit
[![Build Status](https://travis-ci.org/rover/orbit.svg?branch=master)](https://travis-ci.org/rover/orbit)

Orbit is the [client facing API](/docs) server for the Stellar ecosystem.  It acts as the interface between rover-core and applications that want to access the Stellar network. It allows you to submit transactions to the network, check the status of accounts, subscribe to event streams, etc. See [an overview of the Stellar ecosystem](https://www.rover.network/developers/guides/get-started/) for more details.

## Downloading the server
[Prebuilt binaries](https://github.com/laxmicoinofficial/go/releases) of orbit are available on the 
[releases page](https://github.com/laxmicoinofficial/go/releases).

See [the old releases page](https://github.com/rover/orbit/releases) for prior releases

| Platform       | Binary file name                                                                         |
|----------------|------------------------------------------------------------------------------------------|
| Mac OSX 64 bit | [orbit-darwin-amd64](https://github.com/laxmicoinofficial/go/releases/download/orbit-v0.12.0-testing/orbit-v0.12.0-testing-darwin-amd64.tar.gz)      |
| Linux 64 bit   | [orbit-linux-amd64](https://github.com/laxmicoinofficial/go/releases/download/orbit-v0.12.0-testing/orbit-v0.12.0-testing-linux-amd64.tar.gz)       |
| Windows 64 bit | [orbit-windows-amd64.exe](https://github.com/laxmicoinofficial/go/releases/download/orbit-v0.12.0-testing/orbit-v0.12.0-testing-windows-amd64.zip) |

Alternatively, you can [build](#building) the binary yourself.

## Dependencies

Orbit requires go 1.6 or higher to build. See (https://golang.org/doc/install) for installation instructions.

## Building

[mercurial](https://www.mercurial-scm.org/) is used during glide build.

[glide](https://glide.sh/) is used for building orbit.

Given you have a running golang installation, you can install this with:

```bash
curl https://glide.sh/get | sh
```

Next, you must download the source for packages that orbit depends upon. From within the project directory, run:

```bash
glide install
```

Then, simply run `go install github.com/rover/go/services/orbit`.  After successful
completion, you should find `orbit` is present in your `$GOPATH/bin` directory.

More detailed intructions and [admin guide](internal/docs/reference/admin.md). 

## Developing Orbit

See [the development guide](internal/docs/developing.md).

## Contributing
Please see the [CONTRIBUTING.md](./CONTRIBUTING.md) for details on how to contribute to this project.
