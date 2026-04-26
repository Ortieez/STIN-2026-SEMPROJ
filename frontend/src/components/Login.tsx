import React, { useState } from 'react';
import CryptoJS from 'crypto-js';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Loader2 } from "lucide-react";
import { useTranslation } from "@/i18n/LanguageContext";
import { LanguageSwitcher } from "./LanguageSwitcher";

interface LoginProps {
  onLogin: (token: string) => void;
}

const Login: React.FC<LoginProps> = ({ onLogin }) => {
  const { t, language } = useTranslation();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const hashedUsername = CryptoJS.SHA256(username).toString();
      const hashedPassword = CryptoJS.SHA256(password).toString();

      const response = await fetch('http://localhost:3000/login', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Accept-Language': language
        },
        body: JSON.stringify({ username: hashedUsername, password: hashedPassword }),
      });

      const data = await response.json();

      if (response.ok) {
        onLogin(data.token);
      } else {
        setError(data.error || t('login.invalid_credentials'));
      }
    } catch (err) {
      setError(t('login.error_connection'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-[80vh] px-4 space-y-4">
      <div className="w-full max-w-md flex justify-end">
        <LanguageSwitcher />
      </div>
      <Card className="w-full max-w-md shadow-lg border-2">
        <CardHeader className="space-y-1">
          <CardTitle className="text-3xl font-bold text-center">{t('login.title')}</CardTitle>
          <CardDescription className="text-center">
            {t('login.description')}
          </CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            {error && (
              <div className="bg-destructive/15 text-destructive text-sm p-3 rounded-md border border-destructive/20 text-center">
                {error}
              </div>
            )}
            <div className="space-y-2">
              <Label htmlFor="username">{t('login.username')}</Label>
              <Input
                id="username"
                type="text"
                placeholder="admin"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">{t('login.password')}</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
          </CardContent>
          <CardFooter>
            <Button className="w-full" type="submit" disabled={loading}>
              {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {loading ? t('login.logging_in') : t('login.button')}
            </Button>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
};

export default Login;
