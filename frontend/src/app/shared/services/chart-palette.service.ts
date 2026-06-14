import { Injectable } from '@angular/core';

const PALETTE: string[] = [
    '#4E79A7',
    '#A0CBE8', // Blue
    '#F28E2B',
    '#FFBE7D', // Orange
    '#59A14F',
    '#8CD17D', // Green
    '#B6992D',
    '#F1CE63', // Gold
    '#499894',
    '#86BCB6', // Teal
    '#E15759',
    '#FF9D9A', // Red
    '#79706E',
    '#BAB0AC', // Gray
    '#D37295',
    '#FABFD2', // Pink
    '#B07AA1',
    '#D4A6C8', // Purple
    '#9D7660',
    '#D7B5A6', // Brown
    '#17BECF',
    '#9EDAE5', // Cyan
    '#BCBD22', // Chartreuse
    '#393B79', // Dark indigo
    '#CE6DBD', // Magenta
];

@Injectable({ providedIn: 'root' })
export class ChartPaletteService {
    getColor(index: number): string {
        return PALETTE[index % PALETTE.length];
    }

    getColorWithAlpha(index: number, alpha: number): string {
        const hex = PALETTE[index % PALETTE.length];
        const r = parseInt(hex.slice(1, 3), 16);
        const g = parseInt(hex.slice(3, 5), 16);
        const b = parseInt(hex.slice(5, 7), 16);
        return `rgba(${r}, ${g}, ${b}, ${alpha})`;
    }

    colorsForN(n: number): string[] {
        return Array.from({ length: n }, (_, i) => this.getColor(i));
    }
}
