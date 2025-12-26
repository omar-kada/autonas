package mapper

// Mapper is a generic interface for mapping between types T and V.
// It provides methods to convert single instances and slices of types.
type Mapper[T any, V any] interface {
	// Map converts a single source type T to a target type V.
	Map(T) V
}
