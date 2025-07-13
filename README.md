# Open-Sesame
A conceptual local-first access-control system.

Full write-up at https://kieranhardwick.com/

# Internal Structure

- config/                   App configuration (ports, TLS, env)
- model/                    GORM models mapped to SQLite tables
- service/                  Business logic (SetupService, AccessService)
- handlers/                 HTTP handlers and DTOs (management, access, pairing)
- middleware/               HTTP middleware (logging, JSON validation, auth)
- httpserver/               Server bootstrap and middleware wiring

