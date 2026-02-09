"use client";

import { useMemo } from "react";
import { useGetAccounts } from "@/lib/api/generated/accounts/accounts";
import type { Account } from "@/lib/api/models";

export function useAccounts() {
  const { data, isLoading, error } = useGetAccounts();

  const accounts = useMemo(() => data?.data ?? [], [data]);

  const accountMap = useMemo(
    () => new Map(accounts.map((a) => [a.id, a])),
    [accounts]
  );

  const assetAccounts = useMemo(
    () => accounts.filter((a) => a.type === "asset"),
    [accounts]
  );

  const expenseAccounts = useMemo(
    () => accounts.filter((a) => a.type === "expense"),
    [accounts]
  );

  const incomeAccounts = useMemo(
    () => accounts.filter((a) => a.type === "income"),
    [accounts]
  );

  const getAccountName = (id: string | undefined): string => {
    if (!id) return "Unknown";
    return accountMap.get(id)?.name ?? "Unknown";
  };

  const getAccount = (id: string | undefined): Account | undefined => {
    if (!id) return undefined;
    return accountMap.get(id);
  };

  return {
    accounts,
    accountMap,
    assetAccounts,
    expenseAccounts,
    incomeAccounts,
    getAccountName,
    getAccount,
    isLoading,
    error,
  };
}
