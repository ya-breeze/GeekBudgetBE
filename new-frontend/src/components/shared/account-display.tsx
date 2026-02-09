"use client";

import { cn, getApiUrl } from "@/lib/utils";
import { Avatar } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import type { Account } from "@/lib/api/models";

interface AccountDisplayProps {
  account: Account | undefined;
  size?: "sm" | "md";
  showType?: boolean;
  className?: string;
}

const typeColors: Record<string, string> = {
  asset: "bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300",
  income: "bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300",
  expense: "bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300",
};

const avatarColors: Record<string, string> = {
  asset: "bg-blue-500",
  income: "bg-green-500",
  expense: "bg-red-500",
};

function getInitial(name: string): string {
  // If starts with emoji, use it
  const emojiMatch = name.match(
    /^(\p{Emoji_Presentation}|\p{Extended_Pictographic})/u
  );
  if (emojiMatch) return emojiMatch[1];
  return name.charAt(0).toUpperCase();
}

export function AccountDisplay({
  account,
  size = "md",
  showType = false,
  className,
}: AccountDisplayProps) {
  if (!account) {
    return <span className="text-muted-foreground italic">Unknown</span>;
  }

  const sizeClass = size === "sm" ? "h-6 w-6 text-xs" : "h-8 w-8 text-sm";
  const textSize = size === "sm" ? "text-xs" : "text-sm";

  return (
    <div className={cn("flex items-center gap-2", className)}>
      {account.image ? (
        <Avatar className={sizeClass}>
          <img
            src={getApiUrl(`/v1/accounts/${account.id}/image`)}
            alt={account.name}
            className="h-full w-full object-cover rounded-full"
          />
        </Avatar>
      ) : (
        <div
          className={cn(
            sizeClass,
            avatarColors[account.type] ?? "bg-gray-500",
            "flex items-center justify-center rounded-full text-white font-medium shrink-0"
          )}
        >
          {getInitial(account.name)}
        </div>
      )}
      <span className={cn(textSize, "truncate")}>{account.name}</span>
      {showType && (
        <Badge
          variant="secondary"
          className={cn("text-[10px] px-1.5 py-0", typeColors[account.type])}
        >
          {account.type}
        </Badge>
      )}
    </div>
  );
}

export function AccountTypeBadge({
  type,
}: {
  type: string;
}) {
  return (
    <Badge
      variant="secondary"
      className={cn(
        "text-xs capitalize",
        typeColors[type] ?? "bg-gray-100 text-gray-700"
      )}
    >
      {type}
    </Badge>
  );
}
