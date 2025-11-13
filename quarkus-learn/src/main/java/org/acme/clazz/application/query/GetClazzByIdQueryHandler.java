package org.acme.clazz.application.query;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import org.acme.clazz.application.bus.QueryHandler;
import org.acme.clazz.application.mapper.ClazzMapper;
import org.acme.clazz.application.model.ClazzDto;
import org.acme.clazz.domain.exception.ClazzNotFoundException;
import org.acme.clazz.domain.repository.ClazzRepository;

@ApplicationScoped
public class GetClazzByIdQueryHandler implements QueryHandler<GetClazzByIdQuery, ClazzDto> {

    private final ClazzRepository clazzRepository;
    private final ClazzMapper mapper;

    @Inject
    public GetClazzByIdQueryHandler(ClazzRepository clazzRepository, ClazzMapper mapper) {
        this.clazzRepository = clazzRepository;
        this.mapper = mapper;
    }

    @Override
    public ClazzDto handle(GetClazzByIdQuery query) {
        return clazzRepository.findById(query.id())
            .map(mapper::toDto)
            .orElseThrow(() -> new ClazzNotFoundException(query.id()));
    }

    @Override
    public Class<GetClazzByIdQuery> queryType() {
        return GetClazzByIdQuery.class;
    }
}

