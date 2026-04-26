import React, { useState, useEffect } from 'react';
import Settings from './Settings';

interface DashboardProps {
  token: string;
  onLogout: () => void;
}

const Dashboard: React.FC<DashboardProps> = ({ token, onLogout }) => {
  const [view, setView] = useState<'dashboard' | 'settings'>('dashboard');
  const [latestData, setLatestData] = useState<any>(null);
  const [strongest, setStrongest] = useState<any>(null);
  const [weakest, setWeakest] = useState<any>(null);
  const [averageData, setAverageData] = useState<any>(null);
  
  const [fromDate, setFromDate] = useState(new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]);
  const [toDate, setToDate] = useState(new Date().toISOString().split('T')[0]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (view === 'dashboard') {
      fetchDashboardData();
    }
  }, [view]);

  const fetchDashboardData = async () => {
    setLoading(true);
    try {
      const headers = { Authorization: token };
      
      const [latestRes, strongRes, weakRes] = await Promise.all([
        fetch('http://localhost:3000/latest', { headers }),
        fetch('http://localhost:3000/strongest', { headers }),
        fetch('http://localhost:3000/weakest', { headers })
      ]);

      if (latestRes.ok) setLatestData(await latestRes.json());
      if (strongRes.ok) setStrongest(await strongRes.json());
      if (weakRes.ok) setWeakest(await weakRes.json());
      
    } catch (err) {
      console.error('Error fetching dashboard data', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchAverage = async () => {
    try {
      const response = await fetch(`http://localhost:3000/average?from=${fromDate}&to=${toDate}`, {
        headers: { Authorization: token }
      });
      if (response.ok) {
        setAverageData(await response.json());
      }
    } catch (err) {
      console.error('Error fetching average data', err);
    }
  };

  if (view === 'settings') {
    return <Settings token={token} onBack={() => setView('dashboard')} />;
  }

  return (
    <div className="dashboard-view">
      <header className="dashboard-header">
        <h1>Currency Dashboard</h1>
        <div className="header-actions">
          <button onClick={() => setView('settings')}>Settings</button>
          <button className="logout-btn" onClick={onLogout}>Logout</button>
        </div>
      </header>

      {loading ? (
        <div className="loading">Updating data...</div>
      ) : (
        <div className="dashboard-grid">
          {/* Analytics Cards */}
          <div className="stats-row">
            <div className="card stat-card">
              <h4>Strongest Currency</h4>
              {strongest && strongest.data && (
                <div className="stat-val">
                  {Object.entries(strongest.data.rates).map(([k, v]: any) => (
                    <span key={k}>{k}: {v.toFixed(4)}</span>
                  ))}
                </div>
              )}
              <small>Base: {strongest?.data?.base}</small>
            </div>
            <div className="card stat-card">
              <h4>Weakest Currency</h4>
              {weakest && weakest.data && (
                <div className="stat-val">
                  {Object.entries(weakest.data.rates).map(([k, v]: any) => (
                    <span key={k}>{k}: {v.toFixed(4)}</span>
                  ))}
                </div>
              )}
              <small>Base: {weakest?.data?.base}</small>
            </div>
          </div>

          {/* Latest Rates Table */}
          <div className="card full-width">
            <h3>Latest Exchange Rates (Base: {latestData?.data?.base})</h3>
            <div className="table-container">
              <table>
                <thead>
                  <tr>
                    <th>Currency</th>
                    <th>Rate</th>
                  </tr>
                </thead>
                <tbody>
                  {latestData && latestData.data && Object.entries(latestData.data.rates).map(([curr, rate]: any) => (
                    <tr key={curr}>
                      <td>{curr}</td>
                      <td>{rate.toFixed(4)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Average Rates Calculator */}
          <div className="card full-width">
            <h3>Average Rates Over Time</h3>
            <div className="average-controls">
              <div className="date-input">
                <label>From:</label>
                <input type="date" value={fromDate} onChange={(e) => setFromDate(e.target.value)} />
              </div>
              <div className="date-input">
                <label>To:</label>
                <input type="date" value={toDate} onChange={(e) => setToDate(e.target.value)} />
              </div>
              <button onClick={fetchAverage}>Calculate</button>
            </div>

            {averageData && averageData.data && (
              <div className="average-results">
                <h4>Averages for {averageData.data.date}</h4>
                <div className="results-grid">
                  {Object.entries(averageData.data.rates).map(([curr, rate]: any) => (
                    <div key={curr} className="avg-item">
                      <span className="label">{curr}:</span>
                      <span className="value">{rate.toFixed(4)}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;
