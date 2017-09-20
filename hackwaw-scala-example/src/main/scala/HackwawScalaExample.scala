import java.time.Instant

import akka.actor.ActorSystem
import akka.event.Logging
import akka.http.scaladsl.Http
import akka.http.scaladsl.client.RequestBuilding
import akka.http.scaladsl.marshallers.sprayjson.SprayJsonSupport._
import akka.http.scaladsl.model.Uri
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.unmarshalling.Unmarshal
import akka.stream.ActorMaterializer
import com.typesafe.config.ConfigFactory

import scala.concurrent.Future
import scala.concurrent.duration._
import scala.language.postfixOps

object HackwawScalaExample extends App with JsonProtocols {
  implicit val system = ActorSystem()
  implicit val executor = system.dispatcher
  implicit val materializer = ActorMaterializer()

  val config = ConfigFactory.load()
  val logger = Logging(system, getClass)

  val twitterProxyUrl = config.getString("services.twitter-url")
  val slackProxyUrl = config.getString("services.slack-url")

  @volatile var latestTweets: Seq[Tweet] = Seq.empty[Tweet]

  def fetchTweets(from: Instant, to: Instant): Future[Seq[Tweet]] = {
    val params = Uri.Query("from" -> from.toString, "to" -> to.toString)
    val request = RequestBuilding.Get(Uri(s"$twitterProxyUrl/tweets").withQuery(params))
    Http().singleRequest(request).flatMap { response =>
      Unmarshal(response).to[Seq[Tweet]]
    }
  }

  def pushSlackNotifications(notifications: Seq[SlackNotification]): Future[Seq[String]] = {
    Future.sequence(notifications.map { notification =>
      val request = RequestBuilding.Post(s"$slackProxyUrl/push", notification)
      Http().singleRequest(request).flatMap { response =>
        Unmarshal(request).to[String]
      }
    })
  }

  system.scheduler.schedule(5 seconds, 5 seconds) {
    val from = latestTweets.sortBy(_.date).lastOption.map(_.date).getOrElse(Instant.now().minusSeconds(60))
    logger.info("Fetching tweets from {} to now", from)

    fetchTweets(from, Instant.now()).map { tweets =>
      logger.info("Got tweets {}", tweets)
      latestTweets = tweets
      val slackNotifications = tweets.map(SlackNotification.apply)
      logger.info("Pushing notifications {}", slackNotifications)

      pushSlackNotifications(slackNotifications).map { responses =>
        logger.info("Got responses ({}) {}", responses.length, responses)
      }
    }
  }

  val routes = {
    val PageSize = 50

    (get & path("health") & pathEndOrSingleSlash) {
      complete {
        Health(AppState.Operational, DatabaseState.Unnecessary, ServiceState.Operational, ServiceState.Operational)
      }
    } ~
    (get & path("latest") & pathEndOrSingleSlash & parameter('page.as[Int].?)) { pageOpt =>
      complete {
        val page = pageOpt.getOrElse(0)
        latestTweets.slice(page * PageSize, (page + 1) * PageSize).map(SlackNotification.apply)
      }
    }
  }

  Http().bindAndHandle(routes, config.getString("http.interface"), config.getInt("http.port"))
}
