package org.acme.clazz.application.command;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import java.util.concurrent.atomic.AtomicInteger;
import org.acme.clazz.application.bus.CommandHandler;
import org.acme.clazz.application.mapper.ClazzMapper;
import org.acme.clazz.application.model.ClazzDto;
import org.acme.clazz.domain.exception.DuplicatedClazzException;
import org.acme.clazz.domain.model.Clazz;
import org.acme.clazz.domain.repository.ClazzRepository;

@ApplicationScoped
public class CreateClazzCommandHandler implements CommandHandler<CreateClazzCommand, ClazzDto> {

    private final ClazzRepository clazzRepository;
    private final ClazzMapper mapper;
    private final AtomicInteger idGenerator = new AtomicInteger(1000);

    @Inject
    public CreateClazzCommandHandler(ClazzRepository clazzRepository, ClazzMapper mapper) {
        this.clazzRepository = clazzRepository;
        this.mapper = mapper;
    }

    @Override
    public ClazzDto handle(CreateClazzCommand command) {
        var normalizedName = normalize(command.name());
        if (clazzRepository.existsByName(normalizedName)) {
            throw new DuplicatedClazzException(normalizedName);
        }
        var clazz = Clazz.of(idGenerator.incrementAndGet(), normalizedName);
        var saved = clazzRepository.save(clazz);
        return mapper.toDto(saved);
    }

    private static String normalize(String rawName) {
        if (rawName == null) {
            throw new IllegalArgumentException("name must not be null");
        }
        var normalized = rawName.trim();
        if (normalized.isEmpty()) {
            throw new IllegalArgumentException("name must not be blank");
        }
        return normalized;
    }

    @Override
    public Class<CreateClazzCommand> commandType() {
        return CreateClazzCommand.class;
    }
}

