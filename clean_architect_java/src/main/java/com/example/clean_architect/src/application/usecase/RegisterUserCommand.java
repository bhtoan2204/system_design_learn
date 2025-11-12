package com.example.clean_architect.src.application.usecase;

import lombok.Builder;

@Builder
public record RegisterUserCommand(String name, String email) {
}
