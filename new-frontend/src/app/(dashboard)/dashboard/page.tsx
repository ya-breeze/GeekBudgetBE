"use client";

import { PageHeader } from "@/components/shared/page-header";

export default function DashboardPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Dashboard"
        description="Overview of your finances"
      />
      <div className="rounded-lg border bg-card p-8 text-center text-muted-foreground">
        Dashboard coming soon (Phase 3)
      </div>
    </div>
  );
}
