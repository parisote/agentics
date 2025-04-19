package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/parisote/agentics/agentics"
)

type weather struct {
	Current struct {
		Temp float64 `json:"temp_c"`
	}
}

func init() { agentics.RegisterHook("fetchWeather", fetchWeather) }

func fetchWeather(ctx context.Context, c *agentics.Context) error {
	key := os.Getenv("WEATHER_API_KEY")
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=Buenos+Aires", key)
	r, err := c.HTTP.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	var w weather
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	c.Bag.Set("weather_c", w.Current.Temp)
	return nil
}
