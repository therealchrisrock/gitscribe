# GitScribe Meeting Transcription & Action Item System

## Overview

GitScribe is a comprehensive meeting transcription and action item extraction system that integrates with popular meeting platforms (Zoom, Google Meet, Microsoft Teams) and ticketing systems (GitHub Issues, Jira, Linear). The system employs a Domain-Driven Design (DDD) architecture with Go backend and WXT Chrome extension frontend.

## Key Features

- **Multi-Platform Meeting Support**: Zoom, Google Meet, Microsoft Teams, and generic meetings
- **Bot-Based Recording**: Automated bot joins meetings to record conversations
- **AI-Powered Transcription**: Integration with Assembly AI and other transcription services
- **Action Item Extraction**: Intelligent extraction of actionable items from conversations
- **Ticket Integration**: Automatic creation and linking with existing tickets across multiple platforms
- **Extensible Architecture**: Plugin-based system for adding new meeting platforms and ticketing systems

## Architecture Principles

### Domain-Driven Design (DDD)
- **Bounded Contexts**: Clear separation between Meeting, Transcription, Action Item, and Ticketing domains
- **Rich Domain Models**: Business logic encapsulated within domain entities
- **Repository Pattern**: Abstract data access through interfaces
- **Command/Query Separation**: Clear distinction between read and write operations

### Design Patterns
- **Strategy Pattern**: For transcription and ticketing providers
- **Adapter Pattern**: For external service integrations
- **Factory Pattern**: For creating domain entities
- **Observer Pattern**: For event-driven processing

### Key Benefits
- **Loose Coupling**: Easy to add new providers without changing core logic
- **High Cohesion**: Each module has a single, well-defined responsibility
- **Testability**: All dependencies injected through interfaces
- **Extensibility**: Configuration-driven provider selection
- **Maintainability**: Clear separation of concerns across layers

## System Components

### Backend (Go)
- **API Gateway**: HTTP routing, authentication, middleware
- **Domain Modules**: Meeting, Transcription, ActionItem, Ticketing, Integration
- **Shared Infrastructure**: Database, file storage, message queue, logging

### Frontend (Chrome Extension - WXT)
- **Content Scripts**: Meeting detection and interaction
- **Popup Interface**: User controls and status display
- **Background Service**: API communication and event handling

### External Integrations
- **Meeting Platforms**: Zoom API, Google Meet API, Microsoft Teams API
- **Transcription Services**: Assembly AI, Google Speech-to-Text, Zoom transcription
- **Ticketing Systems**: GitHub Issues, Jira, Linear, Asana
- **AI Services**: OpenAI/Claude for action item extraction

## Module Structure

Each domain module follows the DDD layered architecture:

```
modules/{domain}/
├── domain/
│   ├── entities/          # Domain entities and value objects
│   ├── repositories/      # Repository interfaces
│   └── services/          # Domain services
├── application/
│   ├── commands/          # Write operations
│   ├── queries/           # Read operations
│   └── services/          # Application services and orchestration
├── infrastructure/
│   ├── repositories/      # Repository implementations
│   ├── providers/         # External service providers
│   └── services/          # Infrastructure services
└── interfaces/
    └── http/              # HTTP handlers, DTOs, routes
```

## Data Flow

1. **Meeting Detection**: Chrome extension detects meeting URL and type
2. **Bot Deployment**: System deploys appropriate bot to join meeting
3. **Recording**: Bot records audio/video during meeting
4. **Transcription**: Audio processed through transcription service
5. **Action Extraction**: AI analyzes transcript for action items
6. **Ticket Matching**: System searches for existing related tickets
7. **Ticket Creation**: New tickets created or existing ones updated
8. **Notification**: User notified of completed processing

## Next Steps

- [Domain Model Documentation](./domain-model.md)
- [API Documentation](./api-documentation.md)
- [Implementation Guide](./implementation-guide.md)
- [Deployment Guide](./deployment-guide.md) 