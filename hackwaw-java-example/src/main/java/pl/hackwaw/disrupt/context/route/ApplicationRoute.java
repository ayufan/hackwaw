package pl.hackwaw.disrupt.context.route;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import pl.hackwaw.disrupt.context.domain.Tweet;
import pl.hackwaw.disrupt.context.repository.SlackProxyResource;
import pl.hackwaw.disrupt.context.repository.TweetRepository;
import pl.hackwaw.disrupt.context.repository.TwitterProxyResource;

import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.util.List;
import java.util.stream.Collectors;

@Slf4j
@Component
public class ApplicationRoute {

    @Autowired
    private TweetRepository tweetRepository;

    @Autowired
    private TwitterProxyResource twitterProxyResource;

    @Autowired
    private SlackProxyResource slackProxyResource;

    @Scheduled(fixedDelay = 5000, initialDelay = 5000)
    public void doThings() {
        List<Tweet> newTweets = twitterProxyResource.since(latestTweetDate()).stream().map(Tweet::fromProxy).collect(Collectors.toList());
        tweetRepository.save(newTweets);
        newTweets.forEach(it -> slackProxyResource.push(it));
    }

    private LocalDateTime latestTweetDate() {
        return tweetRepository
                .findFirstByOrderByDateDesc()
                .map(Tweet::getDate)
                .map(e -> LocalDateTime.ofInstant(e, ZoneOffset.UTC))
                .orElse(LocalDateTime.now(ZoneOffset.UTC).minusSeconds(60));
    }
}
