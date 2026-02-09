"use client";

import { useState } from "react";
import { toast } from "sonner";
import { Plus } from "lucide-react";
import { PageHeader } from "@/components/shared/page-header";
import { Button } from "@/components/ui/button";
import {
  useAccounts,
  useCreateAccount,
  useUpdateAccount,
  useDeleteAccount,
  useUploadAccountImage,
  useDeleteAccountImage,
} from "@/lib/api/hooks/use-accounts";
import type { Account, AccountNoID } from "@/lib/api/models";
import { AccountsTable } from "./_components/accounts-table";
import { AccountFormDialog } from "./_components/account-form-dialog";
import { AccountDeleteDialog } from "./_components/account-delete-dialog";
import { AccountImageDialog } from "./_components/account-image-dialog";

export default function AccountsPage() {
  const [formDialogOpen, setFormDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [imageDialogOpen, setImageDialogOpen] = useState(false);
  const [selectedAccount, setSelectedAccount] = useState<Account | null>(null);

  // Query
  const { data: accounts = [], isLoading } = useAccounts();

  // Mutations
  const createMutation = useCreateAccount({
    onSuccess: () => {
      toast.success("Account created successfully");
      setFormDialogOpen(false);
    },
    onError: (error: any) => {
      toast.error(error?.message || "Failed to create account");
    },
  });

  const updateMutation = useUpdateAccount({
    onSuccess: () => {
      toast.success("Account updated successfully");
      setFormDialogOpen(false);
      setSelectedAccount(null);
    },
    onError: (error: any) => {
      toast.error(error?.message || "Failed to update account");
    },
  });

  const deleteMutation = useDeleteAccount({
    onSuccess: () => {
      toast.success("Account deleted successfully");
      setDeleteDialogOpen(false);
      setSelectedAccount(null);
    },
    onError: (error: any) => {
      toast.error(error?.message || "Failed to delete account");
    },
  });

  const uploadImageMutation = useUploadAccountImage({
    onSuccess: () => {
      toast.success("Image uploaded successfully");
      setImageDialogOpen(false);
      setSelectedAccount(null);
    },
    onError: (error: any) => {
      toast.error(error?.message || "Failed to upload image");
    },
  });

  const deleteImageMutation = useDeleteAccountImage({
    onSuccess: () => {
      toast.success("Image removed successfully");
      setImageDialogOpen(false);
      setSelectedAccount(null);
    },
    onError: (error: any) => {
      toast.error(error?.message || "Failed to remove image");
    },
  });

  // Handlers
  const handleCreate = () => {
    setSelectedAccount(null);
    setFormDialogOpen(true);
  };

  const handleEdit = (account: Account) => {
    setSelectedAccount(account);
    setFormDialogOpen(true);
  };

  const handleDelete = (account: Account) => {
    setSelectedAccount(account);
    setDeleteDialogOpen(true);
  };

  const handleImageUpload = (account: Account) => {
    setSelectedAccount(account);
    setImageDialogOpen(true);
  };

  const handleFormSubmit = (data: AccountNoID) => {
    if (selectedAccount) {
      updateMutation.mutate({
        id: selectedAccount.id,
        data,
      });
    } else {
      createMutation.mutate(data);
    }
  };

  const handleDeleteConfirm = (replacementAccountId?: string) => {
    if (!selectedAccount) return;

    deleteMutation.mutate({
      id: selectedAccount.id,
      params: replacementAccountId && replacementAccountId !== "__none__"
        ? { replacementAccount: replacementAccountId }
        : undefined,
    });
  };

  const handleImageUploadConfirm = (file: File) => {
    if (!selectedAccount) return;

    uploadImageMutation.mutate({
      id: selectedAccount.id,
      file,
    });
  };

  const handleImageDeleteConfirm = () => {
    if (!selectedAccount) return;

    deleteImageMutation.mutate(selectedAccount.id);
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Accounts"
        description="Manage your asset, income, and expense accounts"
        actions={
          <Button onClick={handleCreate}>
            <Plus className="mr-2 h-4 w-4" />
            New Account
          </Button>
        }
      />

      <AccountsTable
        accounts={accounts}
        isLoading={isLoading}
        onEdit={handleEdit}
        onDelete={handleDelete}
        onImageUpload={handleImageUpload}
      />

      <AccountFormDialog
        open={formDialogOpen}
        onOpenChange={setFormDialogOpen}
        account={selectedAccount ?? undefined}
        isLoading={createMutation.isPending || updateMutation.isPending}
        onSubmit={handleFormSubmit}
      />

      <AccountDeleteDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        account={selectedAccount}
        availableAccounts={accounts}
        isLoading={deleteMutation.isPending}
        onConfirm={handleDeleteConfirm}
      />

      <AccountImageDialog
        open={imageDialogOpen}
        onOpenChange={setImageDialogOpen}
        account={selectedAccount}
        isUploading={uploadImageMutation.isPending}
        isDeleting={deleteImageMutation.isPending}
        onUpload={handleImageUploadConfirm}
        onDelete={handleImageDeleteConfirm}
      />
    </div>
  );
}
