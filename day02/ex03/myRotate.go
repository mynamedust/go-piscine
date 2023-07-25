package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	dir := flag.String("a", "./", "Write valid path for archives.")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Error. There must be at least one argument.")
		return
	}
	createArchives(*dir, args)
}

func createArchives(dir string, files []string) {
	wg := new(sync.WaitGroup)
	for _, file := range files {
		wg.Add(1)
		if !fileValid(file) {
			continue
		}
		go archiveFile(dir, file, wg)
	}
	wg.Wait()
}

func fileValid(file string) bool {
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		fmt.Printf("Error. File %s not exist.\n", file)
		return false
	}
	if ext := filepath.Ext(file); ext != ".log" {
		fmt.Println("Error. Wrong file type.")
		return false
	}
	return true
}

func archiveFile(dir, fileName string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error. File opening failed.")
		return
	}
	defer file.Close()
	stat, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("Error. File stats reading failed.")
		return
	}
	archiveName := createArchiveName(dir, fileName, stat.ModTime())
	archive, err := os.Create(archiveName)
	defer archive.Close()
	if err != nil {
		fmt.Println("Error. Archive creating failed.")
		return
	}
	gzipWriter := gzip.NewWriter(archive)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()
	if err := addFileToArchive(tarWriter, stat, file, fileName); err != nil {
		fmt.Println("Error. File adding failed.")
		return
	}
	fmt.Println(fileName, "added to archive.")
}

func createArchiveName(dir, fileName string, mtime time.Time) string {
	var name strings.Builder
	name.WriteString(dir)
	name.WriteString(fileName[:len(fileName)-4])
	name.WriteString("_")
	name.WriteString(strconv.FormatInt(mtime.Unix(), 10))
	name.WriteString(".tar.gz")
	return name.String()
}

func addFileToArchive(tarWriter *tar.Writer, fileStat fs.FileInfo, file *os.File, fileName string) error {
	header := &tar.Header{
		Name:    filepath.Base(fileName),
		Size:    fileStat.Size(),
		Mode:    int64(fileStat.Mode().Perm()),
		ModTime: fileStat.ModTime(),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}
	_, err := io.Copy(tarWriter, file)
	if err != nil {
		return err
	}
	return nil
}
