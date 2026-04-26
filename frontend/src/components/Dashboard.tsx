import React from 'react';

interface DashboardProps {
  token: string;
  onLogout: () => void;
}

const Dashboard: React.FC<DashboardProps> = ({ onLogout }) => {
  return (
    <div className="dashboard-container">
      <header>
        <h1>Currency Dashboard</h1>
        <button onClick={onLogout}>Logout</button>
      </header>
      <main>
        <p>Welcome! You are logged in.</p>
        {/* Analytics and rates will go here in Phase 4 */}
      </main>
    </div>
  );
};

export default Dashboard;
