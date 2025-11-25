package com.example.clean_architect.src.presentation.http.request;

import com.example.clean_architect.src.application.usecase.RegisterUserCommand;

public record RegisterUserRequest(String name, String email) {
    public RegisterUserCommand toCommand() {
        return RegisterUserCommand.builder().name(name).email(email).build();
    }
}
