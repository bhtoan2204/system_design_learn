package org.acme.clazz.application.command;

import org.acme.clazz.application.bus.Command;
import org.acme.clazz.application.model.ClazzDto;

public record CreateClazzCommand(String name) implements Command<ClazzDto> {
}

