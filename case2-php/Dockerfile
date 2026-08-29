FROM php:8.4-alpine

WORKDIR /var/www

# Install system dependencies & extensions
RUN apk add --no-cache \
    curl \
    libzip-dev \
    oniguruma-dev \
    libpng-dev \
    icu-dev

RUN docker-php-ext-install \
    pdo_mysql \
    mbstring \
    zip \
    bcmath \
    intl

# Install Composer
COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

# Copy application files
COPY . .

EXPOSE 8000

CMD ["php", "artisan", "serve", "--host=0.0.0.0", "--port=8000"]
