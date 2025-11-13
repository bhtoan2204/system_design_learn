package org.acme.clazz.interfaces.rest;

import jakarta.inject.Inject;
import jakarta.validation.Valid;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import java.net.URI;
import java.util.List;
import org.acme.clazz.application.bus.CommandDispatcher;
import org.acme.clazz.application.bus.QueryDispatcher;
import org.acme.clazz.application.command.CreateClazzCommand;
import org.acme.clazz.application.model.ClazzDto;
import org.acme.clazz.application.query.GetAllClazzesQuery;
import org.acme.clazz.application.query.GetClazzByIdQuery;

@Path("/classes")
@Produces(MediaType.APPLICATION_JSON)
public class ClazzResource {

    private final QueryDispatcher queryDispatcher;
    private final CommandDispatcher commandDispatcher;

    @Inject
    public ClazzResource(QueryDispatcher queryDispatcher, CommandDispatcher commandDispatcher) {
        this.queryDispatcher = queryDispatcher;
        this.commandDispatcher = commandDispatcher;
    }

    @GET
    public List<ClazzDto> listClasses() {
        return queryDispatcher.dispatch(new GetAllClazzesQuery());
    }

    @GET
    @Path("/{id}")
    public ClazzDto getClazz(@PathParam("id") Integer id) {
        return queryDispatcher.dispatch(new GetClazzByIdQuery(id));
    }

    @POST
    @Consumes(MediaType.APPLICATION_JSON)
    public Response createClazz(@Valid CreateClazzRequest request) {
        var created = commandDispatcher.dispatch(new CreateClazzCommand(request.name()));
        return Response
            .created(URI.create("/classes/" + created.id()))
            .entity(created)
            .build();
    }
}

