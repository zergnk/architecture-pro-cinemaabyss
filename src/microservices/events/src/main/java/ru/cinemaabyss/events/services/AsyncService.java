package ru.cinemaabyss.events.services;

import jakarta.annotation.PostConstruct;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import ru.cinemaabyss.events.infrastructure.async.MessageConsumer;
import ru.cinemaabyss.events.infrastructure.async.MessageProducer;
import ru.cinemaabyss.events.models.Movie;
import ru.cinemaabyss.events.models.Payment;
import ru.cinemaabyss.events.models.User;

@Service
public class AsyncService {

    @Autowired
    private MessageConsumer messageConsumer;

    @Autowired
    private MessageProducer messageProducer;

    public void createUser(User user) {
        messageProducer.send(user);
    }

    public void createMovie(Movie movie) {
        messageProducer.send(movie);
    }

    public void createPayment(Payment payment) {
        messageProducer.send(payment);
    }

    @PostConstruct
    public void startMessageReceiving() {
        Thread thread = new Thread(new ConsumerRunner());
        thread.start();
    }


    class ConsumerRunner implements Runnable {

        @Override
        public void run() {
            while(true) {
                messageConsumer.receiveAll();
            }
        }
    }
}
