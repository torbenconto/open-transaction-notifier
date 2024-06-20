package main

type Config struct {
	OpenTransactionNotifier struct {
		NotificationMethod string `mapstructure:"NOTIFICATION_METHOD"`
		NotificationInfo   string `mapstructure:"NOTIFICATION_INFO"`
		TimeInterval       string `mapstructure:"TIME_INTERVAL"`
	} `mapstructure:"open_transaction_notifier"`
	Discord struct {
		Message string `mapstructure:"MESSAGE"`
		NoEmbed struct {
			Message string `mapstructure:"MESSAGE"`
		} `mapstructure:"noembed"`
		Embed struct {
			Enabled       bool   `mapstructure:"ENABLED"`
			Title         string `mapstructure:"TITLE"`
			PriceDecimals int    `mapstructure:"PRICE_DECIMALS"`
			Color         string `mapstructure:"COLOR"`
			Fields        struct {
				Ticker        bool `mapstructure:"TICKER"`
				Type          bool `mapstructure:"TYPE"`
				PricePerShare bool `mapstructure:"PRICE_PER_SHARE"`
				Shares        bool `mapstructure:"SHARES"`
				Owner         bool `mapstructure:"OWNER"`
				Date          bool `mapstructure:"DATE"`
				Relationship  bool `mapstructure:"RELATIONSHIP"`
			} `mapstructure:"fields"`
		} `mapstructure:"embed"`
	} `mapstructure:"discord"`
}
