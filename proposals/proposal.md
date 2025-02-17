# HA Bridge

A basic grpc bridge for Home Assisant. Supports initially
device state look up for other purposes

## Tasks

1. Build basic server that does nothing
1. Load token from secret on startup
1. fleet_infra: load habridge into cluster
1. Write proto for supporting getting device state - focus on switches
