package ru.cinemaabyss.events.models;

import lombok.Data;

@Data
public class Movie {
    private Integer movie_id;
    private String title;
    private String action;
    private Integer user_id;
}
