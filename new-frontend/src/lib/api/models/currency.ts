import type { Entity } from './index';

export interface CurrencyNoID {
  name: string;
  symbol?: string;
  decimalPlaces?: number;
}

export type Currency = Entity & CurrencyNoID;
