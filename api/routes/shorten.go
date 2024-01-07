package routes

import (
	"time"
	"url-service/helpers"

	"github.com/asaskevich/govalidator"

	"github.com/gofiber/fiber"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

// * Rate Limit how many attempts left & RateLimitRest to reset the limit.
type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

//* Shorten URL Main logic function
//* Implement rate limiting : we will check whether IP is stored in our db or not. If IP is there in our db then user has used our service then we will have to decrement no of rate limiting by 1. We will allow him to call this API 10 times(Rate) over period 30 min. Means every 30 min it will reset.
//* Check URL is valid or not. We will use govalidator to validate our URL.
//* In domain err : we will not allow localhost url to avoid infinite loop
//* what is Enforcing HTTPS and SSL :
//* 1) The browser forces all communication over HTTPS. The browser prevents the user from using untrusted or invalid certificates.
//* 2) By default we Force SSL to all your visitors. This means if a site visitor loads your old, non-secure web address, http://yourdomainname.com, or they click an old non-secure link, they will be automatically redirected to the secure https://yourdomainname.com.

// ShortenURL ...
func ShortenURL(c *fiber.Ctx) error {
	// check for the incoming request body
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	// check if the input is an actual URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	// check for the domain error
	// users may abuse the shortener by shorting the domain `localhost:3000` itself
	// leading to a inifite loop, so don't accept the domain for shortening
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "haha... nice try",
		})
	}

	// enforce https
	// all url will be converted to https before storing in database
	body.URL = helpers.EnforceHTTP(body.URL)

}
