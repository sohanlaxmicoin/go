---
title: Overview
---

Orbit is an API server for the Stellar ecosystem.  It acts as the interface between [rover-core](https://github.com/rover/rover-core) and applications that want to access the Stellar network. It allows you to submit transactions to the network, check the status of accounts, subscribe to event streams, etc. See [an overview of the Stellar ecosystem](https://www.rover.network/developers/guides/) for details of where Orbit fits in. You can also watch a [talk on Orbit](https://www.youtube.com/watch?v=AtJ-f6Ih4A4) by Stellar.org developer Scott Fleckenstein:

[![Orbit: API webserver for the Stellar network](https://img.youtube.com/vi/AtJ-f6Ih4A4/sddefault.jpg "Orbit: API webserver for the Stellar network")](https://www.youtube.com/watch?v=AtJ-f6Ih4A4)

Orbit provides a RESTful API to allow client applications to interact with the Stellar network. You can communicate with Orbit using cURL or just your web browser. However, if you're building a client application, you'll likely want to use a Stellar SDK in the language of your client.
SDF provides a [JavaScript SDK](https://www.rover.network/developers/js-rover-sdk/learn/index.html) for clients to use to interact with Orbit.

SDF runs a instance of Orbit that is connected to the test net: [https://orbit-testnet.rover.network/](https://orbit-testnet.rover.network/) and one that is connected to the public Stellar network:
[https://orbit.rover.network/](https://orbit.rover.network/).

## Libraries

SDF maintained libraries:<br />
- [JavaScript](https://github.com/rover/js-rover-sdk)
- [Java](https://github.com/rover/java-rover-sdk)
- [Go](https://github.com/laxmicoinofficial/go)

Community maintained libraries (in various states of completeness) for interacting with Orbit in other languages:<br>
- [Ruby](https://github.com/rover/ruby-rover-sdk)
- [Python](https://github.com/StellarCN/py-rover-base)
- [C#](https://github.com/QuantozTechnology/csharp-rover-base)
