# Wintercloak
Wintercloak is a modular file obfuscator. Both a CLI program (`wclk`) and a UI app is available (`wintercloak`).

<img src="https://raw.githubusercontent.com/ma111e/wintercloak/main/readme/preview_home.png" />

## UI
You can switch between the preprocessor and the playground by clicking on the selector in the bottom left corner of the application after selecting an engine.

### Preprocessor
The preprocessor will apply the configured engine to every eligible file within the selected directory.

<img src="https://raw.githubusercontent.com/ma111e/wintercloak/main/readme/preview_preproc.png" />

### Playground
Drag and drop a file to load it onto the playground. The loaded file will be processed by the selected engine when the RUN button is clicked. This button can be used repeatedly to nest obfuscation layers.

<img src="https://raw.githubusercontent.com/ma111e/wintercloak/main/readme/preview_playground.png" />

## CLI
```
Usage:
  wclk [flags]

Flags:
  -e, --engine string           Obfuscator engine to use
  -f, --file string             Source file
  -h, --help                    help for wclk
  -L, --list-data               List all available files that can be used as data sources
  -E, --list-engines            List all available obfuscation engines
  -m, --marker string           Patch marker to use to target specific string with obfuscation (e.g. '-m PATCH' => will replace myvar = PATCH('mystring') with myvar = <obfuscated mystring value>
  -o, --option CLIDataMapping   Pass options to the engine. This flag can be repeated and accepts key-value pairs in key=value format (e.g. '-o 'pattern=??' -o offset=7' (default {})
  -r, --run CLIDataMapping      Set parameters for the run. This flag can be repeated and accepts key-value pairs in key=value format (e.g. '-r useGoBuildTag=true -r goBuildTag=obfuscate' (default {})
  -w, --write                   Write the generated file to disk instead of sending it to stdout. The generated filename is handled by the engine; however the encouraged behaviour is to generate a xxx_gen.yy file next to the source
```

### Examples

Process all the `.go` files in the `./test/data` directory with the `obfuscate` build constraint with the `mask` engine configured to use a repeating pattern of `0x37`
```
./wclk -e mask -r targetDir=./test/data/ -o pattern=0x37 -r goBuildTag=obfuscate -r useGoBuildTag=true
```

Process the `./test/data/main.go` file with the `mask` engine using default parameters and write the obfuscated file as `/test/data/main_gen.go`
```
./wclk -e mask -r targetPath=./test/data/main.go -w
```

## Current state
While I intend to release documentation and most of the obfuscator engines seen on the screenshots, it needs time, which I cannot afford at the moment. PR are welcome.

You should be able to use the `mask` engine as an example to implement your own obfuscator, at least regarding Go source file processing.
