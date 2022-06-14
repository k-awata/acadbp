# acadbp

Batch processing utility using accoreconsole in AutoCAD or LT

## Features

- Bulk conversion of drawing files to specified version of DWG or DXF files
- Running the same script file for each input drawing file
- Publishing drawing files using specified page setup

## Requirements

AutoCAD 2013 or AutoCAD LT 2013 for Windows, or later version (publish command requires 2015 or later)

## Installation

1. If you're using Go:

   ```bat
   go install github.com/k-awata/acadbp@latest
   ```

2. You need to make a `.acadbp.yaml` file in the home directory to specify the path to `accoreconsole.exe` in your AutoCAD installed directory. You can make it with the following command:

   ```bat
   acadbp initconf "C:\Program Files\Autodesk\<your-acad-installed-dir>\accoreconsole.exe"
   ```

## Usage

- Converting all DXF files to DWG in the current directory

  ```bat
  acadbp dwgout *.dxf
  ```

- Downgrading all DWG files to the 2010 format to in the current directory

  ```bat
  acadbp dwgout --format 2010 *.dwg
  ```

- Converting all DWG files to DXF in the current directory

  ```bat
  acadbp dxfout *.dwg
  ```

- Running the script file `example.scr` for each DWG file in the current directory

  ```bat
  acadbp script example.scr *.dwg
  ```

- Publishing PDF files from all DXF files in the current directory by using page setup `Setup1` of `setup.dwg`

  ```bat
  acadbp publish --setup-file setup.dwg --setup-name Setup1 *.dxf
  ```

- Publishing multi sheet PDF

  ```bat
  acadbp publish --setup-file setup.dwg --setup-name Setup1 --multi multi-sheet.pdf *.dxf
  ```

## License

[MIT License](LICENSE)
