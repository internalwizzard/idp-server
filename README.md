![Banner](./docs/images/banner.png)

# InternalWizzard - IDP Server

## Overview

The **IDP Server** is the core backend service of **InternalWizard**, an Internal Developer Platform designed to streamline application lifecycle management.
It acts as the central orchestrator, managing application definitions, authentication, configurations, integrations, and CI/CD pipelines across the ecosystem.

## Core Responsibilities

- Application Management: Handles application creation, configuration, secrets, and environment setup.
- Authentication & Authorization: Manages user identity and access control through OIDC (Keycloak or other providers).
- Integration Management: Connects applications with external services through the Integration Catalog.
- Pipeline Orchestration: Triggers and manages CI/CD pipelines using registered Git Providers.
- Event Dispatching: Publishes and listens to platform events for asynchronous operations and service coordination.

## Architecture

```shell
/idp-server
│
├── cmd/
│   └── idp-server/           # Application entrypoint
│
├── internal/
│   ├── appmanager/           # Application creation, configuration, and secrets
│   ├── auth/                 # Authentication and user session management
│   ├── pipeline/             # CI/CD pipeline integration layer
│   ├── integration/          # Integration with external services via Integration Catalog
│   ├── secrets/              # Secure credentials management (delegated to external Secret Stores)
│   └── shared/               # Common code: middlewares, DTOs, utils
│
└── pkg/                      # Reusable components or SDK clients
```

Each slice contains:

- `domain/` → Core business models
- `application/` → Use cases
- `adapters/` → External interfaces (DB, API, message bus)

## Technology Stack

- Language: Go
- Framework: Echo
- Communication: REST / gRPC
- Storage: PostgreSQL (primary), Redis (cache)

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](./LICENSE)
file for details.
