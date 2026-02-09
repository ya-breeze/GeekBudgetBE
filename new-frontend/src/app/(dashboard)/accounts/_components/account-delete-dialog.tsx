"use client";

import { useState } from "react";
import { ConfirmDialog } from "@/components/shared/confirm-dialog";
import type { Account } from "@/lib/api/models";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { AccountDisplay } from "@/components/shared/account-display";

interface AccountDeleteDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  account: Account | null;
  availableAccounts: Account[];
  isLoading?: boolean;
  onConfirm: (replacementAccountId?: string) => void;
}

export function AccountDeleteDialog({
  open,
  onOpenChange,
  account,
  availableAccounts,
  isLoading,
  onConfirm,
}: AccountDeleteDialogProps) {
  const [replacementAccountId, setReplacementAccountId] = useState<string | undefined>();

  if (!account) return null;

  // Filter out the account being deleted from available replacements
  const replacementOptions = availableAccounts.filter((a) => a.id !== account.id);

  const handleConfirm = () => {
    onConfirm(replacementAccountId);
    setReplacementAccountId(undefined);
  };

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      setReplacementAccountId(undefined);
    }
    onOpenChange(newOpen);
  };

  return (
    <ConfirmDialog
      open={open}
      onOpenChange={handleOpenChange}
      title="Delete Account"
      description={`Are you sure you want to delete "${account.name}"? This action cannot be undone.`}
      confirmText="Delete Account"
      destructive
      isLoading={isLoading}
      onConfirm={handleConfirm}
    >
      {replacementOptions.length > 0 && (
        <div className="space-y-3 py-4">
          <div className="rounded-lg border border-yellow-200 bg-yellow-50 p-3 dark:border-yellow-900 dark:bg-yellow-950">
            <p className="text-sm text-yellow-800 dark:text-yellow-200">
              This account may be referenced in existing transactions. You can optionally select a
              replacement account to reassign those transactions.
            </p>
          </div>
          <div className="space-y-2">
            <Label htmlFor="replacement">Replacement Account (Optional)</Label>
            <Select
              value={replacementAccountId}
              onValueChange={setReplacementAccountId}
            >
              <SelectTrigger id="replacement">
                <SelectValue placeholder="None - delete without replacement" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="__none__">None - delete without replacement</SelectItem>
                {replacementOptions.map((acc) => (
                  <SelectItem key={acc.id} value={acc.id}>
                    <div className="flex items-center gap-2">
                      <AccountDisplay account={acc} size="sm" showType />
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <p className="text-xs text-muted-foreground">
              If a replacement is selected, all transactions using this account will be updated to
              use the replacement instead.
            </p>
          </div>
        </div>
      )}
    </ConfirmDialog>
  );
}
