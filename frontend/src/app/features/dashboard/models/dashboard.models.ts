export interface AssetCard {
    accountId: string;
    accountName: string;
    balance: number;
    currencyId: string;
    currencyName: string;
    trendPercent: number;
    trendDirection: 'up' | 'down' | 'neutral';
    accountImage?: string;
    reconciliationState: 'reconciled' | 'warning' | 'unreconciled' | 'none';
    reconciliationTooltip: string;
}

export interface AssetTotal {
    currencyId: string;
    currencyName: string;
    totalBalance: number;
    trendPercent: number;
    trendDirection: 'up' | 'down' | 'neutral';
    history: number[];
}
