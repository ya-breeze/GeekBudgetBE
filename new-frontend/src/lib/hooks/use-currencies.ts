"use client";

import { useMemo } from "react";
import { useGetCurrencies } from "@/lib/api/generated/currencies/currencies";
import type { Currency } from "@/lib/api/models";

export function useCurrencies() {
  const { data, isLoading, error } = useGetCurrencies();

  const currencies = useMemo(() => data?.data ?? [], [data]);

  const currencyMap = useMemo(
    () => new Map(currencies.map((c) => [c.id, c])),
    [currencies]
  );

  const getCurrencyName = (id: string | undefined): string => {
    if (!id) return "Unknown";
    return currencyMap.get(id)?.name ?? "Unknown";
  };

  const getCurrency = (id: string | undefined): Currency | undefined => {
    if (!id) return undefined;
    return currencyMap.get(id);
  };

  return {
    currencies,
    currencyMap,
    getCurrencyName,
    getCurrency,
    isLoading,
    error,
  };
}
