// AUTHORED protocol guard. Only explicit same-origin pre-effect rejection releases a local fence.
export function definitiveSelectionRejection(response: Pick<Response, "status">, body: unknown): boolean {
    if (!body || typeof body !== "object")
        return false;
    const value = body as Record<string, unknown>;
    return value.operation_effect === "NOT_ATTEMPTED" && ((response.status === 403 && value.code === "SELECTION_NOT_AUTHORIZED") || (response.status === 503 && value.code === "SELECTION_CONFIGURATION_UNAVAILABLE"));
}
