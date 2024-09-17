package main

import "fmt"

type FileCompressionMethod interface {
	Compress(file string) string
}

type ZipCompression struct{}

func (z *ZipCompression) Compress(file string) string {
	return fmt.Sprintf("Compressing %s using ZIP format", file)
}

type TarCompression struct{}

func (t *TarCompression) Compress(file string) string {
	return fmt.Sprintf("Compressing %s using TAR format", file)
}

type RarCompression struct{}

func (r *RarCompression) Compress(file string) string {
	return fmt.Sprintf("Compressing %s using RAR format", file)
}

type Compressor struct {
	fileCompressionMethod FileCompressionMethod
}

func (c *Compressor) SetCompressor(fileCompressionMethod FileCompressionMethod) {
	c.fileCompressionMethod = fileCompressionMethod
}

func (c *Compressor) CompressFile(file string) string {
	return c.fileCompressionMethod.Compress(file)
}

func main() {
	file := "example.txt"

	compressor := &Compressor{}

	zipCompression := &ZipCompression{}
	compressor.SetCompressor(zipCompression)
	fmt.Println(compressor.CompressFile(file))

	tarCompression := &TarCompression{}
	compressor.SetCompressor(tarCompression)
	fmt.Println(compressor.CompressFile(file))

	rarCompression := &RarCompression{}
	compressor.SetCompressor(rarCompression)
	fmt.Println(compressor.CompressFile(file))
}
