import { Component, Input, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { OverlayModule } from '@angular/cdk/overlay';
import { ImageUrlPipe } from '../../pipes/image-url.pipe';

@Component({
    selector: 'app-account-display',
    standalone: true,
    imports: [CommonModule, OverlayModule, ImageUrlPipe],
    template: `
        <div class="account-display" [class]="size">
            @if (image) {
                <img
                    [src]="image | imageUrl"
                    alt=""
                    class="account-icon"
                    cdkOverlayOrigin
                    #trigger="cdkOverlayOrigin"
                    (mouseenter)="isHovered.set(true)"
                    (mouseleave)="isHovered.set(false)"
                />

                <ng-template
                    cdkConnectedOverlay
                    [cdkConnectedOverlayOrigin]="trigger"
                    [cdkConnectedOverlayOpen]="isHovered()"
                    [cdkConnectedOverlayHasBackdrop]="false"
                    [cdkConnectedOverlayPositions]="overlayPositions"
                >
                    <div class="hover-image-hint">
                        <img [src]="image | imageUrl" alt="" class="hint-large-image" />
                    </div>
                </ng-template>
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

    protected readonly isHovered = signal(false);

    protected readonly overlayPositions = [
        {
            originX: 'center' as const,
            originY: 'bottom' as const,
            overlayX: 'center' as const,
            overlayY: 'top' as const,
            offsetY: 8,
        },
        {
            originX: 'center' as const,
            originY: 'top' as const,
            overlayX: 'center' as const,
            overlayY: 'bottom' as const,
            offsetY: -8,
        },
    ];
}
