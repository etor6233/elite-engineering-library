// AUTHORED redirect contract; all payment entry remains on the provider UI.
import { z } from "zod";
export const checkoutID = z.string().min(1).max(200).regex(/^[A-Za-z0-9_-]+$/);
export const checkoutSchema = z.object({
 order_id: checkoutID,
 provider_code: z.enum(["stripe", "mercadopago"]),
 url: z.string().min(1).max(8192),
 expires_at: z.iso.datetime({ offset: true }),
}).strict().superRefine((value,ctx) => {
 let target: URL;
 try { target = new URL(value.url); } catch { ctx.addIssue({ code:"custom", message:"invalid checkout URL" }); return; }
 const hosts = value.provider_code === "stripe" ? ["checkout.stripe.com"] : ["www.mercadopago.com.ar", "sandbox.mercadopago.com.ar"];
 if (target.protocol !== "https:" || target.username || target.password || target.hash || (target.port && target.port !== "443") || !hosts.includes(target.hostname) || /[\r\n]/.test(value.url) || Date.parse(value.expires_at) <= Date.now()) ctx.addIssue({code:"custom",message:"checkout is not available"});
});
