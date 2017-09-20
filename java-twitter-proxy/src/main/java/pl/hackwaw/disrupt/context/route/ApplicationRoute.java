package pl.hackwaw.disrupt.context.route;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import pl.hackwaw.disrupt.context.repository.TweetRepository;
import pl.hackwaw.disrupt.context.repository.TwitterObservable;

import java.util.concurrent.TimeUnit;

@Configuration
public class ApplicationRoute {

    @Autowired
    private TweetRepository tweetRepository;

    @Bean
    public CommandLineRunner initialState() {
        return (args) -> TwitterObservable
                .create()
                .throttleFirst(5L, TimeUnit.SECONDS)
                .subscribe(tweet -> {
                    tweetRepository.save(tweet);
                });
    }
}
