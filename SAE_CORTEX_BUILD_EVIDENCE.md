# SAE Cortex Build Evidence

## 1. Native Build Environment
- **Environment**: WSL2 Ubuntu ext4 ephemeral workspace
- **Dependencies**: Cortex is a Scala-based microservice that, like TheHive, strictly requires Elasticsearch as its datastore layer to persist analyzer job results and taxonomy configurations. It is intended to run as a Docker container.

## 2. Blockers
- **Infrastructure Limitation**: The Docker host environment cannot be bridged natively to the WSL subsystem (`var/run/docker.sock` is unavailable in the execution context). Running Cortex locally without Elasticsearch is fundamentally impossible.

## 3. Results
- Because SAE requires real-world Elasticsearch datastore availability to boot Cortex, the threat analysis stack could not be natively instantiated in this constrained build environment.
