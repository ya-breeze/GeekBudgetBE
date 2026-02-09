"use client";

import { useState, useRef } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Loader2, Upload, X } from "lucide-react";
import type { Account } from "@/lib/api/models";
import { Avatar } from "@/components/ui/avatar";
import { getApiUrl } from "@/lib/utils";

interface AccountImageDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  account: Account | null;
  isUploading?: boolean;
  isDeleting?: boolean;
  onUpload: (file: File) => void;
  onDelete: () => void;
}

export function AccountImageDialog({
  open,
  onOpenChange,
  account,
  isUploading,
  isDeleting,
  onUpload,
  onDelete,
}: AccountImageDialogProps) {
  const [preview, setPreview] = useState<string | null>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  if (!account) return null;

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleUpload = () => {
    if (selectedFile) {
      onUpload(selectedFile);
    }
  };

  const handleClose = () => {
    setPreview(null);
    setSelectedFile(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
    onOpenChange(false);
  };

  const currentImageUrl = account.image
    ? getApiUrl(`/v1/accounts/${account.id}/image`)
    : null;

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Account Image</DialogTitle>
          <DialogDescription>
            Upload an image for {account.name}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Current or Preview Image */}
          <div className="flex justify-center">
            {preview ? (
              <Avatar className="h-32 w-32">
                <img
                  src={preview}
                  alt="Preview"
                  className="h-full w-full object-cover rounded-full"
                />
              </Avatar>
            ) : currentImageUrl ? (
              <Avatar className="h-32 w-32">
                <img
                  src={currentImageUrl}
                  alt={account.name}
                  className="h-full w-full object-cover rounded-full"
                />
              </Avatar>
            ) : (
              <div className="flex h-32 w-32 items-center justify-center rounded-full bg-muted">
                <Upload className="h-12 w-12 text-muted-foreground" />
              </div>
            )}
          </div>

          {/* File Input */}
          <div>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              onChange={handleFileChange}
              className="hidden"
              id="image-upload"
            />
            <Button
              variant="outline"
              className="w-full"
              onClick={() => fileInputRef.current?.click()}
              disabled={isUploading || isDeleting}
            >
              <Upload className="mr-2 h-4 w-4" />
              {selectedFile ? "Choose Different File" : "Choose File"}
            </Button>
            {selectedFile && (
              <p className="mt-2 text-xs text-muted-foreground">
                Selected: {selectedFile.name}
              </p>
            )}
          </div>
        </div>

        <DialogFooter className="flex-col gap-2 sm:flex-row">
          {currentImageUrl && !preview && (
            <Button
              variant="destructive"
              onClick={onDelete}
              disabled={isUploading || isDeleting}
              className="w-full sm:w-auto"
            >
              {isDeleting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              <X className="mr-2 h-4 w-4" />
              Remove Image
            </Button>
          )}
          <div className="flex gap-2 w-full sm:w-auto">
            <Button
              variant="outline"
              onClick={handleClose}
              disabled={isUploading || isDeleting}
              className="flex-1 sm:flex-initial"
            >
              Cancel
            </Button>
            <Button
              onClick={handleUpload}
              disabled={!selectedFile || isUploading || isDeleting}
              className="flex-1 sm:flex-initial"
            >
              {isUploading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Upload
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
