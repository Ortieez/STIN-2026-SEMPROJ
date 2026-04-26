import React, { useState, useEffect } from 'react';
import Settings from './Settings';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { LogOut, Settings as SettingsIcon, TrendingUp, TrendingDown, RefreshCcw, Calculator, Loader2 } from "lucide-react";

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
  const [avgLoading, setAvgLoading] = useState(false);

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
    setAvgLoading(true);
    try {
      const response = await fetch(`http://localhost:3000/average?from=${fromDate}&to=${toDate}`, {
        headers: { Authorization: token }
      });
      if (response.ok) {
        setAverageData(await response.json());
      }
    } catch (err) {
      console.error('Error fetching average data', err);
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
          <h1 className="text-4xl font-extrabold tracking-tight">Currency Dashboard</h1>
          <p className="text-muted-foreground mt-1">Real-time exchange rates and analytics</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={fetchDashboardData}>
            <RefreshCcw className="h-4 w-4 mr-2" />
            Refresh
          </Button>
          <Button variant="outline" size="sm" onClick={() => setView('settings')}>
            <SettingsIcon className="h-4 w-4 mr-2" />
            Settings
          </Button>
          <Button variant="destructive" size="sm" onClick={onLogout}>
            <LogOut className="h-4 w-4 mr-2" />
            Logout
          </Button>
        </div>
      </header>

      {loading ? (
        <div className="flex flex-col items-center justify-center min-h-[400px] space-y-4">
          <Loader2 className="h-12 w-12 animate-spin text-primary" />
          <p className="text-lg font-medium animate-pulse">Fetching latest markets...</p>
        </div>
      ) : (
        <div className="space-y-8">
          {/* Quick Stats Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card className="border-l-4 border-l-green-500 shadow-md">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Strongest Currency</CardTitle>
                <TrendingUp className="h-4 w-4 text-green-500" />
              </CardHeader>
              <CardContent>
                {strongest?.data && Object.entries(strongest.data.rates).map(([k, v]: any) => (
                  <div key={k} className="flex items-baseline gap-2">
                    <span className="text-2xl font-bold">{k}</span>
                    <span className="text-muted-foreground text-lg">{v.toFixed(4)}</span>
                  </div>
                ))}
                <p className="text-xs text-muted-foreground mt-1">Relative to {strongest?.data?.base}</p>
              </CardContent>
            </Card>

            <Card className="border-l-4 border-l-red-500 shadow-md">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Weakest Currency</CardTitle>
                <TrendingDown className="h-4 w-4 text-red-500" />
              </CardHeader>
              <CardContent>
                {weakest?.data && Object.entries(weakest.data.rates).map(([k, v]: any) => (
                  <div key={k} className="flex items-baseline gap-2">
                    <span className="text-2xl font-bold">{k}</span>
                    <span className="text-muted-foreground text-lg">{v.toFixed(4)}</span>
                  </div>
                ))}
                <p className="text-xs text-muted-foreground mt-1">Relative to {weakest?.data?.base}</p>
              </CardContent>
            </Card>
          </div>

          <Tabs defaultValue="rates" className="w-full">
            <TabsList className="grid w-full grid-cols-2 max-w-[400px] mb-4">
              <TabsTrigger value="rates">Market Rates</TabsTrigger>
              <TabsTrigger value="average">Historical Avg</TabsTrigger>
            </TabsList>

            <TabsContent value="rates" className="mt-0">
              <Card className="shadow-lg">
                <CardHeader>
                  <CardTitle>Latest Exchange Rates</CardTitle>
                  <CardDescription>Current market values for your tracked currencies (Base: {latestData?.data?.base})</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="rounded-md border overflow-hidden">
                    <Table>
                      <TableHeader className="bg-muted/50">
                        <TableRow>
                          <TableHead className="font-bold">Currency</TableHead>
                          <TableHead className="text-right font-bold">Exchange Rate</TableHead>
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
                            <TableCell colSpan={2} className="text-center py-8 text-muted-foreground italic">
                              No tracked currencies. Update them in Settings.
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
                  <CardTitle>Average Rates Calculator</CardTitle>
                  <CardDescription>Analyze currency performance over a specific period</CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="flex flex-col md:flex-row items-end gap-4 p-4 bg-muted/20 rounded-xl border">
                    <div className="space-y-2 flex-1 w-full">
                      <Label>Start Date</Label>
                      <Input type="date" value={fromDate} onChange={(e) => setFromDate(e.target.value)} />
                    </div>
                    <div className="space-y-2 flex-1 w-full">
                      <Label>End Date</Label>
                      <Input type="date" value={toDate} onChange={(e) => setToDate(e.target.value)} />
                    </div>
                    <Button className="w-full md:w-auto" onClick={fetchAverage} disabled={avgLoading}>
                      {avgLoading ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : <Calculator className="h-4 w-4 mr-2" />}
                      Calculate
                    </Button>
                  </div>

                  {averageData?.data && (
                    <div className="space-y-4 animate-in slide-in-from-bottom-2 duration-300">
                      <h4 className="text-sm font-semibold text-muted-foreground uppercase tracking-wider px-1">
                        Results for {averageData.data.date}
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
