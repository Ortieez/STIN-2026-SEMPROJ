import React, { useState, useEffect } from 'react';

interface UserSettings {
  baseCurrency: string;
  selectedCurrencies: string[];
}

interface SettingsProps {
  token: string;
  onBack: () => void;
}

const Settings: React.FC<SettingsProps> = ({ token, onBack }) => {
  const [settings, setSettings] = useState<UserSettings>({
    baseCurrency: 'EUR',
    selectedCurrencies: [],
  });
  const [newCurrency, setNewCurrency] = useState('');
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState('');

  useEffect(() => {
    fetchSettings();
  }, []);

  const fetchSettings = async () => {
    try {
      const response = await fetch('http://localhost:3000/settings', {
        headers: { Authorization: token },
      });
      if (response.ok) {
        const data = await response.json();
        setSettings(data);
      }
    } catch (err) {
      console.error('Failed to fetch settings', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setMessage('');
    try {
      const response = await fetch('http://localhost:3000/settings', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: token,
        },
        body: JSON.stringify(settings),
      });
      if (response.ok) {
        setMessage('Settings saved successfully!');
      } else {
        setMessage('Failed to save settings');
      }
    } catch (err) {
      setMessage('Error connecting to server');
    }
  };

  const addCurrency = () => {
    if (newCurrency && !settings.selectedCurrencies.includes(newCurrency.toUpperCase())) {
      setSettings({
        ...settings,
        selectedCurrencies: [...settings.selectedCurrencies, newCurrency.toUpperCase()],
      });
      setNewCurrency('');
    }
  };

  const removeCurrency = (curr: string) => {
    setSettings({
      ...settings,
      selectedCurrencies: settings.selectedCurrencies.filter((c) => c !== curr),
    });
  };

  if (loading) return <div>Loading settings...</div>;

  return (
    <div className="settings-view">
      <h3>User Settings</h3>
      {message && <div className={message.includes('success') ? 'success-msg' : 'error-msg'}>{message}</div>}
      
      <div className="settings-group">
        <label>Base Currency:</label>
        <input
          type="text"
          value={settings.baseCurrency}
          onChange={(e) => setSettings({ ...settings, baseCurrency: e.target.value.toUpperCase() })}
        />
      </div>

      <div className="settings-group">
        <label>Tracked Currencies:</label>
        <div className="currency-list">
          {settings.selectedCurrencies.map((curr) => (
            <span key={curr} className="currency-tag">
              {curr} <button onClick={() => removeCurrency(curr)}>×</button>
            </span>
          ))}
        </div>
        <div className="add-currency">
          <input
            type="text"
            placeholder="Add currency (e.g. USD)"
            value={newCurrency}
            onChange={(e) => setNewCurrency(e.target.value)}
          />
          <button onClick={addCurrency}>Add</button>
        </div>
      </div>

      <div className="settings-actions">
        <button className="secondary" onClick={onBack}>Back to Dashboard</button>
        <button className="primary" onClick={handleSave}>Save Settings</button>
      </div>
    </div>
  );
};

export default Settings;
