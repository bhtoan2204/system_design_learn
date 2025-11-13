package org.acme.clazz.application.query;

import org.acme.clazz.application.bus.Query;
import org.acme.clazz.application.model.ClazzDto;

public record GetClazzByIdQuery(Integer id) implements Query<ClazzDto> {
}

