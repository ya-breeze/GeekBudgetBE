"use client";

import { useMemo } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import { DataTable, SortableHeader } from "@/components/shared/data-table";
import { AccountDisplay, AccountTypeBadge } from "@/components/shared/account-display";
import type { Account } from "@/lib/api/models";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { MoreHorizontal, Pencil, Trash2, Image as ImageIcon } from "lucide-react";

interface AccountsTableProps {
  accounts: Account[];
  isLoading?: boolean;
  onEdit: (account: Account) => void;
  onDelete: (account: Account) => void;
  onImageUpload: (account: Account) => void;
}

export function AccountsTable({
  accounts,
  isLoading,
  onEdit,
  onDelete,
  onImageUpload,
}: AccountsTableProps) {
  const columns = useMemo<ColumnDef<Account>[]>(
    () => [
      {
        id: "account",
        header: ({ column }) => <SortableHeader column={column}>Account</SortableHeader>,
        accessorFn: (row) => row.name,
        cell: ({ row }) => <AccountDisplay account={row.original} />,
      },
      {
        id: "type",
        header: ({ column }) => <SortableHeader column={column}>Type</SortableHeader>,
        accessorKey: "type",
        cell: ({ row }) => <AccountTypeBadge type={row.original.type} />,
      },
      {
        id: "description",
        header: "Description",
        accessorKey: "description",
        cell: ({ row }) => (
          <span className="text-sm text-muted-foreground">
            {row.original.description || "-"}
          </span>
        ),
      },
      {
        id: "actions",
        header: "",
        cell: ({ row }) => (
          <div className="flex justify-end">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8">
                  <MoreHorizontal className="h-4 w-4" />
                  <span className="sr-only">Open menu</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => onEdit(row.original)}>
                  <Pencil className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => onImageUpload(row.original)}>
                  <ImageIcon className="mr-2 h-4 w-4" />
                  {row.original.image ? "Change Image" : "Upload Image"}
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  onClick={() => onDelete(row.original)}
                  className="text-destructive focus:text-destructive"
                >
                  <Trash2 className="mr-2 h-4 w-4" />
                  Delete
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        ),
      },
    ],
    [onEdit, onDelete, onImageUpload]
  );

  const mobileRenderer = (account: Account) => (
    <div className="flex items-start justify-between gap-3">
      <div className="flex items-center gap-2 min-w-0 flex-1">
        <AccountDisplay account={account} />
      </div>
      <div className="flex items-center gap-2 shrink-0">
        <AccountTypeBadge type={account.type} />
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="h-8 w-8">
              <MoreHorizontal className="h-4 w-4" />
              <span className="sr-only">Open menu</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => onEdit(account)}>
              <Pencil className="mr-2 h-4 w-4" />
              Edit
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => onImageUpload(account)}>
              <ImageIcon className="mr-2 h-4 w-4" />
              {account.image ? "Change Image" : "Upload Image"}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={() => onDelete(account)}
              className="text-destructive focus:text-destructive"
            >
              <Trash2 className="mr-2 h-4 w-4" />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );

  return (
    <DataTable
      columns={columns}
      data={accounts}
      isLoading={isLoading}
      emptyMessage="No accounts found. Create your first account to get started."
      mobileRenderer={mobileRenderer}
    />
  );
}
