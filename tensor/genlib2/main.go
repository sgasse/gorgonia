package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"reflect"
)

const genmsg = "GENERATED FILE. DO NOT EDIT"

var (
	gopath, tensorPkgLoc, execLoc, storageLoc string
)

type Kinds struct {
	Kinds []reflect.Kind
}

func init() {
	gopath = os.Getenv("GOPATH")
	tensorPkgLoc = path.Join(gopath, "src/github.com/chewxy/gorgonia/tensor")
	execLoc = path.Join(gopath, "src/github.com/chewxy/gorgonia/tensor/internal/execution")
	storageLoc = path.Join(gopath, "src/github.com/chewxy/gorgonia/tensor/internal/storage")
}

func main() {
	// pipeline("test", "BLAH_1.go", Kinds{allKinds}, generateGenericVecVecArith)
	// pipeline("test", "BLAH_2.go", Kinds{allKinds}, generateGenericMixedArith)
	// pipeline("test", "BLAH_3.go", Kinds{allKinds}, generateEArith)
	// pipeline("test", "BLAH_4.go", Kinds{allKinds}, generateGenericMap)
	// pipeline("test", "BLAH_5.go", Kinds{allKinds}, generateMap)
	// pipeline("test", "BLAH_6.go", Kinds{allKinds}, generateGenericVecVecCmp)
	// pipeline("test", "BLAH_7.go", Kinds{allKinds}, generateGenericMixedCmp)
	// pipeline("test", "BLAH_8.go", Kinds{allKinds}, generateMinMax)
	// pipeline("test", "BLAH_9.go", Kinds{allKinds}, generateStdEngArith)
	// pipeline("test", "BLAH_10.go", Kinds{allKinds}, generateDenseArith)

	// storage
	pipeline(storageLoc, "getset.go", Kinds{allKinds}, generateHeaderGetSet)
	pipeline(tensorPkgLoc, "array_getset.go", Kinds{allKinds}, generateArrayMethods)
	pipeline(tensorPkgLoc, "dense_getset.go", Kinds{allKinds}, generateDenseGetSet)

	// execution
	pipeline(execLoc, "generic_arith_vv.go", Kinds{allKinds}, generateGenericVecVecArith)
	pipeline(execLoc, "generic_arith_mixed.go", Kinds{allKinds}, generateGenericMixedArith)
	pipeline(execLoc, "generic_map.go", Kinds{allKinds}, generateGenericMap)
	pipeline(execLoc, "eng_arith.go", Kinds{allKinds}, generateEArith)
	pipeline(execLoc, "eng_map.go", Kinds{allKinds}, generateMap)

	// level 2 aggregation
	pipeline(tensorPkgLoc, "defaultengine_arith.go", Kinds{allKinds}, generateStdEngArith)

	// level 3 eaggregation
	pipeline(tensorPkgLoc, "dense_arith.go", Kinds{allKinds}, generateDenseArith)
}

func pipeline(pkg, filename string, kinds Kinds, fn func(io.Writer, Kinds)) {
	fullpath := path.Join(pkg, filename)
	f, err := os.Create(fullpath)
	if err != nil {
		log.Printf("fullpath %q", fullpath)
		log.Fatal(err)
	}
	defer f.Close()
	writePkgName(f, pkg)
	fn(f, kinds)

	// gofmt and goimports this stuff
	cmd := exec.Command("goimports", "-w", fullpath)
	if err = cmd.Run(); err != nil {
		log.Fatalf("Go imports failed with %v for %q", err, fullpath)
	}

	cmd = exec.Command("sed", "-i", `s/github.com\/alecthomas\/assert/github.com\/stretchr\/testify\/assert/g`, fullpath)
	if err = cmd.Run(); err != nil {
		log.Fatalf("sed failed with %v for %q", err, fullpath)
	}

	cmd = exec.Command("gofmt", "-s", "-w", fullpath)
	if err = cmd.Run(); err != nil {
		log.Fatalf("Gofmt failed for %q", fullpath)
	}
}