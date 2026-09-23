import { experienceText } from "@/platform/experience/messages";

// AUTHORED presentation labels only. Domain commands and identifiers keep their owners.
const english: Record<string, string> = {
    "Operación": "Operations",
    "Gestión": "Management",
    "Espacios": "Workspaces",
    "Tu espacio de trabajo": "Your workspace",
    "Un lugar para cada tarea.": "A place for every task.",
    "Accesos rápidos": "Quick access",
    "Continuar entrega": "Continue delivery",
    "Revisá, presentá y confirmá.": "Review, present and confirm.",
    "Ver pedidos": "View orders",
    "Abrir agenda": "Open schedule",
    "Atención al cliente": "Customer service",
    "Consultar y dar seguimiento": "Review and follow up",
    "Organizar el día": "Plan the day",
    "Ayuda y consultas": "Help and enquiries",
    "Volver a operaciones": "Back to operations",
};

export function adminChromeText(text: string, language: string): string {
    return language === "en" ? english[text] ?? experienceText(text, "en") : text;
}
