"use client";

import { PrivateReadFailure } from "@/platform/backend/private-read-failure";

export default function PrivatePortalError() {
  return <PrivateReadFailure path="/admin" />;
}
