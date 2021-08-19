# go-fiber-jwt-sample

A simple sample how to use jwt with [fiber](https://github.com/gofiber/fiber).

My problem was that I don't know how to set up the routes correctly.
There were only samples that use the standard middleware functions like:
```
app.use(middleware)
``` 

But in another project I need some different setup. Basically to set the middleware for each route separately.

So here is my setup right now: 

```
app := fiber.New()
	
app.Get("/user", authorizeRequired(), returnUser)

app.Get("/login", login)

app.Get("/hello", authorizeRequired(), hello)

err := app.Listen(":3000")
if err != nil {
	log.Println(err.Error())
}
```

The authorizeRequired function checks the token.
```
func authorizeRequired() func(ctx *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:    []byte(secret),
		SigningMethod: "HS512",
	})
}
```