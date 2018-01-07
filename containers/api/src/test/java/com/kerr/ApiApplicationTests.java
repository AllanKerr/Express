package com.kerr;

import org.cassandraunit.spring.EmbeddedCassandra;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.TestExecutionListeners;
import org.springframework.test.context.junit4.SpringRunner;

@RunWith(SpringRunner.class)
@SpringBootTest
@EmbeddedCassandra
@TestExecutionListeners(listeners = { OrderedCassandraTestExecutionListener.class })
public class ApiApplicationTests {

	@Test
	public void contextLoads() {
	}
}
