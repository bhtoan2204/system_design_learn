package com.example.clean_architect.src.presentation.http;

import com.example.clean_architect.src.application.usecase.RegisterUserCommand;
import com.example.clean_architect.src.application.usecase.RegisterUserResponse;
import com.example.clean_architect.src.application.usecase.RegisterUserUseCase;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/users")
public class UserHttp {

    private final RegisterUserUseCase registerUserUseCase;

    public UserHttp(RegisterUserUseCase registerUserUseCase) {
        this.registerUserUseCase = registerUserUseCase;
    }

    @PostMapping
    public ResponseEntity<RegisterUserResponse> register(@RequestBody RegisterUserRequest request) {
        RegisterUserResponse response = registerUserUseCase.register(RegisterUserCommand.builder()
                .name(request.name())
                .email(request.email())
                .build());
        return ResponseEntity.ok(response);
    }

    public record RegisterUserRequest(String name, String email) {
    }
}
