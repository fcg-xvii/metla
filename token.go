package metla

//type tokenInit func(parser)

type token interface {
	Data() ([]byte, error)
}
