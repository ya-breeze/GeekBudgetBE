/**
 * Custom account hooks that wrap the generated API hooks.
 * This file is NOT generated and won't be overwritten.
 *
 * Use these hooks in your components instead of the generated ones.
 */

import { useQuery, useMutation, useQueryClient, type UseMutationOptions } from "@tanstack/react-query";
import { apiClient } from "../client";
import type { Account, AccountNoID, DeleteAccountParams, UploadAccountImageBody } from "../models";

// Query keys
export const accountKeys = {
  all: ["accounts"] as const,
  lists: () => [...accountKeys.all, "list"] as const,
  list: (filters?: any) => [...accountKeys.lists(), filters] as const,
  details: () => [...accountKeys.all, "detail"] as const,
  detail: (id: string) => [...accountKeys.details(), id] as const,
};

// API functions (correct response types)
const getAccounts = async (): Promise<Account[]> => {
  const { data } = await apiClient.get<Account[]>("/v1/accounts");
  return data;
};

const getAccount = async (id: string): Promise<Account> => {
  const { data } = await apiClient.get<Account>(`/v1/accounts/${id}`);
  return data;
};

const createAccount = async (accountData: AccountNoID): Promise<Account> => {
  const { data } = await apiClient.post<Account>("/v1/accounts", accountData);
  return data;
};

const updateAccount = async (id: string, accountData: AccountNoID): Promise<Account> => {
  const { data } = await apiClient.put<Account>(`/v1/accounts/${id}`, accountData);
  return data;
};

const deleteAccount = async (id: string, params?: DeleteAccountParams): Promise<void> => {
  await apiClient.delete(`/v1/accounts/${id}`, { params });
};

const uploadAccountImage = async (id: string, file: File): Promise<void> => {
  const formData = new FormData();
  formData.append("file", file);
  await apiClient.post(`/v1/accounts/${id}/image`, formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
};

const deleteAccountImage = async (id: string): Promise<void> => {
  await apiClient.delete(`/v1/accounts/${id}/image`);
};

// Query hooks
export function useAccounts() {
  return useQuery({
    queryKey: accountKeys.lists(),
    queryFn: getAccounts,
  });
}

export function useAccount(id: string, enabled = true) {
  return useQuery({
    queryKey: accountKeys.detail(id),
    queryFn: () => getAccount(id),
    enabled: enabled && !!id,
  });
}

// Mutation hooks
export function useCreateAccount(
  options?: UseMutationOptions<Account, Error, AccountNoID>
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createAccount,
    onSuccess: (data, variables, context) => {
      queryClient.invalidateQueries({ queryKey: accountKeys.lists() });
      (options?.onSuccess as any)?.(data, variables, context);
    },
    ...options,
  });
}

export function useUpdateAccount(
  options?: UseMutationOptions<Account, Error, { id: string; data: AccountNoID }>
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }) => updateAccount(id, data),
    onSuccess: (data, variables, context) => {
      queryClient.invalidateQueries({ queryKey: accountKeys.lists() });
      queryClient.invalidateQueries({ queryKey: accountKeys.detail(variables.id) });
      (options?.onSuccess as any)?.(data, variables, context);
    },
    ...options,
  });
}

export function useDeleteAccount(
  options?: UseMutationOptions<void, Error, { id: string; params?: DeleteAccountParams }>
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, params }) => deleteAccount(id, params),
    onSuccess: (data, variables, context) => {
      queryClient.invalidateQueries({ queryKey: accountKeys.lists() });
      queryClient.invalidateQueries({ queryKey: accountKeys.detail(variables.id) });
      (options?.onSuccess as any)?.(data, variables, context);
    },
    ...options,
  });
}

export function useUploadAccountImage(
  options?: UseMutationOptions<void, Error, { id: string; file: File }>
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, file }) => uploadAccountImage(id, file),
    onSuccess: (data, variables, context) => {
      queryClient.invalidateQueries({ queryKey: accountKeys.lists() });
      queryClient.invalidateQueries({ queryKey: accountKeys.detail(variables.id) });
      (options?.onSuccess as any)?.(data, variables, context);
    },
    ...options,
  });
}

export function useDeleteAccountImage(
  options?: UseMutationOptions<void, Error, string>
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteAccountImage,
    onSuccess: (data, variables, context) => {
      queryClient.invalidateQueries({ queryKey: accountKeys.lists() });
      queryClient.invalidateQueries({ queryKey: accountKeys.detail(variables) });
      (options?.onSuccess as any)?.(data, variables, context);
    },
    ...options,
  });
}
