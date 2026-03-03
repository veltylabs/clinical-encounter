# HL7 / FHIR: Análisis y Ruta de Compatibilidad

Este documento sirve como anexo al [`ARCHITECTURE.md`](ARCHITECTURE.md) principal para evaluar el nivel de cumplimiento actual del módulo `clinical-encounter` respecto a los estándares internacionales de interoperabilidad médica (HL7/FHIR) y qué pasos serían necesarios en el futuro para lograr una certificación o integración estricta.

## 1. ¿Qué es HL7 / FHIR?
**HL7** (Health Level Seven) y su iteración moderna **FHIR** (Fast Healthcare Interoperability Resources) son el estándar internacional definitivo para el intercambio de datos médicos. Proveen un modelo de datos universal (Recursos) y una arquitectura API para asegurar que los sistemas hospitalarios, clínicas, laboratorios privados y redes gubernamentales puedan "hablar" entre sí usando una semántica común, sin importar en qué lenguaje o base de datos estén construidos.

## 2. Nivel de Cumplimiento Actual

El módulo actual **no es estrictamente compatible a nivel técnico** (no exporta recursos JSON validados por FHIR ni expone una API REST FHIR), pero **está excelentemente alineado a nivel conceptual e integracional**:

1. **La FSM (Máquina de Estados) del Encuentro:** FHIR utiliza el recurso `Encounter` con un ciclo de vida casi idéntico al implementado aquí. 
   - *Modelo Local:* `created` → `arrived` → `triaged` → `in_progress` → `completed` | `cancelled`.
   - *Modelo FHIR:* `planned` → `arrived` → `triaged` → `in-progress` → `finished` | `cancelled`.
2. **Separación de mediciones clínicas:** En FHIR, lo que denominamos `ClinicalMeasurement` equivale al recurso `Observation`. FHIR exige que las observaciones vivan separadas del "encuentro médico" y simplemente apunten a él mediante referencias. Nuestra arquitectura DB ya cumple con esta separación estricta.
3. **Referencias Planas (Snapshotting):** Usar referencias planas (ej. `patient_id`, `doctor_id`) y *snapshots* es una práctica recomendada en sistemas distribuidos que buscan evitar acoplamientos rígidos con un Master Patient Index (MPI).

## 3. Hoja de Ruta para Certificación / Compatibilidad FHIR

Si en el futuro una entidad de salud o gobierno exige interoperar usando HL7/FHIR, **NO será necesario reescribir la lógica base ni cambiar el modelo de base de datos actual.** En cambio, se debe construir una "capa adaptadora" (Adapter/Translator) que convierta nuestra estructura interna a los complejos recursos FHIR.

Para lograr dicha compatibilidad plena, se deberán abordar las siguientes adaptaciones:

### 3.1 Nomenclaturas exactas (Formato de API y Modelos)
- FHIR requiere documentos JSON con estructuras anidadas muy inflexibles.
- En vez de devolver campos planos como `PatientID`, FHIR espera referencias formales: `{"subject": {"reference": "Patient/123"}}`.
- Los endpoints MCP actuales tendrían que ser envueltos o expuestos secundariamente como una API RESTful estándar (e.g., `GET /Encounter/{id}`).

### 3.2 Terminologías Internacionales (Sistemas de Codificación)
Actualmente, el modelo (`MedicalHistory`) usa campos de texto libre (`Diagnostic string`, `Reason string`). FHIR prohíbe el texto libre como única fuente de verdad médica, requiriendo esquemas internacionales:
- **CIE-10 / CIE-11 / SNOMED CT:** Para codificar los diagnósticos registrados (ej. guardar el código estandarizado `J00` en lugar del texto "Resfriado común").
- **LOINC:** En la tabla `MeasurementType`, para registrar qué tipo de medición se tomó (ej. código universal `85354-9` para Presión Arterial).
- **UCUM:** Unidades de medida clínicas codificadas internacionalmente en lugar de unidades textuales locales (ej. `mm[Hg]`).

### 3.3 El Tiempo como Período
- Nuestro sistema consolida la interacción temporal en un *timestamp* singular (`AttentionAt int64`).
- FHIR maneja el tiempo del consultorio como un tipo `Period` (compuesto por los atributos obligatorios `start` y `end`) para auditar la duración exacta en que el paciente ocupó el recurso o interactuó con el médico. Se requeriría inferir el inicio y el fin basado en los históricos de la FSM o agregar un nuevo campo al esquema en el futuro.

## Conclusión
La madurez de la arquitectura actual es la correcta para un desarrollo ágil sobre WebAssembly. Al respetar principios de separación de dominios y manejo por eventos de estado, sienta las bases perfectas para añadir interoperabilidad HL7/FHIR a futuro a través de traductores independientes, sin sacrificar rendimiento productivo en la primera fase.
