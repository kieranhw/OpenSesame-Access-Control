import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";

interface SectionCardProps {
  title: string;
  children: React.ReactNode;
  icon?: React.ElementType;
  subheader?: React.ReactNode;
  className?: string;
  button?: React.ReactNode;
  bodyPadding?: boolean;
}

export function SectionCard({
  title,
  children,
  subheader,
  icon: Icon,
  className,
  button,
  bodyPadding = true,
}: SectionCardProps) {
  return (
    <Card className={cn("flex w-full flex-col gap-0", className)}>
      <div className="flex justify-between p-6 pt-0">
        <div>
          <h1 className="tracking-narrow flex items-center gap-2 text-xl font-semibold">
            {Icon ? <Icon className="h-6 w-6" /> : null}
            {title ? title : null}
          </h1>
          {subheader ? <div className="text-muted-foreground text-sm">{subheader}</div> : null}
        </div>
        <div>{button ? button : null}</div>
      </div>
      <div className={`px-${bodyPadding ? 6 : 0}`}>{children}</div>
    </Card>
  );
}
