package org.acme.clazz.domain.exception;

public class ClazzNotFoundException extends RuntimeException {

    public ClazzNotFoundException(Integer id) {
        super("Clazz with id %d not found".formatted(id));
    }
}

