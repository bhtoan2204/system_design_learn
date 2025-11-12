package com.example.clean_architect.src.application.usecase;

import com.example.clean_architect.src.domain.entity.User;
import com.example.clean_architect.src.domain.event.UserRegisteredEvent;
import com.example.clean_architect.src.domain.repository.UserRepository;
import com.example.clean_architect.src.domain.service.UserDomainService;
import com.example.clean_architect.src.domain.valueobject.Email;
import org.springframework.stereotype.Service;

import java.util.Optional;

@Service
public class RegisterUserUseCase {

	private final UserRepository userRepository;
	private final UserDomainService userDomainService;

	public RegisterUserUseCase(UserRepository userRepository, UserDomainService userDomainService) {
		this.userRepository = userRepository;
		this.userDomainService = userDomainService;
	}

	public RegisterUserResponse register(RegisterUserCommand command) {
		Email email = Email.of(command.email());
		User candidate = User.create(null, command.name(), email);

		Optional<User> existingUser = userRepository.findByEmail(email);
		userDomainService.ensureUserCanRegister(candidate, existingUser);

		User saved = userRepository.save(candidate);
		UserRegisteredEvent event = new UserRegisteredEvent(saved);

		return RegisterUserResponse.builder()
				.id(saved.getId())
				.name(saved.getName())
				.email(saved.getEmail().getValue())
				.registeredAt(event.occurredOn())
				.build();
	}
}
