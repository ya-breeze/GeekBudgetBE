import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "../../client";
import type { Account, AccountNoID, DeleteAccountParams, UploadAccountImageBody } from "../../models";

// Query keys
export const getGetAccountsQueryKey = () => ["accounts"] as const;
export const getGetAccountQueryKey = (id: string) => ["accounts", id] as const;

// API functions
export const getAccounts = async () => {
  const { data } = await apiClient.get<Account[]>("/v1/accounts");
  return data;
};

export const getAccount = async (id: string) => {
  const { data } = await apiClient.get<{ data: Account }>(`/v1/accounts/${id}`);
  return data;
};

export const createAccount = async (accountData: { data: AccountNoID }) => {
  const { data } = await apiClient.post<{ data: Account }>("/v1/accounts", accountData.data);
  return data;
};

export const updateAccount = async ({ id, data: accountData }: { id: string; data: AccountNoID }) => {
  const { data } = await apiClient.put<{ data: Account }>(`/v1/accounts/${id}`, accountData);
  return data;
};

export const deleteAccount = async ({ id, params }: { id: string; params?: DeleteAccountParams }) => {
  const { data } = await apiClient.delete(`/v1/accounts/${id}`, { params });
  return data;
};

export const uploadAccountImage = async ({ id, data }: { id: string; data: UploadAccountImageBody }) => {
  const formData = new FormData();
  formData.append("file", data.file);
  const response = await apiClient.post(`/v1/accounts/${id}/image`, formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
  return response.data;
};

export const deleteAccountImage = async ({ id }: { id: string }) => {
  const { data } = await apiClient.delete(`/v1/accounts/${id}/image`);
  return data;
};

// Hooks
export function useGetAccounts<TData = Awaited<ReturnType<typeof getAccounts>>, TError = unknown>() {
  return useQuery<Awaited<ReturnType<typeof getAccounts>>, TError, TData>({
    queryKey: getGetAccountsQueryKey(),
    queryFn: () => getAccounts(),
  });
}

export function useGetAccount<TData = Awaited<ReturnType<typeof getAccount>>, TError = unknown>(
  id: string
) {
  return useQuery<Awaited<ReturnType<typeof getAccount>>, TError, TData>({
    queryKey: getGetAccountQueryKey(id),
    queryFn: () => getAccount(id),
    enabled: !!id,
  });
}

export const useCreateAccount = <TError = unknown, TContext = unknown>(options?: {
  mutation?: {
    onSuccess?: (data: Awaited<ReturnType<typeof createAccount>>, variables: { data: AccountNoID }, context: TContext) => void;
    onError?: (error: TError, variables: { data: AccountNoID }, context: TContext | undefined) => void;
  };
}) => {
  return useMutation<
    Awaited<ReturnType<typeof createAccount>>,
    TError,
    { data: AccountNoID },
    TContext
  >({
    mutationFn: createAccount,
    ...options?.mutation,
  });
};

export const useUpdateAccount = <TError = unknown, TContext = unknown>(options?: {
  mutation?: {
    onSuccess?: (data: Awaited<ReturnType<typeof updateAccount>>, variables: { id: string; data: AccountNoID }, context: TContext) => void;
    onError?: (error: TError, variables: { id: string; data: AccountNoID }, context: TContext | undefined) => void;
  };
}) => {
  return useMutation<
    Awaited<ReturnType<typeof updateAccount>>,
    TError,
    { id: string; data: AccountNoID },
    TContext
  >({
    mutationFn: updateAccount,
    ...options?.mutation,
  });
};

export const useDeleteAccount = <TError = unknown, TContext = unknown>(options?: {
  mutation?: {
    onSuccess?: (data: Awaited<ReturnType<typeof deleteAccount>>, variables: { id: string; params?: DeleteAccountParams }, context: TContext) => void;
    onError?: (error: TError, variables: { id: string; params?: DeleteAccountParams }, context: TContext | undefined) => void;
  };
}) => {
  return useMutation<
    Awaited<ReturnType<typeof deleteAccount>>,
    TError,
    { id: string; params?: DeleteAccountParams },
    TContext
  >({
    mutationFn: deleteAccount,
    ...options?.mutation,
  });
};

export const useUploadAccountImage = <TError = unknown, TContext = unknown>(options?: {
  mutation?: {
    onSuccess?: (data: Awaited<ReturnType<typeof uploadAccountImage>>, variables: { id: string; data: UploadAccountImageBody }, context: TContext) => void;
    onError?: (error: TError, variables: { id: string; data: UploadAccountImageBody }, context: TContext | undefined) => void;
  };
}) => {
  return useMutation<
    Awaited<ReturnType<typeof uploadAccountImage>>,
    TError,
    { id: string; data: UploadAccountImageBody },
    TContext
  >({
    mutationFn: uploadAccountImage,
    ...options?.mutation,
  });
};

export const useDeleteAccountImage = <TError = unknown, TContext = unknown>(options?: {
  mutation?: {
    onSuccess?: (data: Awaited<ReturnType<typeof deleteAccountImage>>, variables: { id: string }, context: TContext) => void;
    onError?: (error: TError, variables: { id: string }, context: TContext | undefined) => void;
  };
}) => {
  return useMutation<
    Awaited<ReturnType<typeof deleteAccountImage>>,
    TError,
    { id: string },
    TContext
  >({
    mutationFn: deleteAccountImage,
    ...options?.mutation,
  });
};
