package com.example.user.kafka;

import com.example.user.event.UserEvent;
import lombok.RequiredArgsConstructor;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class KafkaProducer {
    private final KafkaTemplate<String, UserEvent> kafkaTemplate;
    private static final String TOPIC = "user-events";

    public void sendUserEvent(UserEvent event) {
        kafkaTemplate.send(TOPIC, String.valueOf(event.getUserId()), event);
    }
}
