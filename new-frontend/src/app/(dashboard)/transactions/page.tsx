"use client";

import { PageHeader } from "@/components/shared/page-header";

export default function TransactionsPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Transactions"
        description="Manage your financial transactions"
      />
      <div className="rounded-lg border bg-card p-8 text-center text-muted-foreground">
        Transactions coming soon (Phase 2)
      </div>
    </div>
  );
}
