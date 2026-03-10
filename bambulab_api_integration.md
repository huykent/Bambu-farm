# BambuLab API Integration Specification

## For BambuLab Print Farm Manager

Purpose:

Define how the system connects to, monitors, and controls BambuLab printers.

This module integrates with the existing backend services.

Supported printers:

* BambuLab X1
* X1 Carbon
* P1P
* P1S
* A1
* A1 Mini

---

# INTEGRATION MODES

The system must support two connection modes.

## Mode 1 — LAN Control (Preferred)

Direct connection to printer inside LAN.

Protocols:

MQTT
HTTP
WebSocket

Requirements:

* Printer IP
* Access code from printer

Advantages:

* Low latency
* Full telemetry
* Full control

---

## Mode 2 — Cloud API (Fallback)

Use Bambu Cloud service.

Limitations:

* Higher latency
* Limited telemetry
* Requires Bambu account authentication

---

# PRINTER CONNECTION SERVICE

Create a service:

PrinterConnectionManager

Responsibilities:

* maintain connection to printers
* reconnect on failure
* handle telemetry stream
* send commands

Internal structure:

printer_manager/
connection_pool.go
mqtt_client.go
printer_client.go
command_sender.go
telemetry_listener.go

---

# PRINTER REGISTRATION FLOW

When adding a printer:

1. user enters printer IP
2. user enters access code
3. system verifies printer
4. system registers printer

Pseudo flow:

connect_to_printer(ip, access_code)

request_printer_info()

store printer metadata

---

# PRINTER METADATA

Store:

printer_id
serial_number
model
firmware_version
ip_address
access_code
last_seen

---

# TELEMETRY STREAM

Use MQTT subscription.

Topics:

printer/status
printer/temperature
printer/progress
printer/errors

Telemetry fields:

nozzle_temp
bed_temp
print_progress
layer
remaining_time
job_state
filament_type
fan_speed

Example telemetry payload:

{
"progress": 42,
"nozzle_temp": 215,
"bed_temp": 60,
"layer": 120,
"remaining_time": 3600
}

---

# PRINT COMMANDS

Supported commands.

## Start Print

send_gcode_file()

Inputs:

file_url
printer_id

---

## Pause Print

pause_print()

---

## Resume Print

resume_print()

---

## Cancel Print

cancel_print()

---

## Set Temperature

set_nozzle_temp(value)

set_bed_temp(value)

---

# CAMERA STREAM

Bambu printers expose camera via RTSP.

Example:

rtsp://printer-ip:8554/live

System must implement:

CameraProxyService

Responsibilities:

* convert RTSP → WebRTC
* allow browser playback

Architecture:

printer → RTSP → proxy → browser

---

# PRINTER HEALTH MONITOR

Every printer must send heartbeat.

Interval:

5 seconds

If no heartbeat for:

30 seconds

Mark printer as:

OFFLINE

Trigger alert.

---

# ERROR HANDLING

Common printer errors:

filament_runout
bed_leveling_failed
nozzle_clog
print_failed

System must:

store error logs
notify user
trigger alerts

---

# SECURITY REQUIREMENTS

Printer access codes must be:

encrypted in database

Connection rules:

* never expose printer access codes
* isolate printer connections
* rate limit commands

---

# COMMAND QUEUE

All printer commands must go through queue.

Queue type:

Redis / NATS

Reason:

prevent printer overload.

Queue flow:

API → queue → printer_worker → printer

---

# RETRY SYSTEM

If command fails:

retry 3 times

Then:

log error
notify user

---

# FILE UPLOAD

System must support uploading 3D files.

Formats:

.3mf
.gcode

Files stored in:

object storage (MinIO)

Workflow:

upload file
slice if required
send to printer

---

# PRINTER DISCOVERY SUPPORT

Discovery engine must detect Bambu printers.

Detection methods:

mDNS broadcast
known port scan

Ports to scan:

8883
8884

When printer found:

retrieve metadata
add to pending list

---

# OBSERVABILITY

Each printer connection must expose metrics.

Metrics:

printer_online
print_progress
temperature
error_count

Export metrics to:

Prometheus

---

# MODULE OUTPUT REQUIREMENTS

When implementing this module:

Generate:

PrinterConnectionManager
MQTT integration
Command sender
Telemetry listener
Camera proxy

Provide:

API endpoints

GET /printers/:id/status
POST /printers/:id/start
POST /printers/:id/pause
POST /printers/:id/resume
POST /printers/:id/cancel

---

# DEVELOPMENT RULES

1. Each printer runs in isolated goroutine/service.
2. Commands must be asynchronous.
3. Telemetry must stream in real time.
4. Connection must auto-reconnect.

---

# STOP CONDITION

After generating the PrinterConnectionManager:

STOP.

Wait for confirmation before generating:

* camera proxy
* command queue
* telemetry service
