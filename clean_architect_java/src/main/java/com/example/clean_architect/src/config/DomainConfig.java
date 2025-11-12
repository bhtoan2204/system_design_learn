package com.example.clean_architect.src.config;

import com.example.clean_architect.src.domain.service.UserDomainService;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class DomainConfig {

	@Bean
	public UserDomainService userDomainService() {
		return new UserDomainService();
	}
}
