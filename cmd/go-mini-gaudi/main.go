package main

import (
	"encoding/json"
	"flag"
	//"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
	"time"
)

// thread safe output
var fmt = log.New(os.Stdout, "", log.Lmicroseconds)

var g_do_seq *bool = flag.Bool("seq", false, "enable sequential alg seq")
var g_nprocs *int = flag.Int("nprocs", 0, "number of threads")
var g_evtmax *int = flag.Int("evtmax", 10, "number of events")
var g_jobo *string = flag.String("jobo", "cfg.json", "input joboption file")

// type datastore struct {
// 	store map[string]chan interface{}
// }

type datastore map[string]interface{}

func new_datastore() datastore {
	o := make(datastore)
	return o
}

type depgraph map[string]chan int

func new_depgraph() depgraph {
	o := make(depgraph)
	return o
}

type alg struct {
	name  string
	sleep float64
	depg  depgraph
	deps  []string
	store datastore
}

func (a *alg) Initialize() error {
	fmt.Printf(" --> deps: %v\n", a.deps)
	return nil
}

func (a *alg) Execute(ctx int) error {
	for _, dep := range a.deps {
		fmt.Printf(":: [%s] waiting for [%s] (evt: %d)...\n", a.name, dep, ctx)
		v := <-a.depg[dep]
		a.depg[dep] <- v
	}

	// simulate work...
	var f func() int64
	f = func() int64 {
		vv := time.Now()
		for {
			select {
			case delta := <-time.After(time.Duration(a.sleep) * time.Second):
				return int64(delta.Sub(vv))
			}
		}
		return int64(time.Since(vv)) + int64(math.Sqrt(0))
	}
	// f = func() int64 {
	// 	iv := int64(0)
	// 	vv := 0.
	// 	for {
	// 		select {
	// 		case iv = <-time.After(int64(a.sleep * 1e9)):
	// 			return iv
	// 		default:
	// 			vv += math.Sqrt(math.Pi)
	// 		}
	// 	}
	// 	return iv
	// }
	v := f()
	// 

	fmt.Printf(":: [%s] done (evt: %d) (%v)\n", a.name, ctx, v)
	a.depg[a.name] <- 1

	return nil
}

func (a *alg) Finalize() error {
	return nil
}

func init() {
	flag.Parse()
	fmt.Printf("== evtmax:  %d\n", *g_evtmax)
	fmt.Printf("== alg-seq: %v\n", *g_do_seq)
	fmt.Printf("== n-procs: %d\n", *g_nprocs)
	fmt.Printf("== jobo:    %s\n", *g_jobo)
}

func main() {

	runtime.GOMAXPROCS(*g_nprocs)

	fmt.Printf("== hello ==\n")

	cfg, err := os.Open(*g_jobo)
	if err != nil {
		panic(err)
	}
	alg_names := []string{}
	alg_deps := make(map[string][]string)
	alg_sleeps := make(map[string]float64)
	dec := json.NewDecoder(cfg)
	{
		var v map[string]interface{}
		if err := dec.Decode(&v); err != nil {
			panic(err)
		}
		v = v["process"].(map[string]interface{})
		valgs := v["algs"].([]interface{})
		//fmt.Printf("cfg: %v\n", valgs)
		for k, _ := range valgs {
			vv := valgs[k].(map[string]interface{})
			alg_name := vv["@name"].(string)
			alg_deps[alg_name] = []string{}
			alg_names = append(alg_names, alg_name)
			if !*g_do_seq {
				vv_toget := vv["toGet"].([]interface{})
				for kk, _ := range vv_toget {
					vvv := vv_toget[kk].(map[string]interface{})
					label := vvv["label"].(string)
					alg_deps[alg_name] = append(alg_deps[alg_name], label)
				}
			}
			for _,nn := range alg_deps[alg_name] {
				if nn == alg_name {
					panic("detected cycle for alg ["+nn +"]")
				}
			}
			
			if vv_sleep,ok := vv["eventTimes"]; ok {
				alg_sleep := vv_sleep.(float64)
				alg_sleeps[alg_name] = alg_sleep 
			} else {
				alg_sleeps[alg_name] = 0.2
				//panic("no event-time data for ["+alg_name+"]")
			}
			//fmt.Printf("--%d--: %v\n", k, vi)
		}
		//fmt.Printf("::deps: %v\n", alg_deps)
	}
	algs := make([]alg, 0, len(alg_names))

	depg := new_depgraph()
	store := new_datastore()

	for _, n := range alg_names {
		algs = append(
			algs,
			alg{n, alg_sleeps[n], depg, alg_deps[n], store})

		// check if all deps exist
		for _, nn := range alg_deps[n] {
			found := false
			for _, nnn := range alg_names {
				if nnn == nn {
					found = true
					break
				}
			}
			if !found {
				panic("dependency ["+nn+"] for alg ["+n+"] NOT FOUND !")
			}
		}
		depg[n] = make(chan int, 1)
	}

	fmt.Printf("::: #-algs: %d\n", len(algs))

	for _, a := range algs {
		fmt.Printf("--> [init] alg[%s]...\n", a.name)
		if err := a.Initialize(); err != nil {
			panic(err)
		}
	}

	reinit_fct := func() {
		// re-init depg
		for k, _ := range depg {
			depg[k] = make(chan int, 1)
		}
		// re-init store
		for k, _ := range store {
			store[k] = nil
		}
	}

	if *g_do_seq {
		for ievt := 0; ievt < *g_evtmax; ievt++ {
			for i, _ := range algs {
				func(ievt, ialg int) {
					a := algs[ialg]
					fmt.Printf("--> [exec-%d] alg[%s] (%f)...\n", ievt, 
						a.name, a.sleep)
					if err := a.Execute(ievt); err != nil {
						panic(err)
					}
				}(ievt, i)
			}
			reinit_fct()
		}

	} else {
		for ievt := 0; ievt < *g_evtmax; ievt++ {
			var seq sync.WaitGroup
			seq.Add(len(algs))
			for i, _ := range algs {
				go func(ievt, ialg int) {
					a := algs[ialg]
					fmt.Printf("--> [exec-%d] alg[%s] (%f)...\n", ievt, 
						a.name, a.sleep)
					if err := a.Execute(ievt); err != nil {
						panic(err)
					}
					seq.Done()
					fmt.Printf("--> [exec-%d] alg[%s]... [done]\n",
						ievt, a.name)
				}(ievt, i)
			}
			seq.Wait()
			fmt.Printf("--> event [%d] [done]\n", ievt)
			reinit_fct()
		}
	}

	for _, a := range algs {
		fmt.Printf("--> [fini] alg[%s]...\n", a.name)
		if err := a.Finalize(); err != nil {
			panic(err)
		}
	}

	fmt.Printf("== bye.\n")
}
