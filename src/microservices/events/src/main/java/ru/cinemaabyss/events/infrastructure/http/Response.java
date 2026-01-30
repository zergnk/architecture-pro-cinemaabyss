package ru.cinemaabyss.events.infrastructure.http;

import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class Response {
    private Status status;
}
