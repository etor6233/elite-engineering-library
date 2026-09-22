# Mapa de la biblioteca para personas

Alcance planificado. Síntesis para leer en minutos. No certifica que cada recuadro ya esté implementado ni aceptado en producción.

Una persona lee estos tres gráficos. Un agente no los usa como corpus: busca una ruta en `LIBRARY_SEARCH_INDEX.md`, abre ese archivo y una sola fila de estado si hace falta.

```bash
rg -n -F "término" markdown_system/LIBRARY_SEARCH_INDEX.md
```

El índice nombra cada manual, pack y evidencia. Nada queda fuera de la búsqueda. El detalle sigue en su archivo.

## 01. De la biblioteca al sistema

El negocio, un proyecto nuevo o uno existente, y la biblioteca llegan al mismo agente. El agente elige superficies, arma un plan, construye por recorridos y solo sigue si las pruebas pasan.

```mermaid
flowchart TD
  negocio["Vos definís el negocio<br/>Productos o servicios, reglas, presupuesto<br/>y responsabilidades de la marca y franquicias"]
  nuevo["Proyecto nuevo"]
  existente["Proyecto existente<br/>Conservar operación y datos"]
  biblioteca["Biblioteca<br/>Manuales y arquitectura<br/>Fuentes oficiales y licencias<br/>Código reutilizable y adaptadores<br/>Contratos, verificadores y procedimientos"]
  agente["Agente trabajando en el proyecto"]
  superficies["Evaluar las 48 superficies del sistema<br/>Necesarias, opcionales, no aplicables o bloqueadas"]
  plan["Plan concreto del proyecto<br/>Recorridos completos, módulos compatibles,<br/>responsables y criterios de aceptación"]
  construir["Construir o integrar por recorridos<br/>Frontend, backend, datos e integraciones<br/>Ayuda, capacitación y operación"]
  verificar["Verificar el conjunto<br/>Funcionamiento, permisos, seguridad,<br/>rendimiento, errores y recuperación"]
  entorno["Configurar entorno y accesos protegidos<br/>Verificar proveedores y despliegue real"]
  aceptado["Sistema operativo aceptado<br/>con evidencia y recuperación probada"]
  mantener["Mantener y ampliar<br/>con la misma biblioteca y metodología"]

  negocio --> agente
  nuevo --> agente
  existente --> agente
  biblioteca --> agente
  agente --> superficies --> plan --> construir --> verificar
  construir -->|"Falla"| construir
  verificar -->|"Supera las pruebas"| entorno --> aceptado --> mantener
  mantener --> plan
```

## 02. Operación de la red de franquicias

La marca gobierna. Cada franquicia opera con su personal, stock y resultados. Un backend común identifica organización, sucursal y usuario, aplica reglas y registra cada operación. El cliente ve la marca. La franquicia reserva, cobra, entrega y atiende la postventa.

```mermaid
flowchart TD
  marca["Marca / central<br/>Gobierno de la red, compras y políticas<br/>Catálogo, precios, territorios y auditorías"]
  alta["Incorporar franquicia<br/>Acuerdo, territorio, responsables<br/>Permisos, configuración, capacitación<br/>Aceptación y apertura"]
  locales["Franquicias A, B, C<br/>Cada una con su personal, stock,<br/>operaciones y resultados identificados"]
  panel["Panel operativo por rol<br/>Central, franquicias, empleados,<br/>proveedores y logística"]
  backend["Backend común con permisos<br/>Identifica organización, sucursal y usuario<br/>Aplica reglas y registra cada operación"]
  ayuda["Ayuda, capacitación y asistente<br/>Contexto del usuario y permisos<br/>Escalamiento cuando corresponde"]
  control["Control de la red<br/>Auditoría, ventas, costos y márgenes<br/>Stock, campañas y calidad de atención<br/>Regalías y liquidaciones según acuerdos"]
  compras["Compras y abastecimiento<br/>Proveedores, fábricas, órdenes,<br/>producción y transporte"]
  docs["Documentos y costos<br/>Proformas, invoices y packing lists<br/>Extracción validada y revisión cuando corresponda"]
  clientes["Clientes y visitantes"]
  captacion["Captación<br/>SEO, publicidad, redes y mensajes<br/>Origen, cuenta y gasto identificados"]
  stock["Recepción y stock real<br/>Cantidades, calidad, depósitos y ubicaciones<br/>Lotes, series o VIN según el negocio<br/>Transferencias y trazabilidad"]
  publico["Frontend público de la marca<br/>Catálogo y páginas locales<br/>Disponibilidad por sucursal"]
  portal["Portal del cliente<br/>Cotizaciones, pedidos, pagos,<br/>entregas y atención"]
  catalogo["Catálogo y precios<br/>Productos, variantes, contenido,<br/>vigencias y publicación"]
  oportunidades["Clientes y oportunidades<br/>Mensajes, consentimiento, asignación,<br/>seguimiento y cotizaciones"]
  reserva["Reserva y pedido<br/>Franquicia responsable y vendedor identificados"]
  cobro["Cobro y conciliación<br/>Cuenta receptora definida<br/>Estados, devoluciones y duplicados controlados"]
  entrega["Preparación y entrega<br/>Asignación de stock, transporte,<br/>seguimiento y comprobante"]
  postventa["Postventa<br/>Reclamos, devoluciones, garantías,<br/>turnos, repuestos y servicio"]

  marca --> alta --> locales --> panel --> backend
  marca --> backend
  locales --> backend
  panel --> ayuda
  panel --> control
  backend --> compras --> docs --> stock
  clientes --> captacion --> publico
  stock -->|"Disponibilidad publicable"| publico
  publico --> portal
  backend --> catalogo --> reserva
  backend --> oportunidades --> reserva
  stock -->|"Disponibilidad y reserva"| reserva
  reserva --> cobro --> entrega --> postventa
  entrega -->|"Movimientos autorizados"| backend
  postventa --> backend
  ayuda --> backend
  control --> backend
  portal --> backend
```

## 03. Ingeniería, protección y mantenimiento

Todo el sistema anterior se sostiene con experiencia, seguridad, datos, integraciones, infraestructura, pruebas, extensiones y operación. Un incidente no se parchea a ciegas: se diagnostica, se prueba según el riesgo y se despliega con vuelta atrás.

```mermaid
flowchart TD
  todo["Todo el sistema anterior"]
  ux["Experiencia y frontend<br/>Web pública y portales por rol<br/>SEO, accesibilidad, claridad y velocidad<br/>Ayuda y capacitación de la misma versión"]
  sec["Seguridad y privacidad<br/>Identidad, MFA y permisos por recurso<br/>Aislamiento entre organizaciones<br/>Secretos, cifrado, auditoría y retención"]
  datos["Datos y consistencia<br/>Base transaccional y migraciones<br/>Archivos y documentos protegidos<br/>Búsqueda y caché cuando hagan falta"]
  integ["Integraciones y automatización<br/>APIs y contratos versionados<br/>Eventos, trabajos durables y reintentos<br/>Deduplicación, conciliación y notificaciones"]
  infra["Infraestructura y rendimiento<br/>Dominios, DNS, HTTPS y protección web<br/>Cloud y entornos reproducibles<br/>Capacidad, latencia y límites de gasto"]
  pruebas["Construcción y pruebas<br/>Versiones y dependencias fijadas<br/>Licencias, procedencia y artefactos firmados<br/>Pruebas funcionales, E2E, carga y seguridad"]
  ext["Extensiones según necesidad<br/>Móvil, escritorio, offline y periféricos<br/>IoT, dispositivos y telemetría<br/>Analítica, IA, búsqueda documental y GPU"]
  obs["Observabilidad y operación<br/>Métricas, trazas, registros y alertas<br/>Objetivos de servicio y costos"]
  diag["Ante un incidente, actualización o necesidad nueva<br/>El agente diagnostica y prepara el cambio<br/>Consulta evidencia y límites de autorización"]
  riesgo["Pruebas y revisión según riesgo<br/>Actualizar controles, ayuda y capacitación"]
  deploy["Despliegue controlado"]
  continuidad["Continuidad del negocio<br/>Backups y restauración comprobada<br/>Rollback y recuperación ante fallos"]

  todo --> ux
  todo --> sec
  todo --> datos
  todo --> integ
  todo --> infra
  todo --> pruebas
  todo --> ext
  todo --> obs
  obs --> diag --> riesgo --> deploy
  deploy -->|"Si falla"| continuidad
  deploy -->|"Si sale bien"| continuidad
```

El proyecto selecciona las capacidades necesarias. Las reglas específicas y los accesos se verifican en su entorno real.
