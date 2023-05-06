# FHIR Server

> HL7 FHIR server implementation for fructose.dev

## Setup

1. Clone the repo

```bash
git clone https://github.com/hawyar/fhird.git
```

2. Start the server in a container

```bash
make start
```

To stop the server either press `Ctrl+C` or run `make stop` in a separate terminal.

3. Make sure the FHIR server is running

```bash
curl http://localhost:9090/ping
```

## FHIR Server

The RESTful API uses the [FHIR R4](http://hl7.org/fhir/R4/) specification and supports [US Core Profiles]().

Base URL: `http://localhost:9090`

To access the supported resources just append the resource name (PascalCase) to the base URL `http://localhost:9090/api/v1/{resource}` e.g. `http://localhost:9090/api/v1Patient`.

## US Core Profile Resource

List of all supported resources and their endpoints.


| Resource                                                                  |                                Endpoint                               |
|---------------------------------------------------------------------------|:-----------------------------------------------------------------------:|
| [CapabilityStatement](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-capabilitystatement.html) |               /api/{version}/metadata                |
| [Patient](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-patient.html)                          |              /api/{version}/Patient                |
| [Practitioner](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-practitioner.html)                  |            /api/{version}/Practitioner            |
| [PractitionerRole](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-practitionerrole.html)              |       /api/{version}/PractitionerRole       |
| [Organization](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-organization.html)                  |           /api/{version}/Organization           |
| [Encounter](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-encounter.html)                      |              /api/{version}/Encounter              |
| [Location](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-location.html)                          |             /api/{version}/Location             |
| [Observation](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-observation-lab.html)                  | /api/{version}/Observation |
| [Condition](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-condition.html)                          |             /api/{version}/Condition             |
| [Procedure](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-procedure.html)                          |             /api/{version}/Procedure             |
| [DiagnosticReport](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-diagnosticreport-lab.html)        |   /api/{version}/DiagnosticReport  |
| [DocumentReference](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-documentreference.html)          |  /api/{version}/DocumentReference  |
| [Immunization](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-immunization.html)                  |            /api/{version}/Immunization            |
| [Medication](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-medication.html)                      |              /api/{version}/Medication              |
| [MedicationRequest](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-medicationrequest.html)          |   /api/{version}/MedicationRequest  |
| [AllergyIntolerance](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-allergyintolerance.html)        |   /api/{version}/AllergyIntolerance  |
| [CarePlan](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-careplan.html)                          |             /api/{version}/CarePlan             |
| [Goal](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-goal.html)                                  |                  /api/{version}/Goal                  |
| [ServiceRequest](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-servicerequest.html)              |          /api/{version}/ServiceRequest         |
| [DocumentReference](https://www.hl7.org/fhir/us/core/StructureDefinition-us-core-documentreference.html)        |   /api/{version}/DocumentReference  |
| [RelatedPerson](https://www.hl7.org/fhir/us/core/StructureDefinition)                                           |   /api/{version}/RelatedPerson  |
### Other Endpoints

Other endpoints that are not part of the FHIR specification but are useful for the development of FHIR clients.
| Endpoint | Description                                                                                                    |
|----------|--------------------------------------------|
| /ping | Returns `.` to indicate the server is running. |                                                                 

## Examples

Get the FHIR server's capability statement. The capability statement describes
the server's supported resources and operations so a client
may use it as an interface definition when interacting with the server.

```bash
curl -X GET http://localhost:9090/metadata
```

Create patient

```bash
curl -X POST -H "Content-Type: application/json" -d '{"resourceType": "Patient", "name": [{"given": ["John"], "family": "Doe"}]}' http://localhost:9090/Patient
```

Create procedure

```bash
curl -X POST -H "Content-Type: application/json" -d '{"subject":{"reference":"25oYHe8zCfx52wp9S8RKEVjEyTw"}}' http://localhost:9090/Procedure
```
