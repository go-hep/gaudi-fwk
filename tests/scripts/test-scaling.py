#!/usr/bin/env python
import os
import tempfile
if not tempfile.tempdir:
    tempfile.tempdir = "/tmp"
import textwrap
import subprocess

def _save_file(fname, src, show=False):
    if os.path.exists(fname):
        os.remove(fname)
    with open(fname, "w") as f:
        if isinstance(src, (list,tuple)):
            print >> f, "\n".join(src)
        else:
            print >> f, src


_d = os.path.dirname
TESTS_DIR = os.path.dirname(os.path.abspath(__file__))
GAUDIROOT = _d(_d(_d(_d(TESTS_DIR))))
GAUDI_APP = os.path.join(GAUDIROOT, 'cmd', 'go-gaudi','app.py')

print "TESTS_DIR:",TESTS_DIR
print "GAUDI_APP:",GAUDI_APP

JOBO_DIR = os.path.join(GAUDIROOT, 'pkg', 'gaudi', 'tests', 'joboptions')

JOBO_NAME = os.path.join(JOBO_DIR, 'big-jobo.py')
cmd = [ "/usr/bin/time",
        #"-o", "gaudi.profile.out",
        GAUDI_APP, JOBO_NAME ]

nworkers = range(0,5000,100)
nprocs = range(0,48)

nworkers = [0, 10, 20, 30, 50, 100, 200, 500, 1000, 2000, 5000]
nworkers = [5000]
nprocs   = [ 0,  2,  4,  6,  8,
            10, 12, 16, 24,
            26, 28, 30, 40, 48]

for iproc in nprocs:
    for iworker in nworkers:
        jobo = textwrap.dedent(
            """\
            app.props.NbrProcs = %i
            app.svcs['evt-proc'].NbrWorkers = %i
            app.props.EvtMax = 50000
            """ % (iproc, iworker,)
            )
        workdir = tempfile.mkdtemp(prefix='ng-gaudi-go-')
        
        os.chdir(workdir)
        jobo_fname = "./test.py"
        _save_file(jobo_fname, jobo)

        print "::: nworkers: %04i nprocs: %04i" % (iworker, iproc),
        print "    WORK_DIR:",workdir
        icmd = cmd + [jobo_fname]
        #print icmd
        subprocess.check_call(icmd,
                              stderr=open("gaudi.profile.out", "w"),
                              #stderr=subprocess.PIPE,
                              stdout=open('/dev/null','w'))
        
        subprocess.check_call(["cat", "gaudi.profile.out"])
        
