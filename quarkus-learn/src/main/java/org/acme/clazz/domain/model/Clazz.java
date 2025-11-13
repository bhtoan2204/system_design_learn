package org.acme.clazz.domain.model;

import java.util.Objects;

/**
 * Domain aggregate representing a school class.
 */
public final class Clazz {

    private final Integer id;
    private final String name;

    private Clazz(Integer id, String name) {
        this.id = Objects.requireNonNull(id, "id must not be null");
        this.name = normalizeName(name);
    }

    public static Clazz of(Integer id, String name) {
        return new Clazz(id, name);
    }

    private static String normalizeName(String value) {
        Objects.requireNonNull(value, "name must not be null");
        var trimmed = value.trim();
        if (trimmed.isEmpty()) {
            throw new IllegalArgumentException("name must not be blank");
        }
        return trimmed;
    }

    public Integer id() {
        return id;
    }

    public String name() {
        return name;
    }

    public Clazz rename(String newName) {
        return Clazz.of(id, newName);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        Clazz clazz = (Clazz) o;
        return id.equals(clazz.id);
    }

    @Override
    public int hashCode() {
        return Objects.hash(id);
    }

    @Override
    public String toString() {
        return "Clazz{" +
            "id=" + id +
            ", name='" + name + '\'' +
            '}';
    }
}

