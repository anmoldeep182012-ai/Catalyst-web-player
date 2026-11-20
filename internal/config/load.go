/*
 * TgMusicBot - Telegram Music Bot
 * Copyright (c) 2025-2026 Ashok Shau
 *
 * Licensed under GNU GPL v3
 * See https://github.com/AshokShau/TgMusicBot
 */

package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

const (
	defaultDBName         = "CatalystWebApp"
	defaultAPIURL         = "https://api.onegrab.fun"
	defaultDownloadsDir   = "downloads"
	defaultSupportGroup   = "https://t.me/CatalystxSupport"
	defaultSupportChannel = "https://t.me/CatalystxSupport"
	defaultStartImage     = "https://raw.githubusercontent.com/anmoldeep182012-ai/Catalyst-web-player/master/internal/webapp/static/assets/catalyst-logo.jpg"
	defaultMaxFileSize    = int64(500 * 1024 * 1024)
	defaultSongDuration   = int32(3600)
)

var (
	ApiId      = getEnvInt32("API_ID", 0)
	ApiHash    = os.Getenv("API_HASH")
	Token      = os.Getenv("TOKEN")
	DlBotToken = os.Getenv("DL_BOT_TOKEN")

	MongoUri = os.Getenv("MONGO_URI")
	DbName   = getEnv("DB_NAME", defaultDBName)

	ApiUrl = getEnv("API_URL", defaultAPIURL)
	ApiKey = os.Getenv("API_KEY")

	OwnerId  = getEnvInt64("OWNER_ID", 0)
	LoggerId = getEnvInt64("LOGGER_ID", 0)

	Proxy             = os.Getenv("PROXY")
	MaxFileSize       = getEnvInt64("MAX_FILE_SIZE", defaultMaxFileSize)
	SongDurationLimit = getEnvInt32("SONG_DURATION_LIMIT", defaultSongDuration)
	DownloadsDir      = getEnv("DOWNLOADS_DIR", defaultDownloadsDir)

	SupportGroup     = getEnv("SUPPORT_GROUP", defaultSupportGroup)
	SupportChannel   = getEnv("SUPPORT_CHANNEL", defaultSupportChannel)
	StartImg         = getEnv("START_IMG", defaultStartImage)
	PrivacyPolicyURL = getEnv("PRIVACY_POLICY_URL", "")

	Port = getEnv("PORT", "6060")
	WebRTCUDPPortMin = getEnvInt32("WEBRTC_UDP_PORT_MIN", 40000)
	WebRTCUDPPortMax = getEnvInt32("WEBRTC_UDP_PORT_MAX", 49999)
	BotCommandsEnabled = getEnvBool("BOT_COMMANDS_ENABLED", false)

	DEVS        []int64
	CookiesPath []string
	cookiesUrl  = processCookieURLs(os.Getenv("COOKIES_URL"))
)

func LoadEnv() error {
	loadDevelopers()

	if OwnerId != 0 && !slices.Contains(DEVS, OwnerId) {
		DEVS = append(DEVS, OwnerId)
	}

	if err := validate(); err != nil {
		return fmt.Errorf("config is invalid: %w", err)
	}

	if err := os.MkdirAll(DownloadsDir, 0o755); err != nil {
		return errors.New("failed to create downloads directory")
	}

	if len(cookiesUrl) == 0 {
		return nil
	}

	if err := os.MkdirAll(cookiesDr, 0o750); err != nil {
		return errors.New("failed to create cookies directory")
	}

	go saveAllCookies(cookiesUrl)
	return nil
}

func validate() error {
	required := []struct {
		name  string
		valid bool
	}{
		{"API_ID", ApiId > 0},
		{"API_HASH", ApiHash != ""},
		{"TOKEN", Token != ""},
		{"MONGO_URI", MongoUri != ""},
		{"OWNER_ID", OwnerId > 0},
	}

	var missing []string

	for _, field := range required {
		if !field.valid {
			missing = append(missing, field.name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf(
			"missing required configuration: %s",
			strings.Join(missing, ", "),
		)
	}

	if MaxFileSize <= 0 {
		return fmt.Errorf("MAX_FILE_SIZE must be greater than 0")
	}

	if SongDurationLimit <= 0 {
		return fmt.Errorf("SONG_DURATION_LIMIT must be greater than 0")
	}

	if WebRTCUDPPortMin < 1 || WebRTCUDPPortMax > 65535 || WebRTCUDPPortMax < WebRTCUDPPortMin {
		return fmt.Errorf("WEBRTC_UDP_PORT_MIN/MAX must be a valid UDP port range (1-65535)")
	}

	return nil
}
