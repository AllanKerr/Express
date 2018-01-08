package com.kerr.repositories;

import com.kerr.domain.Search;
import com.kerr.domain.SearchKey;
import org.springframework.data.cassandra.repository.TypedIdCassandraRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface SearchRepository extends TypedIdCassandraRepository<Search, SearchKey> {

}
