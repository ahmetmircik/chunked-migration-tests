#!/usr/bin/env bash

set -e

hzVersion=${1:-$(hazelcastVersion)}
user=${2:-ec2-user}

#aws-create --count 3 --instanceType m6g.12xlarge	
aws-create --count 3 --instanceType m6g.16xlarge

ops="$ops -Dhazelcast.enterprise.license.key=${HAZELCAST_EE_KEY}"
ops="$ops -Dhazelcast.partition.count=101"
ops="$ops -Dhazelcast.partition.max.parallel.migrations=5"
ops="$ops -Dhazelcast.partition.migration.chunks.enabled=false"
ops="$ops -Dhazelcast.partition.migration.chunks.max.migrating.data.in.mb=50"
ops="$ops -ea"

mOps="$mOps -Dlog4j.configuration=file:./log4j.properties -Dhazelcast.logging.type=log4j"

hz memberOps "-Xms1G -Xmx1G ${ops} ${mOps}"
hz clientOps "-Xms1G -Xmx1G ${ops}"
hz cluster -ee -size M1 -v ${hzVersion} -boxes a.box -user ${user} -upcwd log4j.properties

hz driver Member
hz run multiLoad 

echo "finished multi load"

hz cluster -ee -size M2 -v ${hzVersion} -boxes a.box -user ${user} -upcwd log4j.properties

hz run untilClusterSize3 untilClusterSafe

hz download
hz wipe
aws-terminate
