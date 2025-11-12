package com.example.clean_architect.src.infrastructure.repository;

import com.example.clean_architect.src.domain.entity.User;
import com.example.clean_architect.src.domain.repository.UserRepository;
import com.example.clean_architect.src.domain.valueobject.Email;
import org.springframework.stereotype.Repository;

import java.util.Map;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;

@Repository
public class InMemoryUserRepository implements UserRepository {

	private final Map<UUID, User> usersById = new ConcurrentHashMap<>();
	private final Map<String, UUID> idsByEmail = new ConcurrentHashMap<>();

	@Override
	public Optional<User> findByEmail(Email email) {
		UUID id = idsByEmail.get(email.getValue());
		return Optional.ofNullable(id).map(usersById::get);
	}

	@Override
	public Optional<User> findById(UUID id) {
		return Optional.ofNullable(usersById.get(id));
	}

	@Override
	public User save(User user) {
		usersById.put(user.getId(), user);
		idsByEmail.put(user.getEmail().getValue(), user.getId());
		return user;
	}
}
