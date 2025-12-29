package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	showVersion     = flag.Bool("version", false, "print the version and exit")
	appPackage      = flag.String("app_package", "project/pkg/app", "app package e.g. project/pkg/app")
	omitempty       = flag.Bool("omitempty", true, "omit if google.api is empty")
	omitemptyPrefix = flag.String("omitempty_prefix", "", "omit if google.api is empty")
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-http %v\n", release)
		return
	}

	opts := protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}

	opts.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f, *appPackage, *omitempty, *omitemptyPrefix)
		}
		return nil
	})
}
