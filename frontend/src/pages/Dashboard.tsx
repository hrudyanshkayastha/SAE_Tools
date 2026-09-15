import { useEffect, useState } from 'react';
import { api } from '../api';
import { Shield, Server, Activity, AlertTriangle, Play, Pause } from 'lucide-react';

export default function Dashboard() {
  const [stats, setStats] = useState({ events: 0, incidents: 0, decisions: 0 });
  const [health, setHealth] = useState<any>(null);
  const [ready, setReady] = useState<any>(null);
  const [polling, setPolling] = useState(true);

  const loadData = async () => {
    try {
      const [ev, inc, dec, h, r] = await Promise.all([
        api.getEvents().catch(() => ({ data: [] })),
        api.getIncidents().catch(() => ({ data: [] })),
        api.getDecisions().catch(() => ({ data: [] })),
        api.getHealth().catch((e) => ({ data: e.response?.data || { status: 'error' } })),
        api.getReady().catch((e) => ({ data: e.response?.data || { status: 'error' } }))
      ]);
      setStats({
        events: ev.data?.length || 0,
        incidents: inc.data?.length || 0,
        decisions: dec.data?.length || 0,
      });
      setHealth(h.data);
      setReady(r.data);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    loadData();
    if (!polling) return;
    const interval = setInterval(loadData, 5000);
    return () => clearInterval(interval);
  }, [polling]);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">SAE Overview</h2>
        <button 
          onClick={() => setPolling(!polling)}
          className={`flex items-center gap-2 px-3 py-1.5 rounded-md text-sm font-medium ${polling ? 'bg-gray-700 hover:bg-gray-600 text-white' : 'bg-blue-600 hover:bg-blue-500 text-white'}`}
        >
          {polling ? <Pause className="w-4 h-4" /> : <Play className="w-4 h-4" />}
          {polling ? 'Auto-refresh On' : 'Auto-refresh Off'}
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-gray-800 p-6 rounded-lg border border-gray-700 flex items-center gap-4">
          <div className="p-3 bg-blue-500/20 rounded-lg text-blue-500">
            <Activity className="w-6 h-6" />
          </div>
          <div>
            <p className="text-sm text-gray-400">Total Telemetry</p>
            <p className="text-2xl font-bold">{stats.events}</p>
          </div>
        </div>
        <div className="bg-gray-800 p-6 rounded-lg border border-gray-700 flex items-center gap-4">
          <div className="p-3 bg-yellow-500/20 rounded-lg text-yellow-500">
            <AlertTriangle className="w-6 h-6" />
          </div>
          <div>
            <p className="text-sm text-gray-400">Active Correlations</p>
            <p className="text-2xl font-bold">{stats.incidents}</p>
          </div>
        </div>
        <div className="bg-gray-800 p-6 rounded-lg border border-gray-700 flex items-center gap-4">
          <div className="p-3 bg-green-500/20 rounded-lg text-green-500">
            <Shield className="w-6 h-6" />
          </div>
          <div>
            <p className="text-sm text-gray-400">AI Decisions</p>
            <p className="text-2xl font-bold">{stats.decisions}</p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-gray-800 rounded-lg border border-gray-700 p-6">
          <h3 className="text-lg font-medium mb-4 flex items-center gap-2">
            <Server className="w-5 h-5 text-gray-400" />
            System Health
          </h3>
          <div className="space-y-4">
            <div className="flex justify-between items-center p-3 bg-gray-900 rounded border border-gray-700">
              <span className="text-sm font-medium">REST API Listener</span>
              <span className={`px-2 py-1 text-xs font-bold rounded ${health?.status === 'ok' ? 'bg-green-500/20 text-green-500' : 'bg-red-500/20 text-red-500'}`}>
                {health?.status === 'ok' ? 'ONLINE' : 'UNAVAILABLE'}
              </span>
            </div>
            <div className="flex justify-between items-center p-3 bg-gray-900 rounded border border-gray-700">
              <span className="text-sm font-medium">PostgreSQL (Operational Store)</span>
              <span className={`px-2 py-1 text-xs font-bold rounded ${ready?.status === 'ready' ? 'bg-green-500/20 text-green-500' : 'bg-red-500/20 text-red-500'}`}>
                {ready?.status === 'ready' ? 'ONLINE' : 'UNAVAILABLE'}
              </span>
            </div>
            <div className="flex justify-between items-center p-3 bg-gray-900 rounded border border-gray-700">
              <span className="text-sm font-medium">ClickHouse (Telemetry Lake)</span>
              <span className="px-2 py-1 text-xs font-bold rounded bg-orange-500/20 text-orange-500">
                BLOCKED
              </span>
            </div>
            <div className="flex justify-between items-center p-3 bg-gray-900 rounded border border-gray-700">
              <span className="text-sm font-medium">Shuffle / Orchestrator</span>
              <span className="px-2 py-1 text-xs font-bold rounded bg-orange-500/20 text-orange-500">
                BLOCKED
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
