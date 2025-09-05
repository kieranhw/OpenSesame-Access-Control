# OpenSesame

A personal, non-commercial project developed for learning and exploration with Raspberry PI hardware.

OpenSesame is a local‑first access‑control system for privacy‑respecting smarthomes, designed to run offline by default.

<img src="https://github.com/user-attachments/assets/fb00dade-5de0-4e60-9489-595643aaf6df" alt="os-rpi" width="480" />

## What it does
- Manage users, credentials, and device permissions locally
- Orchestrate entry events (e.g., keypad/RFID → relay lock)
- Provide real‑time status and audit trails
- Integrate with external clients (control systems, mobile apps) without compromising local‑first principles

## Architecture
- Edge devices (keypads, RFID, locks) connect to a central hub over the local network
- The hub stores configuration and state locally (SQLite) and exposes well‑defined interfaces for management and clients
- Automatic device discovery and registration via mDNS
- Optional external clients demonstrate integrations

## Modules

### os-device (Raspberry PI Pico)
- CircuitPython firmware for edge devices: keypads, RFID readers, and relay lock controllers
- Auto‑discovery and self‑registration with the hub via mDNS

### os-hub (Rasberry PI)
- Central hub written in Go that orchestrates access logic and device coordination
- Local SQLite storage for configuration and state
- Real‑time status updates via long‑polling; clean interfaces for external clients

### os-management
- Web‑based management UI (React/Next.js)
- Install and configure os‑devices, manage user credentials, define access → entry pairings, and review activity

### os-mobile
- React Native client showcasing external integrations (e.g., control systems, mobile access apps)
- Serves as a reference for third‑party client implementations

A fully descriptive write up coming soon!
