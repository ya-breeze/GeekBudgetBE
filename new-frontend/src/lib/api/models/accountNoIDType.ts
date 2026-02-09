export type AccountNoIDType = typeof AccountNoIDType[keyof typeof AccountNoIDType];

export const AccountNoIDType = {
  expense: 'expense',
  income: 'income',
  asset: 'asset',
} as const;
