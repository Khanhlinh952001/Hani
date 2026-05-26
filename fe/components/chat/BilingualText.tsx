import { cn } from "@/lib/utils";

type Props = {
  ko: string;
  vi?: string;
  className?: string;
};

export function BilingualText({ ko, vi, className }: Props) {
  return (
    <div className={cn("space-y-0", className)}>
      <p>{ko}</p>
      {vi ? <p className="bilingual-vi">{vi}</p> : null}
    </div>
  );
}
