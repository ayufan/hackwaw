package pl.hackwaw.disrupt.context.domain;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import twitter4j.Status;

import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.Id;
import java.time.Instant;

@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
@Entity
public class Tweet {

    @Id
    private Long id;

    @Column(nullable = false)
    private String body;

    @Column(nullable = false)
    private Instant date;

    public static Tweet fromStatus(Status status) {
        return Tweet.builder()
                .body(status.getText())
                .id(status.getId())
                .date(status.getCreatedAt().toInstant())
                .build();
    }
}
