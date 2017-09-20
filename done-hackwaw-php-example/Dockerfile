FROM php:7.0.4

# Install required dependencies
RUN apt-get update && apt-get install -y \
	git-core \
	unzip

# Install composer
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer

# Add description of application dependencies
# We add this first, in order to speed up subsequent builds
ADD composer.* /hackwaw/php/
WORKDIR /hackwaw/php/
RUN ["composer", "install"]

# Add application sources
ADD . /hackwaw/php/

CMD ["php", "app/app.php"]
EXPOSE 8080