package main

import (
	"github.com/airgap-solution/go-pkg/mux"
)

func main() {
	r := mux.NewRouter(mux.Config{})
	mux.HandleRoute(r, "/price", func(struct{}) (string, mux.Error) {})
}
