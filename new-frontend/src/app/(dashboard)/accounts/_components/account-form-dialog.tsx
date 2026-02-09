"use client";

import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { ChevronDown } from "lucide-react";
import type { Account, AccountNoID } from "@/lib/api/models";
import { AccountNoIDType } from "@/lib/api/models";

const accountFormSchema = z.object({
  name: z.string().min(1, "Name is required"),
  description: z.string().optional(),
  type: z.enum([AccountNoIDType.asset, AccountNoIDType.income, AccountNoIDType.expense]),
  showInDashboardSummary: z.boolean().default(true),
  hideFromReports: z.boolean().default(false),
  bankInfo: z.object({
    accountId: z.string().optional(),
    bankId: z.string().optional(),
  }).optional(),
  openingDate: z.string().optional(),
  closingDate: z.string().optional(),
  ignoreUnprocessedBefore: z.string().optional(),
});

type AccountFormValues = z.infer<typeof accountFormSchema>;

interface AccountFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  account?: Account;
  isLoading?: boolean;
  onSubmit: (data: AccountNoID) => void;
}

export function AccountFormDialog({
  open,
  onOpenChange,
  account,
  isLoading,
  onSubmit,
}: AccountFormDialogProps) {
  const [bankInfoOpen, setBankInfoOpen] = useState(false);
  const [advancedOpen, setAdvancedOpen] = useState(false);

  const form = useForm<AccountFormValues>({
    resolver: zodResolver(accountFormSchema),
    defaultValues: {
      name: "",
      description: "",
      type: AccountNoIDType.asset,
      showInDashboardSummary: true,
      hideFromReports: false,
      bankInfo: {
        accountId: "",
        bankId: "",
      },
      openingDate: "",
      closingDate: "",
      ignoreUnprocessedBefore: "",
    },
  });

  // Reset form when account changes or dialog opens
  useEffect(() => {
    if (open) {
      if (account) {
        form.reset({
          name: account.name,
          description: account.description || "",
          type: account.type,
          showInDashboardSummary: account.showInDashboardSummary ?? true,
          hideFromReports: account.hideFromReports ?? false,
          bankInfo: {
            accountId: account.bankInfo?.accountId || "",
            bankId: account.bankInfo?.bankId || "",
          },
          openingDate: account.openingDate || "",
          closingDate: account.closingDate || "",
          ignoreUnprocessedBefore: account.ignoreUnprocessedBefore || "",
        });
        setBankInfoOpen(!!(account.bankInfo?.accountId || account.bankInfo?.bankId));
        setAdvancedOpen(!!(account.openingDate || account.closingDate || account.ignoreUnprocessedBefore));
      } else {
        form.reset({
          name: "",
          description: "",
          type: AccountNoIDType.asset,
          showInDashboardSummary: true,
          hideFromReports: false,
          bankInfo: { accountId: "", bankId: "" },
          openingDate: "",
          closingDate: "",
          ignoreUnprocessedBefore: "",
        });
        setBankInfoOpen(false);
        setAdvancedOpen(false);
      }
    }
  }, [account, open, form]);

  const handleSubmit = (values: AccountFormValues) => {
    const payload: AccountNoID = {
      name: values.name,
      description: values.description || undefined,
      type: values.type,
      showInDashboardSummary: values.showInDashboardSummary,
      hideFromReports: values.hideFromReports,
      bankInfo:
        values.bankInfo?.accountId || values.bankInfo?.bankId
          ? {
              accountId: values.bankInfo.accountId || undefined,
              bankId: values.bankInfo.bankId || undefined,
            }
          : undefined,
      openingDate: values.openingDate || undefined,
      closingDate: values.closingDate || undefined,
      ignoreUnprocessedBefore: values.ignoreUnprocessedBefore || undefined,
    };
    onSubmit(payload);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{account ? "Edit Account" : "Create Account"}</DialogTitle>
          <DialogDescription>
            {account
              ? "Update the account details below."
              : "Add a new account to track your finances."}
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input placeholder="e.g., Cash, Bank Account, Groceries" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Type</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select account type" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value={AccountNoIDType.asset}>Asset</SelectItem>
                      <SelectItem value={AccountNoIDType.income}>Income</SelectItem>
                      <SelectItem value={AccountNoIDType.expense}>Expense</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormDescription>
                    Assets track what you own, income tracks where money comes from, expenses track where it goes.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Optional notes about this account"
                      className="resize-none"
                      rows={2}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="space-y-3">
              <FormField
                control={form.control}
                name="showInDashboardSummary"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between rounded-lg border p-3">
                    <div className="space-y-0.5">
                      <FormLabel>Show in Dashboard Summary</FormLabel>
                      <FormDescription className="text-xs">
                        Display this account in the dashboard overview.
                      </FormDescription>
                    </div>
                    <FormControl>
                      <Switch checked={field.value} onCheckedChange={field.onChange} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="hideFromReports"
                render={({ field }) => (
                  <FormItem className="flex items-center justify-between rounded-lg border p-3">
                    <div className="space-y-0.5">
                      <FormLabel>Hide from Reports</FormLabel>
                      <FormDescription className="text-xs">
                        Exclude this account from reports and budget calculations.
                      </FormDescription>
                    </div>
                    <FormControl>
                      <Switch checked={field.value} onCheckedChange={field.onChange} />
                    </FormControl>
                  </FormItem>
                )}
              />
            </div>

            {/* Bank Info Section */}
            <Collapsible open={bankInfoOpen} onOpenChange={setBankInfoOpen}>
              <CollapsibleTrigger asChild>
                <Button
                  variant="outline"
                  className="w-full justify-between"
                  type="button"
                >
                  Bank Information (Optional)
                  <ChevronDown
                    className={`h-4 w-4 transition-transform ${
                      bankInfoOpen ? "rotate-180" : ""
                    }`}
                  />
                </Button>
              </CollapsibleTrigger>
              <CollapsibleContent className="space-y-4 pt-4">
                <FormField
                  control={form.control}
                  name="bankInfo.accountId"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Bank Account ID</FormLabel>
                      <FormControl>
                        <Input placeholder="e.g., 123456789/0100" {...field} />
                      </FormControl>
                      <FormDescription className="text-xs">
                        Account number used for bank imports
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="bankInfo.bankId"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Bank ID</FormLabel>
                      <FormControl>
                        <Input placeholder="e.g., fio, kb, revolut" {...field} />
                      </FormControl>
                      <FormDescription className="text-xs">
                        Bank identifier for automatic imports
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </CollapsibleContent>
            </Collapsible>

            {/* Advanced Options */}
            <Collapsible open={advancedOpen} onOpenChange={setAdvancedOpen}>
              <CollapsibleTrigger asChild>
                <Button
                  variant="outline"
                  className="w-full justify-between"
                  type="button"
                >
                  Advanced Options
                  <ChevronDown
                    className={`h-4 w-4 transition-transform ${
                      advancedOpen ? "rotate-180" : ""
                    }`}
                  />
                </Button>
              </CollapsibleTrigger>
              <CollapsibleContent className="space-y-4 pt-4">
                <FormField
                  control={form.control}
                  name="openingDate"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Opening Date</FormLabel>
                      <FormControl>
                        <Input type="date" {...field} />
                      </FormControl>
                      <FormDescription className="text-xs">
                        Account is ignored before this date
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="closingDate"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Closing Date</FormLabel>
                      <FormControl>
                        <Input type="date" {...field} />
                      </FormControl>
                      <FormDescription className="text-xs">
                        Account is ignored after this date
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="ignoreUnprocessedBefore"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Ignore Unprocessed Before</FormLabel>
                      <FormControl>
                        <Input type="date" {...field} />
                      </FormControl>
                      <FormDescription className="text-xs">
                        Skip unprocessed transactions older than this date
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </CollapsibleContent>
            </Collapsible>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isLoading}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {account ? "Save Changes" : "Create Account"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
