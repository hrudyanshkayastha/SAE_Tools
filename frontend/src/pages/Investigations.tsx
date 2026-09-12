import { useEffect, useState } from 'react';
import { api } from '../api';
import { Network, BrainCircuit, ShieldCheck, ArrowRight } from 'lucide-react';

export default function Investigations() {
  const [decisions, setDecisions] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.getDecisions()
      .then((res) => setDecisions(res.data || []))
      .catch((err) => console.error(err))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">AI Investigations & Decisions</h2>
      
      {loading ? (
        <div className="text-gray-400">Loading AI decisions...</div>
      ) : decisions.length === 0 ? (
        <div className="bg-gray-800 p-8 rounded-lg border border-gray-700 text-center text-gray-400">
          DATA NOT AVAILABLE. No AI decisions found in the operational store.
        </div>
      ) : (
        <div className="space-y-6">
          {decisions.map((dec, i) => (
            <div key={i} className="bg-gray-800 rounded-lg border border-gray-700 p-6 overflow-hidden relative">
              <div className="absolute top-0 right-0 p-4 opacity-10">
                <BrainCircuit className="w-32 h-32" />
              </div>
              
              <h3 className="text-lg font-mono font-bold mb-4">{dec.correlation_id}</h3>
              
              <div className="relative z-10 grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                <div className="bg-gray-900 p-4 rounded-lg border border-gray-700">
                  <div className="flex items-center gap-2 mb-2 text-gray-400">
                    <BrainCircuit className="w-4 h-4" />
                    <span className="text-sm font-bold uppercase">Risk Score</span>
                  </div>
                  <p className={`text-3xl font-bold ${dec.risk_score >= 8 ? 'text-red-500' : dec.risk_score >= 5 ? 'text-yellow-500' : 'text-green-500'}`}>
                    {dec.risk_score} / 10
                  </p>
                </div>

                <div className="bg-gray-900 p-4 rounded-lg border border-gray-700">
                  <div className="flex items-center gap-2 mb-2 text-gray-400">
                    <Network className="w-4 h-4" />
                    <span className="text-sm font-bold uppercase">LangGraph Output</span>
                  </div>
                  <p className="text-sm text-gray-200">{dec.validation}</p>
                </div>

                <div className="bg-gray-900 p-4 rounded-lg border border-gray-700">
                  <div className="flex items-center gap-2 mb-2 text-gray-400">
                    <ShieldCheck className="w-4 h-4" />
                    <span className="text-sm font-bold uppercase">Policy Engine</span>
                  </div>
                  <p className="text-sm font-medium">
                    {dec.action === 'block_ip' || dec.action === 'isolate_host' 
                      ? <span className="text-red-500">BLOCKED: Requires Human Auth</span>
                      : <span className="text-green-500">Permitted: Within Safe Bounds</span>}
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-4 bg-gray-900 p-4 rounded-lg border border-gray-700 relative z-10">
                <div className="flex-1">
                  <p className="text-xs text-gray-500 uppercase tracking-wider mb-1">Recommended Action</p>
                  <p className="text-lg font-bold text-white">{dec.action}</p>
                </div>
                <ArrowRight className="text-gray-500" />
                <div className="flex-1">
                  <p className="text-xs text-gray-500 uppercase tracking-wider mb-1">Response Router Status</p>
                  <p className="text-sm font-medium text-orange-500">BLOCKED - Orchestrator Unavailable</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
