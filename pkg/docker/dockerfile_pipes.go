package docker

import "fmt"

func PHP(version string) Pipe {
	return func(d *Dockerfile) {
		d.Body += fmt.Sprintf("FROM php:%s\n", version)
		d.Body += `RUN apk --update add mysql-client curl git libxml2-dev libzip-dev libintl icu icu-dev zip \
    && docker-php-ext-install pdo_mysql zip xml opcache pcntl intl
`
	}
}

func ApkAdd(packages ...string) Pipe {
	return func(d *Dockerfile) {
		if len(packages) == 0 {
			return
		}

		d.Body += "RUN apk --update add"

		for _, p := range packages {
			d.Body += fmt.Sprintf(" %s", p)
		}
		d.Body += "\n"
	}
}

func Workdir(dir string) Pipe {
	return func(d *Dockerfile) {
		d.Body += fmt.Sprintln("WORKDIR " + dir)
	}
}

func Composer() Pipe {
	return func(d *Dockerfile) {
		d.Body += `ENV COMPOSER_HOME ./.composer
COPY --from=composer:1.9.3 /usr/bin/composer /usr/bin/composer

COPY composer.json composer.json
COPY composer.lock composer.lock

RUN composer install --no-dev --no-autoloader --no-scripts
`
	}
}

func ComposerAutoload(d *Dockerfile) {
	d.Body += fmt.Sprintln("RUN composer dump-autoload --optimize")
}

func Preload() Pipe {
	return func(d *Dockerfile) {
		d.Body += `ADD .breezedev/opcache.ini "$PHP_INI_DIR/conf.d/opcache.ini"
ADD .breezedev/preload.php preload.php
`
	}
}

type CopyOp struct {
	From  string
	To    string
	Chown string
}

func (op *CopyOp) String() string {
	s := "COPY "
	if op.Chown != "" {
		s += fmt.Sprintf("--chown=%s ", op.Chown)
	}

	s += fmt.Sprintf("%s %s", op.From, op.To)

	return s
}

type CopyOption func(op *CopyOp)

func Copy(from string, to string, opts ...CopyOption) Pipe {
	operation := &CopyOp{
		From: from,
		To:   to,
	}

	for _, opt := range opts {
		opt(operation)
	}

	return func(d *Dockerfile) {
		d.Body += fmt.Sprintln(operation.String())
	}
}

func Chown(user string, group string) CopyOption {
	return func(op *CopyOp) {
		op.Chown = fmt.Sprintf("%s:%s", user, group)
	}
}
