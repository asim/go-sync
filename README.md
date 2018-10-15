# Go Sync

Go Sync is a synchronization framework for distributed systems.

## Overview

Distributed systems by their very nature are decoupled and independent. In most cases they must honour 2 out of 3 letters of the CAP theorem 
e.g Availability and Partitional tolerance but sacrificing consistency. In the case of microservices we often offload this concern to 
an external database or eventing system. Go Sync provides a framework for synchronization which can be used in the application by the developer.

We offer three primitives; Lock, Leader and KV

- Lock - distributed locking 
- Leader - leader election
- KV - distributed key value


