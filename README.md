# Compare

Compare creates a static report based on a given directory structure containing images. The images are presented by their category, these correspond to the directory that they were found. Currently only the directory structure shown below is supported, deeper nesting as well as files in the main directory will be ignored.

![Overview](media/overview.png?raw=true "Comparison Overview")
![Detail](media/detail.png?raw=true "Comparison Detail")

## Usage

To use the report generator simply move to the localcation of the images (to use relative paths which are important for the html report) and call compare with the following command `compare [flags] path/to/images`.

The flags can be any combination of:

- --title, -t: to specify a report title
- --html-out: to specify the output path for html reports
- --reporter, -f: to specify a reporter to use (currently only html)
- --ext-pattern: to specify an alternative file extension pattern
- --walker, -w: to specity a filesystem walker (currently only local)

## Directory structure

The directory structure for report generation must follow the structure shown below. The main report directory must contain subdirectories for each comparison category. The category directories may have one or many images that should be compared.

```
MyAppComponents
|
+--UserNameView
|  |  Device1.png
|  |  Device2.png
|  |  ...
|
+--PasswordView
|  |  Device1.png
|  |  Device2.png
|  |  ...
|
...
```
## Reporters

Currently the only supported reporter type is the html reporter that generates a static html report. The html reporter creates a simple responsive html page adding a row of images for each category. Images are loaded lazily to prevent lags during loading. Clicking on one of the images shows an enlarged version for better inspection.

## Filesystem walkers

Currently the only supported filewalker type is local. The local filewalker walks over all files matching a given pattern in the first level subdirectories of the given main directory.