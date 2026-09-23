// Release-bound PUBLIC operational guidance, shared with inline help. Not tenant documents.
import { STATUS_VIEW_VERSION } from "@/platform/notifications/status-contract";
export const QUOTE_GUIDE = {
  id: "quote-acceptance-view", version: "1.0.0", title: "Aceptar una cotización",
  paragraphs: ["Revisá el importe, la moneda y la vigencia. Aceptar solicita crear un pedido con esos datos del servidor.", "Aceptar una cotización no confirma el pago, el stock ni la entrega.", "Si se corta la conexión, usá Actualizar estado o volvé a abrir Mis cotizaciones. No repitas la aceptación hasta comprobar el estado.", "Un pedido se confirma aquí cuando la lectura muestra su referencia. Si el problema continúa, contactá al soporte de la empresa con la referencia de cotización y organización, sólo por un canal autorizado; no compartas contraseñas ni tokens.", "Práctica de esta versión en el entorno de capacitación: reconocer un resultado incierto, recuperar la lectura y ubicar el pedido sin reenviar. No practiques creando pedidos en producción."]
} as const;
export const NOTIFICATION_GUIDE = {
  id: "whatsapp-status-view", version: STATUS_VIEW_VERSION.split("/")[1], title: "Consultar el estado de WhatsApp",
  paragraphs: ["Compartí esta referencia sólo con soporte autorizado de tu organización. No adjuntes teléfonos, tokens ni el contenido del mensaje.", "Si el resultado es incierto o contradictorio: preservá el intento, consultá su evidencia y no reenvíes ni borres el historial. La entrega informada no confirma una venta ni la aceptación del turno."]
} as const;

// Extracted verbatim from the existing versioned inline guidance, V380.
export const OPERATIONAL_GUIDES = {
  "availability-cancel-view": {
    "id": "availability-cancel-view",
    "version": "1.0.0",
    "title": "Cancelar un intervalo",
    "summary": "Ayuda para cancelar un intervalo",
    "paragraphs": [
      "Cancelar este registro no cancela ni reprograma citas de clientes. El servidor comprueba los permisos, la versión y las restricciones existentes. Revisá recurso, fechas y tipo antes de actuar.",
      "Si se pierde la respuesta, Consultar intervalo sólo lee el servidor. Una cancelación consultada no demuestra quién la solicitó; la auditoría autorizada conserva al actor. La referencia de esta pestaña guarda sólo la versión, no motivos, horarios ni datos personales.",
      "Si la consulta falla, el intervalo no aparece en la ventana consultada o permanece activo, conservamos el bloqueo: solicitá al soporte autorizado revisar la referencia, hora y estado. No borres la referencia ni abras otra pestaña para forzar un reintento.",
      "Práctica en pruebas: cancelar un intervalo sin citas, perder la respuesta y consultar Cancelado sin reenviar. No practicar con la agenda real. Reactivar o crear otro intervalo es una operación distinta."
    ],
    "inlineLead": false
  },
  "availability-create-view": {
    "id": "availability-create-view",
    "version": "1.0.0",
    "title": "Registrar un intervalo",
    "summary": "Ayuda para registrar un intervalo",
    "paragraphs": [
      "Conservamos sólo la referencia de operación. Una respuesta perdida no demuestra rechazo: consultá antes de volver a actuar. La consulta muestra el registro aunque quede fuera del rango de la agenda. Preparar otro intervalo no reenvía ni cancela el anterior y sólo se habilita tras recuperarlo. Ante referencia ilegible o consulta fallida, pedí revisión autorizada sin borrar datos. Practicá con datos sintéticos el envío, pérdida de respuesta, consulta y preparación explícita de otro intervalo."
    ],
    "inlineLead": false
  },
  "checklist-completion-view": {
    "id": "checklist-completion-view",
    "version": "1.0.0",
    "title": "Presentar una entrega",
    "summary": "Ayuda para presentar una entrega",
    "paragraphs": [
      "La presentación conserva respuestas inmutables. Guardamos sólo identidades, versiones y una huella; no guardamos las respuestas ni evidencias en la referencia del navegador. Una respuesta perdida puede ocultar un cambio correcto: consultá sin reenviarlo. El estado actual puede avanzar a aceptado o rechazado; consultar o preparar otro formulario no lo cambia. Si el contenido no coincide, pedí revisión. Practicá con datos sintéticos antes de operar."
    ],
    "inlineLead": false
  },
  "checklist-publication-view": {
    "id": "checklist-publication-view",
    "version": "1.0.0",
    "title": "Publicar una checklist",
    "summary": "Ayuda para publicar una checklist",
    "paragraphs": [
      "Una versión publicada es inmutable. Conservamos sólo su ID, versión y huella de contenido; no guardamos el título ni los textos en la referencia del navegador. Una respuesta perdida puede ocultar una publicación correcta: consultá sin volver a enviarla. Si no coincide el contenido o no se puede recuperar, pedí revisión sin borrar la referencia. Preparar otra versión no cambia ni vuelve a publicar la anterior. Practicá este recorrido con datos sintéticos antes de usarlo en operación."
    ],
    "inlineLead": false
  },
  "delivery-resolution-view": {
    "id": "delivery-resolution-view",
    "version": "1.0.0",
    "title": "Resolver una discrepancia",
    "summary": "Ayuda para resolver una discrepancia",
    "paragraphs": [
      "Una respuesta perdida no demuestra rechazo. Conservamos en esta pestaña sólo la versión y la decisión pendiente, sin notas ni credenciales. Consultar vuelve a leer el servidor sin reenviar la acción.",
      "Si la consulta falla, el registro no aparece o hay otra decisión, pedí revisión al responsable autorizado. No borres la referencia ni abras otra pestaña para forzar un reintento. Informá la referencia de la entrega y la hora, sin compartir notas privadas ni tokens.",
      "Práctica: en el entorno de prueba, recuperá una resolución tras perder la respuesta y comprobá su estado antes del siguiente paso. Preparada o autorizada no significa entrega física, reembolso ni cambio finalizados."
    ],
    "inlineLead": false
  },
  "handover-read-view": {
    "id": "handover-read-view",
    "version": "1.0.0",
    "title": "Consultar entregas",
    "summary": "Ayuda para consultar entregas",
    "paragraphs": [
      "Si la información no está disponible, volvé a consultar antes de aceptar o rechazar. No crees otra solicitud para recuperar una lectura.",
      "Si el problema continúa, contactá al soporte autorizado e indicá la acción y la hora. No envíes contraseñas, tokens ni información de otras personas.",
      "Práctica en pruebas: interrumpí la consulta, verificá que no haya acciones de entrega disponibles y restablecé la conexión. La próxima consulta debe recuperar el estado sin enviar una aceptación automática."
    ],
    "inlineLead": true
  },
  "lead-command-view": {
    "id": "lead-command-view",
    "version": "1.0.0",
    "title": "Actualizar un lead",
    "summary": "Ayuda para actualizar un lead",
    "paragraphs": [
      "Una respuesta perdida no demuestra que el cambio haya sido rechazado. Consultar lee el servidor sin reenviar comandos. Una versión posterior permite revisar el estado real; no prueba que haya sido modificado exclusivamente por tu solicitud.",
      "Conservamos sólo versión y tipo de acción en esta pestaña, no responsables ni datos de contacto. Si falta el lead, falla su consulta o la versión no avanzó, pedí revisión al soporte autorizado; no borres la referencia ni abras otra pestaña para forzar un reintento.",
      "Práctica: recuperá un cambio en el entorno de prueba, comprobá responsable y estado y continuá con la versión consultada. No asumas una venta o un cliente convertido por una asignación."
    ],
    "inlineLead": false
  },
  "operation-sections-view": {
    "id": "operation-sections-view",
    "version": "1.0.0",
    "title": "Recuperar una sección",
    "summary": "Ayuda para recuperar una sección",
    "paragraphs": [
      "Una consulta fallida no significa que no existan registros. Volvé a consultar antes de decidir sobre esa sección; las secciones independientes conservan sus propios permisos y estados.",
      "La consulta no reenvía acciones. Si el problema persiste, informá al soporte autorizado la sección, la hora y la operación esperada; no compartas tokens, contraseñas ni datos personales en capturas.",
      "Práctica: ante una consulta fallida, identificá la sección, recuperala mediante consulta y verificá su estado antes de actuar."
    ],
    "inlineLead": false
  },
  "order-operations-view": {
    "id": "order-operations-view",
    "version": "1.1.0",
    "title": "Continuar un pedido",
    "summary": "Ayuda para continuar un pedido",
    "paragraphs": [
      "Elegí una unidad por su serie y reservá una sola vez. El servidor vuelve a comprobar organización, variante, disponibilidad y versiones.",
      "Si se pierde la respuesta o alguien actuó primero, usá Actualizar pedidos. No crees otro pedido ni cambies identificadores para forzar la operación.",
      "Pago: registrar una solicitud guarda una intención del total; no confirma dinero recibido ni libera entrega. Una solicitud existente se consulta, no se duplica. Si falla, no borres la clave del navegador: consultá el estado y escalá al responsable antes de intentar otro pago.",
      "Práctica segura: en pruebas, dos operadores solicitan el pago del mismo pedido; sólo debe existir una intención inicial. Simulá respuesta perdida y recuperá el registro con Actualizar pedidos. No practiques con pedidos reales.",
      "Para soporte, comunicá la referencia del pedido o solicitud, acción y hora al responsable autorizado. No adjuntes tokens ni datos de clientes. Cobro, conciliación y entrega mantienen sus gates separados."
    ],
    "inlineLead": true
  },
  "quote-create-view": {
    "id": "quote-create-view",
    "version": "1.0.0",
    "title": "Emitir una cotización",
    "summary": "Ayuda para emitir una cotización",
    "paragraphs": [
      "El servidor determina precio y moneda. La referencia conserva sólo un identificador de operación. Una respuesta perdida no permite emitir otra: consultá el resultado. Si no aparece, conservá la referencia y pedí revisión al responsable autorizado. Practicá emisión, respuesta perdida y consulta con datos sintéticos; no borres la referencia para repetir una cotización."
    ],
    "inlineLead": false
  },
  "resource-create-view": {
    "id": "resource-create-view",
    "version": "1.0.0",
    "title": "Registrar un recurso",
    "summary": "Ayuda para registrar un recurso",
    "paragraphs": [
      "Conservamos sólo la referencia de operación. Una respuesta perdida no demuestra rechazo: consultá antes de volver a actuar. La consulta identifica el registro por la referencia de esta operación, sin buscar por nombre. Preparar otro recurso no reenvía ni cancela el anterior y sólo se habilita tras recuperarlo. Ante referencia ilegible o consulta fallida, pedí revisión autorizada sin borrar datos. Practicá con datos sintéticos el envío, pérdida de respuesta, consulta y preparación explícita de otro recurso."
    ],
    "inlineLead": false
  },
  "return-operations-view": {
    "id": "return-operations-view",
    "version": "1.0.0",
    "title": "Recibir y decidir devoluciones",
    "summary": "Ayuda para recibir y decidir devoluciones",
    "paragraphs": [
      "Recepción, decisión y solicitudes conservan evidencia inmutable. La referencia del navegador guarda sólo operación, IDs y huella; no guarda serie, notas ni tokens. Una respuesta perdida puede ocultar un cambio correcto. Consultá el caso exacto aunque la lista falle o no lo muestre. Si la evidencia no coincide, pedí revisión. Continuar operando actualiza el caso consultado y habilita una acción nueva; no reenvía la anterior. Una devolución con reembolso crea cuatro solicitudes; un cambio crea tres y no solicita emisión fiscal. Practicá con datos sintéticos antes de operar."
    ],
    "inlineLead": false
  },
  "slot-create-view": {
    "id": "slot-create-view",
    "version": "1.0.0",
    "title": "Registrar un turno",
    "summary": "Ayuda para registrar un turno",
    "paragraphs": [
      "Conservamos sólo la referencia de operación. Una respuesta perdida no demuestra rechazo: consultá antes de volver a actuar. La consulta muestra el turno por su referencia, aunque esté cerrado, completo o fuera del listado público. Publicar no confirma una reserva; siguen vigentes la jornada y las restricciones de agenda. Preparar otro turno no reenvía ni cancela el anterior y sólo se habilita tras recuperarlo. Ante referencia ilegible o consulta fallida, pedí revisión autorizada sin borrar datos. Practicá con datos sintéticos el envío, pérdida de respuesta, consulta y preparación explícita de otro turno."
    ],
    "inlineLead": false
  }
} as const;

// Same-release operational guidance. AUTHORED integration text; no corporate attribution.
export const RELEASE_GUIDES = {
  "supply-role-view": {
    "id": "supply-role-view",
    "version": "1.0.0",
    "title": "Guía de suministro",
    "summary": "Guía de suministro · versión 1.0.0",
    "paragraphs": [
      "Compras registra proveedor, destino, fábrica, importe y cantidades. La fábrica confirma la orden, registra cada serie y pide revisión de calidad. Otra persona autorizada decide su liberación o rechazo con evidencia.",
      "El despacho identifica un envío y sus series. Recepción registra sólo lo recibido; queda en cuarentena hasta la decisión de otra persona autorizada. Las cantidades y versiones se comprueban en cada operación.",
      "Si se pierde una respuesta, consultá el resultado pendiente. Esa consulta no vuelve a enviar. No borres referencias para forzar un reintento. En capacitación, practicá con datos sintéticos el recorrido y la recuperación."
    ],
    "inlineLead": false
  },
  "warranty-role-view": {
    "id": "warranty-role-view",
    "version": "1.0.0",
    "title": "Guía de garantía",
    "summary": "Guía de garantía · versión 1.0.0",
    "paragraphs": [
      "Leé los términos vinculados a la cotización antes de acusar su recepción. La garantía activa conserva esos términos vendidos. Un reclamo necesita entrega y atención registradas.",
      "Diagnóstico, plan y aprobación preceden al trabajo. La persona revisora es diferente de quien propone; calidad es diferente de quien ejecuta. El cliente acepta y la fábrica concilia el servicio. La conciliación es un acuse interno, no un pago.",
      "Si se pierde una respuesta, consultá el resultado pendiente: no vuelve a enviar. Conservá su referencia y pedí revisión autorizada si sigue incierto. Practicá con datos sintéticos en capacitación."
    ],
    "inlineLead": false
  },
  "network-role-view": {
    "id": "network-role-view",
    "version": "1.0.0",
    "title": "Guía de red y acuerdos",
    "summary": "Guía de red y acuerdos · versión 1.0.0",
    "paragraphs": [
      "Registrá la organización bajo un padre autorizado y usá la referencia devuelta para consultar su estado. Crear una organización no concede accesos a personas. Los nuevos accesos se asignan mediante el proveedor de identidad autorizado.",
      "La estructura activa permite preparar acuerdos y sucursales. No equivale a autorización de despliegue o producción. Territorio, versión de términos y fechas deben proceder del acuerdo revisado; no se generan condiciones comerciales aquí.",
      "Antes de operar, completá la capacitación de tu rol. La evaluación no concede permisos. Si una respuesta se pierde, recuperá su resultado y luego consultá el estado actual. No repitas una escritura incierta.",
      "Para cerrar una organización, revisá primero sus sucursales y acuerdos activos. No se eliminan los registros ni se revocan sesiones desde este formulario."
    ],
    "inlineLead": false
  },
  "help-cms-view": {
    "id": "help-cms-view",
    "version": "1.0.0",
    "title": "Guía de contenidos",
    "summary": "Guía de contenidos · versión 1.0.0",
    "paragraphs": [
      "El contenido pertenece a una organización y un idioma. Guardá un borrador, revisalo y publicalo con el permiso correspondiente. Archivar retira el artículo de la consulta de lectores y conserva su historial.",
      "Las versiones publicadas no se editan. Para reemplazarlas, prepará un artículo nuevo y archivá el anterior después de revisar la sustitución. El texto se muestra literalmente; no ejecuta HTML.",
      "Publicar aquí no modifica cursos aprobados ni concede permisos. La incorporación a capacitación requiere revisión y una nueva versión del curso."
    ],
    "inlineLead": false
  },
  "catalog-role-view": {
    "id": "catalog-role-view",
    "version": "1.0.0",
    "title": "Guía de catálogo y publicación",
    "summary": "Guía de catálogo · versión 1.0.0",
    "inlineLead": false,
    "paragraphs": [
      "Creá modelo, variantes y precios con importes y vigencias explícitos. Conservá las referencias devueltas, incorporá la imagen PNG y prepará una versión del catálogo.",
      "La versión requiere las revisiones legal, técnica, de imagen y de publicación. Cada decisión refiere al contenido exacto y su evidencia. Publicar comprueba la generación vigente; una revisión desactualizada no autoriza contenido nuevo.",
      "Consultar un resultado pendiente recupera la operación sin repetirla. Antes de publicar otra versión o volver a una anterior, consultá la publicación vigente. La publicación del catálogo no confirma stock, cobro ni habilitación comercial. Practicá únicamente con datos sintéticos."
    ]
  },
  "training-role-view": {
    "id": "training-role-view",
    "version": "1.0.0",
    "title": "Guía de capacitación y evaluación",
    "summary": "Guía de capacitación · versión 1.0.0",
    "inlineLead": false,
    "paragraphs": [
      "Elegí la guía del rol e iniciá una práctica. Registrá la lectura de cada lección y presentá respuestas para revisión humana. Elegir una guía no cambia los permisos de tu sesión.",
      "Una persona autorizada diferente revisa las respuestas contra la versión presentada. Una evaluación favorable conserva evidencia, pero no asigna accesos ni reemplaza la aceptación del responsable de la operación.",
      "Ante una respuesta incierta, consultá el intento guardado. Un curso nuevo no cambia la evidencia de intentos anteriores; prepará otra práctica para la versión activa. Los artículos privados del CMS requieren incorporación explícita a un curso revisado."
    ]
  }
} as const;

export const ALL_GUIDES = [QUOTE_GUIDE, NOTIFICATION_GUIDE, ...Object.values(OPERATIONAL_GUIDES), ...Object.values(RELEASE_GUIDES)] as const;
export type GuideID = typeof ALL_GUIDES[number]["id"];
