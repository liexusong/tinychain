package p2p

import (
	ma "github.com/multiformats/go-multiaddr"
	"github.com/spf13/viper"
	"github.com/libp2p/go-libp2p-crypto"
)



type Config struct {
	seeds         []ma.Multiaddr // Seed nodes for initialization
	routeFilePath string         // Store route table
	privKey       crypto.PrivKey // Private key of peer
	port          int            // Listener port
}

func LoadConfigFromFile(path string, configName string) (*Config, error) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Cannot load config file: %s", err)
		return nil, err
	}

	var seedAddrs []ma.Multiaddr
	for _, addr := range viper.GetStringSlice("p2p.seeds") {
		multiAddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			log.Fatalf("Invalid multiaddr in config file:%s", err)
			return nil, err
		}
		seedAddrs = append(seedAddrs, multiAddr)
	}

	privKey, err := B64DecodePrivKey(viper.GetString("network.privkey"))
	if err != nil {
		log.Fatalf("Invalid private key:%s", err)
		return nil, err
	}

	routeFilePath := viper.GetString("p2p.routeFile")
	port := viper.GetInt("network.port")

	return &Config{
		seeds:         seedAddrs,
		routeFilePath: routeFilePath,
		privKey:       privKey,
		port:          port,
	}, nil
}
