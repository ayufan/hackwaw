package pl.hackwaw.disrupt.context.endpoint;

import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;
import pl.hackwaw.disrupt.context.domain.SlackProxyRequest;
import pl.hackwaw.disrupt.context.domain.SlackRequest;

@RestController
@RequestMapping
@Slf4j
public class TweetRest {

    private RestTemplate template = new RestTemplate();

    @CrossOrigin
    @RequestMapping(value = "/push", method = RequestMethod.POST)
    private void push(@RequestBody SlackProxyRequest proxyRequest) {
        SlackRequest slackRequest = SlackRequest.fromProxy(proxyRequest);
        log.info("Pushing to slack", slackRequest);
        String responseEntity = template.postForObject("https://hooks.slack.com/services/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", slackRequest, String.class);
        log.info("Slack response:", responseEntity);
    }

}
