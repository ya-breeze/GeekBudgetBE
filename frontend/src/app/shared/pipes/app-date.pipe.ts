import { Pipe, PipeTransform } from '@angular/core';
import { DatePipe } from '@angular/common';

@Pipe({
    name: 'appDate',
    standalone: true,
})
export class AppDatePipe implements PipeTransform {
    private readonly datePipe = new DatePipe('en-US');

    transform(value: any, format = 'dd/MM/yyyy', notSetLabel = 'Not set'): string {
        if (!value) {
            return notSetLabel;
        }

        const dateStr = value.toString();
        if (dateStr.startsWith('0001-01-01')) {
            return notSetLabel;
        }

        const date = new Date(value);
        if (isNaN(date.getTime())) {
            return notSetLabel;
        }

        // Check for 1/1/1 or similar if the input was already a Date object
        if (date.getFullYear() <= 1) {
            return notSetLabel;
        }

        return this.datePipe.transform(value, format) || notSetLabel;
    }
}
