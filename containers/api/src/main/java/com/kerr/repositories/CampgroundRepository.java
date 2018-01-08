package com.kerr.repositories;

import com.kerr.domain.Campground;
import com.kerr.domain.CampgroundKey;
import org.springframework.data.cassandra.repository.TypedIdCassandraRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface CampgroundRepository extends TypedIdCassandraRepository<Campground, CampgroundKey> {

}
