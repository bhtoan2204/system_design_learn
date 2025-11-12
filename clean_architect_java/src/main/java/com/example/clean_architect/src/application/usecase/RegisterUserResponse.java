package com.example.clean_architect.src.application.usecase;

import lombok.Builder;

import java.time.Instant;
import java.util.UUID;

@Builder
public record RegisterUserResponse(UUID id, String name, String email, Instant registeredAt) {
}
