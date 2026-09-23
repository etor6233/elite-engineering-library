# Capacitación con evaluación humana — referencia local

El alumno inicia un intento, lee la versión fijada, declara su lectura y presenta
respuestas. Un revisor autorizado distinto registra su evaluación. Cada efecto
usa audit.event, approval.request/decision y platform.outbox_event existentes.
Una evaluación no crea personas, certificaciones ni permisos. El rol del curso
orienta el contenido; la autorización depende de los permisos de la sesión.

Materializar el perfil completo y aplicar migraciones hasta 0075. No incorporar
0064 de candidatos históricos. Este pack conserva todos los kinds de aprobación
anteriores. Hay dos guards físicos y tres índices de identidad/consulta. La
baja con historial se rechaza; nunca borrar evidencia para poder degradar.

Activación explícita del host Go, deshabilitada de forma predeterminada:

| Variable | Valor de referencia |
|---|---|
| TRAINING_ENABLED | true |
| TRAINING_PROFILE_FILE | Ruta absoluta a deploy/training/reference.profile.json |
| TRAINING_CONTENT_FILE | Ruta absoluta a training_content/help.bundle.json |
| TRAINING_PROFILE_ID | reference-onboarding |
| TRAINING_PROFILE_REVISION | 1 |
| TRAINING_PROFILE_SHA256 | 9e897f609970a1e46c16944a6e6ec9abfc2a1ba8c2e4a3a4d0385095cb540907 |
| TRAINING_TENANT_ID | 50f38793-8a22-4f6b-983f-81dd0fca8208 |
| TRAINING_ORGANIZATION_ID | store-1 |

Son identidades sintéticas, sin cuentas ni secretos. Al materializar otro negocio,
fijar sus identidades y programa de capacitación, producir una nueva revisión del
perfil y calcular su SHA-256; no editar intentos ni evaluaciones previas. El host
rechaza paths relativos, archivos no regulares y hashes/identidades inconsistentes.
El proceso normal de identidad debe emitir training:learn o training:review
para la organización; estas variables no conceden autorización.

El contenido contiene 15 guías públicas de la misma release. SHA fuente:
17d210d1ad3bb0c80bdb94a8d8623cc2b8ea71a4fc94824882ef9362c189c4fd; SHA bundle: 3e636fde7b3e6cbdc109de837024424d0b7a35315b8b0a54edd874fec4a042fb.
El exportador del pack web ejecuta solamente el módulo de ayuda admitido en el
proyecto. No usarlo con código de terceros sin admisión. Los antiguos intentos
conservan el snapshot y se pueden consultar; una revisión nueva no los autoriza
a seguir escribiendo. No introducir documentos privados en este bundle público.

Verificación local: TestTrainingDurableBindingsAndAtomicity, perfil, guards,
host, downgrade y TestTrainingHumanBrowserPostgres. Las pruebas de PG requieren
ELITE_TRAINING_DATABASE_URL a una base aislada con las migraciones aplicadas.
El fixture de navegador usa ELITE_TRAINING_BROWSER=1, ELITE_WEB_ROOT y
ELITE_NODE_BIN absolutos, Next previamente compilado y el browser gate ya fijado.
No se contactan proveedores live. Ejecutar el fixture solamente sobre datos
sintéticos aislados; no es un runner sobre una base productiva.

Los tests prueban cuatro hechos/cuatro eventos/una solicitud/una decisión y cero
grants/recursos. Se rechazan aprobación genérica sin evaluación, duplicación de
assessment por intento, mutación de historial, autoevaluación y acceso cruzado.
Ver reconstruction_evidence/TRAINING_CONNECTED_RELEASE_V402.md en la biblioteca.
