package org.acme.clazz.application.bus;

import jakarta.enterprise.inject.Instance;
import jakarta.inject.Inject;
import jakarta.inject.Singleton;
import java.util.HashMap;
import java.util.Map;

@Singleton
public class QueryDispatcher {

    private final Map<Class<?>, QueryHandler<?, ?>> handlerMap = new HashMap<>();

    @Inject
    public QueryDispatcher(Instance<QueryHandler<?, ?>> handlers) {
        for (QueryHandler<?, ?> handler : handlers) {
            handlerMap.put(handler.queryType(), handler);
        }
    }

    @SuppressWarnings("unchecked")
    public <R, Q extends Query<R>> R dispatch(Q query) {
        var handler = (QueryHandler<Q, R>) handlerMap.get(query.getClass());
        if (handler == null) {
            throw new IllegalStateException("No query handler registered for " + query.getClass().getName());
        }
        return handler.handle(query);
    }
}

