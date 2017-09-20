package pl.hackwaw.disrupt;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import springfox.documentation.builders.ApiInfoBuilder;
import springfox.documentation.builders.RequestHandlerSelectors;
import springfox.documentation.service.ApiInfo;
import springfox.documentation.spi.DocumentationType;
import springfox.documentation.spring.web.plugins.Docket;
import springfox.documentation.swagger2.annotations.EnableSwagger2;
import static com.google.common.base.Predicates.not;

@Configuration
@EnableSwagger2
public class SpringFoxConfiguration {

    @Bean
    public Docket springFoxConfig() {
        return new Docket(DocumentationType.SWAGGER_2)
                .groupName("API")
                .select()
                .apis(RequestHandlerSelectors.basePackage("pl.hackwaw"))
                .build()
                .apiInfo(metadata("API"));
    }

    private ApiInfo metadata(String name) {
        return new ApiInfoBuilder()
                .title(name)
                .build();
    }
}