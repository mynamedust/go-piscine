package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type FlagsType struct {
	isFlags, dFlag, fFlag, slFlag bool
	extFlag, root                 string
}

func main() {
	flags := new(FlagsType)
	dFlag := flag.Bool("d", false, "Find the directories")
	fFlag := flag.Bool("f", false, "Find the files")
	slFlag := flag.Bool("sl", false, "Find the symlinks")
	extFlag := flag.String("ext", "", "Write the correct file extension")
	flag.Parse()
	flags.dFlag = *dFlag
	flags.fFlag = *fFlag
	flags.slFlag = *slFlag
	flags.extFlag = *extFlag
	if flags.extFlag != "" && !flags.fFlag {
		fmt.Println("Error. Flag -ext work only with -f flag.")
		return
	}
	if flags.dFlag || flags.fFlag || flags.slFlag {
		flags.isFlags = true
	}
	args := flag.Args()
	if len(args) > 1 || len(args) == 0 {
		fmt.Println("Error. Write correct number of arguments.")
		return
	}
	flags.root = args[0]
	searchAndPrint(*flags)
}

func searchAndPrint(flags FlagsType) {
	err := filepath.Walk(flags.root, flags.walkFunc)
	if err != nil {
		fmt.Println("Directory reading error:", err)
	}
}

func (flags FlagsType) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.Name() == flags.root || info.Name()+"/" == flags.root {
		return nil
	}
	if flagsCheck(flags, info) == false {
		return nil
	}
	infoPrint(flags, info)
	return nil
}

func flagsCheck(flags FlagsType, info os.FileInfo) bool {
	if flags.isFlags {
		if (info.IsDir() && !flags.dFlag) ||
			(!info.IsDir() && info.Mode()&os.ModeSymlink != 0 && !flags.slFlag) ||
			(!info.IsDir() && !flags.fFlag && info.Mode()&os.ModeSymlink == 0) {
			return false
		}
	}
	return true
}

func infoPrint(flags FlagsType, info os.FileInfo) {
	if info.Mode()&os.ModeSymlink != 0 {
		symPrint(info)
	} else if !info.IsDir() {
		filePrint(flags.extFlag, info)
	} else {
		fmt.Println(info.Name())
	}
}

func symPrint(info os.FileInfo) {
	link, err := os.Readlink(info.Name())
	if err != nil {
		fmt.Printf("Error reading symlink %s: %s\n", info.Name(), err)
		return
	}
	_, err = os.Stat(link)
	if err != nil {
		link = "[broken]"
	}
	fmt.Println(info.Name() + " -> " + link)
}

func filePrint(extFlag string, info os.FileInfo) {
	if extFlag != "" {
		ext := filepath.Ext(info.Name())
		if len(ext) >= 1 {
			ext = ext[1:]
		}
		if ext != extFlag {
			return
		}
	}
	fmt.Println(info.Name())
}
