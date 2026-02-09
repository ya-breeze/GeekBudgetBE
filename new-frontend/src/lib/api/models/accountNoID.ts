import type { AccountNoIDType } from './accountNoIDType';
import type { BankAccountInfo } from './bankAccountInfo';

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
}
