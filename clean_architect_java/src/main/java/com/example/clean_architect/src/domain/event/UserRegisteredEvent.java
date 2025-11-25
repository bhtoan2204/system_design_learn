package com.example.clean_architect.src.domain.event;

import java.time.Instant;
import java.util.UUID;

import com.example.clean_architect.src.domain.entity.User;

public class UserRegisteredEvent implements DomainEvent {

	private final UUID userId;
	private final String email;
	private final Instant occurredOn;

	public UserRegisteredEvent(User user) {
		this.userId = user.getId();
		this.email = user.getEmail().getValue();
		this.occurredOn = Instant.now();
	}

	public UUID getUserId() {
		return userId;
	}

	public String getEmail() {
		return email;
	}

	@Override
	public Instant occurredOn() {
		return occurredOn;
	}
}
