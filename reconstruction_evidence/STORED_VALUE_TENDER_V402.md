# V402 — stored value conectado a cobro y entrega, candidato

Alcance T2802 local: candidato externo gift-tender-candidate. No se incorpora todavía a los packs, no cierra T2802, no cambia82/1039 ni45/48. Quedan HTTP/host/frontend, expiración de aprobación tras espera de outbox y packaging/admisión final. Estos pendientes son código/gates locales, no falta de credenciales.

La derivación aislada OdooCommunity19.0@99edb6dd82b7b560930c00b03b694ba700785370 conserva LGPL3, fuentes exactas, avisos, funciones/oráculos e IPC fijados. Fuente comercial ADAPTED sólo en su módulo Python; Go/SQL/DTO/perfiles y esta composición son AUTHORED glue. No se presenta el ERP ni el algoritmo completo como equivalente. La compra origen para emitir/acumular es fixture confirmado explícito, conforme al estado seleccionado; la compra destino usa los writers reales quote→accept→order→stock. No se falsifica captured para demostrar tender.

La prueba previa stored-value-pg-third aprobó persistencia de emisión, reserva concurrente sin doble gasto, canje, rechazo, reversión y loyalty. Se conserva por correspondencia de partes sin modificar. La nueva composición mantiene el importe bruto del pedido y expone payment.order_funding: base sin stored value devuelve el bruto, módulo opcional proyecta sus entradas aprobadas. Pago/SDK/reconciliación y handover leen el restante después del lock del pedido. Migraciones nuevas0066funding,0067gift/loyalty y0068entrega; la antigua0063 era candidata no publicada y se conserva en el staging original. No se editaron migraciones canónicas anteriores.

Gift parcial: bruto123456minor, gift5000, SDKStripe cobra118456. Loyalty: bruto123456, descuento6173 derivado, SDK cobra117283. Ambos usan HTTPfixture real, callback firmado, inbox/job, GET de reconciliación, handover/checklist/aceptación y recibo comercial. Gift prueba respuesta perdida con un solo POST. Una devolución parcial posterior invalida vigencia y conserva historia. La regresión cash ordinaria también pasa porque cambiaron las mismas fronteras.

FAIL823: la primera propuesta permitía el importe ajustado con un perfil histórico provider-only. Un test SDK/PG reprodujo ese efecto no seleccionado. Corrección candidata: revisión3 y opción FULL_ORDER_WITH_STORED_VALUE expresas, hash-bound; revisiones1/2 preservan semántica original. RED retenido; GREEN rechaza el perfil anterior sin crear preparación. Un script de edición se detuvo en una firma de helper normalizada por gofmt; se reanudó sólo su cola, sin contar ese intento como test.

Cobertura total: writer de observación de financiación en payment.local_funding_receipt, distinto de payment_attempt/provider_observation. Doce requests concurrentes producen un recibo/evento; replay, actor, payload, permiso, cancelación bajo outbox bloqueado y GET durable. SQL conserva FK y exclusión provider-payment XOR funding-receipt. Perfil3 con proveedor desactivado prepara la entrega, recupera contexto/resultado, completa checklist/aceptación y emite recibo comercial. Resultado observado:1funding,1funding_event,1release,0payments,0provider_observations. El cambio de estado cancelado es un fixture negativo explícito, no prueba de un endpoint de cancelación posterior a entrega física. Invalida vigencia; inversa aprobada del canje restaura2000puntos de gift sin cambiar el recibo histórico.

Última ejecución gift-zero-funding-first: cuatro tests superiores y dos subcasos gift/loyalty PASS,65migraciones, sin skips; PostgreSQL propio detenido. Compilación separada ejecutó cero tests. No SCA global, build portable, UX o producción acreditados por este resultado. Siguiente: terminar el acceso HTTP/host y portal, negativa de TTL/commit, fuente/avisos/G0–G8 y reconstruir esta revisión exacta antes de componer.

| Recibo en staging | SHA-256 |
|---|---|
| `gift-tender-connected-manifest.json` | `a9d4ce3ed73df06c9690e060e017a2572f2e7a981b65485757452adc43d26d4d` |
| `STORED_VALUE_TENDER_ADMISSION.candidate.md` | `acebd468fe7e640496ab2f50ebdbd64a24c61da4ad3e98340272b4f17abdb1a7` |
| `gift-loyalty-focal-receipt.json` | `7846ee5029ff6c489d67334f1e3cf915f458e22bb6e4198b691aaa57a6a6e5fa` |
| `stored-value-pg-third/result.json` | `250ca91a74d228a6532df2fc71e71c8909eaf0870b734fdbb1998f145e8d09fc` |
| `gift-tender-pg-first/result.json` | `6114c8078bd2654f038e0185d89ea5bbc83b91458e23e4b6b396fdbaa263ceae` |
| `gift-tender-policy-red/result.json` | `1f51e7d7cfa14ce54c7c4e7da463fd89ee7060b3483747b816547a4258b6979b` |
| `gift-tender-policy-red/connected.log` | `60181e12cf18319532bfc1dc78632adf503c6cbe0c15ad36e13d4eded7b0e968` |
| `gift-tender-fixed-green/result.json` | `0c3aa4f3f43323a50844f18594a6f9ea374d4e1bbd6ab5aef809a9cd8c558bd8` |
| `gift-zero-funding-first/result.json` | `d7bf066fe4e08f33dc98d6e4b132bfddead442fcc8c6a586d7f52ac1994cf5ef` |
| `gift-zero-funding-first/connected.log` | `f5606f763159a2580cf3cc16c4cdb438f4b963ec9f51f8d7f700b91f00fab15a` |
