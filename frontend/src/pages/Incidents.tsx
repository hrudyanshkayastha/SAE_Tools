import { useEffect, useState } from 'react';
import { api } from '../api';

export default function Incidents() {
  const [incidents, setIncidents] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.getIncidents()
      .then((res) => setIncidents(res.data || []))
      .catch((err) => console.error(err))
      .finally(() => setLoading(false));
  }, []);

  const getSeverityColor = (sev: string) => {
    switch(sev?.toLowerCase()) {
      case 'critical': return 'bg-red-500/20 text-red-500 border-red-500/30';
      case 'high': return 'bg-orange-500/20 text-orange-500 border-orange-500/30';
      case 'medium': return 'bg-yellow-500/20 text-yellow-500 border-yellow-500/30';
      default: return 'bg-blue-500/20 text-blue-500 border-blue-500/30';
    }
  };

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold">Correlated Incidents</h2>
      
      {loading ? (
        <div className="text-gray-400">Loading incidents...</div>
      ) : incidents.length === 0 ? (
        <div className="bg-gray-800 p-8 rounded-lg border border-gray-700 text-center text-gray-400">
          DATA NOT AVAILABLE. No active correlations found.
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4">
          {incidents.map((inc, i) => (
            <div key={inc.correlation_id || i} className="bg-gray-800 rounded-lg border border-gray-700 p-5 shadow-lg">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h3 className="text-lg font-mono font-bold text-white mb-1">{inc.correlation_id}</h3>
                  <p className="text-sm text-gray-400">Target/Observable: <span className="font-mono text-gray-200">{inc.target}</span></p>
                </div>
                <span className={`px-2 py-1 text-xs font-bold rounded border ${getSeverityColor(inc.severity)}`}>
                  {inc.severity || 'Unknown'}
                </span>
              </div>
              <div className="grid grid-cols-2 gap-4 mt-4">
                <div className="bg-gray-900 p-3 rounded border border-gray-700">
                  <p className="text-xs text-gray-500 uppercase tracking-wider mb-1">Event Count</p>
                  <p className="text-xl font-bold">{inc.events_count}</p>
                </div>
                <div className="bg-gray-900 p-3 rounded border border-gray-700">
                  <p className="text-xs text-gray-500 uppercase tracking-wider mb-1">Status</p>
                  <p className="text-sm font-medium">{inc.status}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
