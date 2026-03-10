##workflow

# BambuLab Print Farm Manager

## Module-by-Module Workflow for Antigravity

Goal: Generate a production-grade remote management system for multiple BambuLab printers.

This workflow forces Antigravity to generate code in **small independent modules** to reduce failure rate.

---

# GLOBAL RULES

You are a team of:

* software architects
* backend engineers
* frontend engineers
* DevOps engineers
* security engineers

System requirements:

* scalable
* secure
* modular
* observable
* containerized

Rules:

1. Never generate more than 1500 lines per step
2. Every module must compile independently
3. Each module must have tests
4. Use Clean Architecture
5. All configs via ENV
6. API must follow OpenAPI spec
7. Every module must STOP and wait confirmation

---

# MODULE 1 — PROJECT BOOTSTRAP

Goal: Initialize project structure.

Backend:

Go or Node.js (Fastify)

Folder structure:

backend/
cmd/
internal/
domain/
repository/
service/
api/
pkg/

Frontend:

Next.js

frontend/
app/
components/
hooks/
store/
services/

Infrastructure:

docker-compose.yml

Output required:

* full folder structure
* base server
* health check endpoint
* env configuration
* logging system

STOP AND WAIT CONFIRMATION

---

# MODULE 2 — AUTHENTICATION SYSTEM

Features:

* user login
* organization support
* JWT authentication
* RBAC roles

Tables:

users
organizations
roles
permissions

API:

POST /auth/login
POST /auth/register
GET /auth/me

Security:

bcrypt password hashing
JWT refresh token
rate limiting

Output:

* auth service
* middleware
* database schema

STOP

---

# MODULE 3 — PRINTER REGISTRY

Goal:

Allow adding and managing multiple printers.

Tables:

printers
printer_status
printer_logs

Fields:

printer_id
name
ip_address
access_token
model
status
firmware_version

API:

GET /printers
POST /printers
GET /printers/:id
DELETE /printers/:id

Output:

* printer registry service
* CRUD endpoints

STOP

---

# MODULE 4 — DISCOVERY ENGINE

Features:

LAN auto discovery for BambuLab printers.

Methods:

mDNS
network scan
known IP pool

Functions:

scan_network()
detect_bambu_printer()
register_printer()

Output:

discovery daemon
printer auto-add system

STOP

---

# MODULE 5 — TELEMETRY SERVICE

Goal:

Collect printer data.

Data:

temperature
progress
job_status
error_codes
filament

Tables:

printer_metrics

Use WebSocket streaming.

Output:

telemetry collector
metrics storage

STOP

---

# MODULE 6 — PRINT JOB MANAGER

Features:

submit print
pause
resume
cancel

Tables:

print_jobs
print_history

Queue:

Redis / NATS queue.

Output:

job scheduler
queue system

STOP

---

# MODULE 7 — REALTIME SERVICE

Goal:

Live updates.

Use WebSocket gateway.

Events:

printer_status_update
job_progress
temperature_update
alerts

Output:

realtime server

STOP

---

# MODULE 8 — CAMERA STREAM

Features:

view printer camera.

Methods:

RTSP proxy
WebRTC

Output:

camera proxy service

STOP

---

# MODULE 9 — ALERT SYSTEM

Alerts:

print failure
temperature anomaly
offline printer

Delivery:

email
telegram
web notifications

Output:

alert service

STOP

---

# MODULE 10 — DEPLOYMENT

Generate:

docker-compose
production env
CI/CD pipeline

Output:

deployment scripts

STOP
