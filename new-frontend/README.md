# GeekBudget Frontend (Next.js 15)

Modern Next.js 15 frontend for GeekBudget personal finance management.

## Tech Stack

- **Next.js 15.5** - React framework with App Router
- **React 19** - UI library
- **TypeScript 5.5** - Type safety
- **Tailwind CSS 3.4** - Utility-first styling
- **shadcn/ui** - Component library built on Radix UI
- **TanStack Query** - Data fetching and caching
- **React Hook Form** - Form management
- **Zod** - Schema validation
- **Sonner** - Toast notifications
- **next-themes** - Dark mode support

## Getting Started

### Prerequisites

- Node.js 20+ and npm
- Backend API running on `http://localhost:8080` (or configure `NEXT_PUBLIC_API_URL`)

### Installation

```bash
# Install dependencies
npm install --legacy-peer-deps

# Note: --legacy-peer-deps is needed due to cmdk requiring React 18
# while we use React 19. This is safe and the app works correctly.
```

### Development

```bash
# Start development server
npm run dev

# Open http://localhost:3000
```

### Build

```bash
# Create production build
npm run build

# Start production server
npm start
```

### Linting

```bash
npm run lint
```

## Environment Variables

Create a `.env.local` file (already created):

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Project Structure

```
src/
├── app/                    # Next.js App Router pages
│   ├── (dashboard)/        # Dashboard layout group
│   │   ├── accounts/       # Accounts page ✅ COMPLETE
│   │   ├── budget/         # Budget page (Phase 5)
│   │   ├── dashboard/      # Dashboard page (Phase 3)
│   │   ├── transactions/   # Transactions page (Phase 2)
│   │   └── unprocessed/    # Unprocessed page (Phase 4)
│   ├── globals.css         # Global styles & Tailwind
│   ├── layout.tsx          # Root layout
│   ├── page.tsx            # Home (redirects to /dashboard)
│   └── providers.tsx       # React Query & Theme providers
├── components/
│   ├── shared/             # Reusable components
│   │   ├── account-display.tsx
│   │   ├── amount-display.tsx
│   │   ├── confirm-dialog.tsx
│   │   ├── data-table.tsx
│   │   ├── month-navigator.tsx
│   │   └── page-header.tsx
│   └── ui/                 # shadcn/ui components (15 components)
└── lib/
    ├── api/                # API client & generated hooks
    │   ├── client.ts       # Axios instance
    │   ├── models/         # TypeScript types
    │   └── generated/      # API hooks
    │       ├── accounts/
    │       └── currencies/
    ├── hooks/              # Custom hooks
    │   ├── use-accounts.ts
    │   └── use-currencies.ts
    └── utils/
        ├── format.ts       # formatAmount, formatDate
        └── utils.ts        # cn() utility

```

## Implementation Status

| Phase | View | Status |
|---|---|---|
| 0 | Shared Components | ✅ **DONE** |
| 1 | Accounts | ✅ **DONE** |
| 2 | Transactions | 📋 Pending |
| 3 | Dashboard | 📋 Pending |
| 4 | Unprocessed | 📋 Pending |
| 5 | Budget | 📋 Pending |

### Phase 1: Accounts ✅

The Accounts view is fully implemented with:

- **Full CRUD** - Create, read, update, delete accounts
- **Image uploads** - Account avatars with upload/delete
- **Advanced options** - Bank info, opening/closing dates
- **Form validation** - React Hook Form + Zod schemas
- **Responsive design** - Mobile cards, desktop table
- **Type badges** - Color-coded asset/income/expense
- **Toast notifications** - Success/error feedback
- **Query invalidation** - Real-time updates

## Features

### Built-in

- 🌓 **Dark mode** - Automatic system detection + manual toggle
- 📱 **Responsive** - Mobile-first design with breakpoints
- ⚡ **Fast** - Static generation + client-side navigation
- ♿ **Accessible** - WCAG compliant via Radix UI
- 🎨 **Themeable** - CSS variables for easy customization

### Dashboard Layout

- Collapsible sidebar navigation
- Mobile-friendly hamburger menu
- Theme toggle in sidebar
- Active route highlighting

## API Integration

The app uses React Query for data fetching with:

- Automatic caching and background refetching
- Optimistic updates
- Query invalidation on mutations
- Error handling

### API Routing

See [API_ROUTING.md](./API_ROUTING.md) for detailed routing configuration (development vs production).

## Development Notes

### API Client Stubs

The current API client uses stub implementations. When the OpenAPI spec is available, regenerate with:

```bash
# From main repo
make generate

# Copy generated files to new-frontend/src/lib/api/
```

### Next Steps

1. **Phase 2: Transactions** - Most complex view with filters, tags, movements
2. **Phase 3: Dashboard** - Asset cards, expense heatmap, charts
3. **Phase 4: Unprocessed** - Transaction processing workflow
4. **Phase 5: Budget** - Budget matrix with inline editing

### Known Issues

- ESLint warnings about `<img>` vs `<Image />` - cosmetic, not blocking
- cmdk peer dependency warning - safe to ignore with `--legacy-peer-deps`

## Contributing

When adding new pages:

1. Create route in `src/app/(dashboard)/[page-name]/`
2. Add to navigation in `src/app/(dashboard)/layout.tsx`
3. Create `_components/` for page-specific components
4. Use shared components from `src/components/shared/`
5. Add API hooks if needed
6. Test responsive design at mobile widths

## License

See main GeekBudget repository for license information.
