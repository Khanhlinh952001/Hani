import * as React from "react";

import { cn } from "@/lib/utils";

function Input({ className, type, ...props }: React.ComponentProps<"input">) {
  return (
    <input
      type={type}
      data-slot="input"
      className={cn(
        "h-10 w-full rounded-lg border border-input bg-card/90 px-3 text-sm text-foreground shadow-sm shadow-primary/5 outline-none transition-[border-color,box-shadow] placeholder:text-muted-foreground focus-visible:border-primary/55 focus-visible:ring-1 focus-visible:ring-primary/15 disabled:cursor-not-allowed disabled:opacity-50",
        className
      )}
      {...props}
    />
  );
}

export { Input };
