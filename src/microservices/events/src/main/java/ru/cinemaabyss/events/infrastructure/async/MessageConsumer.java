package ru.cinemaabyss.events.infrastructure.async;

import lombok.extern.java.Log;
import org.apache.kafka.clients.consumer.Consumer;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.support.serializer.JacksonJsonDeserializer;
import org.springframework.kafka.support.serializer.JsonDeserializer;
import org.springframework.stereotype.Component;

import java.time.Duration;
import java.util.*;


@Component
@Log
public class MessageConsumer {
    @Value("${bootstrap.servers}")
    private String server;

    @Value("${topic.movie.events}")
    private String topicMovieEvents;

    @Value("${topic.user.events}")
    private String topicUserEvents;

    @Value("${topic.payment.events}")
    private String topicPaymentEvents;

    public void receiveAll() {
        receiveMovie();
        receiveUser();
        receivePayment();
    }

    private void receiveMovie() {
        receiveMessage("group-id-test", topicMovieEvents);
    }

    private void receiveUser() {
        receiveMessage("group-id-test", topicUserEvents);
    }

    private void receivePayment() {
        receiveMessage("group-id-test", topicPaymentEvents);
    }


    private void receiveMessage(String groupId, String topic) {
        Properties props = new Properties();
        props.put("bootstrap.servers", server);
        props.put("group.id", groupId);
        props.put("enable.auto.commit", "true");
        props.put("key.deserializer", StringDeserializer.class.getName());
        props.put("value.deserializer", JacksonJsonDeserializer.class.getName());
        props.put(JsonDeserializer.TRUSTED_PACKAGES.toString(), "ru.cinemaabyss.events.models");

        try (Consumer<String, Object> consumer = new KafkaConsumer<>(props)) {
            consumer.subscribe(Collections.singletonList(topic));
            ConsumerRecords<String, Object> records = consumer.poll(Duration.ofMillis(1000));

            for(ConsumerRecord<String, Object> record : records) {
                log.info("Key: " + record.key() + ", value: " + record.value());
            }
        }
    }

}
