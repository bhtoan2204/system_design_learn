package org.acme.clazz.infrastructure.persistence;

import jakarta.annotation.PostConstruct;
import jakarta.inject.Singleton;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.stream.IntStream;
import org.acme.clazz.domain.model.Clazz;
import org.acme.clazz.domain.repository.ClazzRepository;

@Singleton
public class InMemoryClazzRepository implements ClazzRepository {

    private final List<Clazz> classes = new CopyOnWriteArrayList<>();

    @PostConstruct
    void seedData() {
        if (classes.isEmpty()) {
            IntStream.rangeClosed(1, 3)
                .mapToObj(i -> Clazz.of(i, "Class " + i))
                .forEach(classes::add);
        }
    }

    @Override
    public List<Clazz> findAll() {
        return List.copyOf(classes);
    }

    @Override
    public Optional<Clazz> findById(Integer id) {
        return classes.stream().filter(clazz -> clazz.id().equals(id)).findFirst();
    }

    @Override
    public Clazz save(Clazz clazz) {
        classes.removeIf(existing -> existing.id().equals(clazz.id()));
        classes.add(clazz);
        return clazz;
    }

    @Override
    public boolean existsByName(String name) {
        return classes.stream().anyMatch(clazz -> clazz.name().equalsIgnoreCase(name));
    }
}

