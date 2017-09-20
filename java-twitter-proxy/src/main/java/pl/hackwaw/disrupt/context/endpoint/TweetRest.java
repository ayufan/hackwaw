package pl.hackwaw.disrupt.context.endpoint;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Pageable;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import pl.hackwaw.disrupt.context.domain.Tweet;
import pl.hackwaw.disrupt.context.repository.TweetRepository;

import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.util.List;
import java.util.stream.Collectors;

import static com.google.common.base.Preconditions.checkNotNull;

@RestController
@RequestMapping
public class TweetRest {

    @Autowired
    private TweetRepository tweetRepository;

    @CrossOrigin
    @RequestMapping(value = "/tweets", method = RequestMethod.GET)
    private List<Tweet> tweets(@RequestParam("from") @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME) LocalDateTime from,
                               @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME) @RequestParam("to") LocalDateTime to,
                               Pageable pageRequest
    ) {
        return tweetRepository.findAllByDateBetweenOrderByDateDesc(checkNotNull(from).toInstant(ZoneOffset.UTC), checkNotNull(to).toInstant(ZoneOffset.UTC), pageRequest).collect(Collectors.toList());
    }

}
