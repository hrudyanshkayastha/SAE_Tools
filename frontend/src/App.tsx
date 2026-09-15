import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { Activity, ShieldAlert, Cpu, Network, Search, HardDrive } from 'lucide-react';
import Dashboard from './pages/Dashboard';
import Events from './pages/Events';
import Incidents from './pages/Incidents';
import Investigations from './pages/Investigations';
import Architecture from './pages/Architecture';

function App() {
  return (
    <Router>
      <div className="flex h-screen bg-gray-900 text-gray-100 overflow-hidden font-sans">
        {/* Sidebar */}
        <aside className="w-64 bg-gray-800 border-r border-gray-700 flex flex-col">
          <div className="p-4 border-b border-gray-700 flex items-center gap-3">
            <ShieldAlert className="text-blue-500 w-8 h-8" />
            <h1 className="text-xl font-bold tracking-wider">SAE</h1>
          </div>
          <nav className="flex-1 overflow-y-auto p-4 space-y-2">
            <Link to="/" className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-700 transition-colors">
              <Activity className="w-5 h-5 text-gray-400" />
              <span>Dashboard</span>
            </Link>
            <Link to="/events" className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-700 transition-colors">
              <Search className="w-5 h-5 text-gray-400" />
              <span>Live Events</span>
            </Link>
            <Link to="/incidents" className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-700 transition-colors">
              <ShieldAlert className="w-5 h-5 text-gray-400" />
              <span>Incidents</span>
            </Link>
            <Link to="/investigations" className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-700 transition-colors">
              <Cpu className="w-5 h-5 text-gray-400" />
              <span>AI Investigations</span>
            </Link>
            <Link to="/architecture" className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-700 transition-colors">
              <Network className="w-5 h-5 text-gray-400" />
              <span>Integrations</span>
            </Link>
          </nav>
          <div className="p-4 border-t border-gray-700 text-xs text-gray-500 flex items-center gap-2">
            <HardDrive className="w-4 h-4" />
            <span>SAE Security Fabric</span>
          </div>
        </aside>

        {/* Main Content */}
        <main className="flex-1 flex flex-col h-full overflow-hidden">
          <header className="h-14 bg-gray-800 border-b border-gray-700 flex items-center px-6">
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
              <span className="text-sm font-medium text-gray-300">SYSTEM ONLINE</span>
            </div>
          </header>
          <div className="flex-1 overflow-y-auto p-6">
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/events" element={<Events />} />
              <Route path="/incidents" element={<Incidents />} />
              <Route path="/investigations" element={<Investigations />} />
              <Route path="/architecture" element={<Architecture />} />
            </Routes>
          </div>
        </main>
      </div>
    </Router>
  );
}

export default App;
