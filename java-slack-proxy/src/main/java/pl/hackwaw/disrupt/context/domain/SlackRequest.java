package pl.hackwaw.disrupt.context.domain;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;

@Getter
@NoArgsConstructor
@AllArgsConstructor
public class SlackRequest {

    private String username;
    private String icon_url;
    private String text;

    public static SlackRequest fromProxy(SlackProxyRequest proxyRequest) {
        String slackMessageBody = String.format("Tweet with id '%s', at '%s' with message '%s'", proxyRequest.getTweetId(), proxyRequest.getDate(), proxyRequest.getText());
        return new SlackRequest(proxyRequest.getTeam(), proxyRequest.getIcon_url(), slackMessageBody);
    }

}
