package domain

import (
	"testing"
)

// MockEntity for testing
type MockEntity struct {
	BaseEntity
	Name        string
	Description string
	Active      bool
}

// MockRepository for testing - now implements RepositoryModel
type MockRepository struct {
	BaseRepositoryModel
	Name        string
	Description *string
	Active      *bool
}

// TableName returns the table name for testing
func (MockRepository) TableName() string {
	return "mock_entities"
}

// MockMapper implements DomainMapper for testing
type MockMapper struct {
	BaseDomainMapper
}

func (m *MockMapper) ToRepository(entity *MockEntity) MockRepository {
	repo := MockRepository{
		Name:        entity.Name,
		Description: m.StringToPointer(entity.Description),
		Active:      &entity.Active,
	}

	// Set repository model fields
	repo.SetID(entity.GetID())
	repo.CreatedAt = entity.GetCreatedAt()
	repo.UpdatedAt = entity.GetUpdatedAt()

	return repo
}

func (m *MockMapper) ToDomain(repo MockRepository) *MockEntity {
	entity := &MockEntity{
		Name:        repo.Name,
		Description: m.PointerToString(repo.Description),
		Active:      repo.Active != nil && *repo.Active,
	}
	entity.SetID(repo.GetID())
	entity.CreatedAt = repo.CreatedAt
	entity.UpdatedAt = repo.UpdatedAt
	return entity
}

func (m *MockMapper) ToRepositoryList(entities []*MockEntity) []MockRepository {
	result := make([]MockRepository, len(entities))
	for i, entity := range entities {
		result[i] = m.ToRepository(entity)
	}
	return result
}

func (m *MockMapper) ToDomainList(repos []MockRepository) []*MockEntity {
	result := make([]*MockEntity, len(repos))
	for i := range repos {
		result[i] = m.ToDomain(repos[i])
	}
	return result
}

func TestBaseDomainMapper_StringToPointer(t *testing.T) {
	mapper := BaseDomainMapper{}

	// Test non-empty string
	str := "test"
	ptr := mapper.StringToPointer(str)
	if ptr == nil {
		t.Error("Expected non-nil pointer for non-empty string")
	}
	if *ptr != str {
		t.Errorf("Expected %s, got %s", str, *ptr)
	}

	// Test empty string
	emptyStr := ""
	emptyPtr := mapper.StringToPointer(emptyStr)
	if emptyPtr != nil {
		t.Error("Expected nil pointer for empty string")
	}
}

func TestBaseDomainMapper_PointerToString(t *testing.T) {
	mapper := BaseDomainMapper{}

	// Test non-nil pointer
	str := "test"
	result := mapper.PointerToString(&str)
	if result != str {
		t.Errorf("Expected %s, got %s", str, result)
	}

	// Test nil pointer
	result = mapper.PointerToString(nil)
	if result != "" {
		t.Errorf("Expected empty string for nil pointer, got %s", result)
	}
}

func TestMockMapper_Bidirectional(t *testing.T) {
	mapper := &MockMapper{}

	// Create original entity
	original := &MockEntity{
		Name:        "Test Entity",
		Description: "Test Description",
		Active:      true,
	}
	original.SetID("test-id")

	// Convert to repository model and back
	repo := mapper.ToRepository(original)
	converted := mapper.ToDomain(repo)

	// Verify round-trip conversion
	if converted.GetID() != original.GetID() {
		t.Errorf("Expected ID %s, got %s", original.GetID(), converted.GetID())
	}
	if converted.Name != original.Name {
		t.Errorf("Expected Name %s, got %s", original.Name, converted.Name)
	}
	if converted.Description != original.Description {
		t.Errorf("Expected Description %s, got %s", original.Description, converted.Description)
	}
	if converted.Active != original.Active {
		t.Errorf("Expected Active %v, got %v", original.Active, converted.Active)
	}
}

func TestRepositoryModel_TableName(t *testing.T) {
	repo := MockRepository{}

	// Verify TableName method works
	expectedTableName := "mock_entities"
	actualTableName := repo.TableName()

	if actualTableName != expectedTableName {
		t.Errorf("Expected table name %s, got %s", expectedTableName, actualTableName)
	}
}

func TestRepositoryModel_IDMethods(t *testing.T) {
	repo := MockRepository{}

	// Test SetID and GetID
	testID := "test-repo-id"
	repo.SetID(testID)

	if repo.GetID() != testID {
		t.Errorf("Expected ID %s, got %s", testID, repo.GetID())
	}
}
