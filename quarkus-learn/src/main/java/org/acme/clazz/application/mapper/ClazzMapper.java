package org.acme.clazz.application.mapper;

import jakarta.enterprise.context.ApplicationScoped;
import java.util.List;
import java.util.stream.Collectors;
import org.acme.clazz.application.model.ClazzDto;
import org.acme.clazz.domain.model.Clazz;

@ApplicationScoped
public class ClazzMapper {

    public ClazzDto toDto(Clazz clazz) {
        return new ClazzDto(clazz.id(), clazz.name());
    }

    public List<ClazzDto> toDtoList(List<Clazz> classes) {
        return classes.stream().map(this::toDto).collect(Collectors.toList());
    }
}

