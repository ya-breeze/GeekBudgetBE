// Entity base type
export interface Entity {
  id: string;
  createdAt?: string;
  updatedAt?: string;
  deletedAt?: string;
}

// Account types
export type AccountNoIDType = "asset" | "income" | "expense";

export const AccountNoIDType = {
  asset: "asset" as const,
  income: "income" as const,
  expense: "expense" as const,
};

export interface BankAccountInfoBalancesItem {
  currency?: string;
  balance?: number;
}

export interface BankAccountInfo {
  accountId?: string;
  bankId?: string;
  balances?: BankAccountInfoBalancesItem[];
}

export interface AccountNoID {
  name: string;
  description?: string;
  type: AccountNoIDType;
  bankInfo?: BankAccountInfo;
  showInDashboardSummary?: boolean;
  hideFromReports?: boolean;
  image?: string;
  ignoreUnprocessedBefore?: string;
  openingDate?: string;
  closingDate?: string;
  showInReconciliation?: boolean;
}

export type Account = Entity & AccountNoID & { showInDashboardSummary: boolean };

// Currency types
export interface CurrencyNoID {
  name: string;
  symbol?: string;
  decimalPlaces?: number;
}

export type Currency = Entity & CurrencyNoID;

// API Response wrapper
export interface ApiResponse<T> {
  data: T;
}

// Delete params
export interface DeleteAccountParams {
  replacementAccount?: string;
}

// Upload body
export interface UploadAccountImageBody {
  file: File;
}

// Export all types
export * from "./account";
export * from "./accountNoID";
export * from "./accountNoIDType";
export * from "./currency";
export * from "./deleteAccountParams";
export * from "./uploadAccountImageBody";
