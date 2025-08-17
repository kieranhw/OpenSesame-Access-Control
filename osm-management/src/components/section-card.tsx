import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";

interface SectionCardProps {
  title: string;
  children: React.ReactNode;
  icon?: React.ElementType;
  subheader?: React.ReactNode;
  className?: string;
}

export function SectionCard({ title, children, subheader, icon: Icon, className }: SectionCardProps) {
  return (
    <Card className={cn("flex w-full flex-col gap-6 p-6", className)}>
      <div className="space-y-2">
        <h1 className="tracking-narrow flex items-center gap-2 text-xl font-semibold">
          {Icon && <Icon className="h-6 w-6" />}
          {title}
        </h1>
        {subheader && <div className="text-muted-foreground text-sm">{subheader}</div>}
      </div>
      {children}
    </Card>
  );
}
