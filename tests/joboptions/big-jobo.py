app.props.EvtMax = 10000
app.props.OutputLevel = 1
#app.props.NbrProcs = 1

app.svcs += Svc(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/kernel/evtproc:evtproc",
    "evt-proc",
    OutputLevel=Lvl.INFO,
    NbrWorkers=5000
    )

app.algs += Alg(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg1:alg1",
    "alg1",
    OutputLevel=Lvl.ERROR
    )
app.algs += Alg(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg1:alg2",
    "alg2",
    OutputLevel=Lvl.ERROR
    )
app.algs += Alg(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg2:alg1",
    "alg_one",
    OutputLevel=Lvl.ERROR
    )

app.svcs += Svc(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg1:svc1",
    name="svc1",
    OutputLevel=Lvl.ERROR
    )
app.svcs += Svc(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg2:svc2",
    "svc2"
    )

app.svcs += Svc(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/kernel/datastore:datastoresvc",
    "evt-store"
    )
app.svcs += Svc(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/kernel/datastore:datastoresvc",
    "det-store"
    )

app.algs += Alg(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg2:alg_adder",
    "adder_1",
    OutputLevel=Lvl.ERROR,
    Val=0.
    )
app.algs += Alg(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg2:alg_dumper",
    "dumper-1",
    ExpectedValue=1
    )

app.algs += Alg(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg2:alg_adder",
    "adder_2",
    OutputLevel=Lvl.ERROR,
    Val=3.
    )
app.algs += Alg(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg2:alg_dumper",
    "dumper-2",
    ExpectedValue=2
    )

app.algs += Alg(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg2:alg_adder",
    "adder_3"
    )
app.algs += Alg(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg2:alg_dumper",
    "dumper-3",
    ExpectedValue=3
    )

app.algs += Alg(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg2:alg_dumper",
    "dumper",
    NbrJets = "njets",
    ExpectedValue=3,
    OutputLevel=Lvl.ERROR
    )

if 1:
    for i in xrange(500):
        app.algs += Alg(
            "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg2:alg_adder",
            "addr--%04i" % i,
            SimpleCounter="my_counter"
            )
        app.algs += Alg(
            "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg2:alg_dumper",
            "dump--%04i" % i,
            SimpleCounter="my_counter",
            ExpectedValue=i+1
            )
    
app.toolsvc += Tool(
    "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/tests/pkg1:tool1",
    name="tool1"
    )

