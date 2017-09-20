package pl.hackwaw.disrupt.context.domain.health;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;

@Data
@Builder
@AllArgsConstructor
public class Health {
    State app;
    State database;
    State twitter;
    State slack;
}

