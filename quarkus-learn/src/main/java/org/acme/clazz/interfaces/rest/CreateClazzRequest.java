package org.acme.clazz.interfaces.rest;

import jakarta.validation.constraints.NotBlank;

public record CreateClazzRequest(@NotBlank String name) {
}

