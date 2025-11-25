package com.example.clean_architect.src.application.usecase;

import java.time.Instant;
import java.util.UUID;

import lombok.Builder;

@Builder
public record RegisterUserResponse(UUID id, String name, String email, Instant registeredAt) {
}
