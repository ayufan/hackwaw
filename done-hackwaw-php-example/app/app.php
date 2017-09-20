<?php

use React\Http\Request;
use React\Http\Response;
use React\HttpClient\Response as HttpResponse;
use React\HttpClient\Request as HttpRequest;

require_once __DIR__ . '/../vendor/autoload.php';

$twitterUrl = getenv('TWITTER_URL');
$slackUrl   = getenv('SLACK_URL');

define('STATUS_OPERATIONAL', 'OPERATIONAL');
define('STATUS_SLOW', 'SLOW');
define('STATUS_ERROR', 'ERROR');
define('STATUS_DOWN', 'DOWN');
define('STATUS_UNNECESSARY', 'UNNECESSARY');


$health          = new stdClass();
$health->app     = STATUS_OPERATIONAL;
$health->twitter = STATUS_OPERATIONAL;
$health->slack   = STATUS_OPERATIONAL;

define('TWITTER_API_DATE_FORMAT', 'Y-m-d\TH:i:s.u\Z');
define('TEAM_ID', 'hackwaw-php-example');

$messages = new SplDoublyLinkedList();

$loop   = React\EventLoop\Factory::create();
$socket = new React\Socket\Server($loop);

$dispatcher = FastRoute\simpleDispatcher(function (FastRoute\RouteCollector $r) use ($messages, $health)
{
	$r->addRoute('GET', '/health', function (Request $request, Response $response) use ($health)
	{
		$response->writeHead(200, ['Content-Type' => 'application/json']);
		$response->end(json_encode($health));
	});

	$r->addRoute('GET', '/latest', function (Request $request, Response $response) use ($messages)
	{
		$response->writeHead(200, ['Content-Type' => 'application/json']);
		$response->end(json_encode(iterator_to_array($messages)));
	});
});

$http = new React\Http\Server($socket);
$http->on('request', function (Request $request, Response $response) use ($dispatcher)
{
	$routeInfo = $dispatcher->dispatch($request->getMethod(), $request->getPath());
	switch ($routeInfo[0])
	{
		case FastRoute\Dispatcher::NOT_FOUND:
			$response->writeHead(404);
			$response->end('Not found');
			break;
		case FastRoute\Dispatcher::METHOD_NOT_ALLOWED:
			$response->writeHead(405);
			$response->end('Method not allowed');
			break;
		case FastRoute\Dispatcher::FOUND:
			call_user_func_array($routeInfo[1], [$request, $response]);
			break;
	}
});

$dnsResolverFactory = new React\Dns\Resolver\Factory();
$dnsResolver        = $dnsResolverFactory->createCached('8.8.8.8', $loop);

$factory       = new React\HttpClient\Factory();
$twitterClient = $factory->create($loop, $dnsResolver);
$slackClient   = $factory->create($loop, $dnsResolver);

$loop->addPeriodicTimer(1, function () use ($twitterUrl, $slackUrl, $twitterClient, $slackClient, $messages, $health)
{
	if (empty($last))
	{
		$last = DateTimeImmutable::createFromFormat('U.u', microtime(true))->modify('-1 hour');
	}

	$now = DateTimeImmutable::createFromFormat('U.u', microtime(true));

	$twitterEndpoint = sprintf('%s/tweets?from=%s&to=%s', $twitterUrl, $last->format(TWITTER_API_DATE_FORMAT), $now->format(TWITTER_API_DATE_FORMAT));

	$last = $now;

	$request = $twitterClient->request('GET', $twitterEndpoint);
	$request->on('response', function (HttpResponse $response) use ($messages, $slackClient, $slackUrl, $health)
	{
		$response->on('data', function ($data, HttpResponse $response) use ($messages, $slackClient, $slackUrl, $health)
		{
			$tweets = json_decode($data, false, 512, JSON_UNESCAPED_UNICODE);
			if (null === $tweets)
			{
				$health->twitter = STATUS_ERROR;
			}

			if (is_array($tweets))
			{
				$health->twitter = STATUS_OPERATIONAL;
				foreach ($tweets as $tweet)
				{
					$messages->unshift($tweet);
					$slackRequest = $slackClient->request('POST', $slackUrl.'/push', [
						'Content-Type' => 'application/json',

					]);

					$slackRequest->write(json_encode([
						'team'     => TEAM_ID,
						'tweetId'  => $tweet->id,
						'icon_url' => 'http://www.veryicon.com/icon/ico/System/Arcade%20Daze/Mario.ico',
						'text'     => $tweet->body,
						'date'     => $tweet->date,
					]));

					$slackRequest->on('end', function ($error, HttpResponse $response = null) use ($health)
					{
						if ($error || null === $response)
						{
							$health->slack = STATUS_ERROR;
						}
						else if ($response->getCode() != 200)
						{
							$health->slack = STATUS_ERROR;
						}
						else
						{
							$health->slack = STATUS_OPERATIONAL;
						}
					});

					$slackRequest->end();
				}
			}
		});
	});
	$request->end();
});

$socket->listen(8080, '0.0.0.0');
$loop->run();
