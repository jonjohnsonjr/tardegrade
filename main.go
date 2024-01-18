package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"os"
	"path"
)

func main() {
	if err := mainE(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func mainE(args []string) error {
	keep := make(map[string]struct{}, len(args))

	for _, arg := range args {
		keep[arg] = struct{}{}

		for dir := path.Dir(arg); dir != "." && dir != "/"; dir = path.Dir(dir) {
			keep[dir] = struct{}{}
		}
	}

	br := bufio.NewReader(os.Stdin)
	zr, err := gzip.NewReader(br)
	if err != nil {
		return err
	}
	zr.Multistream(false)
	tr := tar.NewReader(zr)
	zw := gzip.NewWriter(os.Stdout)
	tw := tar.NewWriter(zw)

	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			if _, err := br.Peek(512); errors.Is(err, io.EOF) {
				// We are actually done.
				break
			}
			if err := zr.Reset(br); err != nil {
				return err
			}

			zr.Multistream(false)
			if err := zw.Close(); err != nil {
				return err
			}
			zw.Reset(os.Stdout)
			tr = tar.NewReader(zr)
			continue
		} else if err != nil {
			return err
		}

		if _, ok := keep[hdr.Name]; ok {
			tw.WriteHeader(hdr)
			if _, err := io.Copy(tw, tr); err != nil {
				return err
			}
			if err := tw.Flush(); err != nil {
				return err
			}
		}
	}

	return errors.Join(tw.Close(), zw.Close())
}
