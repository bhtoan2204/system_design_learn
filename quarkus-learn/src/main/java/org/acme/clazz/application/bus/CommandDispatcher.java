package org.acme.clazz.application.bus;

import jakarta.enterprise.inject.Instance;
import jakarta.inject.Inject;
import jakarta.inject.Singleton;
import java.util.HashMap;
import java.util.Map;

@Singleton
public class CommandDispatcher {

    private final Map<Class<?>, CommandHandler<?, ?>> handlerMap = new HashMap<>();

    @Inject
    public CommandDispatcher(Instance<CommandHandler<?, ?>> handlers) {
        for (CommandHandler<?, ?> handler : handlers) {
            handlerMap.put(handler.commandType(), handler);
        }
    }

    @SuppressWarnings("unchecked")
    public <R, C extends Command<R>> R dispatch(C command) {
        var handler = (CommandHandler<C, R>) handlerMap.get(command.getClass());
        if (handler == null) {
            throw new IllegalStateException("No command handler registered for " + command.getClass().getName());
        }
        return handler.handle(command);
    }
}

