package org.acme.clazz.interfaces.rest.exception;

import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.ext.ExceptionMapper;
import jakarta.ws.rs.ext.Provider;
import org.acme.clazz.domain.exception.ClazzNotFoundException;

@Provider
public class ClazzNotFoundExceptionMapper implements ExceptionMapper<ClazzNotFoundException> {

    @Override
    public Response toResponse(ClazzNotFoundException exception) {
        return Response.status(Response.Status.NOT_FOUND)
            .entity(new ErrorResponse(exception.getMessage()))
            .build();
    }
}

