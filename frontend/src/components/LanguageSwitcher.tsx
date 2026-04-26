import { Button } from "@/components/ui/button";
import { useTranslation } from "@/i18n/LanguageContext";

export const LanguageSwitcher = () => {
  const { language, setLanguage } = useTranslation();

  return (
    <div className="flex gap-1 bg-muted p-1 rounded-md">
      <Button
        variant={language === 'en' ? 'secondary' : 'ghost'}
        size="sm"
        className="h-8 text-xs px-2"
        onClick={() => setLanguage('en')}
      >
        EN
      </Button>
      <Button
        variant={language === 'cs' ? 'secondary' : 'ghost'}
        size="sm"
        className="h-8 text-xs px-2"
        onClick={() => setLanguage('cs')}
      >
        CZ
      </Button>
    </div>
  );
};
