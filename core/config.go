package core

import (
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Build Information, overriden on build
var (
	Version   = "DEV"
	BuildDate = "DEV"
	StartTime time.Time
)

func initConfig() {
	if Version == "DEV" {
		logger = zerolog.New(zerolog.NewConsoleWriter()).
			With().
			Timestamp().
			Caller().
			Logger().
			Level(zerolog.DebugLevel)
	} else {
		logger = zerolog.New(zerolog.NewConsoleWriter()).
			With().
			Timestamp().
			Logger().
			Level(zerolog.InfoLevel)
	}

	home, err := homedir.Dir()

	if err != nil {
		logger.Fatal().Msg("Couldn't detect your home directory")
	} else {
		viper.SetConfigName(".thor")
		viper.SetConfigType("toml")
		viper.AddConfigPath(home)
	}

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal().Err(err).Msg("Couldn't read config file")
	}

	StartTime = time.Now().UTC()
	logger.Debug().Msg("Initialized")
}

// GetConfig gets a value from the KV database given its key, if present
func GetConfig(key string) (string, error) {
	db, err := badger.Open(
		badger.DefaultOptions(viper.GetString("core.settings_db")).
			WithLoggingLevel(badger.ERROR),
	)
	if err != nil {
		return "", err
	}

	defer db.Close()

	var value string
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		value = string(val)
		return nil
	})

	return value, err
}

// SetConfig sets a value in the KV database given its key, overwrites if present
func SetConfig(key, value string) error {
	db, err := badger.Open(
		badger.DefaultOptions(viper.GetString("core.settings_db")).
			WithLoggingLevel(badger.ERROR),
	)
	if err != nil {
		return err
	}

	defer db.Close()

	return db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		return err
	})
}

// GetKeyList returns all keys in the database
func GetKeyList() ([]string, error) {
	db, err := badger.Open(
		badger.DefaultOptions(viper.GetString("core.settings_db")).
			WithLoggingLevel(badger.ERROR),
	)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	keys := []string{}

	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			k := string(it.Item().Key())
			keys = append(keys, k)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return keys, nil
}
