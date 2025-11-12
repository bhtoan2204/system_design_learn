package com.example.clean_architect.src.domain.repository;

import com.example.clean_architect.src.domain.entity.User;
import com.example.clean_architect.src.domain.valueobject.Email;

import java.util.Optional;
import java.util.UUID;

public interface UserRepository {

	Optional<User> findByEmail(Email email);

	Optional<User> findById(UUID id);

	User save(User user);
}
