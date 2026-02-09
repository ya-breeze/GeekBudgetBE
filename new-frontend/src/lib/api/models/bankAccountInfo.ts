import type { BankAccountInfoBalancesItem } from './bankAccountInfoBalancesItem';

export interface BankAccountInfo {
  accountId?: string;
  bankId?: string;
  balances?: BankAccountInfoBalancesItem[];
}
