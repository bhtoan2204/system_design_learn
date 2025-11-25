package com.example.clean_architect.src.domain.service;

import java.util.Optional;

import com.example.clean_architect.src.domain.entity.User;

public class UserDomainService {

	public void ensureUserCanRegister(User candidate, Optional<User> existingUser) {
		existingUser.ifPresent(existing -> {
			throw new IllegalStateException(
					"User with email %s already exists".formatted(existing.getEmail().getValue()));
		});

		if (candidate.getName().equalsIgnoreCase("admin")) {
			throw new IllegalArgumentException("Reserved names cannot be used");
		}
	}
}
