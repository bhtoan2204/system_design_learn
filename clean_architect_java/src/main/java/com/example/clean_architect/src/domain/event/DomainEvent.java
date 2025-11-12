package com.example.clean_architect.src.domain.event;

import java.time.Instant;

public interface DomainEvent {

	Instant occurredOn();
}
