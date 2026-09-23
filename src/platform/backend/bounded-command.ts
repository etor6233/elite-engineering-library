// AUTHORED shared body bound, extracted without semantic change from the handover BFF.
export async function boundedRequestBytes(request: Request, maxBytes: number) {
 if (!Number.isSafeInteger(maxBytes) || maxBytes < 1 || maxBytes > 1048576) throw new Error("invalid body budget");
 if (!request.body) throw new Error("invalid handover body");
 const reader = request.body.getReader(); const chunks: Uint8Array[] = []; let bytes = 0; let timer: ReturnType<typeof setTimeout> | undefined;
 const deadline = new Promise<never>((_, reject) => { timer = setTimeout(() => reject(new Error("body deadline")), 2500); });
 try {
  while (true) {
   const part = await Promise.race([reader.read(), deadline]); if (part.done) break;
   bytes += part.value.byteLength; if (bytes > maxBytes) throw new Error("BODY_TOO_LARGE"); chunks.push(part.value);
  }
  const data = new Uint8Array(bytes); let offset = 0; for (const chunk of chunks) { data.set(chunk, offset); offset += chunk.byteLength; }
  return data;
 } finally { clearTimeout(timer); void reader.cancel().catch(() => {}); reader.releaseLock(); }
}

// Existing text budgets and UTF-8 behavior remain unchanged.
export async function boundedCommandBody(request: Request, maxBytes = 4096) {
 if (!Number.isSafeInteger(maxBytes) || maxBytes < 1 || maxBytes > 65536) throw new Error("invalid body budget");
 return new TextDecoder("utf-8", { fatal: true }).decode(await boundedRequestBytes(request,maxBytes));
}
