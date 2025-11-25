package com.example.clean_architect.src.domain.entity;

import java.util.Objects;
import java.util.UUID;

import com.example.clean_architect.src.domain.valueobject.Email;

import lombok.AccessLevel;
import lombok.Builder;
import lombok.Getter;

@Getter
public class User {

	private final UUID id;
	private String name;
	private Email email;

	@Builder(access = AccessLevel.PRIVATE)
	private User(UUID id, String name, Email email) {
		this.id = id == null ? UUID.randomUUID() : id;
		setName(name);
		setEmail(email);
	}

	public static User create(UUID id, String name, Email email) {
		return User.builder().id(id).name(name).email(email).build();
	}

	public void changeName(String newName) {
		setName(newName);
	}

	private void setName(String value) {
		if (value == null || value.trim().isEmpty()) {
			throw new IllegalArgumentException("Name cannot be blank");
		}
		this.name = value.trim();
	}

	private void setEmail(Email email) {
		this.email = Objects.requireNonNull(email, "email must not be null");
	}
}
