package com.example.user.event;

import lombok.Data;

@Data
public class UserEvent {
    private String eventType; // CREATED, UPDATED, DELETED
    private Long userId;
    private String username;
    private String email;
}
