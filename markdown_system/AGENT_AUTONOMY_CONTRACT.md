# Agent Autonomy Contract

## Propósito

Maximizar avance útil sin convertir los controles de calidad en esperas ceremoniales.

## El agente avanza sin preguntar

- inspección, diseño, código, configuración no secreta y tests dentro del workspace;
- selección de packs compatibles y materialización reversible;
- instalación de dependencias fijadas y licenciadas requeridas por el plan;
- refactors internos cubiertos por tests;
- ejecución de build, análisis, benchmark local y smoke tests;
- creación de adapters, mocks, sandboxes y migraciones aún no aplicadas a producción;
- corrección de defectos detectados por gates;
- documentación de supuestos reversibles.

## El agente solicita decisión sólo ante

- reglas comerciales, precios, impuestos, financiación o contabilidad no descubiertas;
- interpretación legal/regulatoria o tratamiento sensible irreversible;
- gasto externo, compra, publicación o contacto con terceros;
- credenciales reales o acceso a producción;
- pérdida/destrucción material de datos;
- cambio de alcance que altere el producto pedido;
- dos opciones irreversibles con trade-off humano genuino.

## Gates no bloqueantes

Si falla un gate local, el agente corrige y repite. Si falta una herramienta, usa alternativa segura o deja comando/evidencia pendiente sin fingir. Si un pack está condicionado, prueba sus condiciones; no pide permiso sólo por la etiqueta.

## Libertad arquitectónica

Los packs son aceleradores, no prohibiciones. El agente puede sustituir un pack cuando el blueprint lo requiera y registre:

- alternativa y razón;
- compatibilidad/migración;
- licencia/procedencia;
- tests y evidencia igual o superior;
- rollback.

