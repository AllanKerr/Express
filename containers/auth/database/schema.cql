
CREATE KEYSPACE IF NOT EXISTS "default" WITH REPLICATION = {
    'class' : 'SimpleStrategy',
    'replication_factor' : 3
};
