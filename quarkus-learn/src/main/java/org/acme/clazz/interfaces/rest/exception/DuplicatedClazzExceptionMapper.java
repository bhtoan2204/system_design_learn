package org.acme.clazz.interfaces.rest.exception;

import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.ext.ExceptionMapper;
import jakarta.ws.rs.ext.Provider;
import org.acme.clazz.domain.exception.DuplicatedClazzException;

@Provider
public class DuplicatedClazzExceptionMapper implements ExceptionMapper<DuplicatedClazzException> {

    @Override
    public Response toResponse(DuplicatedClazzException exception) {
        return Response.status(Response.Status.CONFLICT)
            .entity(new ErrorResponse(exception.getMessage()))
            .build();
    }
}

