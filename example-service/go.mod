module example

go 1.24.0

replace github.com/korableg/space-307-meetup/lib => ../lib

replace github.com/korableg/space-307-meetup/db => ../db

require github.com/korableg/space-307-meetup/lib v0.0.0-00010101000000-000000000000

require (
	github.com/burntcarrot/heaputil v0.0.0-20230927162808-497024fb706a // indirect
	github.com/korableg/space-307-meetup/db v0.0.0-00010101000000-000000000000 // indirect
)
