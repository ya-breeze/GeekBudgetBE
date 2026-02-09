"use client";

import { cn } from "@/lib/utils";
import { formatAmount } from "@/lib/utils/format";

interface AmountDisplayProps {
  amount: number;
  currencyName?: string;
  showCurrency?: boolean;
  colored?: boolean;
  size?: "sm" | "md" | "lg";
  className?: string;
}

export function AmountDisplay({
  amount,
  currencyName,
  showCurrency = true,
  colored = true,
  size = "md",
  className,
}: AmountDisplayProps) {
  const isPositive = amount > 0;
  const isNegative = amount < 0;

  const sizeClasses = {
    sm: "text-sm",
    md: "text-base",
    lg: "text-lg font-semibold",
  };

  const colorClass = colored
    ? isPositive
      ? "text-green-600 dark:text-green-400"
      : isNegative
      ? "text-red-600 dark:text-red-400"
      : "text-muted-foreground"
    : "";

  return (
    <span className={cn(sizeClasses[size], colorClass, "font-mono", className)}>
      {formatAmount(amount)}
      {showCurrency && currencyName && (
        <span className="ml-1 text-muted-foreground font-sans">
          {currencyName}
        </span>
      )}
    </span>
  );
}
