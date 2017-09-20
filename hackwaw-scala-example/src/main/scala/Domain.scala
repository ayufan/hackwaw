import java.time.Instant

sealed trait AppState
object AppState {
  case object Operational extends AppState { override def toString = "OPERATIONAL" }
  case object Slow extends AppState { override def toString = "SLOW" }
  case object Error extends AppState { override def toString = "ERROR" }
}

sealed trait DatabaseState
object DatabaseState {
  case object Operational extends DatabaseState { override def toString = "OPERATIONAL" }
  case object Down extends DatabaseState { override def toString = "DOWN" }
  case object Slow extends DatabaseState { override def toString = "SLOW" }
  case object Error extends DatabaseState { override def toString = "ERROR" }
  case object Unnecessary extends DatabaseState { override def toString = "UNNECESSARY" }
}

sealed trait ServiceState
object ServiceState {
  case object Operational extends ServiceState { override def toString = "OPERATIONAL" }
  case object Down extends ServiceState { override def toString = "DOWN" }
  case object Slow extends ServiceState { override def toString = "SLOW" }
  case object Error extends ServiceState { override def toString = "ERROR" }
}

case class Health(app: AppState, database: DatabaseState, twitter: ServiceState, slack: ServiceState)

case class Tweet(id: Long, body: String, date: Instant)

case class SlackNotification(team: String, tweetId: Long, `icon_url`: String, text: String, date: String)

case object SlackNotification {
  def apply(tweet: Tweet): SlackNotification = {
    SlackNotification("Scala4Ever", tweet.id, "http://fruzenshtein.com/wp-content/uploads/2013/09/scala-logo.png", tweet.body,tweet.date.toString)
  }
}