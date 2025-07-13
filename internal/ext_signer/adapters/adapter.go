package adapters

type ExternalSignerAdapter interface {
	StartSign(srcPath string) (string, error)
	WaitSignature(ID string, dstPath string) error
}
