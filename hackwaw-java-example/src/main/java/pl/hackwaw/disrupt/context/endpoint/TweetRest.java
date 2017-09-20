package pl.hackwaw.disrupt.context.endpoint;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;
import pl.hackwaw.disrupt.context.domain.Tweet;
import pl.hackwaw.disrupt.context.domain.health.Health;
import pl.hackwaw.disrupt.context.domain.health.State;
import pl.hackwaw.disrupt.context.repository.TweetRepository;

import java.util.List;

@RestController
@RequestMapping
public class TweetRest {

    @Autowired
    private TweetRepository tweetRepository;

    @CrossOrigin
    @RequestMapping(value = "/latest", method = RequestMethod.GET)
    private List<Tweet> latest(Pageable pageRequest) {
        return tweetRepository.findAllByOrderByDateDesc(pageRequest).getContent();
    }

    @CrossOrigin
    @RequestMapping(value = "/health", method = RequestMethod.GET)
    private Health health() {
        return Health.builder()
                .app(State.OPERATIONAL)
                .database(State.UNNECESSARY)
                .slack(State.SLOW)
                .twitter(State.OPERATIONAL)
                .build();
    }

}
