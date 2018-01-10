package com.kerr.repositories;

import com.kerr.domain.Search;
import com.kerr.domain.SearchKey;
import java.util.List;
import org.springframework.data.cassandra.repository.TypedIdCassandraRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface SearchRepository extends TypedIdCassandraRepository<Search, SearchKey> {

    List<Search> findAllByKey_UserId(String userId);
}
