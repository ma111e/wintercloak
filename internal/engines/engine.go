package engines

type Engine interface {
	Run(RunOptions) ([]byte, error)
	Apply(string) ([]byte, error)
	WritePatchedFile(patched []byte) error
}
