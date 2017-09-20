import java.time.Instant

import spray.json.{DefaultJsonProtocol, JsString, JsValue, JsonFormat, deserializationError}

import scala.util.Try

trait JsonProtocols extends DefaultJsonProtocol {
  implicit val appStateJsonFormat = new JsonFormat[AppState] {
    override def write(obj: AppState): JsValue = JsString(obj.toString)

    override def read(json: JsValue): AppState = ??? // unnecessary
  }
  implicit val databaseStateJsonFormat = new JsonFormat[DatabaseState] {
    override def write(obj: DatabaseState): JsValue = JsString(obj.toString)

    override def read(json: JsValue): DatabaseState = ??? // unnecessary
  }
  implicit val serviceStateJsonFormat = new JsonFormat[ServiceState] {
    override def write(obj: ServiceState): JsValue = JsString(obj.toString)

    override def read(json: JsValue): ServiceState = ??? // unnecessary
  }
  implicit val healthJsonFormat = jsonFormat4(Health.apply)

  implicit val instantJsonFormat = new JsonFormat[Instant] {
    override def read(json: JsValue): Instant = json match {
      case JsString(value) => Try(Instant.parse(value)).getOrElse(deserializationErrorMessage)
      case _ => deserializationErrorMessage
    }

    override def write(obj: Instant): JsValue = JsString(obj.toString)

    private def deserializationErrorMessage = deserializationError(s"Expecting date in ISO format, ex. ${Instant.now().toString}")
  }

  implicit val tweetJsonFormat = jsonFormat3(Tweet.apply)

  implicit val slackNotificationJsonFormat = jsonFormat5(SlackNotification.apply)
}
