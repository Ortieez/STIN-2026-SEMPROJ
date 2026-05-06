import React, { useState, useEffect, useCallback } from 'react';
import Settings from './Settings';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { LogOut, Settings as SettingsIcon, TrendingUp, TrendingDown, RefreshCcw, Calculator, Loader2, AlertCircle, Lock, AlertTriangle } from "lucide-react";
import { useTranslation } from "@/i18n/LanguageContext";
import { LanguageSwitcher } from "./LanguageSwitcher";

interface DashboardProps {
  token: string;
  onLogout: () => void;
}

const Dashboard: React.FC<DashboardProps> = ({ token, onLogout }) => {
  const { t, language } = useTranslation();
  const [view, setView] = useState<'dashboard' | 'settings'>('dashboard');
  const [latestData, setLatestData] = useState<any>(null);
  const [strongest, setStrongest] = useState<any>(null);
  const [weakest, setWeakest] = useState<any>(null);
  const [averageData, setAverageData] = useState<any>(null);
  
  const [fromDate, setFromDate] = useState(new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]);
  const [toDate, setToDate] = useState(new Date().toISOString().split('T')[0]);
  const [loading, setLoading] = useState(false);
  const [avgLoading, setAvgLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [avgError, setAvgError] = useState<string | null>(null);
  const [userSettings, setUserSettings] = useState<{ baseCurrency: string, selectedCurrencies: string[] } | null>(null);
  /* @ts-ignore */
  const API_URL = import.meta.env.VITE_API_URL || '';

  const fetchDashboardData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const headers: HeadersInit = { 
        'Authorization': token,
        'Accept-Language': language
      };

      const fetchWithAuth = async (url: string) => {
        const res = await fetch(url, { headers });
        if (res.status === 401) {
          onLogout();
          throw new Error('Unauthorized');
        }
        if (!res.ok) throw new Error(`Error ${res.status}`);
        return res.json();
      };

      // Fetch settings first
      const settings = await fetchWithAuth(`${API_URL}/settings`);
      setUserSettings(settings);

      // Construct query parameters
      const params = new URLSearchParams();
      if (settings.baseCurrency) {
        params.append('base', settings.baseCurrency);
      }
      if (settings.selectedCurrencies && settings.selectedCurrencies.length > 0) {
        params.append('symbols', settings.selectedCurrencies.join(','));
      }
      const queryString = params.toString() ? `?${params.toString()}` : '';

      const [latest, strong, weak] = await Promise.all([
        fetchWithAuth(`${API_URL}/latest${queryString}`),
        fetchWithAuth(`${API_URL}/strongest${queryString}`),
        fetchWithAuth(`${API_URL}/weakest${queryString}`)
      ]);

      setLatestData(latest);
      setStrongest(strong);
      setWeakest(weak);
    } catch (err: any) {
      if (err.message !== 'Unauthorized') {
        console.error('Error fetching dashboard data', err);
        setError(t('dashboard.error_connection'));
      }
    } finally {
      setLoading(false);
    }
  }, [token, language, onLogout, t, API_URL]);

  useEffect(() => {
    if (view === 'dashboard') {
      fetchDashboardData();
    }
  }, [view, fetchDashboardData]);

  const hasCurrencies = latestData?.data && Object.keys(latestData.data.rates).length > 0;

  const fetchAverage = async () => {
    setAvgError(null);
    
    if (!hasCurrencies) {
      setAvgError(t('dashboard.avg_calc.no_currencies_error'));
      return;
    }

    setAvgLoading(true);
    try {
      const params = new URLSearchParams();
      params.append('from', fromDate);
      params.append('to', toDate);
      
      if (userSettings) {
        if (userSettings.baseCurrency) {
          params.append('base', userSettings.baseCurrency);
        }
        if (userSettings.selectedCurrencies && userSettings.selectedCurrencies.length > 0) {
          params.append('symbols', userSettings.selectedCurrencies.join(','));
        }
      }

      const response = await fetch(`${API_URL}/average?${params.toString()}`, {
        headers: { 
          'Authorization': token,
          'Accept-Language': language
        }
      });
      if (response.status === 401) {
        onLogout();
        return;
      }
      if (response.ok) {
        setAverageData(await response.json());
      } else {
        const data = await response.json();
        setAvgError(data.error || 'Failed to calculate');
      }
    } catch (err) {
      console.error('Error fetching average data', err);
      setAvgError(t('dashboard.error_connection'));
    } finally {
      setAvgLoading(false);
    }
  };

  if (view === 'settings') {
    return <Settings token={token} onBack={() => setView('dashboard')} />;
  }

  return (
    <div className="container mx-auto py-8 px-4 max-w-6xl animate-in fade-in duration-500">
      <header className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 mb-8 border-b pb-6">
        <div>
          <h1 className="text-4xl font-extrabold tracking-tight">{t('dashboard.title')}</h1>
          <p className="text-muted-foreground mt-1">{t('dashboard.subtitle')}</p>
        </div>
        <div className="flex flex-wrap items-center gap-3">
          <LanguageSwitcher />
          <div className="h-8 w-[1px] bg-border mx-1 hidden sm:block"></div>
          <Button variant="outline" size="sm" onClick={fetchDashboardData} disabled={loading}>
            <RefreshCcw className={`h-4 w-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
            {t('dashboard.refresh')}
          </Button>
          <Button variant="outline" size="sm" onClick={() => setView('settings')}>
            <SettingsIcon className="h-4 w-4 mr-2" />
            {t('dashboard.settings')}
          </Button>
          <Button variant="destructive" size="sm" onClick={onLogout}>
            <LogOut className="h-4 w-4 mr-2" />
            {t('dashboard.logout')}
          </Button>
        </div>
      </header>

      {loading && !latestData ? (
        <div className="flex flex-col items-center justify-center min-h-[400px] space-y-4">
          <Loader2 className="h-12 w-12 animate-spin text-primary" />
          <p className="text-lg font-medium animate-pulse">{t('dashboard.fetching')}</p>
        </div>
      ) : error ? (
        <div className="flex flex-col items-center justify-center min-h-[400px] space-y-6 text-center animate-in zoom-in-95 duration-300">
          <div className="bg-destructive/10 p-6 rounded-full">
            <AlertCircle className="h-16 w-12 text-destructive" />
          </div>
          <div className="space-y-2">
            <h3 className="text-2xl font-bold tracking-tight">{error}</h3>
            <p className="text-muted-foreground max-w-md mx-auto">
              Please check your network or ensure the backend server is running on port 3000.
            </p>
          </div>
          <Button size="lg" onClick={fetchDashboardData}>
            <RefreshCcw className="h-4 w-4 mr-2" />
            {t('dashboard.refresh')}
          </Button>
        </div>
      ) : (
        <div className="space-y-8">
          {hasCurrencies && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 animate-in slide-in-from-top-4 duration-500">
              <Card className="border-l-4 border-l-green-500 shadow-md">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">{t('dashboard.strongest')}</CardTitle>
                  <TrendingUp className="h-4 w-4 text-green-500" />
                </CardHeader>
                <CardContent>
                  {strongest?.data && Object.entries(strongest.data.rates).map(([k, v]: any) => (
                    <div key={k} className="flex items-baseline gap-2">
                      <span className="text-2xl font-bold">{k}</span>
                      <span className="text-muted-foreground text-lg">{v.toFixed(4)}</span>
                    </div>
                  ))}
                  <p className="text-xs text-muted-foreground mt-1">{t('dashboard.relative_to')} {strongest?.data?.base}</p>
                </CardContent>
              </Card>

              <Card className="border-l-4 border-l-red-500 shadow-md">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">{t('dashboard.weakest')}</CardTitle>
                  <TrendingDown className="h-4 w-4 text-red-500" />
                </CardHeader>
                <CardContent>
                  {weakest?.data && Object.entries(weakest.data.rates).map(([k, v]: any) => (
                    <div key={k} className="flex items-baseline gap-2">
                      <span className="text-2xl font-bold">{k}</span>
                      <span className="text-muted-foreground text-lg">{v.toFixed(4)}</span>
                    </div>
                  ))}
                  <p className="text-xs text-muted-foreground mt-1">{t('dashboard.relative_to')} {weakest?.data?.base}</p>
                </CardContent>
              </Card>
            </div>
          )}

          <Tabs defaultValue="rates" className="w-full">
            <TabsList className="grid w-full grid-cols-2 max-w-[400px] mb-4 relative">
              <TabsTrigger value="rates">{t('dashboard.tabs.rates')}</TabsTrigger>
              
              {!hasCurrencies ? (
                <Tooltip delayDuration={0}>
                  <TooltipTrigger asChild>
                    <div className="flex items-center justify-center">
                      <TabsTrigger value="average" disabled className="w-full opacity-50 cursor-not-allowed">
                        <Lock className="h-3 w-3 mr-2" />
                        {t('dashboard.tabs.average')}
                      </TabsTrigger>
                    </div>
                  </TooltipTrigger>
                  <TooltipContent side="top" className="bg-destructive text-destructive-foreground">
                    <p className="font-semibold">{t('dashboard.avg_calc.no_currencies_tooltip')}</p>
                  </TooltipContent>
                </Tooltip>
              ) : (
                <TabsTrigger value="average">{t('dashboard.tabs.average')}</TabsTrigger>
              )}
            </TabsList>

            <TabsContent value="rates" className="mt-0">
              <Card className="shadow-lg">
                <CardHeader>
                  <CardTitle>{t('dashboard.latest.title')}</CardTitle>
                  <CardDescription>{t('dashboard.latest.description', { base: latestData?.data?.base || '...' })}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="rounded-md border overflow-hidden">
                    <Table>
                      <TableHeader className="bg-muted/50">
                        <TableRow>
                          <TableHead className="font-bold">{t('dashboard.latest.col_currency')}</TableHead>
                          <TableHead className="text-right font-bold">{t('dashboard.latest.col_rate')}</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {latestData?.data && Object.entries(latestData.data.rates).length > 0 ? (
                          Object.entries(latestData.data.rates).map(([curr, rate]: any) => (
                            <TableRow key={curr} className="hover:bg-muted/30">
                              <TableCell className="font-medium">{curr}</TableCell>
                              <TableCell className="text-right font-mono">{rate.toFixed(4)}</TableCell>
                            </TableRow>
                          ))
                        ) : (
                          <TableRow>
                            <TableCell colSpan={2} className="text-center py-12 text-muted-foreground">
                              <div className="flex flex-col items-center space-y-2">
                                <AlertCircle className="h-8 w-8 text-muted-foreground/50" />
                                <p className="italic text-lg">{t('dashboard.latest.empty')}</p>
                                <Button variant="link" onClick={() => setView('settings')} className="text-primary font-bold uppercase tracking-wider">
                                  {t('dashboard.settings')}
                                </Button>
                              </div>
                            </TableCell>
                          </TableRow>
                        )}
                      </TableBody>
                    </Table>
                  </div>
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="average" className="mt-0">
              <Card className="shadow-lg">
                <CardHeader>
                  <CardTitle>{t('dashboard.avg_calc.title')}</CardTitle>
                  <CardDescription>{t('dashboard.avg_calc.description')}</CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  {avgError && (
                    <div className="flex items-center gap-2 p-4 bg-destructive/10 text-destructive rounded-lg border border-destructive/20 animate-in slide-in-from-top-2">
                      <AlertTriangle className="h-5 w-5 shrink-0" />
                      <p className="text-sm font-medium">{avgError}</p>
                    </div>
                  )}

                  <div className="flex flex-col md:flex-row items-end gap-4 p-4 bg-muted/20 rounded-xl border">
                    <div className="space-y-2 flex-1 w-full">
                      <Label>{t('dashboard.avg_calc.from')}</Label>
                      <Input type="date" value={fromDate} onChange={(e) => setFromDate(e.target.value)} />
                    </div>
                    <div className="space-y-2 flex-1 w-full">
                      <Label>{t('dashboard.avg_calc.to')}</Label>
                      <Input type="date" value={toDate} onChange={(e) => setToDate(e.target.value)} />
                    </div>
                    <Button className="w-full md:w-auto" onClick={fetchAverage} disabled={avgLoading}>
                      {avgLoading ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <Calculator className="h-4 w-4 mr-2" />}
                      {t('dashboard.avg_calc.calculate')}
                    </Button>
                  </div>

                  {averageData?.data && (
                    <div className="space-y-4 animate-in slide-in-from-bottom-2 duration-300">
                      <h4 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider px-1">
                        {t('dashboard.avg_calc.results_for')} {averageData.data.date}
                      </h4>
                      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
                        {Object.entries(averageData.data.rates).map(([curr, rate]: any) => (
                          <div key={curr} className="p-3 bg-secondary/50 border rounded-lg flex flex-col items-center">
                            <span className="text-xs text-muted-foreground mb-1 font-bold">{curr}</span>
                            <span className="text-lg font-bold">{rate.toFixed(4)}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </div>
      )}
    </div>
  );
};

export default Dashboard;
