package domain

type Config struct {
	ServerPort             string
	JWTSigningKey          []byte
	JWTTokenExpirationHour int
	DBConnectString        string
	DataStorePath          string
}

func (c *Config) Default() {
	c.ServerPort = ":8080"
	c.JWTSigningKey = []byte("SecretYouShouldHide")
	c.JWTTokenExpirationHour = 2
	c.DBConnectString = "postgres://postgres@db:5432/elotus?sslmode=disable"
	c.DataStorePath = "./tmp"
}
