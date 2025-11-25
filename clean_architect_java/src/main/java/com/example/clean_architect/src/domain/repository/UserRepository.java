package com.example.clean_architect.src.domain.repository;

import java.util.Optional;
import java.util.UUID;

import com.example.clean_architect.src.domain.entity.User;
import com.example.clean_architect.src.domain.valueobject.Email;

public interface UserRepository {

	Optional<User> findByEmail(Email email);

	Optional<User> findById(UUID id);

	User save(User user);
}
