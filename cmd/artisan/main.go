package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chenyusolar/aidecms/cmd/artisan/commands"
	"github.com/chenyusolar/aidecms/cmd/artisan/commands/generator"
	"github.com/chenyusolar/aidecms/cmd/artisan/stats"
	"github.com/chenyusolar/aidecms/config"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// Initialize database
	config.InitDB()

	// Check if running as CLI or web server
	if len(os.Args) > 1 && os.Args[1] == "artisan" {
		runArtisan()
		return
	}

	// Otherwise run as web server
	runWebServer()
}

func runArtisan() {
	startTime := time.Now()

	// Initialize database
	config.InitDB()

	// 确保日志目录存在
	if err := os.MkdirAll("storage/logs", 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		return
	}

	// 初始化日志
	logFile, err := os.OpenFile("storage/logs/artisan.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	// Skip "artisan" argument if present
	cmdArgs := os.Args[1:]
	if len(cmdArgs) > 0 && cmdArgs[0] == "artisan" {
		cmdArgs = cmdArgs[1:]
	}

	logger.Printf("[%s] Command started: %v",
		time.Now().Format("2006-01-02 15:04:05"),
		cmdArgs)

	if len(cmdArgs) < 1 {
		showHelp()
		return
	}

	command := cmdArgs[0]
	args := cmdArgs[1:]

	switch command {
	case "ai:setup", "ai:chat", "ai:completion", "ai:models", "ai:test", "ai:config":
		// 将 ai: 前缀的命令传递给AI命令处理器
		aiArgs := []string{strings.TrimPrefix(command, "ai:")}
		aiArgs = append(aiArgs, args...)
		commands.HandleAICommand(aiArgs)
	case "make:command":
		generator.NewCommand().Handle(args)
	case "make:controller":
		generator.NewCommand().Handle(args)
	case "make:model":
		generator.NewCommand().Handle(args)
	case "make:middleware":
		generator.NewCommand().Handle(args)
	case "migrate":
		commands.Migrate(args)
	case "help":
		showHelp()
	case "stats:show":
		commands.ShowStats(args)
	case "stats:reset":
		commands.ResetStats(args)
	case "stats:export":
		commands.ExportStats(args)
	case "stats:chart":
		stats.GenerateChart(commands.GetCommandStats())
	case "stats:cleanup":
		if len(args) > 0 {
			threshold, err := time.ParseDuration(args[0])
			if err != nil {
				fmt.Printf("Invalid duration format: %v\n", err)
				return
			}
			stats.CleanupOldStats(commands.GetCommandStats(), time.Now().Add(-threshold))
		} else {
			stats.CleanupOldStats(commands.GetCommandStats(), time.Now().Add(-30*24*time.Hour)) // Default: 30 days
		}
	case "stats:check":
		if len(args) > 0 {
			threshold, err := time.ParseDuration(args[0])
			if err != nil {
				fmt.Printf("Invalid duration format: %v\n", err)
				return
			}
			stats.CheckForAnomalies(commands.GetCommandStats(), threshold)
		} else {
			stats.CheckForAnomalies(commands.GetCommandStats(), 5*time.Second) // Default: 5 seconds
		}
	case "alert:setup":
		commands.SetupEmailAlert(args)
	case "alert:test":
		commands.SendTestEmail(args)
	case "queue:process":
		commands.ProcessQueue()
	case "queue:status":
		commands.ShowQueueStatus(args)
	case "queue:retry":
		commands.RetryFailedJobs(args)
	case "queue:clean":
		commands.CleanQueue(args)
	case "queue:priority":
		commands.SetPriority(args)
	case "queue:stats":
		commands.ShowQueueStats(args)
	case "schedule:work":
		commands.ScheduleWork(args)
	case "schedule:run":
		commands.ScheduleRun(args)
	case "schedule:list":
		commands.ScheduleList(args)
	case "event:test":
		commands.EventTest(args)
	case "event:list":
		commands.EventList(args)
	case "event:stats":
		commands.EventStats(args)
	case "ratelimit":
		commands.RateLimitCommand(args)
	case "health":
		commands.HealthCommand(args)
	case "web3":
		commands.Web3Command(args)
	case "exchange":
		commands.ExchangeCommandWrapper(args)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
	}

	duration := time.Since(startTime)
	logger.Printf("[%s] Command completed in %v",
		time.Now().Format("2006-01-02 15:04:05"),
		duration)

	// 记录命令使用统计
	commands.RecordCommandUsage(command, duration)
}

func runWebServer() {
	h := server.Default(server.WithHostPorts("0.0.0.0:8080"))

	// Register routes
	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.String(200, "Welcome to AideCMS")
	})

	// Start server
	h.Spin()
}

func showHelp() {
	fmt.Println("AideCMS Artisan Tool")
	fmt.Println("Usage: go run . artisan <command> [arguments]")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  make:command <name>\tCreate a new Artisan command")
	fmt.Println("  make:controller <name>\tCreate a new controller")
	fmt.Println("  make:model <name>\tCreate a new model")
	fmt.Println("  make:middleware <name>\tCreate a new middleware")
	fmt.Println("  migrate\t\tRun database migrations")
	fmt.Println("  help\t\t\tShow this help message")
	fmt.Println("\nAI commands:")
	fmt.Println("  ai:setup <provider> <api_key>\tSetup AI configuration")
	fmt.Println("  ai:chat <message> [model]\tChat with AI")
	fmt.Println("  ai:completion <prompt> [model]\tText completion")
	fmt.Println("  ai:models\t\t\tList available models")
	fmt.Println("  ai:test [model]\t\tTest AI connection")
	fmt.Println("  ai:config <action>\t\tManage AI configurations")
	fmt.Println("\nStatistics commands:")
	fmt.Println("  stats:show\t\tShow command usage statistics")
	fmt.Println("  stats:reset\t\tReset command statistics")
	fmt.Println("  stats:export <format>\tExport statistics (json/csv)")
	fmt.Println("  stats:chart\t\tGenerate usage chart")
	fmt.Println("  stats:cleanup\tClean up old statistics")
	fmt.Println("  stats:check\t\tCheck for performance anomalies")
	fmt.Println("\nAlert commands:")
	fmt.Println("  alert:setup <file>\tSetup email alert configuration")
	fmt.Println("  alert:test\t\tSend test email")
	fmt.Println("\nQueue commands:")
	fmt.Println("  queue:process\t\tProcess email queue")
	fmt.Println("  queue:status\t\tShow queue status")
	fmt.Println("  queue:retry\t\tRetry failed jobs")
	fmt.Println("  queue:clean\t\tClean old jobs")
	fmt.Println("  queue:priority\tSet job priority")
	fmt.Println("  queue:stats\t\tShow queue statistics")
	fmt.Println("\nSchedule commands:")
	fmt.Println("  schedule:work\t\tStart scheduler workers")
	fmt.Println("  schedule:run\t\tRun scheduled tasks once")
	fmt.Println("  schedule:list\t\tList scheduled tasks")
	fmt.Println("\nEvent commands:")
	fmt.Println("  event:test\t\tTest event system")
	fmt.Println("  event:list\t\tList registered events")
	fmt.Println("  event:stats\t\tShow event statistics")
	fmt.Println("\nRate Limit commands:")
	fmt.Println("  ratelimit demo\tRun rate limiting demonstration")
	fmt.Println("\nHealth Check commands:")
	fmt.Println("  health demo\t\tRun health check demonstration")
	fmt.Println("\nWeb3 commands:")
	fmt.Println("  web3 init\t\tInitialize Web3 clients")
	fmt.Println("  web3 chains\t\tList supported chains")
	fmt.Println("  web3 balance <chain> <address>\tGet address balance")
	fmt.Println("  web3 transaction <chain> <hash>\tGet transaction details")
	fmt.Println("  web3 block <chain>\tGet latest block number")
	fmt.Println("  web3 wallet <chain> <address>\tGet wallet info")
	fmt.Println("  web3 validate <chain> <address>\tValidate address format")
	fmt.Println("\nExchange commands:")
	fmt.Println("  exchange list\t\tList supported exchanges")
	fmt.Println("  exchange balance <exchange> <currency>\tGet balance")
	fmt.Println("  exchange balances <exchange>\tGet all balances")
	fmt.Println("  exchange price <exchange> <pair>\tGet trading pair price")
	fmt.Println("  exchange compare <pair>\tCompare prices across exchanges")
	fmt.Println("  exchange balance-all <currency>\tGet balance across all exchanges")
}
