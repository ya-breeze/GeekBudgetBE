import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ImageUrlPipe } from '../../pipes/image-url.pipe';

@Component({
    selector: 'app-account-display',
    standalone: true,
    imports: [CommonModule, ImageUrlPipe],
    template: `
        <div class="account-display" [class]="size">
            @if (image) {
                <img [src]="image | imageUrl" alt="" class="account-icon" />
            }
            <span class="account-name">{{ name }}</span>
        </div>
    `,
    styleUrls: ['./account-display.component.scss'],
})
export class AccountDisplayComponent {
    @Input({ required: true }) name!: string;
    @Input() image?: string | null;
    @Input() size: 'small' | 'medium' | 'large' = 'small';
}
