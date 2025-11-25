package com.example.clean_architect.src.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import com.example.clean_architect.src.domain.service.UserDomainService;

@Configuration
public class DomainConfig {

	@Bean
	public UserDomainService userDomainService() {
		return new UserDomainService();
	}
}
