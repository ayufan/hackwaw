package pl.hackwaw.disrupt.context.domain;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;

@Getter
@NoArgsConstructor
@AllArgsConstructor
public class SlackProxyRequest {

    private String team;
    private String tweetId;
    private String icon_url;
    private String text;
    private String date;

}
