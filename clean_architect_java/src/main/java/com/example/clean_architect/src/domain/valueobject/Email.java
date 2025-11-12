package com.example.clean_architect.src.domain.valueobject;

import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.ToString;

import java.util.regex.Pattern;

@Getter
@EqualsAndHashCode(of = "value")
@ToString
public final class Email {

	private static final Pattern EMAIL_PATTERN = Pattern.compile("^[A-Za-z0-9+_.-]+@[A-Za-z0-9.-]+$");

	private final String value;

	private Email(String value) {
		this.value = value;
	}

	public static Email of(String rawValue) {
		String sanitized = sanitize(rawValue);
		if (!EMAIL_PATTERN.matcher(sanitized).matches()) {
			throw new IllegalArgumentException("Invalid email format");
		}
		return new Email(sanitized);
	}

	private static String sanitize(String value) {
		if (value == null) {
			throw new IllegalArgumentException("Email cannot be null");
		}
		String trimmed = value.trim().toLowerCase();
		if (trimmed.isEmpty()) {
			throw new IllegalArgumentException("Email cannot be empty");
		}
		return trimmed;
	}
}
