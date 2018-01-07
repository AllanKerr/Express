package com.kerr;

import java.util.Enumeration;
import javax.servlet.http.HttpServletRequest;
import javax.validation.Valid;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@SpringBootApplication
@RestController
public class ApiApplication {

	public static void main(String[] args) {
		SpringApplication.run(ApiApplication.class, args);
	}

	@RequestMapping("/")
	public String index(HttpServletRequest request) {


		Enumeration<String> headers = request.getHeaderNames();
		while (headers.hasMoreElements()) {

			String name = headers.nextElement();
			System.out.println(name + " : " + request.getHeader(name));

		}
		return "Hello World";
	}

	@RequestMapping(value = "{configKey}/**")
	public String index(HttpServletRequest request, @PathVariable String configKey) {

		Enumeration<String> headers = request.getHeaderNames();
		while (headers.hasMoreElements()) {

			String name = headers.nextElement();
			System.out.println(name + " : " + request.getHeader(name));

		}

		return "Found: " + configKey;
	}


	@RequestMapping(value = "/search", method = RequestMethod.POST)
	public String addSearch(@Valid @RequestBody SearchRequest search) {

		return search.getPark();
	}
}
