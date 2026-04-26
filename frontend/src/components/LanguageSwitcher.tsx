import { Button } from "@/components/ui/button";
import { useTranslation } from "@/i18n/LanguageContext";
import { cn } from "@/components/lib/utils";

export const LanguageSwitcher = () => {
  const { language, setLanguage } = useTranslation();

  return (
    <div className="flex items-center gap-1 bg-muted/50 p-1 rounded-lg border shadow-sm">
      <Button
        variant={language === 'en' ? 'default' : 'ghost'}
        size="sm"
        className={cn(
          "h-7 text-[10px] font-bold px-3 transition-all duration-200",
          language === 'en' ? "shadow-sm" : "text-muted-foreground hover:text-foreground"
        )}
        onClick={() => setLanguage('en')}
      >
        ENGLISH
      </Button>
      <Button
        variant={language === 'cs' ? 'default' : 'ghost'}
        size="sm"
        className={cn(
          "h-7 text-[10px] font-bold px-3 transition-all duration-200",
          language === 'cs' ? "shadow-sm" : "text-muted-foreground hover:text-foreground"
        )}
        onClick={() => setLanguage('cs')}
      >
        ČEŠTINA
      </Button>
    </div>
  );
};
