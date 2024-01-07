package routes

import (
	"url-service/database"

	"github.com/gofiber/fiber"
	"github.com/redis/go-redis/v9"
)

//* Once we use Shorten() function to shorten the URL which we will give, But the shorten url that we are going to use it should redirect to the actual original link. So we will save that long original link into our database, and creating short url link,  and when somebody uses this shortened URL link  and then we will take user to that actual URL link stored in our database. So, that's our RESOLVE() function . Means after creating shortened URL we will also have to RESOLVE() that url.

//* If we didn't find the URL in redis db then we will return url not found, But if everything went very well we will just pass the value by redirecting

// ResolveURL ...
func ResolveURL(c *fiber.Ctx) error {
	// get the short from the url
	url := c.Params("url")
	// query the db to find the original URL, if a match is found
	// increment the redirect counter and redirect to the original URL
	// else return error message
	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short not found on database",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to DB",
		})
	}
	// increment the counter
	rInr := database.CreateClient(1)
	defer rInr.Close()
	_ = rInr.Incr(database.Ctx, "counter")
	// redirect to original URL
	return c.Redirect(value, 301)

}
