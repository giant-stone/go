package ghttpcache

type Coder interface {
	Compress(data []byte) (rs []byte, err error)
	Decompress(data []byte) (rs []byte, err error)
}
