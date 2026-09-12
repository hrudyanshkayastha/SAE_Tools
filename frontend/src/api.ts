import axios from 'axios';

const API_BASE = 'http://localhost:8080';

export const api = {
  getHealth: () => axios.get(`${API_BASE}/health`),
  getReady: () => axios.get(`${API_BASE}/ready`),
  getEvents: () => axios.get(`${API_BASE}/events`),
  getIncidents: () => axios.get(`${API_BASE}/incidents`),
  getInvestigations: () => axios.get(`${API_BASE}/investigations`),
  getDecisions: () => axios.get(`${API_BASE}/decisions`),
};
