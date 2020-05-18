package docker

type Dockerfile struct {
	Body string
}

type Pipe func(d *Dockerfile)

// dockerfile := &Dockerfile{}
// dockerfile.Build(
//     PHP("7.4-fpm-alpine"),
//     Composer(),
//     Copy(".", "/var/www/html", Chown("www-data:www-data")),
//     Preload(),
//
// )
func (d *Dockerfile) Build(pipes ...Pipe) *Dockerfile {
	for _, p := range pipes {
		p(d)
	}

	return d
}

func (d *Dockerfile) String() string {
	return d.Body
}
