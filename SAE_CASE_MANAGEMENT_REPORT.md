# SAE Case Management Runtime Report

## Component 1: TheHive (SAE Case Management)
1. **Upstream Identity**: TheHive Project
2. **Integrated Commit**: d390a0316ccd7780df6a49eccc773e5c28b43cc2
3. **Integrated Version**: 4.1.24-1
4. **License**: AGPL-3.0 (Archived OSS release)
5. **Runtime Dependencies**: Java 8, Elasticsearch 7.x, Cassandra 3.x (or compatible), minimum 4-6GB RAM total.
6. **Execution Goal**: Spin up TheHive alongside Cortex, receive SAE incident, generate Case, trigger analyzer.
7. **Runtime Constraints**: 
   - Public distributions of TheHive 3 and 4 have ended and Docker images for OSS branch are largely unmaintained.
   - The local environment is currently missing a running Docker Desktop Linux engine (
pipe:////./pipe/dockerDesktopLinuxEngine unavailable in the WSL host layer), making it physically impossible to provision the database (Cassandra) and search (Elasticsearch) dependencies required to run the Java applications.
8. **Final Status**: **BLOCKED** (Runtime Environment Limitation)

---

## Component 2: Cortex (SAE Threat Analysis)
1. **Upstream Identity**: TheHive Project
2. **Integrated Commit**: 9f1bc90ae92d4ba5e843439389f234e25489db1e
3. **Integrated Version**: 4.1.0-1
4. **License**: AGPL-3.0 (Archived OSS release)
5. **Runtime Dependencies**: Java 8, Elasticsearch 7.8.1 (as per docker/cortex/docker-compose.yml), Docker socket access (/var/run/docker.sock) for executing analyzer images.
6. **Execution Goal**: Process an observable pushed from a TheHive Case via SAE, return analysis taxonomy.
7. **Runtime Constraints**: 
   - Same as TheHive. The local environment's lack of Docker Engine access prevents spinning up Elasticsearch and the Cortex daemon. Furthermore, Cortex relies heavily on spawning transient Docker containers for its internal analyzers, which is fundamentally blocked without a healthy Docker socket.
8. **Final Status**: **BLOCKED** (Runtime Environment Limitation)

---

## Conclusion
The integrations for TheHive and Cortex are correctly implemented in the sae-core/internal/thehive and sae-core/internal/cortex packages. The code builds correctly and all associated unit tests pass. However, true end-to-end physical runtime verification cannot proceed until a minimum reproducible stack (Cassandra + Elasticsearch + TheHive + Cortex) can be provisioned. We are keeping the status **BLOCKED** and will not fabricate case creation or analyzer results.
