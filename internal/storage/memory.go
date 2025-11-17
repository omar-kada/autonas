package storage

// MemoryStorage uses memory to store data (to be used mainly for testing)
type MemoryStorage struct{}

// NewMemoryStorage instanciates a new memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

// GetCurrentStacks returns currently stored stacks
func (m *MemoryStorage) GetCurrentStacks() []string {
	return []string{"test"}
}
