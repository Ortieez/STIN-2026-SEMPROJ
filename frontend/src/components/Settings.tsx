import React, { useState, useEffect } from 'react';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { X, Plus, Save, ArrowLeft, Loader2 } from "lucide-react";
import { useTranslation } from "@/i18n/LanguageContext";

interface UserSettings {
  baseCurrency: string;
  selectedCurrencies: string[];
}

interface SettingsProps {
  token: string;
  onBack: () => void;
}

const Settings: React.FC<SettingsProps> = ({ token, onBack }) => {
  const { t, language } = useTranslation();
  const [settings, setSettings] = useState<UserSettings>({
    baseCurrency: 'EUR',
    selectedCurrencies: [],
  });
  const [newCurrency, setNewCurrency] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState<{ text: string, type: 'success' | 'error' } | null>(null);

  useEffect(() => {
    fetchSettings();
  }, []);

  const fetchSettings = async () => {
    try {
      const response = await fetch('http://localhost:3000/settings', {
        headers: { 
          Authorization: token,
          'Accept-Language': language
        },
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
    setMessage(null);
    setSaving(true);
    try {
      const response = await fetch('http://localhost:3000/settings', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: token,
          'Accept-Language': language
        },
        body: JSON.stringify(settings),
      });
      const data = await response.json();
      if (response.ok) {
        setMessage({ text: data.message || t('settings.success'), type: 'success' });
      } else {
        setMessage({ text: data.error || t('settings.error'), type: 'error' });
      }
    } catch (err) {
      setMessage({ text: t('settings.error_connection'), type: 'error' });
    } finally {
      setSaving(false);
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

  if (loading) return (
    <div className="flex items-center justify-center min-h-[50vh]">
      <Loader2 className="h-8 w-8 animate-spin text-primary" />
    </div>
  );

  return (
    <div className="max-w-2xl mx-auto py-8 px-4">
      <Card className="border-2">
        <CardHeader>
          <div className="flex items-center gap-2 mb-2">
            <Button variant="ghost" size="icon" onClick={onBack}>
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <CardTitle className="text-2xl">{t('settings.title')}</CardTitle>
          </div>
          <CardDescription>
            {t('settings.description')}
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {message && (
            <div className={`p-3 rounded-md border text-sm text-center animate-in fade-in zoom-in-95 ${
              message.type === 'success' 
                ? 'bg-green-500/10 text-green-600 border-green-500/20' 
                : 'bg-destructive/10 text-destructive border-destructive/20'
            }`}>
              {message.text}
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="base-currency">{t('settings.base_currency')}</Label>
            <Input
              id="base-currency"
              placeholder="EUR"
              value={settings.baseCurrency}
              onChange={(e) => setSettings({ ...settings, baseCurrency: e.target.value.toUpperCase() })}
              className="max-w-[200px]"
            />
          </div>

          <div className="space-y-4">
            <Label>{t('settings.tracked_currencies')}</Label>
            <div className="flex flex-wrap gap-2 p-4 bg-muted/30 rounded-lg border border-dashed min-h-[60px]">
              {settings.selectedCurrencies.length === 0 && (
                <span className="text-sm text-muted-foreground italic">{t('settings.empty_currencies')}</span>
              )}
              {settings.selectedCurrencies.map((curr) => (
                <Badge key={curr} variant="secondary" className="px-3 py-1 text-sm flex gap-2">
                  {curr}
                  <button onClick={() => removeCurrency(curr)} className="hover:text-destructive transition-colors">
                    <X className="h-3 w-3" />
                  </button>
                </Badge>
              ))}
            </div>
            
            <div className="flex gap-2">
              <Input
                placeholder={t('settings.add_placeholder')}
                value={newCurrency}
                onChange={(e) => setNewCurrency(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && addCurrency()}
                className="max-w-[200px]"
              />
              <Button variant="outline" onClick={addCurrency} type="button">
                <Plus className="h-4 w-4 mr-2" />
                {t('settings.add')}
              </Button>
            </div>
          </div>
        </CardContent>
        <CardFooter className="flex justify-between border-t p-6 bg-muted/10">
          <Button variant="ghost" onClick={onBack}>{t('settings.cancel')}</Button>
          <Button onClick={handleSave} disabled={saving}>
            {saving ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Save className="mr-2 h-4 w-4" />}
            {saving ? t('settings.saving') : t('settings.save')}
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
};

export default Settings;
