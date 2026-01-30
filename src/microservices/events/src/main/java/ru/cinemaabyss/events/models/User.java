package ru.cinemaabyss.events.models;

import lombok.Data;

import java.time.LocalDateTime;

@Data
public class User {
    private Integer user_id;
    private String username;
    private String action;
    private LocalDateTime timestamp;
}
