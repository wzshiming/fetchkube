package main

import (
	"log"
	"os"

	"github.com/wzshiming/fetchkube"
)

func main() {
	conf, err := fetchkube.StartingConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx := os.Getenv("CONTEXT")
	if ctx != "" {
		conf.CurrentContext = ctx
	}

	ctxs, err := fetchkube.Contexts(conf)
	log.Printf("use context %q in %q", conf.CurrentContext, ctxs)

	conf, err = fetchkube.FetchOnlyDefault(conf)
	if err != nil {
		log.Fatal(err)
		return
	}

	conf, err = fetchkube.ResolveLocalPaths(conf)
	if err != nil {
		log.Fatal(err)
		return
	}
	conf, err = fetchkube.InlineData(conf)
	if err != nil {
		log.Fatal(err)
		return
	}
	f, err := fetchkube.Encode(conf)
	if err != nil {
		log.Fatal(err)
		return
	}
	os.Stdout.Write(f)
}
