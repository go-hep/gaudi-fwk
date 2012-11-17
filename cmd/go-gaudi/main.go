package main

import (
	"flag"
	"fmt"
	"os"
	//"reflect"

	"github.com/sbinet/go-gaudi/pkg/gaudi/kernel"
)

var (
	bin  = os.Getenv("GOBIN")
	arch = map[string]string{
		"amd64": "6",
		"386":   "8",
		"arm":   "5",
	}[os.Getenv("GOARCH")]
)

type jobOptions struct {
	pkgs []string
	defs []string
	code []interface{}
	exec string
}

var jobOptName *string = flag.String(
	"jobo", 
	"jobOptions.igo",
	"path to a jobOptions file")

func handle_icomponent(c kernel.IComponent) {
	fmt.Printf(":: handle_icomponent(%s)...\n", c.CompName())
}

func main() {
	flag.Parse()

	fmt.Print("::: gaudi\n")
	fmt.Printf("::: getting options from [%s]...\n", *jobOptName)

	app := kernel.NewAppMgr()
	iapp := app.(kernel.IComponent)
	fmt.Printf(" -> created [%s/%s]\n", 
		iapp.CompType(), 
		iapp.CompName())

	// {
	// 	t := reflect.TypeOf(iapp)
	// 	fmt.Printf("type of [%s]\n", t)
	// 	newt := reflect.New(t)
	// 	fmt.Printf("type of    t: [%s] pkg: [%s] name: [%s]\n", 
	// 		reflect.TypeOf(t),
	// 		t.PkgPath(),
	// 		t.Name())
	// 	fmt.Printf("type of newt: [%s]\n", reflect.TypeOf(newt))
	// }
	handle_icomponent(app)
	fmt.Printf("%s\n", app)

	println("::: configure...")
	sc := app.Configure()
	println("::: configure... [", sc, "]")
	println("::: run...")
	sc = app.Run()
	println("::: run... [", sc, "]")
	fmt.Print("::: bye.\n")
}
