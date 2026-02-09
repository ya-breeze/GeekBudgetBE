import { useQuery } from "@tanstack/react-query";
import { apiClient } from "../../client";
import type { Currency } from "../../models";

// Query keys
export const getGetCurrenciesQueryKey = () => ["currencies"] as const;

// API functions
export const getCurrencies = async () => {
  const { data } = await apiClient.get<{ data: Currency[] }>("/v1/currencies");
  return data;
};

// Hooks
export function useGetCurrencies<TData = Awaited<ReturnType<typeof getCurrencies>>, TError = unknown>() {
  return useQuery<Awaited<ReturnType<typeof getCurrencies>>, TError, TData>({
    queryKey: getGetCurrenciesQueryKey(),
    queryFn: () => getCurrencies(),
  });
}
