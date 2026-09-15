import { useEffect, useState } from 'react';
import { api } from '../api';

export default function Events() {
  const [events, setEvents] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.getEvents()
      .then((res) => {
        setEvents(res.data || []);
      })
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
      <h2 className="text-2xl font-bold">Live Events (OCSF)</h2>
      
      {loading ? (
        <div className="text-gray-400">Loading events...</div>
      ) : events.length === 0 ? (
        <div className="bg-gray-800 p-8 rounded-lg border border-gray-700 text-center text-gray-400">
          DATA NOT AVAILABLE. Backend returned empty events list.
        </div>
      ) : (
        <div className="bg-gray-800 rounded-lg border border-gray-700 overflow-hidden">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-gray-900 border-b border-gray-700">
                <th className="p-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Timestamp</th>
                <th className="p-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Activity</th>
                <th className="p-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Severity</th>
                <th className="p-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Message</th>
                <th className="p-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Correlation ID</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-700">
              {events.map((ev, i) => (
                <tr key={ev.event_id || i} className="hover:bg-gray-750 transition-colors">
                  <td className="p-4 text-sm font-mono text-gray-300">
                    {new Date(ev.time).toLocaleString() !== 'Invalid Date' ? new Date(ev.time).toLocaleString() : 'N/A'}
                  </td>
                  <td className="p-4 text-sm font-medium">{ev.activity_name || 'Unknown'}</td>
                  <td className="p-4">
                    <span className={`px-2 py-1 text-xs font-bold rounded border ${getSeverityColor(ev.severity)}`}>
                      {ev.severity || 'Unknown'}
                    </span>
                  </td>
                  <td className="p-4 text-sm text-gray-300 truncate max-w-xs" title={ev.message}>
                    {ev.message || '-'}
                  </td>
                  <td className="p-4 text-sm font-mono text-gray-400">
                    {ev.correlation_id || '-'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
