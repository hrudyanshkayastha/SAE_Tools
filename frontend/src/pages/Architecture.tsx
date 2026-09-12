import { CheckCircle2, XCircle } from 'lucide-react';

export default function Architecture() {
  const tools = [
    { name: 'Wazuh', cap: 'Endpoint Security', status: 'FULLY VERIFIED', icon: <CheckCircle2 className="w-5 h-5 text-green-500" /> },
    { name: 'Zeek', cap: 'Network Intelligence', status: 'FULLY VERIFIED', icon: <CheckCircle2 className="w-5 h-5 text-green-500" /> },
    { name: 'Suricata', cap: 'Network Detection', status: 'FULLY VERIFIED', icon: <CheckCircle2 className="w-5 h-5 text-green-500" /> },
    { name: 'Falco', cap: 'Container Security', status: 'PARTIALLY VERIFIED', icon: <CheckCircle2 className="w-5 h-5 text-yellow-500" /> },
    { name: 'KubeArmor', cap: 'CloudGuard', status: 'PARTIALLY VERIFIED', icon: <CheckCircle2 className="w-5 h-5 text-yellow-500" /> },
    { name: 'Trivy', cap: 'Vulnerability Intel', status: 'FULLY VERIFIED', icon: <CheckCircle2 className="w-5 h-5 text-green-500" /> },
    { name: 'ScoutSuite', cap: 'Cloud Assessment', status: 'PARTIALLY VERIFIED', icon: <CheckCircle2 className="w-5 h-5 text-yellow-500" /> },
    { name: 'Shuffle', cap: 'Response Orchestrator', status: 'RUNTIME BLOCKED', icon: <XCircle className="w-5 h-5 text-red-500" /> },
    { name: 'TheHive', cap: 'Case Management', status: 'RUNTIME BLOCKED', icon: <XCircle className="w-5 h-5 text-red-500" /> },
    { name: 'Cortex', cap: 'Threat Analysis', status: 'RUNTIME BLOCKED', icon: <XCircle className="w-5 h-5 text-red-500" /> },
    { name: 'Ollama', cap: 'Local AI Runtime', status: 'FULLY VERIFIED', icon: <CheckCircle2 className="w-5 h-5 text-green-500" /> },
    { name: 'LangGraph', cap: 'AI Reasoning Engine', status: 'FULLY VERIFIED', icon: <CheckCircle2 className="w-5 h-5 text-green-500" /> },
    { name: 'Garak', cap: 'AI Security Testing', status: 'FULLY VERIFIED', icon: <CheckCircle2 className="w-5 h-5 text-green-500" /> },
  ];

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">SAE Unified Architecture</h2>
      
      <div className="bg-gray-800 p-6 rounded-lg border border-gray-700">
        <h3 className="text-lg font-bold mb-4">Pipeline Flow</h3>
        <div className="font-mono text-sm text-gray-300 whitespace-pre overflow-x-auto p-4 bg-gray-900 rounded border border-gray-700">
{`          ┌──────────────────────────┐
          │     13 SECURITY TOOLS    │
          └────────────┬─────────────┘
                       ↓
                 OCSF NORMALIZER
                       ↓
                  REDIS STREAM
                       ↓
                 POSTGRESQL
                       ↓
              CORRELATION ENGINE
                       ↓
                LANGGRAPH
                       ↓
                  OLLAMA
                       ↓
                 RISK ENGINE
                       ↓
                POLICY ENGINE
                       ↓
                  DECISION`}
        </div>
      </div>

      <div className="bg-gray-800 p-6 rounded-lg border border-gray-700">
        <h3 className="text-lg font-bold mb-4">Capabilities Matrix</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {tools.map((t, i) => (
            <div key={i} className="bg-gray-900 p-4 rounded border border-gray-700 flex items-center justify-between">
              <div>
                <p className="font-bold text-gray-200">{t.name}</p>
                <p className="text-xs text-gray-500">{t.cap}</p>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-xs font-mono">{t.status}</span>
                {t.icon}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
