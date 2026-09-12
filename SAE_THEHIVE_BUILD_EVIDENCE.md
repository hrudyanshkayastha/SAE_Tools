# SAE TheHive Build Evidence

## 1. Native Build Environment
- **Environment**: WSL2 Ubuntu ext4 ephemeral workspace
- **Dependencies**: TheHive relies heavily on Scala/sbt, Cassandra/ScyllaDB (for storage), and Elasticsearch (for indexing) usually orchestrated via Docker Compose for production deployments.

## 2. Blockers
- **Infrastructure Limitation**: The underlying Docker host (required to spin up Cassandra/ES for TheHive) is detached from the WSL subsystem (no `/var/run/docker.sock`). While the `sbt` application could technically be compiled natively, a functioning case management system requires its heavy backing datastores to be active.

## 3. Results
- Because SAE requires real-world Docker datastore availability to boot TheHive, the case management stack could not be natively instantiated here.
