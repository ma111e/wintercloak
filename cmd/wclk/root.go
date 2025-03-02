package main

import (
	"fmt"
	"github.com/gookit/color"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/fs"
	"path/filepath"
	"strings"
	"wintercloak/internal/datasources"
	"wintercloak/internal/engines"
	"wintercloak/internal/utils/engineutils"
)

var (
	RootCmd = &cobra.Command{
		Use:   "wclk",
		Short: "Wintercloak preprocessor",
		PreRun: func(cmd *cobra.Command, _ []string) {
			// godotenv's autoload is used
			viper.SetEnvPrefix("WCLK")
			viper.AutomaticEnv()

			viper.BindPFlag("patch-empty", cmd.Flags().Lookup("patch-empty"))
			viper.BindEnv("patch-empty")
			viper.BindPFlag("list-data", cmd.Flags().Lookup("list-data"))
			viper.BindEnv("list-data")
			viper.BindPFlag("list-engines", cmd.Flags().Lookup("list-engines"))
			viper.BindEnv("list-engines")
			viper.BindPFlag("engine", cmd.Flags().Lookup("engine"))
			viper.BindEnv("engine")
			viper.BindPFlag("file", cmd.Flags().Lookup("file"))
			viper.BindEnv("file")
			viper.BindPFlag("marker", cmd.Flags().Lookup("marker"))
			viper.BindEnv("marker")
			viper.BindPFlag("write", cmd.Flags().Lookup("write"))
			viper.BindEnv("write")

			viper.BindPFlag("option", cmd.Flags().Lookup("option"))
		},
	}

	rawEngineCLIArgs CLIDataMapping = make(map[string]string)
	rawRunCLIArgs    CLIDataMapping = make(map[string]string)
)

func init() {
	RootCmd.Run = Run
	RootCmd.Flags().BoolP("list-data", "L", false, "List all available files that can be used as data sources")
	RootCmd.Flags().BoolP("list-engines", "E", false, "List all available obfuscation engines")
	RootCmd.Flags().BoolP("write", "w", false, "Write the generated file to disk instead of sending it to stdout. The generated filename is handled by the engine; however the encouraged behaviour is to generate a xxx_gen.yy file next to the source")
	RootCmd.Flags().StringP("engine", "e", "", "Obfuscator engine to use")
	RootCmd.Flags().StringP("file", "f", "", "Source file")
	RootCmd.Flags().StringP("marker", "m", "", "Patch marker to use to target specific string with obfuscation (e.g. '-m PATCH' => will replace myvar = PATCH('mystring') with myvar = <obfuscated mystring value>")

	RootCmd.MarkFlagRequired("engine")

	RootCmd.Flags().VarP(&rawEngineCLIArgs, "option", "o", "Pass options to the engine. This flag can be repeated and accepts key-value pairs in key=value format (e.g. '-o 'pattern=??' -o offset=7'")
	RootCmd.Flags().VarP(&rawRunCLIArgs, "run", "r", "Set parameters for the run. This flag can be repeated and accepts key-value pairs in key=value format (e.g. '-r useGoBuildTag=true -r goBuildTag=obfuscate'")
}

func Run(_ *cobra.Command, _ []string) {
	engineName := viper.GetString("engine")
	listData := viper.GetBool("list-data")
	listEngines := viper.GetBool("list-engines")

	if listData {
		fmt.Println("Available data sources:")
		dirs, err := datasources.ListDataSources()
		if err != nil {
			log.Fatalln(err)
		}

		for _, path := range dirs {
			fmt.Printf("\t* %s\n", path.Name())
		}
		return
	}

	if listEngines {
		fmt.Println("Available obfuscation engines:")
		available := engines.Available.List()
		for _, eng := range available {
			desc := eng.ShortDescription
			if len(eng.LongDescription) > 0 {
				desc = eng.LongDescription
			}

			fmt.Printf("\t* %s\n", eng.Name)
			fmt.Printf("\t\t%s\n", desc)
			fmt.Printf("\t\tOptions:\n")
			for _, arg := range eng.Args {
				fmt.Printf("\t\t\t+ '%s': %s\n", arg.Name, arg.Description)
			}
		}
		return
	}

	runOptions := engines.NewRunOptionsFromCLI(rawRunCLIArgs)

	// Marker & file shorthand
	if viper.IsSet("marker") {
		runOptions.PatchMarker = viper.GetString("marker")
	}

	if viper.IsSet("file") {
		runOptions.TargetPath = viper.GetString("file")
	}

	var targetFiles []string

	if len(runOptions.TargetPath) > 0 {
		targetFiles = append(targetFiles, runOptions.TargetPath)
	}

	if len(runOptions.TargetDir) > 0 {
		err := filepath.Walk(runOptions.TargetDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			targetFiles = append(targetFiles, path)
			return nil
		})

		if err != nil {
			log.Fatalf("failed to gather files to preprocess: %s", err)
		}
	}

	engineArgs := engines.Available.GetArgsFor(engineName)
	if engineArgs == nil {
		log.Fatalf("failed to fetch args for engine %s", engineName)
	}

	engineArgs.InitDefault()
	for name, value := range rawEngineCLIArgs {
		arg := engineArgs.GetByName(name)
		if arg == nil {
			log.Printf("unknown arg '%s' for engine %s", name, engineName)
			continue
		}

		arg.Value = value
	}
	for _, targetFile := range targetFiles {
		eng, err := engineutils.NewEngine(strings.ToLower(engineName), engineArgs)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Processing %s\n", color.New(color.FgWhite, color.Bold).Sprintf(targetFile))
		runOptions.TargetPath = targetFile

		patched, err := eng.Run(runOptions)
		if err != nil {
			log.Error(err)
			continue
		}

		if viper.GetBool("write") {
			err = eng.WritePatchedFile(patched)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			fmt.Println(string(patched))
		}
	}
}
