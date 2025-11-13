package org.acme.clazz.domain.repository;

import java.util.List;
import java.util.Optional;
import org.acme.clazz.domain.model.Clazz;

public interface ClazzRepository {

    List<Clazz> findAll();

    Optional<Clazz> findById(Integer id);

    Clazz save(Clazz clazz);

    boolean existsByName(String name);
}

