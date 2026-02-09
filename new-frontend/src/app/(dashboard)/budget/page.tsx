"use client";

import { PageHeader } from "@/components/shared/page-header";

export default function BudgetPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Budget"
        description="Track your budget and spending"
      />
      <div className="rounded-lg border bg-card p-8 text-center text-muted-foreground">
        Budget coming soon (Phase 5)
      </div>
    </div>
  );
}
