package played

import (
	Provider "GameApplication/entity/entityProvider"
)

func getProvider(db string) Provider.Provider {
	return Provider.Provider{db}
}
