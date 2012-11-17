#!/bin/bash

jobo=./data/rdotoesd-notrigger.json
nprocs='0 1 2 4 6 8 10 12 14 16 18 20 22 24 32 56'
evtmax=100

echo "::: parameters..."
echo ":: jobo:   [$jobo]"
echo ":: nprocs: [$nprocs]"
echo ":: evtmax: [$evtmax]"

# echo "::: flat-sequence :::"
# for i in 0 16
# do
#     echo "::: nprocs: $i -- "`date`
#     time go-mini-gaudi -evtmax ${evtmax} -jobo ${jobo} -seq -nprocs $i >| log.sseq.${i}
#     echo "::: nprocs:$i:$?"
#     echo ""
# done

echo ""
echo "::: parallel-sequence :::"
for i in $nprocs
do
    echo "::: nprocs: $i -- "`date`
    time go-mini-gaudi -evtmax ${evtmax} -jobo ${jobo} -nprocs $i >| log.pseq.${i}
    echo "::: nprocs:$i:$?"
    echo ""
done
