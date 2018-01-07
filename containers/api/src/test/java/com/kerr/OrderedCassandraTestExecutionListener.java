package com.kerr;

import org.cassandraunit.spring.CassandraUnitDependencyInjectionTestExecutionListener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.core.Ordered;

public class OrderedCassandraTestExecutionListener extends CassandraUnitDependencyInjectionTestExecutionListener {

    private static final Logger logger = LoggerFactory
            .getLogger(OrderedCassandraTestExecutionListener.class);

    @Override
    public int getOrder() {
        return Ordered.HIGHEST_PRECEDENCE;
    }
}
