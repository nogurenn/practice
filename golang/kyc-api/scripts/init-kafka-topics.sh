#/bin/sh

# blocks until kafka is reachable
kafka-topics --bootstrap-server kafka-broker-1:29092 --list
kafka-topics --bootstrap-server kafka-broker-2:29092 --list
kafka-topics --bootstrap-server kafka-broker-3:29092 --list

echo -e 'Creating kafka topics'
kafka-topics --bootstrap-server kafka-broker-1:29092 --create --if-not-exists --topic kyc-verification-requests --replication-factor 3 --partitions 1

echo -e 'Successfully created the following topics:'
kafka-topics --bootstrap-server kafka-broker-1:29092 --list
