# THIRD-PARTY NOTICES

SAE incorporates and integrates with various third-party open-source software. This file contains the applicable attribution and license notices for these components.

## 1. Integrated Security Providers

The SAE Security Fabric integrates with the following third-party sensors, orchestration tools, and runtimes:

*   **Wazuh** (GPLv2) - Used via the Wazuh adapter.
*   **Zeek** (BSD 3-Clause) - Used via the Zeek adapter.
*   **Suricata** (GPLv2) - Used via the Suricata adapter.
*   **Falco** (Apache 2.0) - Used via the Falco adapter.
*   **KubeArmor** (Apache 2.0) - Used via the KubeArmor adapter.
*   **Trivy** (Apache 2.0) - Used via the Trivyligence adapter.
*   **ScoutSuite** (GPLv2) - Used via the ScoutSuite adapter.
*   **Shuffle** (MIT) - Used via the Shuffle adapter.
*   **TheHive** (AGPL 3.0) - Used via the TheHive adapter.
*   **Cortex** (AGPL 3.0) - Used via the Cortex adapter.
*   **Ollama** (MIT) - Used as the local execution runtime for SAE AI reasoning.
*   **LangGraph** (MIT) - Used as the stateful orchestrator for SAE Investigation.
*   **Garak** (Apache 2.0) - Used via the Garak adapter.

## 2. Infrastructure & Libraries

**Go Backend Dependencies:**
*   `github.com/redis/go-redis/v9` - License: BSD 2-Clause
*   `github.com/lib/pq` - License: MIT
*   `github.com/google/uuid` - License: BSD 3-Clause

**Frontend Dependencies:**
*   `react` & `react-dom` - License: MIT (Copyright Meta Platforms, Inc.)
*   `vite` - License: MIT (Copyright Evan You)
*   `tailwindcss` - License: MIT (Copyright Tailwind Labs, Inc.)
*   `lucide-react` - License: ISC (Copyright Lucide Contributors)
*   `axios` - License: MIT

## 3. Disclaimers

The above trademarks, logos, and project names are the property of their respective owners. SAE does not claim ownership over any third-party software, libraries, or sensors. They remain independently licensed under their respective open-source agreements. SAE only provides proprietary telemetry fabric, correlation, policy engine, and user interfaces on top of these providers.
