package ru.cinemaabyss.events.infrastructure.http;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.*;
import org.springframework.web.bind.annotation.*;
import ru.cinemaabyss.events.models.Movie;
import ru.cinemaabyss.events.models.Payment;
import ru.cinemaabyss.events.models.User;
import ru.cinemaabyss.events.services.AsyncService;

@RestController
@RequestMapping(path="api/events/", produces = MediaType.APPLICATION_JSON_VALUE)
public class EventsController {

    @Autowired
    private AsyncService asyncService;

    @GetMapping(path = "health")
    public ResponseEntity<String> getHealth() {
        return  ResponseEntity
                .status(HttpStatus.OK)
                .body("{\"status\": true}");
    }

    @PostMapping(path = "movie", consumes = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<Response> createMovie(@RequestBody Movie movie) {
        asyncService.createMovie(movie);
        return  ResponseEntity
                .status(HttpStatus.CREATED)
                .body(new Response(Status.success));
    }

    @PostMapping(path = "user", consumes = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<Response> createUser(@RequestBody User user) {
        asyncService.createUser(user);
        return  ResponseEntity
                .status(HttpStatus.CREATED)
                .body(new Response(Status.success));
    }

    @PostMapping(path = "payment", consumes = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<Response> createPayment(@RequestBody Payment payments) {
        asyncService.createPayment(payments);
        return  ResponseEntity
                .status(HttpStatus.CREATED)
                .body(new Response(Status.success));
    }
}
