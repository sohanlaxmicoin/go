---
title: XDR
---

**XDR**, also known as _External Data Representation_, is used extensively in
the Rover Network, especially in the core protocol.  The ledger, transactions, results,
history, and even the messages passed between computers running rover-core
are encoded using XDR.

XDR is specified in [RFC 4506](http://tools.ietf.org/html/rfc4506.html).

Since XDR is a binary format and not known as widely as JSON for example, we try
to hide most of it from Orbit.  Instead, we opt to interpret the XDR for you
and present the values as JSON attributes.  That said, we also expose the XDR
to you so you can get access to the raw, canonical data.

In general, Orbit will encode the XDR structures in base64 so that they can be 
transmitted within a json body.  You should decode the base64 string
into a byte stream, then decode the XDR into an in-memory data structure.

## .X files

Data structures in XDR are specified in an _interface definition file_ (IDL).
The IDL files used for the Rover Network are available
[on GitHub](https://github.com/rover/rover-core/tree/master/src/xdr).
