//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"

	"github.com/fopina/wireguard-mitm-gui/assets"
	"github.com/shurcooL/vfsgen"
)

func main() {
	// need to change directory so assets.Assets "static" works for both dev and vfsgen...
	// gap in vfsgen or misuse? https://github.com/shurcooL/vfsgen/issues/77
	err := os.Chdir("..")
	if err != nil {
		log.Fatalln(err)
	}
	err = vfsgen.Generate(assets.Assets, vfsgen.Options{
		PackageName:  "assets",
		BuildTags:    "!dev",
		VariableName: "Assets",
		Filename:     "assets/assets_vfsdata.go",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
