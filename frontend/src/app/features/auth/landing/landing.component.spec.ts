import { ComponentFixture, TestBed } from '@angular/core/testing';
import { LandingComponent } from './landing.component';
import { provideRouter } from '@angular/router';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

describe('LandingComponent', () => {
    let component: LandingComponent;
    let fixture: ComponentFixture<LandingComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [LandingComponent, NoopAnimationsModule],
            providers: [provideRouter([])],
        }).compileComponents();

        fixture = TestBed.createComponent(LandingComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });

    it('should have features', () => {
        expect(component.features.length).toBe(6);
    });

    it('should have steps', () => {
        expect(component.steps.length).toBe(4);
    });

    it('should have concepts', () => {
        expect(component.concepts.length).toBe(6);
    });
});
