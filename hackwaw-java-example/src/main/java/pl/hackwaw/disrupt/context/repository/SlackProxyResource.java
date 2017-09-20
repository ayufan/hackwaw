package pl.hackwaw.disrupt.context.repository;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
import pl.hackwaw.disrupt.context.domain.ProxySlackTweet;
import pl.hackwaw.disrupt.context.domain.Tweet;

import static java.lang.String.format;

@Slf4j
@Service
public class SlackProxyResource {

    private RestTemplate template = new RestTemplate();

    @Value("${SLACK_URL}")
    private String slackUrl;

    public void push(Tweet tweet) {
        log.info("Pushing to slack", tweet);
        String responseEntity = template.postForObject(format("%s/push", slackUrl), new ProxySlackTweet(tweet), String.class);
        log.info("Slack response:", responseEntity);
    }

}

