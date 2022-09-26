/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/cobra"
)

var Config string
var cfg ConfigDatabase

type ConfigDatabase struct {
	AllowedUserList []string `yaml:"AllowedUserList" env:"ALLOWED_USER_LIST" env-separator:","`
	ChannelName     string   `yaml:"ChannelName" env:"CHANNEL_NAME"`
	Token           string   `yaml:"Token" env:"O_AUTH_TOKEN"`
	StartText       string   `yaml:"StartText"`
	EndText         string   `yaml:"EndText"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "twitch-bot",
	Short: "A twitch rating bot",
	Long:  `A twitch rating bot.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.twitch-bot.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
func initConfig() {
	err := cleanenv.ReadConfig(Config, &cfg)
	if err != nil {
		panic(err)
	}
}
