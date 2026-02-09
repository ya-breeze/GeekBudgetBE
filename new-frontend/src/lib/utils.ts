import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// Get the correct API URL for fetching resources (like images)
// In dev: points to localhost:8080
// In production: uses /api prefix (rewritten by Next.js)
export function getApiUrl(path: string): string {
  const isProduction = process.env.NODE_ENV === "production";
  const baseUrl = isProduction ? "/api" : (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080");
  return `${baseUrl}${path}`;
}
