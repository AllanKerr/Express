FROM gcr.io/google-samples/cassandra:v11

COPY schema.cql /schema.cql

CMD ["/sbin/dumb-init", "/bin/bash", "/run.sh", "cqlsh -f /schema.cql"]
