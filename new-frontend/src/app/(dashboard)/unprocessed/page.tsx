"use client";

import { PageHeader } from "@/components/shared/page-header";

export default function UnprocessedPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Unprocessed Transactions"
        description="Review and process imported transactions"
      />
      <div className="rounded-lg border bg-card p-8 text-center text-muted-foreground">
        Unprocessed transactions coming soon (Phase 4)
      </div>
    </div>
  );
}
