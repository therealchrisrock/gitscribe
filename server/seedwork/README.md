# Seedwork - Shared DDD Building Blocks

This package provides shared abstractions and utilities for implementing Domain-Driven Design (DDD) patterns across all modules.

## Domain vs Repository Separation

A key architectural principle in this codebase is the clear separation between domain entities and repository models:

### BaseEntity (Domain Layer)
```go
// ✅ Pure domain object - NO database concerns
type BaseEntity struct {
    ID        string    `json:"id"`          // No gorm tags!
    CreatedAt time.Time `json:"created_at"`  // No gorm tags!
    UpdatedAt time.Time `json:"updated_at"`  // No gorm tags!
}
```

**Characteristics:**
- 🚫 **NO GORM tags** - keeps domain pure
- ✅ **Business logic focused** - only domain concerns
- ✅ **Persistence ignorant** - doesn't know about databases
- ✅ **Testable** - no external dependencies

### BaseRepositoryModel (Infrastructure Layer)
```go
// ✅ Infrastructure object - WITH database concerns
type BaseRepositoryModel struct {
    ID        string         `json:"id" gorm:"primaryKey;type:varchar(128)"`
    CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

**Characteristics:**
- ✅ **HAS GORM tags** - handles persistence concerns
- ✅ **Database aware** - knows about tables, indices, constraints
- ✅ **ORM integration** - works directly with GORM
- ✅ **Soft deletion** - includes DeletedAt for safe deletes

### Why This Separation Matters

**❌ Wrong - Domain with GORM tags:**
```go
// Violates DDD - domain should not know about persistence
type Meeting struct {
    BaseEntity                    // This used to have gorm tags - BAD!
    Title string `gorm:"not null"` // Domain tied to database - BAD!
}
```

**✅ Right - Clean separation:**
```go
// Domain entity - pure business logic
type Meeting struct {
    entities.BaseEntity  // Clean, no gorm tags
    Title string        // Pure domain field
}

// Repository model - handles persistence  
type Meeting struct {
    domain.BaseRepositoryModel           // Has gorm tags
    Title string `gorm:"not null"`      // Database constraints
}
```

**Benefits of this approach:**
1. **Domain Independence**: Business logic isn't tied to database structure
2. **Technology Flexibility**: Can switch ORMs without changing domain
3. **Clean Testing**: Domain tests don't need database setup
4. **Clear Responsibility**: Each layer has distinct concerns
5. **Maintainability**: Database changes don't affect business rules

## Data Mappers

The `DomainMapper` abstraction provides standardized bidirectional conversion between domain entities and repository models.

### Problem Solved

Without data mappers, each application service would need manual, repetitive mapping code:

```go
// ❌ Manual mapping - error prone and repetitive
repoMeeting := &repositories.Meeting{
    ID:            meeting.GetID(),
    UserID:        meeting.UserID,
    Title:         meeting.Title,
    Type:          string(meeting.Type),
    Status:        string(meeting.Status),
    // ... many more fields
    BotJoinURL:    nil, // Forgot to handle this!
    RecordingPath: nil, // And this!
}
if meeting.BotJoinURL != "" {
    repoMeeting.BotJoinURL = &meeting.BotJoinURL
}
// etc...
```

### Solution

The `DomainMapper` interface provides clean, reusable mapping with type-safe constraints:

```go
// ✅ Clean mapper usage with type constraints
type MeetingMapper struct {
    domain.BaseDomainMapper
}

func (m *MeetingMapper) ToRepository(meeting *entities.Meeting) repositories.Meeting {
    repo := repositories.Meeting{
        UserID:        meeting.UserID,
        Title:         meeting.Title,
        Type:          string(meeting.Type),
        Status:        string(meeting.Status),
        StartTime:     meeting.StartTime,
        EndTime:       meeting.EndTime,
        MeetingURL:    meeting.MeetingURL,
        BotJoinURL:    m.StringToPointer(meeting.BotJoinURL),
        RecordingPath: m.StringToPointer(meeting.RecordingPath),
    }
    
    // Set base repository model fields
    repo.SetRepositoryID(meeting.GetID())
    repo.SetRepositoryCreatedAt(meeting.GetCreatedAt())
    repo.UpdatedAt = meeting.GetUpdatedAt()
    
    return repo
}
```

### Repository Model Interface

All repository models must implement the `RepositoryModel` interface:

```go
type RepositoryModel interface {
    GetID() string
    SetID(id string)
    TableName() string
}
```

Repository models should embed `BaseRepositoryModel` and implement `TableName()`:

```go
type Meeting struct {
    domain.BaseRepositoryModel  // Provides ID, CreatedAt, UpdatedAt with GORM tags
    UserID        string        `json:"user_id" gorm:"not null"`
    Title         string        `json:"title" gorm:"not null"`
    // ... other fields with appropriate gorm tags
}

// TableName returns the database table name for meetings
func (Meeting) TableName() string {
    return "meetings"
}
```

### Interface

```go
type DomainMapper[D Entity, R RepositoryModel] interface {
    ToRepository(domain D) R
    ToDomain(repo R) D
    ToRepositoryList(domains []D) []R
    ToDomainList(repos []R) []D
}
```

### Base Utilities

The `BaseDomainMapper` provides common conversion utilities:

```go
type BaseDomainMapper struct{}

// String/pointer conversions
func (BaseDomainMapper) StringToPointer(s string) *string
func (BaseDomainMapper) PointerToString(s *string) string

// Time/pointer conversions  
func (BaseDomainMapper) TimeToPointer(t time.Time) *time.Time
func (BaseDomainMapper) PointerToTime(t *time.Time) time.Time
```

### Usage in Application Services

```go
type MeetingService struct {
    meetingRepo   repositories.MeetingRepository
    meetingMapper *MeetingMapper
}

func (s *MeetingService) CreateMeeting(ctx context.Context, cmd commands.CreateMeetingCommand) (*entities.Meeting, error) {
    // Domain logic
    meeting, err := entities.CreateMeeting(cmd.UserID, cmd.Title, cmd.Type, cmd.MeetingURL)
    if err != nil {
        return nil, fmt.Errorf("domain validation failed: %w", err)
    }

    // Clean conversion with type safety
    repoMeeting := s.meetingMapper.ToRepository(meeting)
    
    // Persistence
    if err := s.meetingRepo.SaveMeeting(ctx, &repoMeeting); err != nil {
        return nil, fmt.Errorf("failed to persist meeting: %w", err)
    }

    return meeting, nil
}
```

### Benefits

1. **DRY Principle**: Eliminates duplicate mapping code
2. **Type Safety**: `RepositoryModel` constraint prevents mapping errors
3. **Consistency**: Standardized patterns across all modules
4. **Maintainability**: Single place to update mapping logic
5. **Testability**: Isolated mapper testing
6. **Null Safety**: Built-in utilities for pointer conversions
7. **Interface Compliance**: Repository models must implement common contract
8. **ORM Integration**: `TableName()` method supports ORM table mapping

### Type Safety Improvements

The `RepositoryModel` interface provides several advantages over `any`:

**Before (loose typing)**:
```go
type DomainMapper[D Entity, R any] interface {
    ToRepository(domain D) *R  // Could be any type
    ToDomain(repo *R) D
}
```

**After (constrained typing)**:
```go
type DomainMapper[D Entity, R RepositoryModel] interface {
    ToRepository(domain D) R   // Must implement RepositoryModel
    ToDomain(repo R) D
}
```

This ensures:
- ✅ All repository models have consistent ID handling
- ✅ All repository models define their table names
- ✅ Compile-time verification of repository model compliance
- ✅ Better IntelliSense/IDE support
- ✅ Prevents accidental use of non-repository types

### Validation Support

For domain entities requiring validation:

```go
validator := NewMeetingValidator()
validatedMapper := domain.NewValidatedMapper(meetingMapper, validator)

// This will validate before conversion
repoMeeting, err := validatedMapper.ToRepository(meeting)
if err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

### Guidelines

1. **Domain Entities**: 
   - ❌ **NO GORM tags** - keep domain pure
   - ✅ **Embed `BaseEntity`** for common domain fields
   - ✅ **Focus on business logic** only

2. **Repository Models**:
   - ✅ **USE GORM tags** - handle persistence concerns
   - ✅ **Embed `BaseRepositoryModel`** for common persistence fields
   - ✅ **Implement `TableName()`** to define database table names

3. **Mapping**:
   - ✅ **Embed `BaseDomainMapper`** in your concrete mappers
   - ✅ **Use helper methods** for pointer conversions
   - ✅ **Keep mapping logic simple** - complex transformations belong in domain
   - ✅ **Test bidirectional conversion** to ensure data integrity

4. **Architecture**:
   - ✅ **Maintain clear layer separation** - domain vs infrastructure
   - ✅ **Handle all fields** - don't leave mappings incomplete
   - ✅ **Implement `RepositoryModel`** for type safety

### Example Implementation

See `server/modules/meeting/application/services/meeting_mapper.go` for a complete example of implementing domain mappers for Meeting, Participant, and BotSession entities. 

### Migration Guide

When updating existing repository models:

1. **Clean Domain Entities**:
   ```go
   type Meeting struct {
       entities.BaseEntity  // No gorm tags
       Title string        // No gorm tags
   }
   ```

2. **Add GORM tags to Repository Models**:
   ```go
   type Meeting struct {
       domain.BaseRepositoryModel           // Has gorm tags
       Title string `gorm:"not null"`      // Add database constraints
   }
   ```

3. **Implement TableName method**:
   ```go
   func (Meeting) TableName() string {
       return "meetings"
   }
   ```

4. **Update mapper signatures**:
   ```go
   // Before
   func (m *Mapper) ToRepository(entity *Entity) *RepoModel
   
   // After  
   func (m *Mapper) ToRepository(entity *Entity) RepoModel
   ```

5. **Update mapper implementations**:
   ```go
   func (m *Mapper) ToRepository(entity *Entity) RepoModel {
       repo := RepoModel{/* fields */}
       repo.SetID(entity.GetID())
       repo.CreatedAt = entity.GetCreatedAt()
       repo.UpdatedAt = entity.GetUpdatedAt()
       return repo
   }
   ``` 