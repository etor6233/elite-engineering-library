import { createElement } from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { expect, test, vi } from "vitest";
import AdminError from "@/app/admin/error";
import CustomerError from "@/app/customer/error";
import FactoryError from "@/app/factory/error";

test.each([
  { path: "/admin", component: AdminError },
  { path: "/customer", component: CustomerError },
  { path: "/factory", component: FactoryError }
])("$path recovery never displays the supplied failure or calls a retry", ({ path, component }) => {
  const reset = vi.fn(); const retry = vi.fn();
  const error = Object.assign(new Error("private-bearer customer-secret database-host"), { digest: "private-digest" });
  const html = renderToStaticMarkup(createElement(component, { ...{ error, reset, retry } }));
  expect(html).toContain('role="alert"');
  expect(html).toContain('href="'+path+'"');
  expect(html).toContain("Volver a consultar");
  expect(html).toContain("no confirma el resultado de una operación anterior");
  expect(html).not.toMatch(/private-bearer|customer-secret|database-host|private-digest|<form|<button/);
  expect(reset).not.toHaveBeenCalled(); expect(retry).not.toHaveBeenCalled();
});
