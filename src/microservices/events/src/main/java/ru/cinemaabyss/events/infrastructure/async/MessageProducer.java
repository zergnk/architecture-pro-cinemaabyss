package ru.cinemaabyss.events.infrastructure.async;

import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.Producer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.support.serializer.JacksonJsonSerializer;
import org.springframework.stereotype.Component;
import ru.cinemaabyss.events.models.Movie;
import ru.cinemaabyss.events.models.Payment;
import ru.cinemaabyss.events.models.User;

import java.util.Properties;

@Component
public class MessageProducer {

    @Value("${bootstrap.servers}")
    private String server;

    @Value("${topic.movie.events}")
    private String topicMovieEvents;

    @Value("${topic.user.events}")
    private String topicUserEvents;

    @Value("${topic.payment.events}")
    private String topicPaymentEvents;


    public void send(Object message) {
        Properties props = new Properties();
        props.put("bootstrap.servers",server);
        props.put("key.serializer", org.apache.kafka.common.serialization.StringSerializer.class.getName());
        props.put("value.serializer", JacksonJsonSerializer.class.getName());

        String topic;
        if (Movie.class.equals(message.getClass())) {
            topic = topicMovieEvents;
        } else if (User.class.equals(message.getClass())) {
            topic = topicUserEvents;
        } else if (Payment.class.equals(message.getClass())) {
            topic = topicPaymentEvents;
        } else {
            throw new RuntimeException("Unknown message type.");
        }

        Producer<String, Object> producer = new KafkaProducer<>(props);
        producer.send(new ProducerRecord<>(topic, message.getClass().getName(), message));
        producer.close();
    }
}
