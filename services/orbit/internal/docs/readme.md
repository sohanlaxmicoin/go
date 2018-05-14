---
title: Horizon
---

Horizon is the server for the client facing API for the Stellar ecosystem.  It acts as the interface between [rover-core](https://www.stellar.org/developers/learn/rover-core) and applications that want to access the Stellar network. It allows you to submit transactions to the network, check the status of accounts, subscribe to event streams, etc. See [an overview of the Stellar ecosystem](https://www.stellar.org/developers/guides/) for more details.

You can interact directly with orbit via curl or a web browser but SDF provides a [JavaScript SDK](https://www.stellar.org/developers/js-rover-sdk/learn/) for clients to use to interact with Horizon.

SDF runs a instance of Horizon that is connected to the test net [https://orbit-testnet.stellar.org/](https://orbit-testnet.stellar.org/).

## Libraries

SDF maintained libraries:<br />
- [JavaScript](https://github.com/rover/js-rover-sdk)
- [Java](https://github.com/rover/java-rover-sdk)
- [Go](https://github.com/rover/go)

Community maintained libraries (in various states of completeness) for interacting with Horizon in other languages:<br>
- [Ruby](https://github.com/rover/ruby-rover-sdk)
- [Python](https://github.com/StellarCN/py-rover-base)
- [C# .NET 2.0](https://github.com/QuantozTechnology/csharp-rover-base)
- [C# .NET Core 2.x](https://github.com/elucidsoft/dotnetcore-rover-sdk)
- [C++](https://bitbucket.org/bnogal/stellarqore/wiki/Home)
