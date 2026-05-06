import React, { createContext, useContext, useState, useEffect } from 'react';
import en from './locales/en.json';
import cs from './locales/cs.json';

type Language = 'en' | 'cs';

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (path: string, variables?: Record<string, string>) => string;
}

const translations: Record<Language, any> = { en, cs };

const LanguageContext = createContext<LanguageContextType | undefined>(undefined);

export const LanguageProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [language, setLanguage] = useState<Language>(
    (localStorage.getItem('app_lang') as Language) || 'en'
  );

  useEffect(() => {
    localStorage.setItem('app_lang', language);
    document.documentElement.lang = language;
  }, [language]);

  const t = (path: string, variables?: Record<string, string>) => {
    const keys = path.split('.');
    let result = translations[language];

    for (const key of keys) {
      if (result[key] === undefined) return path;
      result = result[key];
    }

    let finalString = result as string;
    if (variables) {
      Object.entries(variables).forEach(([key, val]) => {
        finalString = finalString.replace(`{{${key}}}`, val);
      });
    }

    return finalString;
  };

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  );
};

export const useTranslation = () => {
  const context = useContext(LanguageContext);
  if (!context) throw new Error('useTranslation must be used within a LanguageProvider');
  return context;
};
