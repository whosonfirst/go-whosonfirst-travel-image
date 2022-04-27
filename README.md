# go-whosonfirst-travel-image

Go package for traveling Who's On First documents and their relations, and making pictures

## Important

This should still be considered experimental. Documentation is incomplete.

## Tools

```
$> make cli
go build -mod vendor -o bin/wof-travel-id-image cmd/wof-travel-id-image/main.go
go build -mod vendor -o bin/wof-belongs-to-image cmd/wof-belongs-to-image/main.go
go build -mod vendor -o bin/wof-travel-filename cmd/wof-travel-filename/main.go
```

### wof-belongs-to-image

```
$> ./bin/wof-belongs-to-image \
	-labels \
	-html \
	-out ./belongs-to \
	-source fs:///usr/local/data/whosonfirst-data/data \
	-include-placetype neighbourhood \
	-belongs-to 85922583 \
	/usr/local/data/whosonfirst-data
```

Time passes and eventually something like this will be produced:

![](docs/images/20180813-sf-belongs-to.png)

_Why are some neighbourhoods rendered blank? I don't know yet..._

### wof-travel-id-image

```
$> ./bin/wof-travel-id-image \
	-source fs:///usr/local/data/whosonfirst-data/data \
	-superseded-by \
	-timings \
	-out test \
	-html \
	1108955789
	
2018/08/07 17:43:03 time to travel feature ID 85633755 182.121771ms
2018/08/07 17:43:03 time to travel feature ID 85632667 307.617697ms
2018/08/07 17:43:03 time to travel feature ID 1108955783 1.042050504s
2018/08/07 17:43:04 time to travel feature ID 85632609 109.83179ms
2018/08/07 17:43:04 time to travel feature ID 1108955785 1.570722402s
2018/08/07 17:43:04 time to travel feature ID 85632373 147.954674ms
2018/08/07 17:43:04 time to travel feature ID 1108955787 2.255208485s
2018/08/07 17:43:04 time to travel feature ID 85633779 117.256202ms
2018/08/07 17:43:04 time to travel feature ID 85633229 294.142877ms
2018/08/07 17:43:04 time to travel feature ID 1126113567 3.544154652s
2018/08/07 17:43:04 time to travel feature ID 1108955789 4.276219004s
```

Which produces this:

![](docs/images/20180807-yugoslavia.png)

All of the images and HTML output are included in [docs/example/1108955789](docs/example/1108955789)

Other examples:

The "20th" in Paris:

![](docs/images/20180807-20th.png)

Leningrad and St. Petersburg in Russia:

![](docs/images/20180807-leningrad-stpetersburg.png)

## See also

* https://github.com/whosonfirst/go-whosonfirst-travel