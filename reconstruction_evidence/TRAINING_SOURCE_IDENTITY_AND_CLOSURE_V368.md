# V368 — identidad de TRL y criterio de cierre funcional

2026-09-09. Mantenimiento/discovery. Checkpoint182,172owners, resume/plan
validados con kit1.3.1 reconstruido antes de continuar. No se repite V367.

## Cierre funcional

La spec0.1.2 explicita el criterio existente de ENGINEERING_EXECUTION_PLAYBOOK
§§22–23 ante la pregunta del usuario: el cierre exige las capacidades prometidas
implementadas, probadas en sus composiciones y condiciones, sin bloqueos que
contradigan el alcance y con el release exacto verificado. Los gates estructurales
no prueban por sí solos comportamiento ni calidad de datos/modelos.

Las reglas/datos/accesos/entorno/aceptación del consumer siguen siendo específicos.
No usar esa separación para diferir un defecto genérico de la biblioteca.
Candidatos/rechazos pueden conservarse como investigación, pero una capability
requerida sin implementación admitida impide el cierre. No hay garantía universal
de cero bugs, de futuras versiones ni de combinaciones nunca verificadas.

## Avance técnico: comparación completa de metadatos Git

V367 dejó pendiente contrastar TRL1.11.0 frente al duplicado accidental1.12.0.
Se resolvió el commit completo de1.12.0 desde la [release oficial](https://github.com/huggingface/trl/releases/tag/v1.12.0)
y se consultaron commit, ambos árboles recursivos y compare de GitHub.

| Observación | Resultado |
|---|---|
| base1.11.0 | d0a2cd58229dd77804e18020830d3e3c5d4b5b72 |
| head1.12.0 | 59c4a8e104413fa9f4ca1a54eaf2ff93c0f299be |
| tree base | 8533a30746599fefc0c3e1235644b48dcbf74d70 |
| tree head | 57eac04907602c09cba52d2ee9965192d3fe7b49 |
| entradas no-directorio por árbol | 580 y580; ninguna respuesta truncada |
| únicas entradas cambiadas | VERSION;6bytes en ambos árboles |
| entradas restantes | 579/579 con path/mode/type/Git SHA/size iguales |
| parentesco y compare | head tiene base como único padre;1commit ahead,0behind |

El [commit exacto](https://github.com/huggingface/trl/commit/59c4a8e104413fa9f4ca1a54eaf2ff93c0f299be)
modifica VERSION de1.11.0 a1.12.0. GitHub informa firma validada; no se ejecutó
verificación criptográfica local. El resultado compara identidades Git observadas
desde API oficial, no bytes extraídos de un archivo de distribución.

Este control acota la identidad del source entre revisiones. Permite continuar
la adquisición considerando la revisión1.12.0 exacta y trasladar sólo la lectura
estática de paths con la misma identidad Git; no hereda pruebas de runtime.
Los hashes diferentes de wheel/sdist observados en V367 permanecen válidos:
los artefactos de publicación no son byte-idénticos por cambiar versión. Falta
adquisición gobernada, inventario del paquete y vínculo publicación/source,
LICENSE/NOTICE, lock completo, SCA, pruebas y pipeline reusable. No se inventa
un digest de archive GitHub a partir del hash de tree ni del metadata JSON.

Las cuatro respuestas originales permanecen en el stage con sus digests; el
resultado compacto siguiente preserva endpoints, hashes, counts, diferencias
y límites del control. El script de investigación local no ejecuta upstream,
instala paquetes, obtiene datos/modelos ni inicia entrenamiento.

Estado43/48controles PASS y5BLOCKED, sin promoción.162packs/1461archivos y
53perfiles sin cambios;807Markdown esperado. Gate estructural183 pendiente.
Stage: %TEMP%/elite-v368-36f0d5728d0645469634d6cc9b11323b.

| Recibo relativo al stage | SHA-256 |
|---|---|
| resume182.log | `9a16ff3fbfd4481e46303ae6722aa794455e0b73a2c9a35c703c77adff231485` |
| plan182.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
| baseline182.json | `ba288454788e1afebfe3ab48dac1890ceb6e74226a9b9b18fcb316516e18077f` |
| compare_sources.py | `e8dd0766aa0be5295161db949d3d81b5d744cce52548348f7d84afbabf2da6cd` |
| head-commit.json | `07b05589820fad974ced659d52f1b3062f37edc802bd48a73412c267e5f96563` |
| base-tree.json | `7b79457b49c815afe2ef46387c0c64806096f1b556b2a8594627be3d23caf059` |
| head-tree.json | `58565a749e6ad76aafc40008bf73a69ec18c7403b25a439a0151402249da9226` |
| compare.json | `e9e61d79cfcf6f8939c69b5da547137797e1509fa116a66f5c9da00aba4a4005` |
| source-identity-result.json | `05f395b746b4dff583ba54468061b727d37355dac0551c844bed45b60ccca2c9` |

## Resultado portable — source-identity-result.json

Base64 de bytes originales; decodificar y comprobar SHA-256 anterior. Es evidencia de investigación, no código ni un pack nuevo.

````base64
ewogICJraW5kIjogImNvbXBhcmlzb24gb2Ygb2ZmaWNpYWwgR2l0IG9iamVjdCBtZXRhZGF0YSwgbm90IHNvdXJjZSBhcmNoaXZlIG9yIGluc3RhbGxlZCBhcnRpZmFjdCBlcXVpdmFsZW5jZSIsCiAgImJhc2VfY29tbWl0IjogImQwYTJjZDU4MjI5ZGQ3NzgwNGUxODAyMDgzMGQzZTNjNWQ0YjViNzIiLAogICJoZWFkX2NvbW1pdCI6ICI1OWM0YThlMTA0NDEzZmE5ZjRjYTFhNTRlYWYyZmY5M2MwZjI5OWJlIiwKICAiYmFzZV90cmVlIjogIjg1MzNhMzA3NDY1OTlmZWZjMGMzZTEyMzU2NDRiNDhkY2JmNzRkNzAiLAogICJoZWFkX3RyZWUiOiAiNTdlYWMwNDkwNzYwMmMwOWNiYTUyZDJlZTk5NjUxOTJkM2ZlN2I0OSIsCiAgImhlYWRfcGFyZW50cyI6IFsKICAgICJkMGEyY2Q1ODIyOWRkNzc4MDRlMTgwMjA4MzBkM2UzYzVkNGI1YjcyIgogIF0sCiAgImhlYWRfZ2l0aHViX3ZlcmlmaWNhdGlvbiI6IHsKICAgICJ2ZXJpZmllZCI6IHRydWUsCiAgICAicmVhc29uIjogInZhbGlkIiwKICAgICJ2ZXJpZmllZF9hdCI6ICIyMDI2LTA4LTI2VDE5OjQ1OjU0WiIKICB9LAogICJiYXNlX2xlYWZfY291bnQiOiA1ODAsCiAgImhlYWRfbGVhZl9jb3VudCI6IDU4MCwKICAidHJ1bmNhdGVkIjogZmFsc2UsCiAgImNoYW5nZWRfbGVhdmVzIjogWwogICAgewogICAgICAicGF0aCI6ICJWRVJTSU9OIiwKICAgICAgImJlZm9yZSI6IHsKICAgICAgICAibW9kZSI6ICIxMDA2NDQiLAogICAgICAgICJ0eXBlIjogImJsb2IiLAogICAgICAgICJzaGEiOiAiMTY5ZjE5YjQ5MDFlYTZiNmVkOGUzZjdmNGY0NmM2ZGI0YmUxMmRiNyIsCiAgICAgICAgInNpemUiOiA2CiAgICAgIH0sCiAgICAgICJhZnRlciI6IHsKICAgICAgICAibW9kZSI6ICIxMDA2NDQiLAogICAgICAgICJ0eXBlIjogImJsb2IiLAogICAgICAgICJzaGEiOiAiMzJiZDkzMmYzNTViMWI1ZjM5NGMyNTg2MWNkZTIxYTgyMzY2NTI0NCIsCiAgICAgICAgInNpemUiOiA2CiAgICAgIH0KICAgIH0KICBdLAogICJub252ZXJzaW9uX2xlYWZfY291bnQiOiA1NzksCiAgImFsbF9ub252ZXJzaW9uX2dpdF9pZGVudGl0aWVzX2VxdWFsIjogdHJ1ZSwKICAiY29tcGFyZV9zdGF0dXMiOiAiYWhlYWQiLAogICJhaGVhZF9ieSI6IDEsCiAgImJlaGluZF9ieSI6IDAsCiAgInRvdGFsX2NvbW1pdHMiOiAxLAogICJjb21wYXJlX2ZpbGVfbmFtZXMiOiBbCiAgICAiVkVSU0lPTiIKICBdLAogICJjb21wYXJlX3BhdGNoIjogWwogICAgewogICAgICAiZmlsZW5hbWUiOiAiVkVSU0lPTiIsCiAgICAgICJwYXRjaCI6ICJAQCAtMSArMSBAQFxuLTEuMTEuMFxuXFwgTm8gbmV3bGluZSBhdCBlbmQgb2YgZmlsZVxuKzEuMTIuMFxuXFwgTm8gbmV3bGluZSBhdCBlbmQgb2YgZmlsZSIKICAgIH0KICBdLAogICJzb3VyY2VfYXJjaGl2ZXNfYWNxdWlyZWQiOiAwLAogICJzaWduYXR1cmVzX3ZlcmlmaWVkX2xvY2FsbHkiOiBmYWxzZSwKICAidHJhaW5pbmdfam9icyI6IDAsCiAgInJlY2VpcHRzIjogWwogICAgewogICAgICAibmFtZSI6ICJoZWFkLWNvbW1pdC5qc29uIiwKICAgICAgInVybCI6ICJodHRwczovL2FwaS5naXRodWIuY29tL3JlcG9zL2h1Z2dpbmdmYWNlL3RybC9jb21taXRzLzU5YzRhOGUxMDQ0MTNmYTlmNGNhMWE1NGVhZjJmZjkzYzBmMjk5YmUiLAogICAgICAic3RhdHVzIjogMjAwLAogICAgICAic2hhMjU2IjogIjA3YjA1NTg5ODIwZmFkOTc0Y2VkNjU5ZDUyZjFiMzA2MmYzN2VkYzgwMmJkNDhhNzM0MTJjMjY3ZTVmOTY1NjMiLAogICAgICAib2JzZXJ2ZWRfYXQiOiAiMjAyNi0wOS0wOVQyMDo1NzozOC42MzQ1NjgrMDA6MDAiCiAgICB9LAogICAgewogICAgICAibmFtZSI6ICJiYXNlLXRyZWUuanNvbiIsCiAgICAgICJ1cmwiOiAiaHR0cHM6Ly9hcGkuZ2l0aHViLmNvbS9yZXBvcy9odWdnaW5nZmFjZS90cmwvZ2l0L3RyZWVzLzg1MzNhMzA3NDY1OTlmZWZjMGMzZTEyMzU2NDRiNDhkY2JmNzRkNzA/cmVjdXJzaXZlPTEiLAogICAgICAic3RhdHVzIjogMjAwLAogICAgICAic2hhMjU2IjogIjdiNzk0NTdiNDljODE1YWZlMmVmNDYzODdjMGM2NDgwNjA5NmYxYjU1NmIyYTg1OTQ2MjdiZTNkMjNjYWYwNTkiLAogICAgICAib2JzZXJ2ZWRfYXQiOiAiMjAyNi0wOS0wOVQyMDo1NzozOS4yMjM0MzUrMDA6MDAiCiAgICB9LAogICAgewogICAgICAibmFtZSI6ICJoZWFkLXRyZWUuanNvbiIsCiAgICAgICJ1cmwiOiAiaHR0cHM6Ly9hcGkuZ2l0aHViLmNvbS9yZXBvcy9odWdnaW5nZmFjZS90cmwvZ2l0L3RyZWVzLzU3ZWFjMDQ5MDc2MDJjMDljYmE1MmQyZWU5OTY1MTkyZDNmZTdiNDk/cmVjdXJzaXZlPTEiLAogICAgICAic3RhdHVzIjogMjAwLAogICAgICAic2hhMjU2IjogIjU4NTY1YTc0OWU2YWQ3NmFhZmM0MDAwOGJmNzNhNjllYzE4Yzc0MDNiMjVhNDM5YTAxNTE0MDIyNDlkYTkyMjYiLAogICAgICAib2JzZXJ2ZWRfYXQiOiAiMjAyNi0wOS0wOVQyMDo1NzozOS43NDc2MjcrMDA6MDAiCiAgICB9LAogICAgewogICAgICAibmFtZSI6ICJjb21wYXJlLmpzb24iLAogICAgICAidXJsIjogImh0dHBzOi8vYXBpLmdpdGh1Yi5jb20vcmVwb3MvaHVnZ2luZ2ZhY2UvdHJsL2NvbXBhcmUvZDBhMmNkNTgyMjlkZDc3ODA0ZTE4MDIwODMwZDNlM2M1ZDRiNWI3Mi4uLjU5YzRhOGUxMDQ0MTNmYTlmNGNhMWE1NGVhZjJmZjkzYzBmMjk5YmUiLAogICAgICAic3RhdHVzIjogMjAwLAogICAgICAic2hhMjU2IjogImU5ZTYxZDc5Y2ZjZjZmODkzOWM2OWI1ZGE1NDcxMzc3OTdlMTUwOWZhMTE2YTY2ZjVjOWRhMDBhYmE0YTQwMDUiLAogICAgICAib2JzZXJ2ZWRfYXQiOiAiMjAyNi0wOS0wOVQyMDo1Nzo0MC4zNzc4MTErMDA6MDAiCiAgICB9CiAgXQp9Cg==
````

## Cierre184

VERIFY_LIBRARY183 PASS162/1461/807/53.163Markdown de packs intactos;48definiciones y estados de tests, requisitos del contrato,readiness JSON y2verifiers intactos. Spec0.1.2 sólo explicita aceptación ya exigida; TEST09 recibe EVID78 parcial. Resultado portable byte-idéntico y4respuestas de metadata con hashes verificados. Checkpoint184 se registra con resume/plan;43/48controles siguen PASS,5BLOCKED. No se cerró un control por esa comparación ni se adquirió/instaló/entrenó un paquete o modelo.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| verify-library183.log | `4a7f6732a9e490d8fe07a8771e514387d4065735bd69807541107b0a756e3e04` |
| closure-verification183.json | `d362e48811ea211daef63ec66406c61c94e92389cae59e95071f35ba1f084392` |
| checkpoint183.log | `6759b12597d5b184cce53fd4c6d2e0beda280b643b1f163498b50c7de0057740` |
| resume183.log | `609c39b8ea5285c41445381f651a9a8c6df997523fa4f5d7cb24073ceb6ca279` |
| plan183.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
