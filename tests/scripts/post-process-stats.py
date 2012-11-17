#!/usr/bin/env python2
"""

::: nworkers: 0010 nprocs: 0048     WORK_DIR: /tmp/ng-gaudi-go-WF1x2A
184.67user 7.69system 0:14.20elapsed 1354%CPU (0avgtext+0avgdata 0maxresident)k
0inputs+0outputs (0major+63157minor)pagefaults 0swaps

"""

import sys
import os
import re

header_pat = re.compile(r'::: nworkers: (?P<nworkers>\d*?) nprocs: (?P<nprocs>\d*?) .*?').match
cpu_pat = re.compile(r'(?P<usr>.*?)user (?P<sys>.*?)system (?P<real>.*?)elapsed (?P<cpu>.*?)%CPU.*?').match

fname = sys.argv[1]
print "::: massaging [%s]..." % (fname,)
f = open(fname, 'r')

# skip first 2 header-lines of the form:
# TESTS_DIR: /storage/binet/dev/go/dev/ng-go-gaudi/gaudi/tests/scripts
# GAUDI_APP: /storage/binet/dev/go/dev/ng-go-gaudi/gaudi/app/app.py

lines = f.readlines()[2:]
headers = lines[::3]
cpus = lines[1::3]

nitems = len(cpus)
assert len(cpus) == len(headers)
print "-->",nitems
nworkers = []
nprocs = []

#import numpy as np
ntuple = []


for i in xrange(nitems):
    h = headers[i]
    hh= header_pat(h).groupdict()
    for k in hh:
        hh[k] = int(hh[k])

    d = cpus[i]
    dd= cpu_pat(d).groupdict()
    real_minutes,real_seconds = dd['real'].split(':')
    dd['real'] = float(real_minutes)*60. + float(real_seconds)
    for k in dd:
        dd[k] = float(dd[k])

    print i, hh,dd
    nworkers += [ hh['nworkers'] ]
    nprocs   += [ hh['nprocs'] ]
    ntuple += [ dict(nworkers=hh['nworkers'], nprocs=hh['nprocs'],
                     usr=dd['usr'],
                     sys=dd['sys'],
                     real=dd['real'],
                     cpu=dd['cpu']) ]

data = {
    'nworkers': [],
    'nprocs': [],
    'usr': [],
    'sys': [],
    'real': [],
    'cpu': [],
}

for n in ntuple:
    for k in data:
        data[k] += [n[k]]

if 0:
    for k in data:
        print k, data[k]

my_data_x = []
my_data_y = []
print "="*80
for d in ntuple:
    if d['nworkers'] != 5000:
        continue
    print d
    my_data_x += [d['nprocs']]
    my_data_y += [d['real']]
print "="*80
print "x=",my_data_x
print "y=",my_data_y




