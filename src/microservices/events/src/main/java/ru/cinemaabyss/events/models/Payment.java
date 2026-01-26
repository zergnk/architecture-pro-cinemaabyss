package ru.cinemaabyss.events.models;

import lombok.Data;

import java.time.LocalDateTime;

@Data
public class Payment {
   private Integer payment_id;
   private Integer user_id;
   private Double amount;
   private String status;
   private LocalDateTime timestamp;
   private String method_type;
}
