module github.com/korableg/space-307-meetup/lib

go 1.23.5

replace github.com/korableg/space-307-meetup/db => ../db

require github.com/korableg/space-307-meetup/db v0.0.0-00010101000000-000000000000

require github.com/burntcarrot/heaputil v0.0.0-20230927162808-497024fb706a // indirect

exclude github.com/burntcarrot/heaputil v1.0.0