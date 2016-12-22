hadoop on Runv
==============

### perf

The cluster statistics:

*. 5 bare metals
*. KVM
*. Linux Bridge
*. hdfs on Runv: 4 VCores, 8G * (1 namenode + 3 datanode)
*. yarn on Runv: 8VCores, 16G * (1 resource manager + 4 node manager)

in house benchmark:

Terasort in 100G data costed:

```
real    601m25.369s
user    0m56.305s
sys     0m21.853s
```

Yahoo benchmark:

In May 2008, Owen O'Malley ran this code on a 910 node cluster and sorted the 10 billion records (1 TB) in 209 seconds (3.48 minutes) to win the annual general purpose (daytona) terabyte sort benchmark.

*. 910 nodes
*. 4 dual core Xeons @ 2.0ghz per a node
*. 4 SATA disks per a node
*. 8G RAM per a node
*. 1 gigabit ethernet on each node
*. 40 nodes per a rack
*. 8 gigabit ethernet uplinks from each rack to the core
*. Red Hat Enterprise Linux Server Release 5.1 (kernel 2.6.18)
*. Sun Java JDK 1.6.0_05-b13

refer link:

*. http://www.michael-noll.com/blog/2011/04/09/benchmarking-and-stress-testing-an-hadoop-cluster-with-terasort-testdfsio-nnbench-mrbench/
*. http://hadoop.apache.org/docs/r2.7.3/api/org/apache/hadoop/examples/terasort/package-summary.html

