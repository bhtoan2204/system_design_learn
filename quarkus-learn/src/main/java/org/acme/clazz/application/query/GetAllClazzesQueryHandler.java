package org.acme.clazz.application.query;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import java.util.List;
import org.acme.clazz.application.bus.QueryHandler;
import org.acme.clazz.application.mapper.ClazzMapper;
import org.acme.clazz.application.model.ClazzDto;
import org.acme.clazz.domain.repository.ClazzRepository;

@ApplicationScoped
public class GetAllClazzesQueryHandler implements QueryHandler<GetAllClazzesQuery, List<ClazzDto>> {

    private final ClazzRepository clazzRepository;
    private final ClazzMapper mapper;

    @Inject
    public GetAllClazzesQueryHandler(ClazzRepository clazzRepository, ClazzMapper mapper) {
        this.clazzRepository = clazzRepository;
        this.mapper = mapper;
    }

    @Override
    public List<ClazzDto> handle(GetAllClazzesQuery query) {
        return mapper.toDtoList(clazzRepository.findAll());
    }

    @Override
    public Class<GetAllClazzesQuery> queryType() {
        return GetAllClazzesQuery.class;
    }
}

