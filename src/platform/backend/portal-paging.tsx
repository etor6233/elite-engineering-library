import{PrivatePagingLinks}from"./private-paging-links";
export type PortalQuery = Record<string, string | string[] | undefined>;
export type PortalPageProps = { searchParams?: Promise<PortalQuery> };
type CursorName = "orders_after" | "cases_after" | "units_after" | "leads_after";

export function portalCursor(value: string | string[] | undefined): string | undefined {
  if (value === undefined || value === "") return undefined;
  if (typeof value !== "string" || value.length > 512) throw new Error("Invalid page cursor");
  return value;
}

export function portalPageHref(path: "/customer" | "/factory" | "/admin", cursors: Partial<Record<CursorName, string | undefined>>): string {
  const query = new URLSearchParams();
  for (const name of path === "/admin" ? ["orders_after", "cases_after", "leads_after"] as const : path === "/customer" ? ["orders_after", "cases_after"] as const : ["units_after"] as const) {
    const value = portalCursor(cursors[name]);
    if (value !== undefined) query.set(name, value);
  }
  return path + (query.size ? "?" + query.toString() : "");
}

export function PortalPaging(props:{label:string;nextHref:string|undefined;firstHref:string|undefined}){return <PrivatePagingLinks {...props}/>;}
