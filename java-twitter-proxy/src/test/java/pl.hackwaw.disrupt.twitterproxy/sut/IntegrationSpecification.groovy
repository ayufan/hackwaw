package pl.hackwaw.disrupt.twitterproxy.sut

import groovy.util.logging.Slf4j
import org.springframework.beans.factory.annotation.Value
import org.springframework.boot.test.SpringApplicationConfiguration
import org.springframework.boot.test.WebIntegrationTest
import pl.hackwaw.disrupt.Application
import spock.lang.Specification

@Slf4j
@SpringApplicationConfiguration(Application)
@WebIntegrationTest('server.port:0')
class IntegrationSpecification extends Specification {

    @Value("\${local.server.port}")
    int serverPort

    Sut sut

    def setup() {
        log.info('App server on port {}', serverPort)
        sut = new Sut(serverPort)
    }
}
