# BambuLab Print Farm Dashboard

## UI Design Specification

UI inspiration:

* Proxmox
* Portainer
* Bambu Studio

Goal:

Modern industrial dashboard for managing printer farms.

---

# DESIGN PRINCIPLES

Style:

dark theme first

colors:

background: #0b0f14
panel: #111827
accent: #3b82f6

Typography:

Inter font

Layout:

grid based

---

# MAIN NAVIGATION

Sidebar layout:

Dashboard
Printers
Jobs
History
Alerts
Analytics
Settings

---

# DASHBOARD PAGE

Overview widgets:

Total printers
Active prints
Failed jobs
Queue size

Farm grid:

Each printer card shows:

printer name
model
status
progress bar
temperature
camera thumbnail

Example layout:

Printer Tile

Printer Name
Status LED
Progress bar
Temp nozzle
Temp bed
Camera preview

Click tile → open printer detail

---

# PRINTER DETAIL PAGE

Sections:

Printer Info
Live Camera
Print Progress
Temperature Chart
Job History

Controls:

Pause
Resume
Cancel

---

# PRINT JOB PAGE

Queue manager UI.

Columns:

job id
printer
status
progress
duration

Controls:

cancel
reassign printer

---

# ALERTS PAGE

Alert feed.

Color levels:

green = info
yellow = warning
red = critical

---

# ANALYTICS PAGE

Charts:

printer uptime
job success rate
material usage

---

# COMPONENT LIBRARY

Use:

Shadcn UI
Tailwind

Components:

Card
DataTable
StatusBadge
ProgressBar
CameraPanel
MetricWidget

---

# UX FEATURES

Keyboard shortcuts

R = refresh

Quick search

Cmd + K

Drag printers to reorder farm layout.

---

# MOBILE SUPPORT

Responsive grid.

Printer tiles collapse to list view.

---

# PERFORMANCE RULES

Use:

virtual lists
websocket updates
lazy loading

---

# FINAL INSTRUCTION FOR ANTIGRAVITY

When implementing UI:

Generate components first.

Then pages.

Then integrate with backend API.

Stop after each page generation.
