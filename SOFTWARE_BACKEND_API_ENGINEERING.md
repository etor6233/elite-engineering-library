# Software, Backend & API Engineering — manual operativo

> **Estado:** núcleo general v1 auditado el 2026-08-21; las 19 lecturas de MIT 6.102, la normativa API/HTTP seleccionada, la evidencia de producción y los cruces con sistemas, redes, datos, GPU, seguridad/SRE/cloud, microestructura y frontend están auditados; las extensiones dependientes de proyecto o manuales futuros se declaran al final.
> **Autoridad académica:** MIT 6.102 *Software Construction*.
> **Autoridad normativa:** RFC 9110–9114, RFC 9457, RFC 8446, RFC 9700, RFC 8725, OpenAPI 3.2.0, JSON Schema 2020-12, gRPC y Protocol Buffers.
> **Autoridad operativa:** *Software Engineering at Google*, Google API Improvement Proposals, Amazon Builders’ Library y OWASP API Security Top 10 2023.
> **Propósito:** convertir requisitos de negocio en software y APIs correctos, comprensibles, evolutivos, seguros, observables y operables, sin confundir un framework con la ingeniería subyacente.

## Cómo usar este manual

Para cada cambio, recorrer esta secuencia:

```text
problema y usuario
→ invariantes y riesgos
→ contrato observable
→ ejemplos y casos límite
→ tests que expresan el contrato
→ diseño interno y límites de módulos
→ implementación mínima
→ revisión humana y automática
→ pruebas de integración/fallo/carga
→ despliegue gradual y observación
→ aprendizaje y actualización del contrato
```

No comenzar por elegir framework, base de datos o patrón. Esas son decisiones subordinadas al contrato, la carga, el riesgo, el equipo y la vida esperada del sistema.

## Procedencia

- `[SPEC]`: comportamiento normativo de una especificación o documentación oficial versionada.
- `[ACADEMIC]`: síntesis fiel de material académico identificado; no es transcripción literal.
- `[CODE]`: comportamiento comprobado en implementación o repositorio oficial.
- `[PROD]`: práctica publicada por operadores o expertos con experiencia de producción.
- `[MEASURED]`: evidencia reproducible en el entorno declarado.
- `[IMPL]`: traducción propia a arquitectura, código, test, gate o runbook.
- `[OPEN]`: cuestión aún no auditada o dependiente del proyecto; no tratar como regla cerrada.

## 1. Contrato de ingeniería de software

### 1.1 Las tres propiedades que gobiernan el diseño

**[ACADEMIC]** MIT 6.102 formula tres objetivos que sobreviven al lenguaje y al framework:

| Propiedad | Pregunta de revisión |
|---|---|
| segura frente a bugs | ¿es correcta hoy y defensiva ante entradas, estados y cambios futuros? |
| fácil de entender | ¿comunica intención, supuestos y límites a otra persona? |
| preparada para cambiar | ¿una modificación razonable queda localizada o exige reescribir? |

Rendimiento, seguridad y usabilidad también son requisitos. Cuando compiten, el trade-off debe quedar explícito y medido; nunca se degrada corrección silenciosamente para obtener velocidad aparente.

**[PROD]** *Software Engineering at Google* distingue programación de ingeniería por tiempo y escala: las prácticas deben responder a cuánto vivirá el software, cuántas personas lo mantendrán y cuántos consumidores dependerán de él.

### 1.2 Definición de terminado

**[IMPL]** Una unidad de backend no está terminada porque “responde 200”. Debe demostrar:

- contrato documentado y validable;
- invariantes de dominio preservados;
- autenticación y autorización probadas;
- semántica de error y retry definida;
- límites de tiempo, tamaño, concurrencia y recursos;
- observabilidad sin secretos ni datos innecesarios;
- compatibilidad o plan de migración;
- tests del camino feliz, bordes y fallos;
- estrategia de despliegue, rollback y recuperación;
- ownership y runbook.

### 1.3 Plantilla de contexto antes de diseñar

```text
usuario/actor:
resultado de negocio:
recursos y operaciones:
invariantes:
datos sensibles:
volumen, concurrencia y distribución de tamaños:
SLO de latencia/disponibilidad/corrección:
consistencia y durabilidad requeridas:
fallos tolerables y no tolerables:
consumidores y ritmo de actualización:
vida esperada y ownership:
restricciones regulatorias/legales:
```

## 2. Especificaciones que permiten razonar

### 2.1 Qué debe decir una especificación

**[ACADEMIC]** Una especificación es la frontera entre cliente e implementador. Debe ser declarativa: describir qué se garantiza, no imponer accidentalmente cómo se logra.

```text
operación:
precondiciones:
resultado/postcondiciones:
efectos observables:
errores distinguidos:
ordenamiento y consistencia:
límites:
idempotencia:
compatibilidad:
```

- Una precondición más débil acepta más entradas y exige más al implementador.
- Una postcondición más fuerte promete más resultados y exige más al implementador.
- Una especificación puede quedar deliberadamente subdeterminada para conservar libertad de implementación.
- Todo comportamiento observable no documentado puede convertirse en dependencia de alguien; minimizar observables accidentales.

### 2.2 Ejemplos, particiones y límites

**[ACADEMIC]** Diseñar tests desde la especificación antes de implementar. Particionar el espacio de entrada por comportamientos distintos y elegir bordes: vacío/no vacío, mínimo/máximo, válido/inválido, presente/ausente, autorizado/no autorizado, nuevo/duplicado, antes/en/durante/después de expiración.

**[IMPL]** Para una operación de creación, como mínimo:

| Partición | Caso |
|---|---|
| identidad | ID generado, ID solicitado válido, colisión |
| datos | mínimos válidos, completos, campo desconocido, tipo/formato inválido |
| autorización | actor permitido, otro tenant, credencial ausente/expirada |
| repetición | primera solicitud, mismo token+mismo payload, mismo token+otro payload |
| concurrencia | dos creaciones lógicamente iguales, cancelación, deadline |
| dependencia | éxito, timeout, respuesta inválida, rechazo, indisponibilidad |

### 2.3 Tipos abstractos e invariantes de representación

**[ACADEMIC]** Un ADT se define por sus operaciones y comportamiento, no por su representación. La función de abstracción conecta representación con valor abstracto; el invariante de representación define qué estados internos son legales.

**[IMPL]** Toda entidad con reglas no triviales debe concentrarlas:

```text
Money:
  invariante: currency válida; minor_units entero dentro de rango
Order:
  invariante: quantity > 0; price cumple tick; transición de estado legal
UserEmail:
  invariante: forma normalizada definida; longitud y dominio acotados
```

No repartir el mismo invariante entre controlador, ORM, job y frontend. Validar en el límite y preservar en el núcleo de dominio; la base de datos añade una última línea de defensa mediante constraints.

**[ACADEMIC]** Clasificar las operaciones de un ADT como creators, producers, observers y mutators. Cada creator/producer establece el invariante; cada observer/mutator debe preservarlo. Cuando sea viable, una función `checkRep()` o assertions equivalentes convierten el invariante documentado en un control ejecutable en construcción, mutación y tests. Desactivarlo en producción sólo después de evaluar costo y pérdida de diagnóstico.

Receta auditada de MIT:

```text
elegir abstracción y mutabilidad
→ escribir operaciones y especificaciones
→ elegir representación
→ escribir AF + RI + argumento contra rep exposure
→ generar tests desde operaciones/particiones
→ implementar primero de forma simple
→ comprobar RI
→ integrar temprano
→ optimizar/rediseñar con evidencia
```

### 2.4 Mutabilidad, aliasing y ownership

- Preferir valores inmutables para mensajes, configuración y value objects.
- No devolver referencias mutables a representación interna.
- Definir quién crea, modifica, comparte y destruye cada objeto o recurso.
- Evitar estado global mutable; hace implícitas dependencias, orden y sincronización.
- Si existe caché o memoización, declarar si forma parte de la semántica o sólo de la implementación.

### 2.5 Igualdad, identidad y hashing

**[ACADEMIC]** Igualdad debe ser una relación de equivalencia: reflexiva, simétrica y transitiva. Distinguir:

- identidad/referencia: es la misma instancia;
- igualdad observacional: hoy expone el mismo valor abstracto;
- igualdad de comportamiento: seguirá siendo indistinguible bajo operaciones futuras.

En valores inmutables, igualdad observacional y de comportamiento coinciden. En objetos mutables pueden divergir. Una clave de mapa/set debe conservar igualdad y hash mientras esté indexada; no usar campos mutables como base sin una política que lo haga seguro. Canonicalización e interning cambian representación, no la definición del valor abstracto.

### 2.6 Tipos, estados imposibles y exhaustividad

- Usar enums/unions selladas para conjuntos finitos; evitar strings mágicos.
- Hacer exhaustivo el manejo de variantes para que una nueva variante produzca error estático o test fallido.
- Distinguir `null`, ausencia y variante vacía; preferir tipos que expresen cada caso.
- No esperar que el tipo estático valide valores externos en runtime: parsear y validar primero.
- Elegir enteros, decimales y timestamps con rango/precisión explícitos; los números JSON/JavaScript no representan exactamente todos los enteros de 64 bits.

## 3. Arquitectura de un servicio

### 3.1 Flujo de una solicitud

```text
transporte
→ terminación TLS / proxy
→ parsing y límites
→ autenticación
→ autorización
→ validación semántica
→ caso de uso / dominio
→ transacción y dependencias
→ serialización
→ métricas, tracing y auditoría
```

El orden exacto depende del riesgo, pero nunca ejecutar trabajo costoso o producir efectos antes de autenticar, autorizar y acotar la entrada cuando sea posible.

### 3.2 Límites de módulo

**[IMPL]** Separación mínima:

| Capa | Responsabilidad | No debe decidir |
|---|---|---|
| transporte | protocolo, parsing, status, headers | reglas de negocio |
| aplicación | orquestar caso de uso, transacción, permisos | detalles del protocolo |
| dominio | invariantes y transiciones | HTTP, SQL, proveedor cloud |
| adaptadores | persistencia, colas, servicios externos | política central de negocio |

Las dependencias apuntan hacia contratos estables. Un repositorio o gateway no debe ocultar semántica crítica: aislamiento de transacción, consistencia, reintentos y fallos siguen siendo explícitos.

### 3.3 Monolito modular, servicios y funciones

**[IMPL]** Comenzar con el límite de despliegue más simple que satisfaga ownership, escalado, aislamiento y cumplimiento. Separar servicios sólo cuando exista una frontera justificable:

- autonomía de equipo y ciclo de entrega;
- escalado o recursos radicalmente distintos;
- aislamiento de fallos o seguridad;
- reglas de datos y consistencia claramente separables.

La distribución añade latencia, fallos parciales, version skew, observabilidad, seguridad y operación. Un llamado local convertido en RPC ya no conserva las mismas garantías.

## 4. Semántica HTTP

### 4.1 Modelo correcto

**[SPEC]** RFC 9110 define HTTP como una familia de protocolos stateless de solicitud/respuesta con interfaz uniforme sobre recursos y representaciones. Las semánticas son comunes a HTTP/1.1, HTTP/2 y HTTP/3; el transporte no redefine el significado de métodos y status.

Separar siempre:

- recurso: concepto identificado;
- representación: estado transferido en un formato;
- método: intención de la solicitud;
- status: resultado genérico para componentes HTTP;
- contenido: información específica de aplicación;
- metadata: campos sobre representación, control, caché o conexión.

### 4.2 Métodos, seguridad e idempotencia

| Método | Intención usual | Safe | Idempotente |
|---|---|---:|---:|
| GET | obtener representación | sí | sí |
| HEAD | GET sin contenido de respuesta | sí | sí |
| OPTIONS | opciones de comunicación | sí | sí |
| PUT | crear/reemplazar estado en URI conocida | no | sí |
| DELETE | eliminar asociación/recurso | no | sí |
| POST | procesar según semántica del recurso | no | no por defecto |
| PATCH | aplicar cambios parciales | no | depende del formato/operación |

**[SPEC]** “Idempotente” significa que múltiples solicitudes idénticas tienen el mismo efecto intencionado que una; no exige respuestas byte a byte iguales ni ausencia de logging. “Safe” significa que el cliente no solicita un cambio de estado, aunque puedan existir efectos incidentales como métricas.

**[IMPL]** No usar `GET` para acciones. No asumir que `PUT` vuelve atómica una cadena de dependencias. Para hacer reintentable una operación naturalmente no idempotente, diseñar un contrato de deduplicación.

### 4.3 Selección de status

| Situación | Status orientativo |
|---|---:|
| lectura/actualización exitosa | 200 |
| creación con URI disponible | 201 + `Location` |
| aceptada para procesamiento diferido | 202 |
| éxito sin representación | 204 |
| representación no modificada por condición | 304 |
| sintaxis/tipo/formato inválido | 400 |
| credencial ausente o inválida | 401 + desafío aplicable |
| actor autenticado sin permiso | 403 |
| recurso no encontrado o deliberadamente oculto | 404 |
| método no permitido | 405 + `Allow` |
| conflicto con estado actual | 409 |
| precondición fallida | 412 |
| media type no soportado | 415 |
| semántica de contenido inválida | 422 |
| límite de tasa | 429; `Retry-After` si aplica |
| error no clasificado del servidor | 500 |
| dependencia/gateway inválido | 502 |
| servicio temporalmente no disponible | 503 |
| timeout de gateway | 504 |

El status describe la clase de resultado. El cuerpo añade un código estable de dominio; no devolver siempre 200 con un `success:false` privado.

### 4.4 Negociación y representación

- Validar `Content-Type` de la solicitud; negociar respuesta mediante `Accept` cuando se soporten formatos.
- Declarar encoding y formato temporal; preferir timestamps RFC 3339/ISO 8601 con offset o UTC y precisión definida.
- Distinguir campo ausente, `null`, vacío y valor por defecto.
- No representar enteros fuera del rango seguro del consumidor como números JSON sin contrato explícito.
- Poner límites a profundidad, cantidad de campos, arrays, strings y cuerpo total.

### 4.5 Caché y condiciones

**[SPEC]** RFC 9111 define almacenamiento, reutilización, freshness y validación. Diseñar explícitamente:

```text
cacheabilidad:
cache key y Vary:
freshness: Cache-Control / Expires
validator: ETag fuerte o débil / Last-Modified
revalidación: If-None-Match / If-Modified-Since
invalidación:
datos privados o compartidos:
```

Usar `no-store` cuando no deba almacenarse; `no-cache` exige validación antes de reutilizar y no significa “no guardar”. No cachear respuestas personalizadas en cachés compartidas sin claves y directivas correctas.

**[IMPL]** Para control optimista de escritura, exponer versión o `ETag` y exigir `If-Match`; una precondición fallida evita lost updates, pero no sustituye una transacción adecuada.

### 4.6 HTTP/1.1, HTTP/2 y HTTP/3

**[SPEC]** RFC 9112 especifica mensajería HTTP/1.1; RFC 9113, HTTP/2 multiplexado; RFC 9114, HTTP/3 sobre QUIC. Elegir y medir según infraestructura:

- multiplexar no elimina límites de aplicación, CPU, pools o dependencias;
- HTTP/2 evita head-of-line entre streams a nivel HTTP, pero TCP aún puede bloquear por pérdida;
- HTTP/3 cambia el transporte, handshake y operación de red; no “arregla” una API lenta;
- proxies, balanceadores y conversiones de versión forman parte del camino real.

### 4.7 Matices de métodos que evitan contratos falsos

**[SPEC]** Pasada completa de métodos RFC 9110 y PATCH RFC 5789:

- Un body en GET/HEAD no tiene semántica general, puede ser rechazado y presenta riesgo de request smuggling entre intermediarios; no diseñar una API que dependa de él.
- PUT pide reemplazar el estado de la URI conocida; si el servidor elige la URI, normalmente corresponde POST. Una respuesta PUT no es cacheable e invalida representaciones almacenadas de la URI objetivo.
- DELETE solicita quitar la asociación entre URI y funcionalidad; no garantiza borrado físico, inmediato o irreversible de todo dato relacionado. El contrato de retención debe decirlo.
- PATCH aplica un documento de cambios identificado por media type; no es safe ni idempotente por defecto. Debe aplicarse atómicamente a los recursos directamente afectados o no aplicar nada.
- Publicar formatos PATCH soportados mediante `Accept-Patch`; usar strong `ETag` + `If-Match` cuando el patch depende de un base-point.
- Una respuesta POST/PATCH sólo puede reutilizarse posteriormente para GET/HEAD bajo condiciones explícitas de freshness y `Content-Location`; no asumirlo como comportamiento común.
- Una respuesta exitosa sólo demuestra que la intención se logró al procesarla; una lectura posterior puede observar cambios concurrentes.

### 4.8 Validators, condiciones y rangos

- `If-Match` usa comparación fuerte y es apropiado para impedir lost updates.
- `If-None-Match` usa comparación débil y sirve para validación de caché o creación condicional (`*`).
- Evaluar precondiciones en la precedencia definida por RFC 9110; no inventar orden cuando llegan varias.
- Un `ETag` identifica una representación según el origin; no asumir que es hash, versión de fila o firma criptográfica.
- `Range`/206 representa partes de una representación, no paginación de una colección.
- Para byte-range no satisfacible, 416 puede incluir `Content-Range: bytes */longitud`; servidores pueden ignorar Range y responder 200 completo.
- Rechazar rangos numerosos/solapados que impliquen abuso de recursos.

## 5. Diseño de recursos y operaciones

### 5.1 Recursos primero

**[PROD]** Google AIP-121/122 organiza APIs alrededor de recursos como sustantivos con nombres canónicos y operaciones estándar. Es una convención útil, no una ley universal de HTTP.

```text
publishers/{publisher}/books/{book}
```

- Identificadores estables y opacos; no filtrar secretos o PII en URI.
- Jerarquía sólo cuando expresa ownership o ámbito real.
- Referenciar otro recurso por identidad, no duplicar su representación sin política de consistencia.
- Usar operaciones estándar para CRUD; acción custom sólo si no encaja limpiamente.

### 5.2 Colecciones

Toda lista no trivial define:

- orden estable y desempate;
- paginación y tamaño máximo;
- semántica del page token y expiración;
- filtros y campos permitidos;
- snapshot/consistencia mientras la colección cambia;
- autorización por elemento;
- límite de costo.

**[IMPL]** Preferir cursor/token opaco a offset en conjuntos grandes o mutables. No codificar estado sensible sin autenticidad/confidencialidad; validar que el token pertenece al filtro, actor y versión correctos.

### 5.3 Actualización parcial

- Distinguir reemplazo completo de patch parcial.
- Definir si `null` limpia, asigna nulo o se ignora.
- Permitir lista explícita de campos actualizables; rechazar campos output-only.
- Con Protobuf, usar field mask cuando el ecosistema lo soporte.
- Aplicar autorización por propiedad, no sólo por endpoint.

### 5.4 Operaciones largas

Para trabajo que supera el presupuesto síncrono:

```text
POST → 202 + operation resource
GET operation → estado/progreso/resultado/error
cancel operation → mejor esfuerzo o garantía declarada
```

Definir retención, idempotencia, cancelación, expiración, ownership, visibilidad y cuándo el efecto se vuelve irreversible.

## 6. Contrato ejecutable

### 6.1 OpenAPI

**[SPEC]** OpenAPI 3.2.0 es una descripción independiente del lenguaje para que humanos y herramientas comprendan una API HTTP sin inspeccionar código o tráfico.

El documento debe contener, según corresponda:

- operations y `operationId` estable;
- parámetros, request bodies y media types;
- respuestas exitosas y cada clase de error;
- schemas y ejemplos;
- security schemes y requirements;
- callbacks/webhooks, links o streaming si existen;
- deprecaciones y versión.

**[IMPL]** Validar el documento con tooling, pero recordar que el schema de OpenAPI no detecta toda violación de la especificación. Hacer lint semántico propio: nombres, paginación, errores, auth, límites e idempotencia.

Detalles de interoperabilidad y seguridad de OpenAPI 3.2:

- Los Schema Objects forman un superset del dialecto JSON Schema 2020-12; fijar `jsonSchemaDialect`/`$schema` cuando la portabilidad importe.
- En `security`, varios Security Requirement Objects son alternativas **OR**; varios schemes dentro del mismo objeto son requisitos **AND**; `{}` habilita acceso anónimo. Testear que la descripción coincide con enforcement real.
- `$ref` y recursos externos pueden apuntar a hosts no confiables; resolver offline o con allowlist, límites y caché controlada.
- Tooling debe detectar ciclos de referencias para evitar agotamiento de recursos.
- Sanitizar Markdown/HTML antes de renderizar documentación generada.
- Diferenciar annotations (`default`, `examples`, `readOnly`, `writeOnly`) de validación efectiva.

### 6.2 JSON Schema

**[SPEC]** La versión publicada vigente es 2020-12, separada en Core y Validation. Declarar dialecto con `$schema`; no asumir que dos validadores implementan los mismos drafts o formatos.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["symbol", "quantity"],
  "properties": {
    "symbol": {"type": "string", "minLength": 1, "maxLength": 32},
    "quantity": {"type": "integer", "minimum": 1}
  },
  "additionalProperties": false
}
```

Schema valida forma; no reemplaza autorización, invariantes entre campos, consulta de estado ni reglas de negocio.

En 2020-12, `format` es annotation por defecto; sólo es assertion completa cuando se declara/soporta el vocabulario correspondiente. `default` tampoco inserta valores durante validación. Por lo tanto, probar configuración y versión exacta del validator, no inferir enforcement del documento.

### 6.3 gRPC y Protocol Buffers

**[SPEC]** gRPC define RPC unary y streaming con status, metadata, deadlines, cancelación, health checking y flow control. Protocol Buffers define el contrato serializado.

Reglas de evolución de Protobuf:

- nunca reutilizar un tag; reservar números y nombres eliminados;
- casi nunca cambiar el tipo de un campo;
- no cambiar defaults de forma incompatible;
- clientes y servidores no se actualizan al mismo tiempo;
- no depender de serialización byte-estable entre builds;
- añadir campos de modo que lectores viejos puedan ignorarlos;
- evaluar compatibilidad binaria y JSON por separado.
- no añadir campos `required`; representar requisito como contrato validado en aplicación;
- incluir valor `UNSPECIFIED = 0` en enums para evolución y presencia;
- no usar boolean si el dominio probablemente tendrá más estados;
- separar mensajes RPC de mensajes persistidos para no congelar ambos ciclos de vida juntos.

**[SPEC]** gRPC puede hacer transparent retries aun sin policy si sabe que la lógica de aplicación no procesó la llamada. Con policy explícita, el RPC se considera committed al recibir response headers; después no lo reintenta automáticamente. Revisar retry, hedging y service config por método y lenguaje, y observar intentos además de calls lógicos.

### 6.4 Elegir REST/HTTP, gRPC o eventos

| Opción | Encaja cuando | Riesgo principal |
|---|---|---|
| HTTP + JSON | interoperabilidad web, consumo humano/tooling, caché | contratos laxos y payload/parse cost |
| gRPC + Protobuf | servicios tipados, streaming, multi-lenguaje controlado | version skew, proxies/debugging, browser edge |
| evento/mensaje | desacoplar tiempo y productores/consumidores | entrega, orden, duplicados, schema evolution |

No elegir por moda. Comparar consumidores, red, latencia, tamaño, evolución, tooling, observabilidad y operación. El manual de redes/distribuidos cerrará entrega, orden y backpressure de eventos.

## 7. Errores como parte de la API

### 7.1 Modelo estable

**[SPEC]** RFC 9457 define `application/problem+json` para detalles machine-readable y reemplaza RFC 7807.

```json
{
  "type": "https://api.example.com/problems/insufficient-balance",
  "title": "Insufficient balance",
  "status": 403,
  "detail": "The operation cannot be completed.",
  "instance": "/operations/01J...",
  "code": "INSUFFICIENT_BALANCE",
  "trace_id": "..."
}
```

- `type`/`code` es identidad estable para código cliente.
- `title` es breve; `detail` es específico de la ocurrencia.
- `instance` identifica la ocurrencia, no expone internals.
- extensiones nuevas deben ser ignorables por consumidores viejos.
- no incluir stack, SQL, filesystem, tokens, claves ni detalles explotables.

Para gRPC, usar status canónico y detalles tipados. **[PROD]** Google AIP-193 exige información machine-readable porque mensajes humanos cambian y no deben parsearse.

### 7.2 Taxonomía de acción

Cada error debe indicar internamente:

```text
culpa: cliente | servicio | dependencia | operación
permanencia: permanente | transitorio | desconocido
retry: nunca | inmediato | backoff | después de condición
efecto: ninguno | posible | confirmado
visibilidad: usuario | desarrollador | operador | seguridad
```

No mapear toda excepción a 500. No revelar si un recurso de otro tenant existe: autorización puede preceder a existencia.

## 8. Deadlines, cancelación, reintentos e idempotencia

### 8.1 Presupuesto temporal end-to-end

**[SPEC]** gRPC no establece deadline por defecto; un cliente podría esperar indefinidamente. El servidor debe detener trabajo derivado cuando vence o se cancela y propagar la cancelación aguas abajo.

```text
deadline_total
  = cola_cliente
  + red
  + cola_servidor
  + cómputo
  + dependencias
  + serialización
  + margen
```

- Elegir timeout desde SLO y distribución medida, no un número copiado.
- Distinguir connect, TLS, request/write, response/read y deadline total.
- Propagar el tiempo restante, no reiniciar el presupuesto en cada salto.
- Cancelar trabajo que ya no puede producir un resultado útil.
- Registrar timeout separado de error del servicio.

### 8.2 Cuándo reintentar

**[PROD]** AWS y gRPC advierten que un retry agrega trabajo precisamente cuando el sistema puede estar fallando. Aplicar sólo si:

1. el error es transitorio y clasificado;
2. la operación es idempotente o deduplicada;
3. queda presupuesto temporal;
4. existe límite de intentos/trabajo;
5. backoff y jitter reducen sincronización;
6. la capa elegida evita multiplicación en cada salto.

```text
delay_i = random(0, min(cap, base * 2^i))   # full jitter, ejemplo
```

**[IMPL]** Medir tasa de requests originales, intentos, éxitos recuperados, amplificación y carga descartada. Un circuito abierto, token bucket o retry budget puede frenar tormentas; debe tener tests de transición y recuperación.

No apilar políticas sin presupuesto global: SDK, proxy, service mesh y aplicación pueden reintentar la misma operación. Un solo owner debe decidir; los demás niveles exponen fallos o consumen una cuota coordinada. Hedging inicia intentos concurrentes y sólo es aceptable con idempotencia, cancelación del perdedor y capacidad medida.

### 8.3 Claves de idempotencia

**[PROD]** AWS describe idempotent APIs para conservar un resultado lógico aun con múltiples intentos.

Contrato mínimo:

```text
scope de la clave: tenant + operación
entropía y longitud:
fingerprint del payload semántico:
resultado guardado:
estado in-progress:
TTL >= ventana máxima de retry:
comportamiento misma clave/mismo payload:
comportamiento misma clave/distinto payload: conflicto
atomicidad entre deduplicación y efecto:
```

No implementar como “buscar y luego insertar” sin atomicidad. Considerar crash entre el efecto y el registro; la solución depende de transacción local, outbox/inbox, reserva previa o reconciliación.

## 9. Identidad, autenticación y autorización

### 9.1 Separar conceptos

| Concepto | Pregunta |
|---|---|
| autenticación | ¿qué identidad/principal presentó evidencia válida? |
| autorización | ¿puede ese principal ejecutar esta acción sobre este objeto y propiedades? |
| auditoría | ¿qué decisión y efecto deben quedar registrados? |

Nunca confiar autorización al frontend, al conocimiento de un ID o a que la ruta sea difícil de adivinar.

### 9.2 OAuth y tokens

OAuth 2.0 delega autorización; no es por sí solo un protocolo de autenticación. **[SPEC]** OpenID Connect Core 1.0 añade una capa de identidad sobre OAuth mediante ID Token y claims de autenticación. No usar un ID Token como access token ni aceptar un access token como prueba genérica de login.

**[SPEC]** RFC 9700 es el BCP vigente de seguridad OAuth 2.0 y depreca modos inseguros. Diseñar con Authorization Code + PKCE para clientes públicos; validar redirect URIs estrictamente; evitar tokens en URLs; rotar/proteger refresh tokens según el modelo; considerar tokens sender-constrained para reducir replay.

**[SPEC]** RFC 8725 exige seleccionar explícitamente algoritmos criptográficamente actuales para JWT y evitar confusión de algoritmos/validadores.

Validación mínima de un token:

- firma y algoritmo allowlist;
- issuer exacto;
- audience destinada a este servicio;
- expiración y not-before con skew acotado;
- tipo/uso del token;
- scopes/claims y tenant;
- key rotation y fallo seguro;
- revocación o vida corta según amenaza.

Para OIDC, además validar firma, `iss`, `aud`, `azp` cuando corresponda, tiempos y el `nonce` ligado a la sesión; consumir metadata/keys sólo desde issuer confiable y con política de rotación/fallo definida.

Un JWT firmado no cifra su contenido y no vuelve confiable cualquier claim; la política del recurso sigue en el servidor.

### 9.3 Autorización en cada objeto y propiedad

**[PROD]** OWASP API Security 2023 destaca fallos de autorización a nivel de objeto, propiedad y función.

**[IMPL]** Para cada operación:

```text
principal = authenticate(credential)
resource = load_in_authorized_scope(principal, canonical_name)
authorize(principal, action, resource)
validate_writable_properties(principal, patch)
apply_domain_transition(resource, patch)
```

Evitar `load(id)` seguido de un check olvidable. Encapsular el ámbito de tenant/owner en repositorio o policy, y testear IDs válidos pertenecientes a otro actor.

## 10. Validación y consumo hostil

### 10.1 Toda frontera es no confiable

Validar también respuestas de APIs externas, mensajes de cola, archivos, cache, datos históricos y configuración. **[PROD]** OWASP incluye el consumo inseguro de APIs: TLS o reputación del proveedor no garantiza payload correcto.

Controles:

- parsing estricto y fail-closed;
- allowlist de campos, enums y formatos;
- límites de bytes, profundidad, elementos, tiempo y compresión;
- canonicalización antes de comparar identidad, paths o firmas;
- queries parametrizadas; encoding contextual de salida;
- rechazo de tipos ambiguos y coerciones inesperadas;
- protección ante duplicate keys, integer overflow y expansión excesiva;
- deserialización a DTO explícito, no a objetos con comportamiento.

### 10.2 SSRF

Si la API recibe destinos o URLs:

- permitir esquemas y destinos explícitos;
- resolver DNS y validar todas las direcciones resultantes;
- bloquear loopback, link-local, metadata, redes privadas y rangos especiales según política;
- controlar redirects y revalidar cada salto;
- usar egress proxy/policy y segmentación;
- limitar puertos, método, tamaño y tiempo;
- tratar DNS rebinding y formatos alternativos de IP.

Un regex de URL no es defensa suficiente.

### 10.3 Consumo de recursos y abuso de flujos

**[PROD]** OWASP separa consumo irrestricto de recursos y acceso irrestricto a flujos sensibles de negocio.

- límites por identidad, tenant, IP/dispositivo y operación según amenaza;
- cuotas de CPU, memoria, bytes, conexiones, concurrencia y costo downstream;
- límites más estrictos para búsqueda, exportación, upload y cómputo costoso;
- anti-automation donde el flujo lo requiere;
- respuesta uniforme que no permita enumeración;
- presupuesto y cancelación para jobs;
- protección contra payload comprimido pequeño que expande enormemente.

### 10.4 Parsing y pequeños lenguajes

**[ACADEMIC]** MIT 6.102 trata grammars como especificaciones declarativas de strings, streams, formatos y wire protocols. Para filtros, expresiones, reglas o configuración no trivial:

```text
bytes/chars
→ tokenizer con límites
→ parser desde gramática versionada
→ AST tipado e inmutable
→ validación semántica/autorización
→ evaluación con presupuesto
```

- No sustituir una gramática recursiva por un regex críptico.
- Acotar profundidad, longitud, cantidad de nodos y costo de evaluación.
- Rechazar trailing input y ambigüedad no declarada.
- Separar syntax tree de modelo de dominio.
- Fuzzear lexer/parser y verificar que errores no filtren internals.
- Un DSL aumenta superficie de ataque y compatibilidad; crearlo sólo si resuelve una familia repetida de problemas mejor que una estructura de datos ordinaria.

El patrón interpreter facilita agregar variantes; visitor concentra nuevas operaciones. Elegir según qué eje cambiará y exigir exhaustividad de variantes.

## 11. Concurrencia y asincronía

### 11.1 Races de dominio

**[ACADEMIC]** MIT 6.102 presenta concurrencia mediante shared memory y message passing y busca diseños que no requieran razonar sobre cada interleaving.

Las races más peligrosas del backend suelen ser semánticas:

- check-then-act;
- lost update;
- doble creación/cobro;
- transición de estado duplicada;
- lease vencida pero aún usada;
- evento emitido sin commit o commit sin evento.

Elegir una garantía explícita: constraint única, compare-and-set/version, lock transaccional, serialización por clave, idempotency record o reconciliación. “Async” no elimina races.

### 11.2 Backpressure y colas

- Toda cola debe ser acotada o tener política explícita de overflow.
- Rechazar temprano puede ser más seguro que aceptar trabajo que vencerá.
- Propagar cancelación y deadlines a tareas derivadas.
- Separar concurrencia máxima de rate; ambas protegen recursos distintos.
- Medir edad de cola, no sólo longitud.
- Definir shutdown: dejar de admitir, drenar con deadline, persistir/reintentar o cancelar.

En streaming gRPC, flow control protege buffers entre sender y receiver, pero no sustituye límites de mensajes ni de trabajo de aplicación. La documentación oficial advierte deadlock si ambos extremos hacen writes síncronos/manuales sin reads; especificar quién lee, quién escribe y cómo progresa cada dirección.

### 11.3 Estado de una operación

```text
RECEIVED → VALIDATED → COMMITTED → PUBLISHED → COMPLETED
             ↘ REJECTED
                         ↘ RECONCILIATION_REQUIRED
```

**[IMPL]** Registrar el punto de no retorno. Una respuesta perdida después de `COMMITTED` es ambigua para el cliente; idempotencia y consulta de estado deben resolverla.

### 11.4 Interleaving, safety y liveness

**[ACADEMIC]** En runtimes cooperativos, cada `await` es un punto donde otra operación puede observar o cambiar estado. El código entre dos `await` puede ser mutuamente excluyente respecto del mismo event loop, pero no respecto de procesos, workers o sistemas externos.

Para cada operación async:

1. marcar cada punto de yield;
2. anotar qué invariantes deben sostenerse antes de ceder control;
3. volver a validar estado/version al reanudar si pudo cambiar;
4. no mantener locks o transacciones a través de espera remota sin prueba de liveness;
5. tratar rejection/cancelación como salida normal del protocolo;
6. nunca hacer busy-wait: bloquea progreso o desperdicia CPU.

Dos obligaciones distintas:

- **safety:** nada malo ocurre; invariantes y postcondiciones siempre se conservan;
- **liveness:** algo bueno eventualmente ocurre; no hay espera infinita, deadlock o starvation bajo supuestos declarados.

Callbacks/listeners deben terminar rápido, evitar ciclos y darse de baja al finalizar; de lo contrario retienen memoria, trabajo y comportamiento obsoleto.

## 12. Datos y transacciones en el límite del servicio

El manual `DATABASE_STORAGE_INTERNALS.md` profundiza WAL, MVCC, aislamiento e índices. En backend son obligatorias estas decisiones:

- unidad de consistencia y límites de transacción;
- invariantes defendidos por constraints;
- nivel de aislamiento requerido y anomalías toleradas;
- estrategia de concurrencia optimista/pesimista;
- migración compatible con versiones de aplicación coexistentes;
- retención, borrado, privacidad y auditoría;
- outbox/inbox o reconciliación para efectos externos.

No mantener una transacción de base abierta durante un RPC remoto salvo diseño extraordinariamente justificado. No afirmar “exactly once” end-to-end sin demostrar identidad, persistencia, atomicidad y recuperación en cada frontera.

Auditoría cerrada el 2026-08-21 con `DATABASE_STORAGE_INTERNALS.md`:

- el deadline incluye acquire del pool, lock wait, query, WAL/commit y serialización; cancellation debe alcanzar DB o el trabajo queda huérfano;
- `40001`, deadlock, unique violation, timeout y commit con respuesta perdida no comparten política automática; clasificar y reintentar la transacción completa sólo cuando sea seguro;
- idempotency record y efecto de negocio deben compartir transacción o tener reconciliación demostrada;
- expand/contract se valida con versiones coexistentes, backfill resumible, WAL/replica lag y rollback que preserve writes nuevos;
- pool size es admission control sobre CPU/I/O/locks; una conexión por request o un pool enorme desplaza saturación;
- cache invalidation, outbox/CDC y replica reads deben declarar orden, cursor/LSN, staleness y recovery.

## 13. Observabilidad y operación

### 13.1 Telemetría mínima por operación

| Señal | Dimensiones controladas |
|---|---|
| tasa/goodput | operación, resultado, versión |
| latencia | end-to-end y dependencias; percentiles |
| errores | código estable, origen, retryable |
| saturación | concurrencia, queue age/depth, pools |
| retries | intentos, éxito recuperado, amplificación |
| dependencias | deadline, status, circuit state |

Evitar labels de alta cardinalidad como user ID, URL cruda o error message. Correlacionar con trace/request ID; no usar el ID como credencial.

### 13.2 Logs y auditoría

Log estructurado mínimo:

```text
timestamp, service, version, operation, request_id, trace_id,
principal_class, tenant_hash/opcional, outcome_code,
duration, dependency_summary, retry_count
```

- redactar tokens, cookies, passwords, secrets y payload sensible;
- separar log operativo de audit log resistente a alteración;
- sincronizar tiempo y declarar retención/acceso;
- registrar decisión de autorización suficientemente para investigar, sin exponer datos innecesarios.

### 13.3 SLO y alertas

Definir SLIs desde experiencia del consumidor: proporción de solicitudes válidas correctas dentro de latencia, frescura o durabilidad requerida. Alertar sobre consumo de error budget y síntomas accionables, no cada excepción aislada.

## 14. Estrategia de pruebas

### 14.1 Pirámide por riesgo, no por dogma

| Nivel | Prueba | Falla que localiza |
|---|---|---|
| unidad | dominio/función en entorno controlado | regla local |
| propiedad/fuzz | invariantes sobre muchas entradas | bordes no enumerados |
| contrato | implementación contra OpenAPI/proto/schema | drift proveedor-consumidor |
| integración | adaptadores reales: DB, broker, auth | semántica entre componentes |
| end-to-end | camino desplegado | configuración y composición |
| carga/resiliencia | concurrencia, fallos, degradación | saturación, tails y recuperación |

**[PROD]** Google enfatiza escribir, ejecutar y reaccionar: un test rojo ignorado destruye confianza. Probar fallos mediante excepciones, latencia y errores RPC controlados; cobertura de líneas no equivale a calidad.

### 14.2 Matriz obligatoria de API

- éxito mínimo y completo;
- request malformado, demasiado grande y media type incorrecto;
- credencial ausente, inválida, expirada y audience incorrecta;
- objeto, propiedad y función de otro actor;
- transición de dominio inválida;
- duplicado e idempotency conflict;
- retry después de timeout con efecto desconocido;
- dependencia lenta, caída, corrupta o inconsistente;
- concurrencia sobre el mismo recurso;
- compatibilidad cliente viejo/servidor nuevo y viceversa;
- migración en coexistencia y rollback;
- logs/metrics sin secretos.

### 14.3 Tests que sobreviven cambios

- Probar comportamiento público, no secuencia interna de llamadas salvo que sea contrato.
- Usar fakes sólo si mantienen semántica relevante; reservar mocks para fronteras estrechas.
- Un bug corregido produce primero un test de regresión que falla por la causa correcta.
- Tests deterministas: controlar reloj, aleatoriedad, IDs, red y scheduling donde sea posible.
- Separar datos de test de lógica para cubrir particiones con claridad.
- Documentar la estrategia de partición y cobertura funcional; una suite pequeña y justificable puede ser mejor que casos redundantes.
- Mantener tests DAMP —descriptivos y significativos— aunque repitan algo de setup; aplicar DRY sólo si la abstracción conserva claridad de cada caso.
- Probar cada operación de un ADT combinando creators, producers, observers y mutators, no métodos aislados sin observar su valor abstracto.

## 15. Aprender de errores y depurar científicamente

**[ACADEMIC]** MIT 6.102 enseña el ciclo estudiar → hipótesis → experimento → repetir, con registro cuando el problema supera unos pocos intentos. El objetivo no es esconder el síntoma, sino aumentar confianza en la causa y la corrección.

### 15.1 Registro causal

```text
síntoma observable:
impacto y alcance:
primera evidencia confiable:
hipótesis actual:
predicción falsable:
experimento mínimo:
resultado observado:
hipótesis aceptada/rechazada:
causa raíz técnica:
condiciones que la hicieron posible:
fix:
test de regresión:
telemetría que detectará recurrencia:
acción sistémica y owner:
```

### 15.2 Método

1. Reproducir o acotar con evidencia; preservar request ID, versión y timeline.
2. Estudiar primero el sistema y el cambio, no modificar al azar.
3. Formular una causa que prediga una observación adicional.
4. Cambiar o instrumentar una variable material.
5. Repetir y mantener audit trail.
6. Reducir el caso hasta aislar la frontera defectuosa.
7. Crear regresión que falle antes del fix y pase después.
8. Buscar la misma clase de error en código vecino.
9. Mejorar diseño, gate o observabilidad para que el sistema aprenda.

No cerrar con “error humano”. Preguntar qué contrato, revisión, test, tipo, permiso o mecanismo permitió que un error individual llegara a producción.

## 16. Compatibilidad y evolución

### 16.1 Tres compatibilidades

**[PROD]** Google AIP-180 separa:

| Tipo | Pregunta |
|---|---|
| source | ¿el código consumidor compila contra la nueva librería? |
| wire | ¿cliente y servidor de versiones distintas se comunican? |
| semántica | ¿el consumidor recibe el comportamiento razonablemente esperado? |

Añadir un campo puede romper consumidores que rechazan desconocidos; cambiar orden, timing, formato de mensaje o longitud puede romper dependencias observables. Compatibilidad se prueba, no se declara sólo por SemVer.

### 16.2 Cambios aditivos y sustractivos

Secuencia segura general:

```text
expand contract/schema
→ desplegar lectores tolerantes
→ desplegar escritores nuevos/dual-write si está demostrado
→ backfill y verificar
→ cambiar lecturas
→ detener escritura vieja
→ observar ventana completa
→ contract/remove
```

- No reutilizar identificadores de campos.
- No renombrar como delete+add sin plan de clientes.
- Mantener resource names estables entre versiones.
- Deprecar con owner, consumidores, fecha, métricas y ruta automatizable.
- Versionar cuando hay ruptura semántica real, no para evitar disciplina de compatibilidad.

### 16.3 Dependencias

**[PROD]** Google trata cada dependencia como una relación continua de confianza y costo, no una importación gratuita.

Registrar:

- propósito y alternativas;
- maintainer/licencia/procedencia;
- versión y lockfile;
- vulnerabilidades y actualización;
- superficie transitiva;
- compatibilidad y rollback;
- comportamiento si desaparece o falla.

## 17. Revisión de diseño y código

### 17.1 Design review antes de la API

Documento corto:

```text
contexto y no-objetivos
contratos e invariantes
modelo de datos y ownership
API y errores
amenazas y abuso
capacidad/SLO
dependencias y fallos
alternativas descartadas
migración, rollout y rollback
preguntas abiertas
```

La revisión de código no es el momento de descubrir que el contrato o la arquitectura nunca fueron acordados.

### 17.2 Code review

**[PROD]** En Google prácticamente todo cambio se revisa antes de integrarse. La revisión combina corrección/comprensión, ownership y legibilidad; automatización elimina trabajo mecánico para que la persona examine decisiones.

Checklist:

- ¿resuelve un problema real con el menor cambio coherente?;
- ¿la descripción explica qué y por qué?;
- ¿preserva invariantes y casos de fallo?;
- ¿contrato, código y tests coinciden?;
- ¿authz ocurre sobre cada objeto/propiedad?;
- ¿timeouts, cancelación, retries e idempotencia son coherentes?;
- ¿hay límites y backpressure?;
- ¿migración y rollback funcionan con versiones coexistentes?;
- ¿telemetría ayuda sin filtrar datos?;
- ¿la complejidad futura está justificada?

Comentarios humanos deben ser concretos y profesionales. Cada comentario se responde o resuelve; no expandir cambios indiscriminadamente.

**[ACADEMIC]** La pasada local de claridad también busca clones/DRY mal aplicado, comentarios que explican `why` y no repiten `what`, fail-fast, números mágicos, variables con una sola función, nombres precisos, globals mutables, funciones que retornan resultados en vez de imprimirlos y special cases que deberían modelarse uniformemente. Refactorizar preserva comportamiento observable y se protege con tests.

## 18. Entrega y operación segura

### 18.1 Pipeline

```text
format/lint/typecheck
→ unit/property/security checks
→ contract/integration
→ build artefacto inmutable + SBOM/provenance
→ entorno de alta fidelidad
→ canary/segmento
→ evaluación automática de SLI
→ promoción gradual
→ rollback o completar
```

**[PROD]** Google recomienda promover el mismo release candidate entre entornos, sin recompilarlo; feature flags y despliegue selectivo reducen blast radius, pero agregan estado que debe retirarse.

### 18.2 Compatibilidad de despliegue

Antes de desplegar:

- aplicación nueva con schema viejo/expandido;
- aplicación vieja con schema nuevo;
- clientes viejos contra servidor nuevo;
- jobs y consumers rezagados;
- rollback después de que datos nuevos fueron escritos;
- cachés, tokens y eventos producidos por ambas versiones.

### 18.3 Runbook mínimo

```text
servicio/owner/on-call:
SLO y dashboards:
dependencias críticas:
señales de saturación:
errores conocidos y clasificación:
cómo reducir tráfico/costo:
cómo desactivar feature:
cómo hacer rollback:
cómo drenar y reiniciar:
cómo reconciliar efectos ambiguos:
escalamiento y evidencia a preservar:
```

## 19. Seguridad API: gate OWASP 2023

**[PROD]** Auditoría explícita de los diez riesgos:

| Riesgo | Evidencia exigida |
|---|---|
| API1 BOLA | consulta y test por objeto/tenant |
| API2 Broken Authentication | flujos, tokens, rotación, rate/abuse tests |
| API3 Broken Object Property Level Authorization | allowlist de lectura/escritura y tests |
| API4 Unrestricted Resource Consumption | límites de tamaño, tiempo, costo y concurrencia |
| API5 Broken Function Level Authorization | policy por operación/rol y deny tests |
| API6 Sensitive Business Flows | threat model de automatización/abuso |
| API7 SSRF | egress policy, resolución/redirect tests |
| API8 Security Misconfiguration | hardening reproducible, headers, defaults, secrets |
| API9 Improper Inventory Management | catálogo de hosts/versiones/endpoints/owners |
| API10 Unsafe Consumption of APIs | validación, timeout, límites y confianza mínima |

Reglas de cierre del gate:

- **API1:** un identificador impredecible sólo agrega defensa en profundidad; cada acción sobre cada objeto debe volver a comprobar autorización y aislamiento de tenant.
- **API2:** login, recuperación, MFA y emisión/renovación de tokens requieren controles de abuso propios. Rate limiting general no basta; lockout y CAPTCHA deben diseñarse sin crear denegación de servicio contra cuentas legítimas.
- **API3:** serializar y aceptar sólo propiedades permitidas por caso de uso; no enlazar automáticamente payload externo con objetos internos privilegiados.
- **API4/API6:** limitar tamaño, duración, memoria, CPU, concurrencia y gasto de terceros, y modelar por separado la automatización abusiva de flujos de negocio válidos.
- **API5:** una capa de autorización consistente debe ser invocada por toda función de negocio; ocultar una ruta en el cliente no es control de acceso.
- **API7:** resolver, validar y restringir destino, esquema, puerto, redirects y egress; volver a validar después de cada resolución o redirección relevante.
- **API8:** hardening y configuración deben ser reproducibles, revisados continuamente y aplicados también a respuestas de error; validar que todo response cumpla su schema y no revele secretos.
- **API9:** inventariar ambientes, hosts, versiones, endpoints, owners, dependencias y flujos de datos, con evidencia de retiro de versiones antiguas.
- **API10:** evaluar postura del proveedor, exigir TLS, validar y limitar toda respuesta externa y usar allowlists para redirects; una API confiable también puede entregar datos defectuosos o comprometidos.

Este gate no sustituye threat modeling ni el manual de seguridad/SRE; garantiza que una API no llegue a producción ignorando clases conocidas.

## 20. Gate profesional antes de producción

La respuesta debe ser demostrable para cada ítem:

1. ¿Existe contrato ejecutable y ejemplos de bordes?
2. ¿Invariantes están centralizados y defendidos también bajo concurrencia?
3. ¿Métodos, status, caché y errores respetan semántica oficial?
4. ¿La identidad del error es machine-readable y estable?
5. ¿Cada llamada tiene deadline y cancelación útil?
6. ¿Todo retry es selectivo, acotado, con jitter e idempotencia demostrada?
7. ¿Autorización cubre función, objeto y propiedad?
8. ¿Inputs y outputs externos tienen schema y límites?
9. ¿Colas, pools y trabajo asíncrono tienen backpressure?
10. ¿Clientes y servidores con version skew están probados?
11. ¿Migración, canary y rollback preservan datos?
12. ¿Logs, métricas y traces permiten explicar un fallo sin filtrar secretos?
13. ¿Las pruebas incluyen fallos y efecto ambiguo, no sólo happy path?
14. ¿Hay owner, SLO, alertas y runbook?
15. ¿Se aprendió de bugs anteriores mediante regresiones y controles sistémicos?

Si una respuesta material es “no sé”, registrar `[OPEN]`; no ocultarla con lenguaje de certeza.

### 20.1 Cruce con microestructura y order entry

`MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md` define la semántica económica; este manual gobierna adapters, APIs y lifecycle distribuido:

- `timeout/5xx/disconnect` puede dejar outcome `UNKNOWN`; no convertirlo en error terminal ni retry ciego;
- client/order/execution IDs, endpoint data source y reconciliation forman el contrato de idempotencia;
- public market data, private execution/drop-copy y account query son planos distintos con auth, freshness y SLO propios;
- REST/WebSocket/FIX/SBE son transports/schemas, no una state machine universal; preservar estados y unknown fields del venue;
- rate limits, deadlines y circuit breakers deben permitir cancel/risk-reducing actions según policy;
- cualquier endpoint capaz de operar permanece fuera de agentes/readers por defecto y detrás de risk/authorization gate independiente.

Cruce recíproco auditado el 2026-08-21: un adapter “limpio” que pierde STP, amend priority, partial fill o ambiguous outcome es incorrecto.

## 21. Contrato para un agente Codex

Al solicitar implementación backend, incluir:

```text
OBJETIVO
- usuario, resultado y no-objetivos

CONTRATO
- recursos, operaciones, invariantes, schemas, errores
- consistencia, idempotencia, orden y compatibilidad

RIESGO
- trust boundaries, datos sensibles, abuso y permisos

OPERACIÓN
- SLO, carga, límites, deadlines, retries y backpressure

CAMBIO
- archivos/módulos permitidos, migración, rollout y rollback

EVIDENCIA
- tests requeridos, benchmark si aplica, métricas y criterios de aceptación
```

Instrucción base al agente:

> Inspecciona primero contratos, consumidores, tests, migraciones y telemetría existentes. No inventes semántica ausente. Mantén el cambio mínimo y compatible; separa validación, autorización, dominio y adaptadores. Explicita efectos bajo timeout, duplicación y concurrencia. Implementa pruebas que fallen por la causa correcta, ejecuta gates relevantes y reporta evidencia, riesgos residuales y rollback.

### 21.1 Formato de entrega del agente

```text
contrato interpretado:
supuestos:
archivos y decisiones:
invariantes preservados:
seguridad:
compatibilidad/migración:
tests ejecutados y resultados:
medición:
observabilidad:
riesgos residuales [OPEN]:
rollback:
```

## 22. Inventario de fuentes

### 22.1 MIT 6.102, Spring 2025

**[ACADEMIC]** Archivo oficial auditado lectura por lectura el 2026-08-21: 19 lecturas, cinco problem sets con fase alpha/review/beta y un proyecto de equipo iterativo. Se contrastaron objetivos, índices, resúmenes y pasajes operativos; ejercicios/evaluaciones no se reproducen.

| # | Lectura | Extracción operativa |
|---:|---|---|
| 1 | Static Checking | tipos, mutabilidad y supuestos explícitos |
| 2 | Testing | test-first, particiones, bordes y regresión |
| 3 | Code Review | comprensión compartida y defectos antes de integrar |
| 4 | Specifications | frontera cliente/implementador |
| 5 | Designing Specifications | fuerza de pre/postcondiciones y libertad de implementación |
| 6 | Abstract Data Types | comportamiento independiente de representación |
| 7 | Abstraction Functions & Rep Invariants | estados representables versus legales |
| 8 | Interfaces & Subtyping | sustitución y contratos entre implementaciones |
| 9 | Functional Programming | inmutabilidad y composición |
| 10 | Equality | identidad, valor y contratos de comparación |
| 11 | Recursive Data Types | estructuras inductivas y operaciones totales |
| 12 | Grammars & Parsing | entradas formales y parsing seguro |
| 13 | Debugging | hipótesis, experimento, repetición y audit trail |
| 14 | Concurrency | interleavings, shared memory y message passing |
| 15 | Promises | composición asíncrona y errores |
| 16 | Mutual Exclusion | seguridad de estado compartido |
| 17 | Callbacks & GUIs | eventos y separación modelo/vista |
| 18 | Message-Passing & Networking | procesos, mensajes y fallos de red |
| 19 | Little Languages | gramáticas y DSL acotados |

### 22.2 Normativa y contratos

| Fuente | Versión verificada | Uso |
|---|---|---|
| RFC 9110 | 2022 | semántica, métodos, condiciones, ranges y status auditados |
| RFC 9111 | 2022 | storage, reuse, freshness, validation e invalidation auditados |
| RFC 9112 | 2022 | framing y HTTP/1.1 auditados al nivel del manual |
| RFC 9113 | 2022 | HTTP/2 auditado al nivel del manual; internals van a redes |
| RFC 9114 | 2022 | HTTP/3 auditado al nivel del manual; QUIC va a redes |
| RFC 5789 | 2010 | PATCH y `Accept-Patch` auditados |
| RFC 9457 | 2023 | Problem Details; obsoleta RFC 7807 |
| RFC 8446 | 2018 | TLS 1.3 |
| RFC 9700 | 2025 | OAuth 2.0 Security BCP |
| RFC 8725 | 2020 | JWT BCP |
| RFC 6585 | 2012 | 428/429/431/511; rate limiting HTTP |
| OpenID Connect Core | 1.0 + Second Errata | autenticación sobre OAuth 2.0 |
| OpenAPI | 3.2.0, 2025-09-19 | descripción, dialecto, security y riesgos de tooling auditados |
| JSON Schema | 2020-12 | core, validation, vocabularios y annotations auditados |
| gRPC docs | verificadas 2026-08-21 | deadline, cancellation, retry, status, health y flow control auditados |
| Protobuf docs | verificadas 2026-08-21 | evolución y compatibilidad de mensajes auditadas |

### 22.3 Producción y seguridad

- *Software Engineering at Google*: tabla de contenidos completa inspeccionada; se auditaron tiempo/escala, documentación, code review, testing, dependencias, deprecación, CI/CD y cambios grandes.
- Google AIPs: catálogo inspeccionado y AIP-121, 122, 131, 155, 158, 180, 193 y 194 auditadas para recursos, métodos, paginación, request identification, freshness, errores, retries y compatibilidad.
- Amazon Builders’ Library: auditados timeouts, retries, backoff+jitter, idempotent APIs, fallos correlacionados, sobrecarga y aislamiento de dependencias.
- OWASP API Security Top 10 2023: los diez riesgos y sus mitigaciones fueron auditados el 2026-08-21 e incorporados como gate mínimo.

### 22.4 Auditoría transversal con sistemas, redes, GPU y seguridad/SRE

Cruce completado el 2026-08-21 con `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md`:

| Decisión backend | Obligación de sistemas |
|---|---|
| SLO de API | distribución p50/p95/p99/p99.9, goodput y carga ofrecida |
| serialización/parsing | CPU, allocations, copias, tamaño y límites medidos |
| concurrencia | threads/event loop, pools, run queue, locks y memoria |
| colas/backpressure | Little (`L=λW`), edad de cola, capacidad y rechazo |
| retry | amplificación de trabajo y degradación bajo saturación |
| caché | working set, hit ratio, invalidación, coherencia y memoria |
| logging/tracing | overhead y cardinalidad medidos; sampling justificable |
| optimización | misma semántica, authz, errores y tests antes/después |

La medición localiza costo; no autoriza debilitar el contrato. Un benchmark de handler sin proxy, TLS, auth, serialización, dependencias y colas sólo prueba ese componente.

Cruce completado el 2026-08-21 con `NETWORKING_DISTRIBUTED_STREAMING.md`:

| Decisión backend | Obligación de red/distribución |
|---|---|
| retry/hedging | sólo con idempotencia/deduplicación, presupuesto y cancelación |
| deadline | incluir DNS/connect/TLS/queue/write/read y propagar tiempo restante |
| RPC remoto | modelar partial I/O, pérdida de respuesta y efecto ambiguo |
| streaming | framing, flow control, límites, cierre y slow-consumer policy |
| eventos | key/partition, orden, delivery, commit, replay y side effects |
| despliegue | reconnect/rebalance/version skew y recuperación de estado |

La red puede duplicar, retrasar, reordenar o perder observaciones; no debe decidir por accidente la semántica de negocio. HTTP/3, gRPC, WebSocket o Kafka sólo son correctos cuando su contrato de aplicación hace explícitas esas consecuencias.

Cruce completado el 2026-08-21 con `GPU_ACCELERATED_COMPUTING.md`:

| Decisión backend | Obligación del trabajo GPU |
|---|---|
| admisión/batching | capacidad y espera acotadas; queue time forma parte de la latencia |
| deadline/cancelación | propagar presupuesto y descartar resultados tardíos; cancelar request no prueba que el device haya parado |
| concurrencia | ownership explícito de context/stream/buffers por worker o mecanismo sincronizado documentado |
| finalización | enqueue no equivale a completion; observar event/sync/error antes de publicar output |
| fallo/OOM | no reutilizar resultado o communicator dudoso; fallback, aislamiento y rollback probados |
| SLO | medir request completo, cold/warm, p99.9, goodput, memoria y calidad; no sólo kernel time |

La GPU es una dependencia asíncrona y finita del servicio. El handler conserva authz, idempotencia, deadline y lifecycle aunque el trabajo device continúe; dynamic batching no puede volver ilimitada la cola ni mezclar tenants o respuestas.

Cruce completado el 2026-08-21 con `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md`:

| Decisión backend/API | Obligación de seguridad/SRE/cloud |
|---|---|
| identidad/autorización | trusted service autoriza cada operación sobre recurso/tenant real; cache y cambios de policy tienen freshness/revocation |
| input y lógica | límites antes de trabajo caro; ASVS aplicable, abuse cases, fuzz y pruebas negativas/sequence/atomicity |
| secrets/datos | workload identity efímera; minimización, redaction, retention/delete y ninguna credencial en payload/log/trace |
| dependencia/tercero | TLS peer esperado, egress/SSRF policy, quota/costo, timeout/circuit y failure/exit plan |
| build/deploy | mismo digest promovido con provenance/policy; config/schema/mixed-version, canary, rollback y desired-vs-deployed |
| operación | SLI/SLO/error budget, bounded overload, alert accionable, runbook/on-call, incident y restore drill |
| agente/LLM | contenido no confiable no concede tools/capabilities; side effects reautorizados y acotados |

El gate OWASP de §19 es un mínimo, no el threat model completo. El proyecto selecciona requisitos ASVS exactos y conserva evidencia; “autenticado”, “interno” o “detrás de cloud” no sustituye autorización, límites ni recovery.

Cruce completado el 2026-08-21 con `FRONTEND_PRODUCT_ENGINEERING_UX.md`:

| Contrato backend/API | Obligación frontend |
|---|---|
| recurso/permiso/tenant | UI puede ocultar o deshabilitar por UX, pero servidor autoriza cada operación real |
| schema/error/deadline | validar runtime; representar loading/empty/partial/error/timeout y mensajes accionables |
| idempotencia/outcome ambiguo | operation ID, pending/reconcile; no repetir mutación por doble click o retry ciego |
| paginación/freshness/cache | cursor y autoridad explícitos; no mezclar páginas/versiones ni presentar stale como live |
| compatibilidad/rollout | server/client coexistentes, feature flags y rollback probados sobre contratos versionados |
| observabilidad | correlación por release/route/operation sin PII, secretos ni detalles DOM accidentales |

El navegador es un cliente no confiable y cancelable. Abort local, cierre de pestaña o navegación no prueban que el servidor no haya cometido; el API ofrece idempotencia, consulta/reconciliación y estados terminales suficientes para que la UI diga la verdad.

Cruce completado el 2026-08-21 con `ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md`:

| Decisión backend | Obligación algorítmica |
|---|---|
| lookup/cache/dedup | justificar hash, árbol, índice o filtro por orden, adversario, memoria y falsos positivos |
| paginación/ranking/top-`k` | orden total y tie-break determinista; no reordenar páginas bajo cursor sin contrato |
| queue/scheduler | capacidad acotada, prioridad/fairness y complejidad de enqueue/dequeue bajo overload |
| validación/parsing | límites de tamaño, profundidad, expansión y peor caso antes de trabajo caro |
| operación compuesta | invariante de dominio y atomicidad; complejidad no sustituye transacción/idempotencia |
| optimización | baseline/oráculo, prueba de corrección y benchmark end-to-end con distribución real |

La estructura interna no debe filtrarse como semántica accidental del API. Un resultado `O(1)` esperado no autoriza memoria ilimitada, orden no determinista, hash flooding ni pérdida de aislamiento por tenant.

Cruce con `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md` cerrado el 2026-08-21: arquitectura decide bounded contexts, topology, ownership y quality/failure scenarios; backend materializa esos límites en specifications, APIs, módulos, lifecycle y tests. Un service boundary no es válido si obliga distributed transactions, chatty calls o despliegues coordinados sin que el beneficio supere el costo. El contrato ejecutable y la telemetría deben coincidir con C4/runtime/ADR.

Cruce con `DATA_ENGINEERING_ANALYTICS.md` cerrado el 2026-08-21: source/serving API declara pagination/cursor, snapshot, rate, schema, event/effective time, correction/delete y idempotency. Un API no autoriza scraping “latest” irreproducible; un mart no sustituye servicio operacional sin SLO/authz. Dataset/run/operation IDs permiten reconciliar ingestión, respuesta y lineage sin exponer PII.

Cruce con `TOOLCHAINS_BUILDS_PACKAGING_FFI.md` cerrado el 2026-08-21: API/schema/codegen/client packages comparten version/compatibility contract; CI prueba artifact instalado, no sólo source tree. Service/container se construye herméticamente, porta SBOM/provenance y se promueve por digest. Binary rollback se coordina con schema/config/protocol; native extension/FFI no puede ocultar crash, leak o thread-unsafety detrás del handler.

Cruce con `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` cerrado el 2026-08-21: el API trata mobile/desktop como clientes intermitentes y version-skewed: idempotency, cursor/version, auth refresh, cache validators, bounded payloads y compatibility window permiten offline outbox, retry y reconcile. Push sólo despierta o señala; no reemplaza la fuente autoritativa. Deep links y app integrity son inputs/señales, no autorización. Rollout del servidor preserva clientes store-delayed y pending operations.

### 22.5 Práctica pública xAI — protobuf, SDK sync/async y streaming

**[CODE]** `xai-org/xai-proto@723dd2aa22d17be35617463837dc47cda008d90e` publica contratos protobuf y `xai-org/xai-sdk-python@4dab6a449736890a70d20bc886221a4219ff6ae7` un SDK gRPC con clientes síncrono/asíncrono, streaming, deadlines, retry configurable, telemetry y tests. Ambos declaran Apache-2.0. Son evidencia de diseño de cliente público, no de internals del servicio.

Patrones transferibles:

- proto es fuente del wire contract; wrappers añaden ergonomía sin ocultar status/deadline;
- sync y async deben conservar semántica, errores, cancelación y cobertura equivalentes;
- streaming separa chunk incremental de respuesta acumulada y mantiene IDs de tool calls;
- retry queda desactivado o configurado explícitamente según idempotencia y presupuesto;
- el objeto raw proto sigue accesible para forward compatibility, pero uso recurrente señala un hueco del wrapper;
- OpenTelemetry registra lifecycle y metadata permitida sin payloads/secrets por defecto.

**Gate para SDK:** golden protobuf, unknown fields, compatibilidad entre revisiones, fragmentación arbitraria de chunks, backpressure/consumer lento, cancelación y deadline, error mapping, retry amplification, sync↔async parity, thread/task safety y redacción de traces. Fijar versión del paquete y del proto; no depender de un alias móvil cuando el workflow exige reproducibilidad.

Fuentes fijadas: [xai-proto](https://github.com/xai-org/xai-proto/tree/723dd2aa22d17be35617463837dc47cda008d90e) y [xai-sdk-python](https://github.com/xai-org/xai-sdk-python/blob/4dab6a449736890a70d20bc886221a4219ff6ae7/README.md).

## 23. Fuentes oficiales

### Academia

- <https://web.mit.edu/6.102/www/sp25/>
- <https://web.mit.edu/6.102/www/sp26/classes/02-testing/>
- <https://web.mit.edu/6.102/www/sp26/classes/05-designing-specs/>
- <https://web.mit.edu/6.102/www/sp26/classes/13-debugging/>
- <https://web.mit.edu/6.102/www/sp26/classes/14-concurrency/>
- <https://web.mit.edu/6.102/www/sp26/classes/16-mutual-exclusion/>

### Especificaciones

- <https://www.rfc-editor.org/rfc/rfc9110.html>
- <https://www.rfc-editor.org/rfc/rfc9111.html>
- <https://www.rfc-editor.org/rfc/rfc9112.html>
- <https://www.rfc-editor.org/rfc/rfc9113.html>
- <https://www.rfc-editor.org/rfc/rfc9114.html>
- <https://www.rfc-editor.org/rfc/rfc5789.html>
- <https://www.rfc-editor.org/rfc/rfc9457.html>
- <https://www.rfc-editor.org/rfc/rfc8446.html>
- <https://www.rfc-editor.org/rfc/rfc9700.html>
- <https://www.rfc-editor.org/rfc/rfc8725.html>
- <https://www.rfc-editor.org/rfc/rfc6585.html>
- <https://openid.net/specs/openid-connect-core-1_0.html>
- <https://spec.openapis.org/oas/v3.2.0.html>
- <https://json-schema.org/specification>
- <https://grpc.io/docs/guides/>
- <https://grpc.io/docs/guides/retry/>
- <https://grpc.io/docs/guides/flow-control/>
- <https://protobuf.dev/best-practices/dos-donts/>

### Ingeniería de producción

- <https://abseil.io/resources/swe-book>
- <https://google.aip.dev/general>
- <https://aws.amazon.com/builders-library/>
- <https://owasp.org/API-Security/editions/2023/en/0x00-toc/>

## 24. Límites y extensiones deliberadamente condicionadas

Este núcleo ya permite diseñar, revisar e implementar servicios y APIs generales. No sustituye:

- internals de red, consenso, entrega y streaming: `NETWORKING_DISTRIBUTED_STREAMING.md`;
- aislamiento, WAL, MVCC, índices y motores: `DATABASE_STORAGE_INTERNALS.md`;
- kernels, streams, memoria device, collectives y profiling: `GPU_ACCELERATED_COMPUTING.md`;
- IAM de infraestructura, supply chain, Kubernetes, SRE e incidentes: `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md`;
- semántica concreta del runtime/framework escogido, que deberá verificarse contra su documentación y versión oficial.

El núcleo general está auditado. Permanecen deliberadamente condicionados:

- los cruces transversales con sistemas/rendimiento, redes/distribuidos, bases de datos, GPU y seguridad/SRE/cloud quedaron cerrados el 2026-08-21;
- las reglas específicas de lenguaje, runtime, framework, proveedor y threat model, que sólo deben añadirse cuando un proyecto declare esas elecciones.
