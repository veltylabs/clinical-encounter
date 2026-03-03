```mermaid
sequenceDiagram
    actor Nurse as Nurse / Staff
    actor Doctor
    participant Reception
    participant Scheduling
    participant clinical-encounter
    participant EventBus

    Scheduling->>clinical-encounter: CreateVisit(patient_id, doctor_id, attention_at, ...)<br/>Status: created
    clinical-encounter-->>Scheduling: MedicalHistory{id, status: created}

    Reception->>clinical-encounter: MarkArrived(id)<br/>Status: arrived
    clinical-encounter->>EventBus: publish(visit.patient_arrived)
    EventBus-->>Nurse: 🔔 Patient registered

    Nurse->>clinical-encounter: AddMeasurement(visit_id, type_id, value, unit)<br/>vitals during triage
    clinical-encounter-->>Nurse: ClinicalMeasurement saved

    Nurse->>clinical-encounter: MarkTriaged(id)<br/>Status: triaged
    clinical-encounter->>EventBus: publish(visit.patient_triaged)
    EventBus-->>Doctor: 🔔 Patient ready for attention

    Doctor->>clinical-encounter: StartVisit(id)<br/>Status: in_progress

    Doctor->>clinical-encounter: AddMeasurement(visit_id, type_id, value, unit)<br/>additional during consultation
    clinical-encounter-->>Doctor: ClinicalMeasurement saved

    Nurse->>clinical-encounter: AddHistoryDetail(visit_id, catalog_item_id, qty, ...)<br/>items / services used
    clinical-encounter-->>Nurse: HistoryDetail saved

    Doctor->>clinical-encounter: CompleteVisit(id)<br/>Status: completed
    clinical-encounter->>EventBus: publish(visit.completed)
    EventBus-->>Scheduling: update reservation status
    EventBus-->>Billing: trigger invoice

    note over clinical-encounter: CancelVisit(id) available from<br/>created / arrived / triaged / in_progress<br/>publishes visit.cancelled
```
