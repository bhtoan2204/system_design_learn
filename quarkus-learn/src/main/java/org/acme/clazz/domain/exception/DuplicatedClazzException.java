package org.acme.clazz.domain.exception;

public class DuplicatedClazzException extends RuntimeException {

    public DuplicatedClazzException(String name) {
        super("Clazz with name '%s' already exists".formatted(name));
    }
}

