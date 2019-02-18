package metla

// Общий интерфейс объекта результирующих данных
type token interface {
	Data() ([]byte, error) // Получение значения
}

// Интерфейс контейнера значения
type value interface {
	Val() interface{}
	Type() valueType
	Data() ([]byte, error)
}
