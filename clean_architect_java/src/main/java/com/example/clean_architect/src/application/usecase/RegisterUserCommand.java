package com.example.clean_architect.src.application.usecase;

import lombok.Builder;

@Builder
public record RegisterUserCommand(String name, String email) {
    public static Builder builder() {
        return new Builder();
    }

    public static class Builder {
        private String name;
        private String email;

        public Builder name(String name) {
            this.name = name;
            return this;
        }

        public Builder email(String email) {
            this.email = email;
            return this;
        }

        public RegisterUserCommand build() {
            return new RegisterUserCommand(name, email);
        }
    }
}
