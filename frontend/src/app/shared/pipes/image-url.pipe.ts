import { Pipe, PipeTransform, inject } from '@angular/core';
import { ApiConfiguration } from '../../core/api/api-configuration';

@Pipe({
    name: 'imageUrl',
    standalone: true,
})
export class ImageUrlPipe implements PipeTransform {
    private readonly config = inject(ApiConfiguration);

    transform(imageId: string | null | undefined): string | null {
        if (!imageId) {
            return null;
        }

        // If it's already a data URL (e.g. preview), return as is
        if (imageId.startsWith('data:')) {
            return imageId;
        }

        // If it's already an absolute URL (unlikely for IDs, but safe check)
        if (imageId.startsWith('http')) {
            return imageId;
        }

        // Construct backend URL
        // rootUrl is typically '/api' or 'http://localhost:8080'.
        // Backend serves images at /images/{id}
        // If rootUrl is '/api' (proxy), result is '/api/images/{id}'
        // If rootUrl has trailing slash, remove it to be clean, though usually not strictly required if browser handles //
        const root = this.config.rootUrl.endsWith('/')
            ? this.config.rootUrl.slice(0, -1)
            : this.config.rootUrl;

        return `${root}/images/${imageId}`;
    }
}
