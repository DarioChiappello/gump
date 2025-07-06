package config

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/DarioChiappello/gump/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigWatcher(t *testing.T) {
	// Create temp folder for tests
	tempDir, err := os.MkdirTemp("", "config_watcher_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Helper function to create new config
	createConfigFile := func(name, content string) string {
		path := filepath.Join(tempDir, name)
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)
		return path
	}

	t.Run("Succesfully watcher created", func(t *testing.T) {
		filePath := createConfigFile("test1.json", `{}`)
		cfg := config.NewConfig()
		watcher, err := config.NewConfigWatcher(cfg, time.Minute, filePath)
		require.NoError(t, err)
		assert.NotNil(t, watcher)
		watcher.Stop()
	})

	t.Run("Changes detection through events", func(t *testing.T) {
		filePath := createConfigFile("event_test.json", `{"app": {"name": "GUMP"}}`)
		cfg := config.NewConfig()
		require.NoError(t, cfg.LoadFromJSON(filePath))

		watcher, err := config.NewConfigWatcher(cfg, time.Hour, filePath)
		require.NoError(t, err)

		reloadCh := make(chan bool, 1)
		watcher.OnReload(func(c *config.Config) {
			reloadCh <- true
		})

		go watcher.Start()
		defer watcher.Stop()

		time.Sleep(100 * time.Millisecond) // Wait for the watcher is ready
		err = os.WriteFile(filePath, []byte(`{"app": {"name": "GUMP_MODIFIED"}}`), 0644)
		require.NoError(t, err)

		select {
		case <-reloadCh:
			name, err := cfg.GetString("app.name")
			require.NoError(t, err)
			assert.Equal(t, "GUMP_MODIFIED", name)
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting event reload")
		}
	})

	t.Run("Watcher detection", func(t *testing.T) {
		filePath := createConfigFile("stop_test.json", `{}`)
		cfg := config.NewConfig()
		watcher, err := config.NewConfigWatcher(cfg, 100*time.Millisecond, filePath)
		require.NoError(t, err)

		stoppedCh := make(chan bool, 1)
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			watcher.Start()
			stoppedCh <- true
		}()

		time.Sleep(200 * time.Millisecond)
		watcher.Stop()

		select {
		case <-stoppedCh:
			// Success
		case <-time.After(500 * time.Millisecond):
			t.Fatal("Timeout waiting detection")
		}
	})

	t.Run("Events for file creation", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "new_config.json")
		cfg := config.NewConfig()
		watcher, err := config.NewConfigWatcher(cfg, 100*time.Millisecond, filePath)
		require.NoError(t, err)

		reloadCh := make(chan bool, 1)
		watcher.OnReload(func(c *config.Config) {
			reloadCh <- true
		})

		go watcher.Start()
		defer watcher.Stop()

		time.Sleep(100 * time.Millisecond)
		err = os.WriteFile(filePath, []byte(`{"new": "value"}`), 0644)
		require.NoError(t, err)

		select {
		case <-reloadCh:
			val, err := cfg.GetString("new")
			require.NoError(t, err)
			assert.Equal(t, "value", val)
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for creation detection")
		}
	})

}
