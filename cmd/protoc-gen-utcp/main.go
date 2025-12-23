package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/dimelords/protoc-gen-utcp/internal/generator"
)

const version = "0.1.0"

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("protoc-gen-utcp %s\n", version)
		return
	}

	var flags flag.FlagSet
	baseURL := flags.String("base_url", "", "Base URL for the API endpoints")
	providerType := flags.String("provider_type", "http", "Provider type (http, grpc, connectrpc)")
	authType := flags.String("auth_type", "bearer", "Authentication type (bearer, api_key, oauth2)")

	opts := protogen.Options{
		ParamFunc: flags.Set,
	}

	opts.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		config := &generator.Config{
			BaseURL:      *baseURL,
			ProviderType: *providerType,
			AuthType:     *authType,
		}

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			if err := generator.GenerateFile(gen, f, config); err != nil {
				return err
			}
		}

		return nil
	})
}
