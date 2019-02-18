package metla

// Общий интерфейс объекта результирующих данных
type token interface {
	Data() ([]byte, error) // Получение значения
	IsExecutable() bool
}

// Интерфейс контейнера значения
type value interface {
	token
	Val() interface{}
}
